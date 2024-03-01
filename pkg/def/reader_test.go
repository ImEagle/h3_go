package def

import (
	"github.com/ImEagle/h3_go/pkg/lod"
	"testing"
)

func TestReader_Read(t *testing.T) {
	type fields struct {
		fileName string
	}
	type args struct {
		data []byte
	}

	lodReader := lod.NewReader("/Users/eagle/Downloads/h3/H3sprite.lod")
	lodReader.LoadMetadata()
	defPayload, err := lodReader.GetFile("ab03_.def")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		args    args
		want    *header
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{data: defPayload},
			want: &header{
				Type:        0x00000000,
				Width:       0x00000020,
				Height:      0x00000020,
				BlocksCount: 0x00000001,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader()
			err := r.Load(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if r.CanGenerateSpriteSheet() {
				// TODO
			}

		})
	}
}
