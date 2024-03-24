package def

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"io"
	"os"
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
			var pixelData []byte

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
				if firstFullWidth > int(imgFrame.FullWidth) {
					imgFrame.FullWidth = uint32(firstFullWidth)
				}

				if firstFullWidth < int(imgFrame.FullWidth) {
					firstFullWidth = int(imgFrame.FullWidth)
				}

				if firstFullHeight > int(imgFrame.FullHeight) {
					imgFrame.FullHeight = uint32(firstFullHeight)
				}

				if firstFullHeight < int(imgFrame.FullHeight) {
					firstFullHeight = int(imgFrame.FullHeight)
					return errors.New("first height smaller than latter one")
				}

				if imgFrame.Width != 0 || imgFrame.Height != 0 {
					if imgFrame.Format == 0 {
						var err error
						pixelData, err = extractFromFormat0(imgFrame, bReader)
						if err != nil {
							return err
						}

					}

				}

				err := debugCreateImageFromBytes(pixelData, r.palette, imgFrame, fmt.Sprintf("image_%d.png", i))
				if err != nil {
					return err
				}

			}

		}
	}

	return nil
}

func extractFromFormat0(imgFrame frame, reader *bytes.Reader) ([]byte, error) {
	pixelData := make([]byte, imgFrame.Width*imgFrame.Height)
	err := binary.Read(reader, binary.LittleEndian, &pixelData)
	if err != nil {
		return nil, err
	}

	return pixelData, nil
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

func debugCreateImageFromBytes(data []byte, palette []byte, imgFrame frame, fileName string) error {

	pal := make([]color.Color, len(palette)/3)
	for i := 0; i < len(palette); i += 3 {
		pal[i/3] = color.RGBA{palette[i], palette[i+1], palette[i+2], 0xff}
	}

	imp := image.NewPaletted(image.Rect(0, 0, int(imgFrame.FullWidth), int(imgFrame.FullHeight)), pal)
	copy(imp.Pix, data)

	imRGB := image.NewRGBA(imp.Bounds())
	draw.Draw(imRGB, imRGB.Bounds(), imp, image.Point{}, draw.Src)

	replaceColor(imRGB, color.RGBA{0, 0, 0, 0}, palette, 0)
	replaceColor(imRGB, color.RGBA{0, 0, 0, 0x40}, palette, 1)
	replaceColor(imRGB, color.RGBA{0, 0, 0, 0x80}, palette, 4)
	replaceColor(imRGB, color.RGBA{0, 0, 0, 0}, palette, 5)
	replaceColor(imRGB, color.RGBA{0, 0, 0, 0x80}, palette, 6)
	replaceColor(imRGB, color.RGBA{0, 0, 0, 0x40}, palette, 7)

	im := image.NewRGBA(image.Rect(0, 0, int(imgFrame.FullWidth), int(imgFrame.FullHeight)))
	draw.Draw(im, im.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)
	draw.Draw(im, im.Bounds(), imRGB, image.Point{int(imgFrame.LeftMargin), int(imgFrame.TopMargin)}, draw.Src)

	// Create or open the file for writing
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating/opening file: %v", err)
	}
	defer file.Close()

	// Save the image to the file in PNG format
	err = png.Encode(file, im)
	if err != nil {
		return fmt.Errorf("error encoding image to PNG: %v", err)
	}

	fmt.Printf("Image successfully saved to %s\n", fileName)
	return nil
}

func replaceColor(img *image.RGBA, newColor color.RGBA, palette []uint8, paletteIndex uint8) {

	for i := 0; i < len(img.Pix); i += 4 {
		if img.Pix[i+3] == paletteIndex {
			img.Pix[i] = newColor.R
			img.Pix[i+1] = newColor.G
			img.Pix[i+2] = newColor.B
			img.Pix[i+3] = newColor.A
		}
	}
}
