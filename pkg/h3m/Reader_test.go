package h3m

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *H3m
		wantErr bool
	}{
		//{
		//	name:    "example_map.h3m",
		//	args:    args{fileName: "example_map.h3m"},
		//	want:    &H3m{Format: "H3M"},
		//	wantErr: false,
		//},
		{
			name:    "generated.h3m",
			args:    args{fileName: "generated.h3m"},
			want:    &H3m{Format: "H3M"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
