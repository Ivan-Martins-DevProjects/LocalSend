package main

import (
	"fmt"
	"os"
	"sync"

	"Ivan-Martins-DevProjects/localsend/cmd/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run ./cmd <caminho-do-arquivo>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	transfer, err := internal.NewTransfer(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo:", err)
		os.Exit(1)
	}
	defer transfer.File.Close()

	jobs := make(chan internal.Chunk)
	results := make(chan internal.Result)

	workers := 8
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			internal.Worker(jobs, results, transfer)
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

	if !internal.ValidateResults(results, transfer) {
		os.Exit(1)
	}
}
