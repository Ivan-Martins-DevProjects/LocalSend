package internal

import (
	"log"
	"os"

	"Ivan-Martins-DevProjects/localsend/cmd/internal/config"
)

type Chunk struct {
	Index  int
	Offset int64
	Size   int64
}

func NewTransfer(path string) (*Transfer, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := info.Size()

	totalChunks := int((fileSize + cfg.BufferSize - 1) / cfg.BufferSize)
	chunks := make([]Chunk, 0, totalChunks)

	log.Println("Processando arquivo...")
	for i := 0; i < totalChunks; i++ {
		// Cada chunk possui cfg.BufferSize bytes.
		// Para descobrir onde um chunk começa no arquivo, multiplicamos
		// seu índice (i) pelo tamanho de um chunk.
		// Ex.: i = 2 e cfg.BufferSize = 1024 => offset = 2048.
		offset := int64(i) * cfg.BufferSize

		log.Printf("Carregado: %d/%d", i+1, totalChunks)
		size := cfg.BufferSize
		remaining := fileSize - offset
		if remaining < cfg.BufferSize {
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
