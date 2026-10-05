package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ListConfig() error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	fmt.Printf("Port: %s\n", cfg.Port)
	fmt.Printf("Workers: %d\n", cfg.Workers)
	fmt.Printf("BufferSize: %d\n", cfg.BufferSize)

	return nil
}

func CreateConfig() error {
	scanner := bufio.NewScanner(os.Stdin)

	actualCfg, err := GetConfig()
	if err != nil {
		return err
	}

	fmt.Print("Informe abaixo as configuraçõe as quais deseja alterar\n\n")

	var port string
	for {
		fmt.Printf("----- Porta atual: %s -----\n", actualCfg.Port)
		fmt.Println("Deseje a nova porta desejada (pressione Return para pular esse item): ")
		fmt.Print("> ")

		scanner.Scan()
		port = strings.TrimSpace(scanner.Text())

		if port == "" {
			fmt.Println("Porta padrão preservada")
			break
		}

		break
	}

	return nil
}
