package h3m

import (
	"github.com/ImEagle/h3_go/pkg/h3m/models"
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
			name:    "test_maps/test_1.h3m",
			args:    args{fileName: "test_maps/test_1.h3m"},
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

func Test_loadPlayersData(t *testing.T) {
	type args struct {
		mapFileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		players []models.Player
	}{
		{
			name:    "test_1.h3m",
			args:    args{mapFileName: "test_maps/test_1.h3m"},
			wantErr: false,
			players: []models.Player{
				{
					SPlayer: &models.SPlayer{
						PlayableByHuman:    true,
						PlayableByComputer: true,
						Behavior:           0,
						SetTowns:           0,
						Towns:              256,
						HasRandomTown:      false,
						HasMainTown:        true,
					},
					TownCoordinates: &models.TownCoordinates{
						CreateHero: true,
						TownType:   models.TownConflux,
						X:          25,
						Y:          19,
						Z:          0,
					},
				},
				{
					SPlayer: &models.SPlayer{
						PlayableByHuman:    true,
						PlayableByComputer: true,
						Behavior:           0,
						SetTowns:           0,
						Towns:              2,
						HasRandomTown:      false,
						HasMainTown:        true,
					},
					TownCoordinates: &models.TownCoordinates{
						CreateHero: true,
						TownType:   models.TownRampart,
						X:          17,
						Y:          7,
						Z:          0,
					},
				},
				{
					SPlayer: &models.SPlayer{
						PlayableByHuman:    false,
						PlayableByComputer: true,
						Behavior:           0,
						SetTowns:           0,
						Towns:              32,
						HasRandomTown:      false,
						HasMainTown:        true,
					},
					TownCoordinates: &models.TownCoordinates{
						CreateHero: true,
						TownType:   models.TownDungeon,
						X:          8,
						Y:          29,
						Z:          0,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h3m, err := Load(tt.args.mapFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadPlayersData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i, expectedPlayer := range tt.players {
				mapPlayer := h3m.Players[i]

				if mapPlayer.PlayableByHuman != expectedPlayer.PlayableByHuman {
					t.Errorf("loadBasicMapParameters.loadPlayersData.PlayableByHuman (Player %d) got = %v, want %v", i, mapPlayer.PlayableByHuman, expectedPlayer.PlayableByHuman)
					return
				}

				if mapPlayer.PlayableByComputer != expectedPlayer.PlayableByComputer {
					t.Errorf("loadBasicMapParameters.loadPlayersData.PlayableByComputer (Player %d) got = %v, want %v", i, mapPlayer.PlayableByComputer, expectedPlayer.PlayableByComputer)
					return
				}

				if mapPlayer.Behavior != expectedPlayer.Behavior {
					t.Errorf("loadBasicMapParameters.loadPlayersData.Behavior (Player %d) got = %v, want %v", i, mapPlayer.Behavior, expectedPlayer.Behavior)
					return
				}

				if mapPlayer.SetTowns != expectedPlayer.SetTowns {
					t.Errorf("loadBasicMapParameters.loadPlayersData.SetTowns (Player %d) got = %v, want %v", i, mapPlayer.SetTowns, expectedPlayer.SetTowns)
					return
				}

				if mapPlayer.Towns != expectedPlayer.Towns {
					t.Errorf("loadBasicMapParameters.loadPlayersData.Towns (Player %d) got = %v, want %v", i, mapPlayer.Towns, expectedPlayer.Towns)
					return
				}

				if mapPlayer.HasRandomTown != expectedPlayer.HasRandomTown {
					t.Errorf("loadBasicMapParameters.loadPlayersData.HasRandomTown (Player %d) got = %v, want %v", i, mapPlayer.HasRandomTown, expectedPlayer.HasRandomTown)
					return
				}

				if mapPlayer.HasMainTown != expectedPlayer.HasMainTown {
					t.Errorf("loadBasicMapParameters.loadPlayersData.HasMainTown (Player %d) got = %v, want %v", i, mapPlayer.HasMainTown, expectedPlayer.HasMainTown)
					return
				}

				if mapPlayer.TownCoordinates.CreateHero != expectedPlayer.TownCoordinates.CreateHero {
					t.Errorf("loadBasicMapParameters.loadPlayersData.TownCoordinates.CreateHero (Player %d) got = %v, want %v", i, mapPlayer.TownCoordinates.CreateHero, expectedPlayer.TownCoordinates.CreateHero)
					return
				}

				if mapPlayer.TownCoordinates.TownType != expectedPlayer.TownCoordinates.TownType {
					t.Errorf("loadBasicMapParameters.loadPlayersData.TownCoordinates.TownType (Player %d) got = %v, want %v", i, mapPlayer.TownCoordinates.TownType, expectedPlayer.TownCoordinates.TownType)
					return
				}

				if mapPlayer.TownCoordinates.X != expectedPlayer.TownCoordinates.X {
					t.Errorf("loadBasicMapParameters.loadPlayersData.TownCoordinates.X (Player %d) got = %v, want %v", i, mapPlayer.TownCoordinates.X, expectedPlayer.TownCoordinates.X)
					return
				}

				if mapPlayer.TownCoordinates.Y != expectedPlayer.TownCoordinates.Y {
					t.Errorf("loadBasicMapParameters.loadPlayersData.TownCoordinates.Y (Player %d) got = %v, want %v", i, mapPlayer.TownCoordinates.Y, expectedPlayer.TownCoordinates.Y)
					return
				}

				if mapPlayer.TownCoordinates.Z != expectedPlayer.TownCoordinates.Z {
					t.Errorf("loadBasicMapParameters.loadPlayersData.TownCoordinates.Z (Player %d) got = %v, want %v", i, mapPlayer.TownCoordinates.Z, expectedPlayer.TownCoordinates.Z)
					return
				}
			}
		})
	}
}

func Test_loadRumors(t *testing.T) {
	type args struct {
		mapFileName string
	}
	tests := []struct {
		name   string
		args   args
		rumors []models.Rumor
	}{
		{
			name: "test_1.h3m",
			args: args{mapFileName: "test_maps/test_1.h3m"},
			rumors: []models.Rumor{
				{
					Name: "Rumor 1",
					Text: "Example rumor 1",
				},
				{
					Name: "Rumor 2",
					Text: "Example rumor 2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h3m, err := Load(tt.args.mapFileName)
			if err != nil {
				t.Errorf("loadRumors() error = %v", err)
				return
			}

			if len(h3m.Rumors) != len(tt.rumors) {
				t.Errorf("loadRumors() got = %v, want %v", h3m.Rumors, tt.rumors)
				return

			}

			if len(h3m.Rumors) != len(tt.rumors) {
				t.Errorf("loadRumors() got = %v, want %v", h3m.Rumors, tt.rumors)
				return
			}

			for i, expectedRumor := range tt.rumors {
				if h3m.Rumors[i].Name != expectedRumor.Name {
					t.Errorf("loadRumors() got = %v, want %v", h3m.Rumors[i].Name, expectedRumor.Name)
					return
				}

				if h3m.Rumors[i].Text != expectedRumor.Text {
					t.Errorf("loadRumors() got = %v, want %v", h3m.Rumors[i].Text, expectedRumor.Text)
					return
				}
			}
		})
	}
}

func Test_loadMapObjects(t *testing.T) {
	type args struct {
		mapFileName string
	}
	tests := []struct {
		name              string
		args              args
		objectsDefinition []models.MapObjectDefinition
	}{
		{
			name: "mapa-teren-zasoby.h3m",
			args: args{mapFileName: "test_maps/mapa-teren-zasoby.h3m"},
			objectsDefinition: []models.MapObjectDefinition{
				{
					MapObjectData: &models.MapObjectData{
						PassableSquares: [6]byte{},
						ActiveSquare:    [6]byte{},
						Landscape:       [2]byte{},
						LandscapeGroup:  [2]byte{},
						Class:           0,
						Number:          0,
						Group:           0,
						OverOrBelow:     0,
						Unknown:         [16]byte{},
					},
					SpriteName: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h3m, err := Load(tt.args.mapFileName)
			if err != nil {
				t.Errorf("loadRumors() error = %v", err)
				return
			}

			// NOP

			if len(h3m.ObjectsDefinition) != len(tt.objectsDefinition) {
				t.Errorf("loadObjectsDefinition() got = %v, want %v", h3m.ObjectsDefinition, tt.objectsDefinition)
				return
			}
		})
	}
}
