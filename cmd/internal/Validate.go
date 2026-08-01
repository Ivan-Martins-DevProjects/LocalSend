package internal

import (
	"fmt"
	"time"
)

func ValidateResults(results <-chan Result, transfer *Transfer) bool {
	start := time.Now()

	processed := make(map[int]struct{}, transfer.TotalChunks)
	duplicated := make([]int, 0)
	var failed []Result

	for result := range results {
		if _, exists := processed[result.ChunkInfo.Index]; exists {
			duplicated = append(duplicated, result.ChunkInfo.Index)
			continue
		}
		processed[result.ChunkInfo.Index] = struct{}{}

		if result.Err != nil {
			failed = append(failed, result)
		}
	}

	var missing []int
	for i := 0; i < transfer.TotalChunks; i++ {
		if _, exists := processed[i]; !exists {
			missing = append(missing, i)
		}
	}

	ok := len(missing) == 0 && len(duplicated) == 0 && len(failed) == 0

	fmt.Println()
	fmt.Println("========== RELATORIO DE VALIDACAO DOS WORKERS ==========")
	fmt.Printf("Arquivo: %s (%d bytes, %d chunks)\n", transfer.FileName, transfer.FileSize, transfer.TotalChunks)
	fmt.Printf("Chunks processados: %d/%d\n", len(processed), transfer.TotalChunks)
	fmt.Printf("Chunks perdidos: %d\n", len(missing))
	for _, index := range missing {
		fmt.Printf("  - chunk %d nao foi processado\n", index)
	}
	fmt.Printf("Chunks duplicados: %d\n", len(duplicated))
	for _, index := range duplicated {
		fmt.Printf("  - chunk %d processado mais de uma vez\n", index)
	}
	fmt.Printf("Chunks com erro: %d\n", len(failed))
	for _, result := range failed {
		fmt.Printf("  - chunk %d: %v\n", result.ChunkInfo.Index, result.Err)
	}
	fmt.Printf("Tempo total: %s\n", time.Since(start).Round(time.Millisecond))

	if ok {
		fmt.Println("RESULTADO: VALIDADO - todos os chunks foram processados corretamente")
	} else {
		fmt.Println("RESULTADO: FALHA - corrija os problemas acima")
	}
	fmt.Println("=========================================================")

	return ok
}
