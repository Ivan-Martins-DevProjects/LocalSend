package internal

import (
	"io"
	"os"
)

type Transfer struct {
	IP          string
	File        *os.File
	FileName    string
	FileSize    int64
	TotalChunks int
	Chunks      []Chunk
}

func (t *Transfer) ReadChunk(chunk Chunk) ([]byte, error) {
	buffer := make([]byte, chunk.Size)

	_, err := t.File.ReadAt(buffer, chunk.Offset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buffer, nil
}
