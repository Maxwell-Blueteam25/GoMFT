package parser

import (
	"os"
)

type JournalReader struct {
	FileHandler *os.File
	Buffer      []byte
}

func NewJournalReader(path string) (*JournalReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &JournalReader{
		FileHandler: file,
		Buffer:      make([]byte, 65536),
	}, nil
}

func (r *JournalReader) ReadChunk() (int, error) {
	n, err := r.FileHandler.Read(r.Buffer)
	return n, err
}
