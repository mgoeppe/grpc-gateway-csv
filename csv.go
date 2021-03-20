package csv

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Marshaler struct {
	runtime.Marshaler

	// RowDelim specifies the delimiter between rows
	RowDelim string
	// FieldDelim specifies the delimiter between fields
	FieldDelim string
	// InnerDelim specifies the delimiter of merged values (slices / maps) within one field.
	InnerDelim string

	// NoHeader suppresses to render the header
	NoHeader bool
}

func (m *Marshaler) initDefaults() {
	if m.RowDelim == "" {
		m.RowDelim = "\n"
	}
	if m.FieldDelim == "" {
		m.FieldDelim = ";"
	}
	if m.InnerDelim == "" {
		m.InnerDelim = "|"
	}
}

// Marshal renders the structure in i as CSV.
//
// If i is a slice or a struct that contains slices each slice
// is rendered to a CSV block. These top-level slices need to
// contain structs otherwise an error is returned.
// Each top-level slices is rendered to one block. The blocks
// are delimited by '---\n'. Empty slices or nil pointers are ignored.
//
// Each csv block consists of a header (if NoHeader option is false) and
// multiple rows delimited by m.RowDelimi. Each row is a 'flat'
// representation of the corresponding slice elements:
//  * struct fields are visible on top-level with own header delimited by m.FieldDelim
//  * nested slices / maps are flatened delimited by m.InnerDelim
func (m *Marshaler) Marshal(i interface{}) ([]byte, error) {
	m.initDefaults()

	v := reflect.ValueOf(i)
	v = followPtr(v)

	slices := []string{}
	var err error
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			v := v.Field(i)
			if v.Kind() == reflect.Slice {
				slices, err = m.marshalSliceAndAppend(slices, v)
				if err != nil {
					return nil, err
				}
			}
		}
	case reflect.Slice:
		slices, err = m.marshalSliceAndAppend(slices, v)
		if err != nil {
			return nil, err
		}
	}
	return []byte(strings.Join(slices, "---\n")), nil

}

func (m *Marshaler) marshalSliceAndAppend(slices []string, v reflect.Value) ([]string, error) {
	s, err := m.marshalSlice(v)
	if err != nil {
		return nil, err
	}
	if s != "" {
		return append(slices, s), nil
	}
	return slices, nil
}

type visit struct {
	addr uintptr
	typ  reflect.Type
	next *visit
}

func (m *Marshaler) marshalSlice(v reflect.Value) (string, error) {
	if v.IsNil() || v.Len() == 0 {
		return "", nil
	}

	res := ""
	// header
	first := followPtr(v.Index(0))
	if first.Type().Kind() != reflect.Struct {
		return "", fmt.Errorf("top-level slice with non struct type: %s", v.Index(0).Type().Kind())
	}
	if !m.NoHeader {
		header := m.marshal(first, true, map[uintptr]*visit{})
		res = res + fmt.Sprintf("%s%s", strings.Join(header, m.FieldDelim), m.RowDelim)
	}

	for i := 0; i < v.Len(); i++ {
		row := m.marshal(v.Index(i), false, map[uintptr]*visit{})
		res = res + fmt.Sprintf("%s%s", strings.Join(row, m.FieldDelim), m.RowDelim)
	}
	return res, nil
}

func (m *Marshaler) marshal(v reflect.Value, header bool, visited map[uintptr]*visit) []string {
	res := []string{}
	v = followPtr(v)

	// break recursion
	if v.CanAddr() {
		addr := v.UnsafeAddr()
		typ := v.Type()
		seen := visited[addr]
		for p := seen; p != nil; p = p.next {
			if p.addr == addr && p.typ == typ {
				return res
			}
		}
		visited[addr] = &visit{addr, typ, seen}
	}

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)
		typ := v.Type().Field(i)
		// old proto version uses XXX_ as field name prefix for internal stuff
		if strings.ToLower(string(typ.Name[0])) == string(typ.Name[0]) || strings.HasPrefix(typ.Name, "XXX_") {
			continue
		}

		switch val.Kind() {
		case reflect.Map:
			if header {
				res = append(res, name(typ))
			} else {
				s := []string{}
				for _, k := range val.MapKeys() {
					// k: struct keys are not supported so far
					if val.MapIndex(k).Kind() == reflect.Struct {
						s = append(s, fmt.Sprintf("%s:%s",
							fmt.Sprintf("%v", k),
							strings.Join(m.marshal(val.MapIndex(k), header, visited), m.InnerDelim),
						))
					} else {
						s = append(s, fmt.Sprintf("%s:%s",
							fmt.Sprintf("%v", k),
							fmt.Sprintf("%v", val.MapIndex(k)),
						))
					}
				}
				res = append(res, strings.Join(s, m.InnerDelim))
			}
		case reflect.Slice:
			if header {
				res = append(res, name(typ))
			} else {
				s := []string{}
				for j := 0; j < val.Len(); j++ {
					if val.Index(j).Kind() == reflect.Struct {
						s = append(s, m.marshal(val.Index(j), header, visited)...)
					} else {
						s = append(s, fmt.Sprintf("%v", val.Index(j)))
					}
				}
				res = append(res, strings.Join(s, m.InnerDelim))
			}
		case reflect.Struct:
			res = append(res, m.marshal(val, header, visited)...)
		case reflect.Ptr:
			if val.Elem().Kind() == reflect.Struct {
				res = append(res, m.marshal(val, header, visited)...)
			}
		default:
			if header {
				res = append(res, name(typ))
			} else {
				res = append(res, fmt.Sprintf("%v", val))
			}
		}
	}
	return res
}

// name evaluates field name to use when marshaling. The following order applies:
// 1. csv tag
// 2. field name
func name(f reflect.StructField) string {
	n := f.Tag.Get("csv")
	if n != "" {
		return n
	}
	return f.Name

}

func followPtr(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

// ContentType returns 'text/csv'.
func (m *Marshaler) ContentType(v interface{}) string {
	return "text/csv"
}
