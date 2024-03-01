package lod

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReader_LoadMetadata(t *testing.T) {
	type fields struct {
		fileName string
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "example",
			fields: fields{
				fileName: "/Users/eagle/Downloads/h3/h3ab_bmp.lod",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.fields.fileName)

			if got := r.LoadMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadMetadata() = %v, want %v", got, tt.want)
			}

			data, err := r.GetFile("sgelwa5.pcx")
			//_, _ = r.GetFile("ara_cobl.pcx")
			fmt.Print(data, err)
		})
	}
}
