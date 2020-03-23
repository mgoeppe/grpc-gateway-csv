package csv

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

type outer struct {
	col1       string
	col2       string
	s          []string      `csv:"slice"`
	m1         map[int]int   `csv:"map1"`
	m2         map[int]inner `csv:"map2"`
	XXX_ignore string
	inner      inner
	innerSlice []inner
}

type inner struct {
	col3  string
	col4  int
	outer *outer
}

func TestMarshaler_Marshal(t *testing.T) {

	v := []outer{
		{
			col1: "a",
			col2: "b",
			s:    []string{"a", "b"},
			m1:   map[int]int{1: 2},
			m2: map[int]inner{
				1: {
					col3: "x",
					col4: 7,
				},
			},
			inner: inner{
				col3: "c",
				col4: 6,
			},
			innerSlice: []inner{
				{
					col3: "u",
					col4: 8,
				},
				{
					col3: "v",
					col4: 9,
				},
			},
		},
		{
			col1: "e",
			col2: "f",
			s:    []string{"a", "b", "c"},
			inner: inner{
				col3: "g",
				col4: 6,
			},
		},
	}
	v[1].inner.outer = &v[1] // introduce cycle

	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			name: "multiple empty slices",
			v: struct {
				a []string
				b []string
			}{
				a: nil,
				b: nil,
			},
			want: "",
		},
		{
			name: "multiple slices w/o struct",
			v: struct {
				a []string
				b []string
			}{
				a: []string{"a", "b"},
				b: []string{"c", "d"},
			},
			wantErr: true,
		},
		{
			name: "multiple slices w/ struct",
			v: struct {
				a []inner
				b []outer
			}{
				a: []inner{
					{
						col3: "a",
						col4: 1,
					},
					{
						col3: "b",
						col4: 2,
					},
				},
				b: []outer{
					{
						col1: "c",
						col2: "d",
					},
					{
						col1: "e",
						col2: "f",
					},
				},
			},
			want: "col3;col4\na;1\nb;2\n---\ncol1;col2;slice;map1;map2;col3;col4;innerSlice\nc;d;;;;;0;\ne;f;;;;;0;\n",
		},
		{
			name: "deep structure",
			v:    v,
			want: "col1;col2;slice;map1;map2;col3;col4;innerSlice\na;b;a|b;1:2;1:x|7;c;6;u|8|v|9\ne;f;a|b|c;;;g;6;\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Marshaler{}
			g, err := m.Marshal(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVMarshaler.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := string(g)
			if diff := pretty.Compare(got, tt.want); diff != "" {
				t.Errorf("CSVMarshaler.Marshal() generate unexpected results:\n%s", diff)
			}
		})
	}
}

func TestMarshaler_ContentType(t *testing.T) {
	want := "text/csv"
	m := &Marshaler{}
	if got := m.ContentType(); got != want {
		t.Errorf("Marshaler.ContentType() = %v, want %v", got, want)
	}
}
