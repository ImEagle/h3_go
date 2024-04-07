package objects

import (
	"encoding/binary"
	"github.com/ImEagle/h3_go/pkg/h3m"
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
		message, err := h3m.ReadString(decompressedMap)
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
