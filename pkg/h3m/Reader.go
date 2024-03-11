package h3m

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/ImEagle/h3_go/pkg/h3m/models"
	"io"
	"os"
)

const (
	HeaderRoE uint32 = 14
	HeaderAB  uint32 = 21
	HeaderSOD uint32 = 28

	HeaderRoEName = "RoE"
	HeaderABName  = "AB"
	HeaderSODName = "SoD"
)

var headerNames = map[uint32]string{
	HeaderRoE: HeaderRoEName,
	HeaderAB:  HeaderABName,
	HeaderSOD: HeaderSODName,
}

func Load(fileName string) (*H3m, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decompressedMap, err := decompressGZIP(f)
	if err != nil {
		return nil, err
	}

	h3m := &H3m{}

	// ----- Load header -----
	var header uint32
	binary.Read(decompressedMap, binary.LittleEndian, &header)

	h3m.Format = headerNames[header]

	err = loadBasicMapParameters(decompressedMap, h3m)
	if err != nil {
		return nil, err
	}
	// ----- ~ Load header ~ -----

	// ----- Load players data -----
	err = loadPlayersData(decompressedMap, h3m)
	if err != nil {
		return nil, err
	}

	// ----- ~ Load players data ~ -----

	// ----- Load victory condition -----
	err = loadVictoryCondition(decompressedMap, h3m)
	// ----- ~ Load victory condition ~ -----

	// ----- Load loss condition -----
	err = loadLossCondition(decompressedMap, h3m)
	// ----- ~ Load loss condition ~ -----

	// ----- Load teams -----
	err = loadTeams(decompressedMap, h3m)
	// ----- ~ Load teams ~ -----

	// ----- Load available heroes -----
	err = loadAvailableHeroes(decompressedMap, h3m)
	// ----- ~ Load available heroes ~ -----

	// ----- NOP -----
	// 4 empty bytes
	decompressedMap.Seek(4, io.SeekCurrent)
	// ----- ~ NOP ~ -----

	// ----- Load custom heroes -----
	err = loadCustomHeroes(decompressedMap, h3m)
	// ----- ~ Load custom heroes ~ -----

	// ----- NOP -----
	// 31 empty bytes
	decompressedMap.Seek(31, io.SeekCurrent)
	// ----- ~ NOP ~ -----

	// ----- Load random artifacts -----
	err = loadRandomArtifacts(decompressedMap, h3m)
	// ----- ~ Load random artifacts ~ -----

	// ----- Rumors -----
	err = loadRumors(decompressedMap, h3m)
	// ----- ~ Rumors ~ -----

	// ----- Hero settings -----
	err = loadHeroSettings(decompressedMap, h3m)
	// ----- ~ Hero settings ~ -----

	return h3m, nil
}

