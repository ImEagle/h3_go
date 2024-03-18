package lod

import (
	"fmt"
	"github.com/ImEagle/h3_go/pkg/def"
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

func TestReader_loadTerrain(t *testing.T) {
	tests := []struct {
		name           string
		spriteFileName string
		assetName      string
		wantError      bool
	}{
		{
			name:           "example",
			spriteFileName: "/Users/eagle/projects/h3_go/pkg/h3m/test_maps/dirt_only.h3m",
			assetName:      "sandtl.def",
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.spriteFileName)
			r.LoadMetadata()

			imgDetails, err := r.GetFile(tt.assetName)
			if (err != nil) != tt.wantError {
				t.Errorf("loadTerrain() error = %v, wantError %v", err, tt.wantError)
				return
			}

			dr := def.NewReader()
			dr.Load(imgDetails)

		})
	}
}
