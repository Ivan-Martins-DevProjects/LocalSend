package internal

import (
	"log"
	"os"
)

const ChunkSize int64 = 4 * 1024 * 1024

type Chunk struct {
	Index  int
	Offset int64
	Size   int64
}

func NewTransfer(path string) (*Transfer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := info.Size()

	totalChunks := int((fileSize + ChunkSize - 1) / ChunkSize)
	chunks := make([]Chunk, 0, totalChunks)

	log.Println("Processando arquivo...")
	for i := 0; i < totalChunks; i++ {
		// Cada chunk possui ChunkSize bytes.
		// Para descobrir onde um chunk começa no arquivo, multiplicamos
		// seu índice (i) pelo tamanho de um chunk.
		// Ex.: i = 2 e ChunkSize = 1024 => offset = 2048.
		offset := int64(i) * ChunkSize

		log.Printf("Carregado: %d/%d", i+1, totalChunks)
		size := ChunkSize
		remaining := fileSize - offset
		if remaining < ChunkSize {
			size = remaining
		}

		chunks = append(chunks, Chunk{
			Index:  i,
			Offset: offset,
			Size:   size,
		})

	}
	response := &Transfer{
		File:        file,
		FileName:    info.Name(),
		FileSize:    info.Size(),
		TotalChunks: totalChunks,
		Chunks:      chunks,
	}

	return response, nil
}
