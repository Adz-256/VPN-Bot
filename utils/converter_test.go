package utils

import (
	"reflect"
	"testing"
	"time"
)

type testStruct struct {
	ID    int       `db:"id"`
	Name  string    `db:"name"`
	Empty time.Time `db:"empty,omitempty"`
}

func TestStructToMap(t *testing.T) {
	type args struct {
		s    any
		full bool
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{name: "TestStructToMap", args: args{s: testStruct{
			ID:   12,
			Name: "jhon",
		}, full: true}, want: map[string]any{"id": 12, "name": "jhon"}, wantErr: false},

		{name: "TestStructToMapFullFalse", args: args{s: testStruct{
			Name: "jhon",
		}, full: false}, want: map[string]any{"name": "jhon"}, wantErr: false},

		{name: "TestStructToMapFullTrueError", args: args{s: testStruct{
			Name: "jhon",
		}, full: true}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StructToMap(tt.args.s, tt.args.full)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
