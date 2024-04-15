package objects

import (
	"encoding/binary"
	"github.com/ImEagle/h3_go/pkg/h3m"
	"io"
)

func ReadRandomMonster(decompressedMap io.ReadSeeker, mapType string) (*RandomMonster, error) {
	var randomMonster RandomMonster

	if mapType != "RoE" {
		err := binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.Identifier)
		if err != nil {
			return nil, err
		}
	}

	err := binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.Count)
	if err != nil {
		return nil, err
	}

	err = binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.Character)
	if err != nil {
		return nil, err
	}

	hasMsg, msg, err := h3m.ReadMessageIfSet(decompressedMap)
	if hasMsg {
		randomMonster.Message = msg
		var resources MonsterResources
		err = binary.Read(decompressedMap, binary.LittleEndian, &resources)
		if err != nil {
			return nil, err
		}
		// Read Artifact
	}

	binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.NeverFlees)
	binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.NotGrown)

	decompressedMap.Seek(2, io.SeekCurrent)

	// TODO: How to check HOTA3 feature
	featureLevelHOTA3 := true
	if featureLevelHOTA3 {
		var featureHOTA3 MonsterFeatureHOTA3
		binary.Read(decompressedMap, binary.LittleEndian, &featureHOTA3)
	}

	return &randomMonster, nil
}

type RandomMonster struct {
	Identifier uint32 // ??
	Count      uint16
	Character  uint8
	Message    string
	Resources  MonsterResources
	NeverFlees bool
	NotGrown   bool
}

type MonsterResources struct {
	Res1 uint32 // Gold?
	Res2 uint32 // Wood?
	Res3 uint32 // Ore?
	Res4 uint32 // Mercury?
	Res5 uint32 // Sulfur
	Res6 uint32 // Crystal
	Res7 uint32 // Gems
}

type MonsterFeatureHOTA3 struct {
	AggressionExact  uint32 // -1 = default, 1-10 = possible values range
	JoinOnlyForMoney bool   // if true, monsters will only join for money
	JoinPercent      uint32 // 100 = default, percent of monsters that will join on succesfull agression check
	UpgradedStack    uint32 // Presence of upgraded stack, -1 = random, 0 = never, 1 = always
	StackCount       uint32 // TODO: check possible values. How many creature stacks will be present on battlefield, -1 = default
}