func loadHeroSettings(decompressedMap io.ReadSeeker, m *H3m) error {
	var heroDetails []*models.Hero

	const HeroesCount = 156
	for i := 0; i < HeroesCount; i++ {
		var customHero bool
		binary.Read(decompressedMap, binary.LittleEndian, &customHero)
		if !customHero {
			continue
		}

		var hd models.Hero

		var customExperience bool
		binary.Read(decompressedMap, binary.LittleEndian, &customExperience)
		if customExperience {
			binary.Read(decompressedMap, binary.LittleEndian, &hd.Experience)
		}

		var customSecondarySkills bool
		binary.Read(decompressedMap, binary.LittleEndian, &customSecondarySkills)

		// Secondary skills
		if customSecondarySkills {
			var secondarySkillsCount uint32
			binary.Read(decompressedMap, binary.LittleEndian, &secondarySkillsCount)

			for j := uint32(0); j < secondarySkillsCount; j++ {
				var secondarySkill models.SecondarySkill
				binary.Read(decompressedMap, binary.LittleEndian, &secondarySkill)
				hd.SecondarySkills = append(hd.SecondarySkills, &secondarySkill)
			}
		}

		// Artifacts
		var customArtifacts bool
		binary.Read(decompressedMap, binary.LittleEndian, &customArtifacts)
		if customArtifacts {
			binary.Read(decompressedMap, binary.LittleEndian, &hd.ArtifactsDetails)

			// Backpack
			var backpackCount uint16
			binary.Read(decompressedMap, binary.LittleEndian, &backpackCount)

			for j := uint16(0); j < backpackCount; j++ {
				var art models.ArtifactId
				binary.Read(decompressedMap, binary.LittleEndian, &art)

				hd.ArtifactsDetails.Backpack = append(hd.ArtifactsDetails.Backpack, &art)
			}
		}

		// Biography
		var customBiography bool
		binary.Read(decompressedMap, binary.LittleEndian, &customBiography)

		if customBiography {
			var err error
			hd.Biography, err = readString(decompressedMap)
			if err != nil {
				return nil
			}
		}

		// Gender
		var customGender bool
		binary.Read(decompressedMap, binary.LittleEndian, &customGender)

		if customGender {
			binary.Read(decompressedMap, binary.LittleEndian, &hd.Gender)
		}

		// Spells
		var customSpells bool
		binary.Read(decompressedMap, binary.LittleEndian, &customSpells)

		if customSpells {
			hd.Spells = make([]byte, 9)
			binary.Read(decompressedMap, binary.LittleEndian, &hd.Spells)
		}

		// Primary skills
		var customPrimarySkills bool
		binary.Read(decompressedMap, binary.LittleEndian, &customPrimarySkills)

		if customPrimarySkills {
			binary.Read(decompressedMap, binary.LittleEndian, &hd.PrimaryAttack)
			binary.Read(decompressedMap, binary.LittleEndian, &hd.PrimaryDefence)
			binary.Read(decompressedMap, binary.LittleEndian, &hd.PrimarySpellPower)
			binary.Read(decompressedMap, binary.LittleEndian, &hd.PrimaryKnowledge)
		}

		heroDetails = append(heroDetails, &hd)
	}

	m.Heroes = heroDetails

	return nil
}

func loadRumors(decompressedMap io.ReadSeeker, m *H3m) error {
	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Rumors: %d\n", currentOffset)

	var rumorsCount uint32
	err = binary.Read(decompressedMap, binary.LittleEndian, &rumorsCount)

	for i := uint32(0); i < rumorsCount; i++ {
		var rumor models.Rumor

		rumor.Name, err = readString(decompressedMap)
		if err != nil {
			return err
		}

		rumor.Text, err = readString(decompressedMap)
		if err != nil {
			return err
		}

		m.Rumors = append(m.Rumors, &rumor)

	}

	return nil
}

func loadRandomArtifacts(decompressedMap io.ReadSeeker, m *H3m) error {
	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Random artifacts: %d\n", currentOffset)

	// 18 bytes for artifacts
	return binary.Read(decompressedMap, binary.LittleEndian, &m.Artifacts)
}

func loadCustomHeroes(decompressedMap io.ReadSeeker, m *H3m) error {
	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Custom heroes offset: %d\n", currentOffset)

	var customHeroesCount uint8
	err = binary.Read(decompressedMap, binary.LittleEndian, &customHeroesCount)

	for i := uint8(0); i < customHeroesCount; i++ {
		var customHero models.CustomHeroes

		err := binary.Read(decompressedMap, binary.LittleEndian, &customHero.Id)
		if err != nil {
			return err
		}

		err = binary.Read(decompressedMap, binary.LittleEndian, &customHero.Portrait)
		if err != nil {
			return err
		}

		customHero.Name, err = readString(decompressedMap)
		if err != nil {
			return err
		}

		err = binary.Read(decompressedMap, binary.LittleEndian, &customHero.CanBeHired)
		if err != nil {
			return err
		}

		m.CustomHeroes = append(m.CustomHeroes, &customHero)
	}

	return nil
}

func loadAvailableHeroes(decompressedMap io.ReadSeeker, m *H3m) error {
	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Teams offset: %d\n", currentOffset)

	var availableHeroes models.AvailableHeroes

	err = binary.Read(decompressedMap, binary.LittleEndian, &availableHeroes)
	if err != nil {
		return err
	}
	return nil
}

func loadTeams(decompressedMap io.ReadSeeker, m *H3m) error {
	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Teams offset: %d\n", currentOffset)

	var numberOfTeams uint8
	err = binary.Read(decompressedMap, binary.LittleEndian, &numberOfTeams)
	if err != nil {
		return err
	}

	if numberOfTeams == 0 {
		return nil
	}

	var teamColors models.TeamColors
	err = binary.Read(decompressedMap, binary.LittleEndian, &teamColors.Red)
	if err != nil {
		return err
	}

	m.TeamColors = &teamColors

	return nil
}

