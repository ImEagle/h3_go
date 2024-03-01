package lod

import (
	"encoding/binary"
	"os"
)

type File struct {
	*os.File
}

func (f *File) ReadInt32() (int32, error) {
	var buf int32
	err := binary.Read(f, binary.LittleEndian, &buf)
	return buf, err
}
