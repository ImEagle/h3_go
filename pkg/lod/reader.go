package lod

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const FILE_NAME_LENGTH = 16

type fileMetadata struct {
	Offset   uint32
	Size     uint32
	Unknown1 uint32
	Csize    uint32
}

func NewReader(fileName string) *Reader {
	return &Reader{
		fileName: fileName,
		metadata: make(map[string]fileMetadata),
	}
}

type Reader struct {
	fileName string

	fileList []string
	metadata map[string]fileMetadata
}

func (r *Reader) GetFileList() []string {
	if r.fileList == nil {
		err := r.LoadMetadata()
		if err != nil {
			panic(err)
		}
	}
	return r.fileList
}

func (r *Reader) LoadMetadata() error {
	f, err := os.OpenFile(r.fileName, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// ----- Read header -----
	header := make([]byte, 4)
	if _, err := f.Read(header); err != nil {
		return err
	}

	// ----- ~ Read header ~ -----

	// ----- Read file count -----
	_fileCountOffset := int64(8) // TODO: Magic number
	_, err = f.Seek(_fileCountOffset, io.SeekStart)
	if err != nil {
		return err
	}

	fileCount, err := readInt32(f)
	if err != nil {
		return err
	}
	fmt.Printf("Total files %d\n", fileCount)
	// ----- ~ Read file count ~ -----

	// ----- Read file list -----
	_fileDataOffset := int64(92) // TODO: Magic number
	_, err = f.Seek(_fileDataOffset, io.SeekStart)

	for i := int32(0); i < fileCount; i++ {
		fileName, err := readString(f, FILE_NAME_LENGTH)
		if err != nil {
			return err
		}

		// ----- Read file metadata -----
		var metadata fileMetadata
		err = binary.Read(f, binary.LittleEndian, &metadata)
		if err != nil {
			return err
		}
		// ----- ~ Read file metadata ~ -----

		r.metadata[fileName] = metadata
	}

	return nil
}

func readInt32(f *os.File) (int32, error) {
	var buf int32
	err := binary.Read(f, binary.LittleEndian, &buf)
	if err != nil {
		return 0, err
	}
	return buf, nil
}

func readString(f *os.File, len int) (string, error) {
	buf := make([]byte, len)
	_, err := f.Read(buf)
	if err != nil {
		return "", err
	}

	nullIndex := bytes.IndexByte(buf, 0)
	fileName := buf[:nullIndex]

	return string(fileName), nil
}
