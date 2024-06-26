package h3m

import "github.com/ImEagle/h3_go/pkg/h3m/models"

type H3m struct {
	Format         string
	HasHeroOnMap   bool
	MapSize        uint32
	HasUnderground bool
	Name           string
	Description    string
	Difficulty     byte

	Players           []*models.Player
	VictoryCondition  *models.VictoryCondition
	LossCondition     *models.LoseCondition
	TeamColors        *models.TeamColors
	CustomHeroes      []*models.CustomHeroes
	Artifacts         [18]byte
	Rumors            []*models.Rumor
	Heroes            []*models.Hero
	LandMap           []models.MapTile
	UndergroundMap    []models.MapTile
	ObjectsDefinition []*models.MapObjectDefinition
}
