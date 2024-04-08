package objects

import (
	"encoding/binary"
	"io"
)

func ReadResource(decompressedMap io.ReadSeeker) (*Resource, error) {
	var resMessage string
	var hasMessage bool
	err := binary.Read(decompressedMap, binary.LittleEndian, &hasMessage)
	if err != nil {
		return nil, err
	}

	if hasMessage {
		message, err := readString(decompressedMap)
		if err != nil {
			return nil, err
		}
		resMessage = message

		var hasGuards bool
		err = binary.Read(decompressedMap, binary.LittleEndian, &hasGuards)
		if err != nil {
			return nil, err
		}

		if hasGuards {
			panic("Not implemented")
		}
	}

	return &Resource{
		Message: resMessage,
	}, nil
}

type Resource struct {
	Message string
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
