package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func buildTestFile(t *testing.T, size int64) (string, []byte) {
	t.Helper()

	content := make([]byte, size)
	for i := int64(0); i < size; i++ {
		content[i] = byte(i*31 + 7)
	}

	path := filepath.Join(t.TempDir(), "test.bin")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("falha ao criar arquivo de teste: %v", err)
	}

	return path, content
}

func runWorkers(t *testing.T, jobs chan Chunk, results chan Result, transfer *Transfer, workerCount int) {
	t.Helper()

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Worker(jobs, results, transfer)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
}

func sendChunks(t *testing.T, jobs chan<- Chunk, chunks []Chunk) {
	t.Helper()

	go func() {
		for _, chunk := range chunks {
			jobs <- chunk
		}
		close(jobs)
	}()
}

func collectResults(t *testing.T, results <-chan Result, totalChunks int) map[int]Result {
	t.Helper()

	processed := make(map[int]Result, totalChunks)
	for result := range results {
		if _, duplicated := processed[result.ChunkInfo.Index]; duplicated {
			t.Errorf("chunk %d processado mais de uma vez", result.ChunkInfo.Index)
		}
		processed[result.ChunkInfo.Index] = result
	}

	return processed
}

func validateCoverage(t *testing.T, processed map[int]Result, totalChunks int) {
	t.Helper()

	for i := 0; i < totalChunks; i++ {
		result, ok := processed[i]
		if !ok {
			t.Errorf("chunk %d nao foi processado", i)
			continue
		}
		if result.Err != nil {
			t.Errorf("chunk %d falhou: %v", i, result.Err)
		}
	}

	if len(processed) != totalChunks {
		t.Errorf("esperados %d chunks processados, obtidos %d", totalChunks, len(processed))
	}
}

func TestWorkers_ProcessaTodosOsChunksUmaVez(t *testing.T) {
	const fileSize int64 = 10*1024*1024 + 12345

	path, _ := buildTestFile(t, fileSize)
	transfer, err := NewTransfer(path)
	if err != nil {
		t.Fatalf("falha ao criar transfer: %v", err)
	}
	defer transfer.File.Close()

	jobs := make(chan Chunk)
	results := make(chan Result)

	runWorkers(t, jobs, results, transfer, 8)
	sendChunks(t, jobs, transfer.Chunks)

	processed := collectResults(t, results, transfer.TotalChunks)
	validateCoverage(t, processed, transfer.TotalChunks)
}

func TestWorkers_IntegridadeDosDados(t *testing.T) {
	const fileSize int64 = 5*1024*1024 + 7

	path, content := buildTestFile(t, fileSize)
	transfer, err := NewTransfer(path)
	if err != nil {
		t.Fatalf("falha ao criar transfer: %v", err)
	}
	defer transfer.File.Close()

	jobs := make(chan Chunk)
	results := make(chan Result)

	runWorkers(t, jobs, results, transfer, 4)
	sendChunks(t, jobs, transfer.Chunks)

	processed := collectResults(t, results, transfer.TotalChunks)
	validateCoverage(t, processed, transfer.TotalChunks)

	for _, chunk := range transfer.Chunks {
		data, err := transfer.ReadChunk(chunk)
		if err != nil {
			t.Errorf("chunk %d nao pode ser relido: %v", chunk.Index, err)
			continue
		}

		expected := content[chunk.Offset : chunk.Offset+chunk.Size]
		if !bytes.Equal(data, expected) {
			t.Errorf("chunk %d com conteudo divergente do arquivo original", chunk.Index)
		}
	}
}

func TestWorkers_ArquivosLimite(t *testing.T) {
	tests := []struct {
		name          string
		size          int64
		totalExpected int
	}{
		{"arquivo vazio", 0, 0},
		{"menor que um chunk", 1024, 1},
		{"exatamente um chunk", ChunkSize, 1},
		{"um chunk e um byte", ChunkSize + 1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, _ := buildTestFile(t, tt.size)
			transfer, err := NewTransfer(path)
			if err != nil {
				t.Fatalf("falha ao criar transfer: %v", err)
			}
			defer transfer.File.Close()

			if transfer.TotalChunks != tt.totalExpected {
				t.Fatalf("esperados %d chunks, obtidos %d", tt.totalExpected, transfer.TotalChunks)
			}

			jobs := make(chan Chunk)
			results := make(chan Result)

			runWorkers(t, jobs, results, transfer, 1)
			sendChunks(t, jobs, transfer.Chunks)

			processed := collectResults(t, results, transfer.TotalChunks)
			validateCoverage(t, processed, transfer.TotalChunks)
		})
	}
}
