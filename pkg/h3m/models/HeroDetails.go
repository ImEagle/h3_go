package models

type SecondarySkill struct {
	Id    uint8
	Level uint8
}

type ArtifactId struct {
	Id uint16
}

type ArtifactDetails struct {
	HeadWear  uint16
	Shoulders uint16
	Neck      uint16
	RightHand uint16
	LeftHand  uint16
	Torso     uint16
	RightRing uint16
	LeftRing  uint16
	Legs      uint16
	Misc1     uint16
	Misc2     uint16
	Misc3     uint16
	Misc4     uint16
	Device1   uint16
	Device2   uint16
	Device3   uint16
	Device4   uint16
	Spellbook uint16
	Misc5     uint16
	Backpack  []*ArtifactId
}

type Hero struct {
	Experience        uint32
	SecondarySkillSet bool
	SecondarySkills   []*SecondarySkill
	HasArtifacts      bool
	ArtifactsDetails  ArtifactDetails
	HasBiography      bool
	BiographyLen      uint32
	Biography         string
	HasCustomGender   bool
	Gender            uint8
	HasSpells         bool
	Spells            []byte
	HasPrimarySkills  bool
	PrimaryAttack     uint8
	PrimaryDefence    uint8
	PrimarySpellPower uint8
	PrimaryKnowledge  uint8
}