func loadVictoryCondition(decompressedMap io.ReadSeeker, m *H3m) error {
	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Victory condition offset: %d\n", currentOffset)

	var victoryCondition models.VictoryCondition
	binary.Read(decompressedMap, binary.LittleEndian, &victoryCondition.Type)

	if victoryCondition.Type != 255 {
		// TODO: victoryCondition
		panic("Victory conditions not implemented")
	}

	m.VictoryCondition = &victoryCondition

	return nil
}

func loadLossCondition(decompressedMap io.ReadSeeker, m *H3m) error {
	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Loss condition offset: %d\n", currentOffset)

	var lossCondition models.LoseCondition
	binary.Read(decompressedMap, binary.LittleEndian, &lossCondition.Type)

	if lossCondition.Type != 255 {
		panic("Loss conditions not implemented")
	}

	m.LossCondition = &lossCondition

	return nil
}

func loadBasicMapParameters(decompressedMap io.Reader, h3m *H3m) error {
	var err error

	binary.Read(decompressedMap, binary.LittleEndian, &h3m.HasHeroOnMap)
	binary.Read(decompressedMap, binary.LittleEndian, &h3m.MapSize)
	binary.Read(decompressedMap, binary.LittleEndian, &h3m.HasUnderground)

	h3m.Name, err = readString(decompressedMap)
	if err != nil {
		return err
	}

	h3m.Description, err = readString(decompressedMap)
	if err != nil {
		return err

	}

	binary.Read(decompressedMap, binary.LittleEndian, &h3m.Difficulty)
	return nil
}

func loadPlayersData(decompressedMap io.ReadSeeker, h3m *H3m) error {
	//  Red, Blue, Tan, Green, Orange, Purple, Teal, Pink
	//  0,   1,    2,   3,     4,     5,      6,    7

	var heroesMaxLevel uint8 // max level of each hero separately ?

	currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	fmt.Printf("Current offset: %d\n", currentOffset)

	binary.Read(decompressedMap, binary.LittleEndian, &heroesMaxLevel)

	for i := 0; i < 8; i++ {
		var player models.SPlayer
		var fullPlayer models.Player

		currentOffset, err := decompressedMap.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		fmt.Printf("Player %d Current offset: %d\n", i, currentOffset)
		binary.Read(decompressedMap, binary.LittleEndian, &player)

		fullPlayer.SPlayer = &player

		if player.HasMainTown {
			var townCoordinates models.TownCoordinates
			binary.Read(decompressedMap, binary.LittleEndian, &townCoordinates)
			fullPlayer.TownCoordinates = &townCoordinates
		}

		// Player main hero
		var haveRandomHero bool
		binary.Read(decompressedMap, binary.LittleEndian, &haveRandomHero)

		var heroType uint8
		binary.Read(decompressedMap, binary.LittleEndian, &heroType)

		var heroFace uint8
		binary.Read(decompressedMap, binary.LittleEndian, &heroFace)

		heroName, err := readString(decompressedMap.(io.Reader))
		if err != nil {
			return err
		}

		if heroName == "" {
			h3m.Players = append(h3m.Players, &fullPlayer)
			continue
		}

		var fullHero models.MainHero
		fullHero.Name = heroName

		var unknown uint8
		binary.Read(decompressedMap, binary.LittleEndian, &unknown)

		var heroesCount uint32
		binary.Read(decompressedMap, binary.LittleEndian, &heroesCount)

		for j := uint32(0); j < heroesCount; j++ {
			var heroId uint8
			binary.Read(decompressedMap, binary.LittleEndian, &heroId)

			heroName, err := readString(decompressedMap)
			if err != nil {
				return err
			}

			fullHero.Heroes = append(fullHero.Heroes, &models.HeroDetails{
				Id:   heroId,
				Name: heroName,
			})

		}

		h3m.Players = append(h3m.Players, &fullPlayer)
	}

	return nil
}

func decompressGZIP(gzipedR io.Reader) (io.ReadSeeker, error) {
	r, err := gzip.NewReader(gzipedR)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	out := bytes.NewReader(b)

	return out, nil
}

func readString(r io.Reader) (string, error) {
	var len uint32
	binary.Read(r, binary.LittleEndian, &len)

	text := make([]byte, len)
	_, err := r.Read(text)
	if err != nil {
		return "", err
	}
	return string(text), nil
}
