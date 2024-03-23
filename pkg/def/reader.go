package def

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const SpriteSheetType = 0x42

type frame struct {
	Size       uint32
	Format     uint32
	FullWidth  uint32
	FullHeight uint32
	Width      uint32
	Height     uint32
	LeftMargin uint32
	TopMargin  uint32
}

type header struct {
	Type        uint32
	Width       uint32
	Height      uint32
	BlocksCount uint32
}

func NewReader() *Reader {
	return &Reader{
		header:  header{},
		palette: make([]byte, 256*3),
		blocks:  make([]block, 0),
	}
}

type Reader struct {
	header  header
	palette []byte
	blocks  []block
}

func (r *Reader) Load(data []byte) error {
	bReader := bytes.NewReader(data)

	// Load header
	if err := binary.Read(bReader, binary.LittleEndian, &r.header); err != nil {
		return err
	}

	err := binary.Read(bReader, binary.LittleEndian, &r.palette)
	if err != nil {
		return err
	}

	for i := uint32(0); i < r.header.BlocksCount; i++ {
		blck, err := readBlock(bReader)
		if err != nil {
			return err
		}

		r.blocks = append(r.blocks, *blck)
	}

	r.fetchImages(bReader)

	debug := 1
	debug += 1

	return nil
}

func (r *Reader) fetchImages(bReader *bytes.Reader) error {

	firstFullWidth := -1
	firstFullHeight := -1

	for i := uint32(0); i < r.header.BlocksCount; i++ {
		// different block = different frames?

		for _, offset := range r.blocks[i].Offsets {
			var imgFrame frame

			bReader.Seek(int64(offset), io.SeekStart)
			binary.Read(bReader, binary.LittleEndian, &imgFrame)

			if (imgFrame.LeftMargin > imgFrame.FullWidth) || (imgFrame.TopMargin > imgFrame.FullHeight) {
				return errors.New("margins are higher than dimensions")
			}

			// https://gitlab.mister-muffin.de/josch/lodextract/src/branch/main/defextract.py#L92
			if firstFullWidth == -1 && firstFullHeight == -1 {
				firstFullWidth = int(imgFrame.FullWidth)
				firstFullHeight = int(imgFrame.FullHeight)
			} else {
				// TODO
			}

		}
	}

	return nil
}

func (r *Reader) CanGenerateSpriteSheet() bool {
	return r.header.Type == SpriteSheetType
}

type block struct {
	Id      uint32
	Count   uint32
	Width   uint32
	Height  uint32
	Names   []string
	Offsets []uint32
}

func readBlock(bReader *bytes.Reader) (*block, error) {
	blck := block{
		Names: make([]string, 0),
	}

	err := binary.Read(bReader, binary.LittleEndian, &blck.Id)
	if err != nil {
		return nil, err
	}

	err = binary.Read(bReader, binary.LittleEndian, &blck.Count)
	if err != nil {
		return nil, err
	}

	err = binary.Read(bReader, binary.LittleEndian, &blck.Width)
	if err != nil {
		return nil, err
	}

	err = binary.Read(bReader, binary.LittleEndian, &blck.Height)
	if err != nil {
		return nil, err
	}

	for i := int32(0); i < int32(blck.Count); i++ {
		name, err := readString(bReader, 13) // #TODO: Check if 13 is correct; Default string length is 16
		if err != nil {
			return nil, err
		}

		blck.Names = append(blck.Names, name)
	}

	blck.Offsets = make([]uint32, blck.Count)
	err = binary.Read(bReader, binary.LittleEndian, &blck.Offsets)
	if err != nil {
		return nil, err
	}

	return &blck, nil
}

func readString(f io.Reader, len int) (string, error) {
	buf := make([]byte, len)
	_, err := f.Read(buf)
	if err != nil {
		return "", err
	}

	nullIndex := bytes.IndexByte(buf, 0)
	fileName := buf[:nullIndex]

	return string(fileName), nil
}
