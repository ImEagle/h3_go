package lod

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"os"
)

func isPCX(data []byte) bool {
	if len(data) < 12 {
		return false
	}

	size := binary.LittleEndian.Uint32(data[:4])
	width := binary.LittleEndian.Uint32(data[4:8])
	height := binary.LittleEndian.Uint32(data[8:12])

	return size == width*height || size == width*height*3
}

func pcxToImage(data []byte) image.Image {
	if len(data) < 12 {
		return nil
	}

	size := binary.LittleEndian.Uint32(data[:4])
	width := binary.LittleEndian.Uint32(data[4:8])
	height := binary.LittleEndian.Uint32(data[8:12])

	if size == width*height {
		img := image.NewPaletted(image.Rect(0, 0, int(width), int(height)), make([]color.Color, 255))
		palette := make([]color.Color, 256)

		for i := 0; i < 256; i++ {
			offset := 12 + int(width*height) + i*3
			r := data[offset]
			g := data[offset+1]
			b := data[offset+2]
			palette[i] = color.RGBA{R: r, G: g, B: b, A: 255}
		}

		img.Palette = palette
		for y := 0; y < int(height); y++ {
			for x := 0; x < int(width); x++ {
				index := data[12+y*int(width)+x]
				img.Set(x, y, palette[index])
			}
		}

		return img
	} else if size == width*height*3 {
		img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		copy(img.Pix, data[12:])
		return img
	}

	return nil
}

func imageToBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func imageToFile(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
