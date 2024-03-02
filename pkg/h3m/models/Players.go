package models

type TownCoordinates struct {
	CreateHero bool
	TownType   uint8
	X          uint8
	Y          uint8
	Z          uint8
}

type Player struct {
	*SPlayer
	TownCoordinates *TownCoordinates
}

type SPlayer struct {
	//HeroesMasteryLevelCap uint8
	PlayableByHuman    bool
	PlayableByComputer bool
	Behavior           uint8
	SetTowns           uint8
	Towns              uint16
	HasRandomTown      bool
	HasMainTown        bool
}
