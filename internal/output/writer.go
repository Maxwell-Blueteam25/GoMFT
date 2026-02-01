package output

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONWriter struct {
	file    *os.File
	encoder *json.Encoder
	mu      sync.Mutex
}

func NewJSONWriter(path string) (*JSONWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &JSONWriter{
		file:    f,
		encoder: json.NewEncoder(f),
	}, nil
}

func (w *JSONWriter) Write(data interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.encoder.Encode(data)
}

func (w *JSONWriter) Close() error {
	return w.file.Close()
}
