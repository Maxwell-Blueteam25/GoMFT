package parser

import (
	"encoding/binary"
	"io"
	"os"
)

type MftReader struct {
	FileHandler *os.File
	Buffer      []byte
}

func NewMftReader(path string) (*MftReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &MftReader{
		FileHandler: file,
		Buffer:      make([]byte, 1024),
	}, nil
}

func (r *MftReader) Next() ([]byte, error) {
	_, err := io.ReadFull(r.FileHandler, r.Buffer)
	if err != nil {
		return nil, err
	}

	e := r.fixup(r.Buffer)
	if e != nil {
		return nil, e
	}
	return r.Buffer, nil
}

func (r *MftReader) fixup(data []byte) error {

	usaOffset := binary.LittleEndian.Uint16(data[4:6])
	usaSize := binary.LittleEndian.Uint16(data[6:8])

	expectedSeqNum := binary.LittleEndian.Uint16(data[usaOffset : usaOffset+2])

	for i := 1; i < int(usaSize); i++ {

		sectorEnd := (i * 512) - 2

		if sectorEnd+2 > len(data) {
			break
		}

		actualSeqNum := binary.LittleEndian.Uint16(data[sectorEnd : sectorEnd+2])
		if actualSeqNum != expectedSeqNum {
			return os.ErrInvalid
		}

		savedBytesStart := int(usaOffset) + (i * 2)
		copy(data[sectorEnd:sectorEnd+2], data[savedBytesStart:savedBytesStart+2])
	}

	return nil
}
