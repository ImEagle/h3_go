package models

const (
	TownCastle     = 0
	TownRampart    = 1
	TownTower      = 2
	TownInferno    = 3
	TownNecropolis = 4
	TownDungeon    = 5
	TownStronghold = 6
	TownFortress   = 7
	TownConflux    = 8 // Armageddon's Blade

	// Different game version?
	// TownCove    = 9  // Horn of the Abyss
	// TownFactory = 10 // Horn of the Abyss
)

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
}

type HeroDetails struct {
	Id   uint8
	Name string
}

type MainHero struct {
	Name string

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
