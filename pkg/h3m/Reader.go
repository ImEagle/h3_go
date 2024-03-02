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

	return h3m, nil
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
		fmt.Printf("Player %d Current offset: %d\n", i, currentOffset)
		binary.Read(decompressedMap, binary.LittleEndian, &player)

		fullPlayer.SPlayer = &player

		if player.HasMainTown {
			var townCoordinates models.TownCoordinates
			binary.Read(decompressedMap, binary.LittleEndian, &townCoordinates)
			fullPlayer.TownCoordinates = &townCoordinates
		}

		// Player main hero
		var fullHero models.MainHero
		var hero models.SMainHero
		binary.Read(decompressedMap, binary.LittleEndian, &hero)

		fullHero.SMainHero = &hero

		fullHero.Name, err = readString(decompressedMap)
		if err != nil {
			return err
		}

		binary.Read(decompressedMap, binary.LittleEndian, &fullHero.Unknown)

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
