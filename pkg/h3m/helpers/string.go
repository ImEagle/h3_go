package helpers

import (
	"encoding/binary"
	"io"
)

func ReadString(r io.Reader) (string, error) {
	var len uint32
	binary.Read(r, binary.LittleEndian, &len)

	text := make([]byte, len)
	_, err := r.Read(text)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func ReadMessageIfSet(r io.Reader) (bool, string, error) {
	var hasMsg bool
	err := binary.Read(r, binary.LittleEndian, &hasMsg)
	if err != nil {
		return false, "", err
	}

	if hasMsg {
		message, err := ReadString(r)
		if err != nil {
			return false, "", err
		}
		return true, message, nil
	}

	return false, "", nil
}
