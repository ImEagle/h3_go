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

func Test_loadHeader(t *testing.T) {
	type args struct {
		mapFileName string
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		expectedHeader string
	}{
		{
			name:           "SoD_map",
			args:           args{mapFileName: "test_maps/test_1.h3m"},
			wantErr:        false,
			expectedHeader: "SoD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h3m, err := Load(tt.args.mapFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if h3m.Format != tt.expectedHeader {
				t.Errorf("loadHeader() got = %v, want %v", h3m.Format, tt.expectedHeader)
				return
			}
		})
	}
}

func Test_loadBasicMapParameters(t *testing.T) {
	type basicMapParameters struct {
		HasHeroOnMap   bool
		MapSize        uint32
		HasUnderground bool
		Name           string
		Description    string
	}

	type args struct {
		mapFileName string
	}
	tests := []struct {
		name                string
		args                args
		wantErr             bool
		expectedBasicParams basicMapParameters
	}{
		{
			name:    "test_1.h3m",
			args:    args{mapFileName: "test_maps/test_1.h3m"},
			wantErr: false,
			expectedBasicParams: basicMapParameters{
				HasHeroOnMap:   true,
				MapSize:        36,
				HasUnderground: false,
				Name:           "Test Map 1",
				Description:    "Map Type SoD\nTest Map description 1\nDifficulty level - Low\nNo teams\nNo special wins\nNo special loose",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h3m, err := Load(tt.args.mapFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadBasicMapParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if h3m.HasHeroOnMap != tt.expectedBasicParams.HasHeroOnMap {
				t.Errorf("loadBasicMapParameters.HasHeroOnMap got = %v, want %v", h3m.HasHeroOnMap, tt.expectedBasicParams.HasHeroOnMap)
				return
			}

			if h3m.MapSize != tt.expectedBasicParams.MapSize {
				t.Errorf("loadBasicMapParameters.MapSize got = %v, want %v", h3m.MapSize, tt.expectedBasicParams.MapSize)
				return
			}

			if h3m.HasUnderground != tt.expectedBasicParams.HasUnderground {
				t.Errorf("loadBasicMapParameters.HasUnderground got = %v, want %v", h3m.HasUnderground, tt.expectedBasicParams.HasUnderground)
				return
			}

			if h3m.Name != tt.expectedBasicParams.Name {
				t.Errorf("loadBasicMapParameters.Name got = %v, want %v", h3m.Name, tt.expectedBasicParams.Name)
				return
			}

			if h3m.Description != tt.expectedBasicParams.Description {
				t.Errorf("loadBasicMapParameters.Description got = %v, want %v", h3m.Description, tt.expectedBasicParams.Description)
				return
			}
		})
	}
}
