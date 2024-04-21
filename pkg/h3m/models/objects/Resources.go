package objects

import (
	"encoding/binary"
	"github.com/ImEagle/h3_go/pkg/h3m/helpers"
	"io"
)

func ReadResource(decompressedMap io.ReadSeeker) (*Resource, error) {
	var hasMessage bool
	err := binary.Read(decompressedMap, binary.LittleEndian, &hasMessage)
	if err != nil {
		return nil, err
	}

	hasMsg, msg, err := helpers.ReadMessageIfSet(decompressedMap)
	if hasMsg {
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
		Message: msg,
	}, nil
}

type Resource struct {
	Message string
}
