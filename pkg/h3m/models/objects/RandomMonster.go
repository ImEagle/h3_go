package objects

import (
	"encoding/binary"
	"io"
)

func ReadRandomMonster(decompressedMap io.ReadSeeker, mapType string) (*RandomMonster, error) {
	var randomMonster RandomMonster

	err := binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.Count)
	if err != nil {
		return nil, err
	}

	err = binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.Character)
	if err != nil {
		return nil, err
	}

	var hasMsg bool
	err = binary.Read(decompressedMap, binary.LittleEndian, &hasMsg)
	if err != nil {
		return nil, err
	}

	if hasMsg {
		message, err := readString(decompressedMap)
		if err != nil {
			return nil, err
		}
		randomMonster.Message = message
	}

	// ReadResources???
	err = binary.Read(decompressedMap, binary.LittleEndian, &randomMonster.Resources)
	if err != nil {
		return nil, err
	}

	if mapType == "RoE" {
		decompressedMap.Seek(1, io.SeekCurrent)
	} else {
		decompressedMap.Seek(2, io.SeekCurrent)
	}

	// ????
	decompressedMap.Seek(2, io.SeekCurrent)

	return &randomMonster, nil
}

type RandomMonster struct {
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
