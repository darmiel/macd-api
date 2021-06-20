package csv

import (
	"reflect"
	"testing"
)

func Test_escapedRuneSplit(t *testing.T) {
	type args struct {
		str string
		sep rune
	}
	tests := []struct {
		name     string
		args     args
		wantResp []string
	}{
		{"Normal", args{"Hello|World", '|'}, []string{"Hello", "World"}},
		{"Empty Input", args{"", '|'}, []string{}},
		{"Escaped", args{"Hello\\|World", '|'}, []string{"Hello|World"}},
		{"Multi", args{"Hello|World|What's|Up?", '|'}, []string{"Hello", "World", "What's", "Up?"}},
		{"Open End", args{"Hello|World|What's|Up?|", '|'}, []string{"Hello", "World", "What's", "Up?", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResp := escapedRuneSplit(tt.args.str, tt.args.sep); !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("escapedRuneSplit() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestParseUnmarshal(t *testing.T) {
	const data = `Name|Age
Max|20
John|31
Jim|420`

	type T struct {
		Name string `csv:"Name"`
		Age  int    `csv:"Age"`
	}

	expected := []*T{{"Max", 20}, {"John", 31}, {"Jim", 420}}

	parse, err := Parse([]byte(data), &ParseOptions{
		Separator: '|',
	})
	if err != nil {
		t.Error(err)
		return
	}

	var out []*T
	for _, row := range parse.Rows {
		trg := new(T)
		if err = row.Unmarshal(trg); err != nil {
			t.Error(err)
			return
		}
		out = append(out, trg)
	}

	if !reflect.DeepEqual(expected, out) {
		t.Errorf("out = %v, want expected = %v", out, expected)
	}
}
