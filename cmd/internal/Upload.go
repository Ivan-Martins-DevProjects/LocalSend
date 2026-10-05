package internal

import (
	"fmt"
	"net"
)

func UploadChunk(chunk Chunk, data []byte, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	buffer := make([]byte, 4096)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao receber:", err)
		return nil
	}

	response := fmt.Sprintf("Resposta: %s\n", buffer[:n])
	_, err = conn.Write([]byte(response))
	if err != nil {
		return err
	}
	return nil
}
