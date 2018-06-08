package commander

import (
	"reflect"
	"testing"
)

func TestArgument_GetValue(t *testing.T) {
	type fields struct {
		Name          string
		Type          string
		OriginalValue string
	}
	tests := []struct {
		name      string
		fields    fields
		parameter string
		want      interface{}
	}{
		{
			name: "String",
			fields: fields{
				Type: "String",
			},
			parameter: "testname",
			want:      "testname",
		},
		{
			name: "Int64",
			fields: fields{
				Type: "Int64",
			},
			parameter: "42",
			want:      int64(42),
		},
		{
			name: "-Int64",
			fields: fields{
				Type: "Int64",
			},
			parameter: "-42",
			want:      int64(-42),
		},
		{
			name: "Uint64",
			fields: fields{
				Type: "Uint64",
			},
			parameter: "42",
			want:      uint64(42),
		},
		{
			name: "-Uint64",
			fields: fields{
				Type: "Uint64",
			},
			parameter: "-42",
			want:      uint64(0),
		},
		{
			name: "FilePath [exists]",
			fields: fields{
				Type: "FilePath",
			},
			parameter: "argument_test.go",
			want:      "argument_test.go",
		},
		{
			name: "FilePath [not exists]",
			fields: fields{
				Type: "FilePath",
			},
			parameter: "/ajshdjkashdjashdjasd",
			want:      "",
		},
		{
			name: "StringArray[]",
			fields: fields{
				Type: "StringArray[]",
			},
			parameter: "one,two,three",
			want:      []string{"one", "two", "three"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Argument{
				Name:          tt.fields.Name,
				Type:          tt.fields.Type,
				OriginalValue: tt.fields.OriginalValue,
			}
			a.SetValue(tt.parameter)
			if got := a.Value; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Argument.Value = %v, want %v", got, tt.want)
			}
		})
	}
}
