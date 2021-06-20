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
