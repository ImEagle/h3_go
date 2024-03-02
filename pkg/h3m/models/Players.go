package models

type TownCoordinates struct {
	CreateHero bool
	TownType   uint8
	X          uint8
	Y          uint8
	Z          uint8
}

type SMainHero struct {
	HaveRandomHero bool
	HeroType       uint8
	HeroFace       uint8
}

type HeroDetails struct {
	Id   uint8
	Name string
}

type MainHero struct {
	*SMainHero
	Name    string
	Unknown uint8

	Heroes []*HeroDetails
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
