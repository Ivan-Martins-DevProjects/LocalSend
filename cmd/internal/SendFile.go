package internal

import (
	"fmt"
	"os"
	"sync"
)

func SendFile(filepath, address string) {
	filePath := filepath

	transfer, err := NewTransfer(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo:", err)
		os.Exit(1)
	}
	defer transfer.File.Close()

	jobs := make(chan Chunk)
	results := make(chan Result)

	workers := 8
	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()
			Worker(jobs, results, transfer, address)
		}()
	}

	go func() {
		for _, chunk := range transfer.Chunks {
			jobs <- chunk
		}

		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	if !ValidateResults(results, transfer) {
		os.Exit(1)
	}
}
