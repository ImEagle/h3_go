package lod

import "encoding/binary"

func isPCX(data []byte) bool {
	if len(data) < 12 {
		return false
	}

	size := binary.LittleEndian.Uint32(data[:4])
	width := binary.LittleEndian.Uint32(data[4:8])
	height := binary.LittleEndian.Uint32(data[8:12])

	return size == width*height || size == width*height*3
}
