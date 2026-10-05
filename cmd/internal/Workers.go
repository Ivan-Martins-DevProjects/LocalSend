package internal

type Result struct {
	ChunkInfo Chunk
	Err       error
}

func Worker(jobs <-chan Chunk, results chan<- Result, transfer *Transfer, address string) {
	for chunk := range jobs {
		data, err := transfer.ReadChunk(chunk)
		if err != nil {
			results <- Result{ChunkInfo: chunk, Err: err}
			continue
		}

		err = UploadChunk(chunk, data, address)
		results <- Result{
			ChunkInfo: chunk,
			Err:       err,
		}
	}
}
