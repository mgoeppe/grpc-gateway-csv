package csv

import (
	"github.com/kylelemons/godebug/pretty"
	"testing"
)

type outer struct {
	col1  string
	col2  string
	s     []string    `csv:"slice"`
	m     map[int]int `csv:"map"`
	inner inner
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
			m:    map[int]int{1: 2, 3: 4},
			inner: inner{
				col3: "c",
				col4: 6,
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
			name: "deep structure",
			v:    v,
			want: "col1;col2;slice;map;col3;col4\na;b;a|b;1:2|3:4;c;6\ne;f;a|b|c;;g;6\n",
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
