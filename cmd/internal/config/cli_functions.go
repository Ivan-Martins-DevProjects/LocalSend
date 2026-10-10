package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	apperror "Ivan-Martins-DevProjects/localsend/cmd/internal/app_errors"
)

func ListConfig() error {
	var cfg *Config

	if err := FindFile(path); err != nil {
		if err.Error() == apperror.FILE_NOT_FOUND {
			err = defaultConfig.SaveConfig()
			if err != nil {
				return err
			}
			cfg = defaultConfig
		}

		return err
	}

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

	for {
		fmt.Printf("----- Porta atual: %s -----\n", actualCfg.Port)
		fmt.Println("Digite a nova porta desejada (pressione Return para pular esse item): ")
		fmt.Print("> ")

		scanner.Scan()
		port := strings.TrimSpace(scanner.Text())
		if port == "" || fmt.Sprintf(":%s", port) == actualCfg.Port {
			fmt.Printf("\nPorta preservada\n\n")
			break
		}

		actualCfg.Port = fmt.Sprintf(":%s", port)
		fmt.Printf("\nPorta alterada\n\n")
		break
	}

	for {
		fmt.Printf("----- Número de workers: %d -----\n", actualCfg.Workers)
		fmt.Println("Deseje a quantidade de workers desejada (pressione Return para pular esse item): ")
		fmt.Print("> ")

		scanner.Scan()
		workers := strings.TrimSpace(scanner.Text())
		if workers == "" {
			fmt.Printf("\nNúmero de workers preservados\n\n")
			break
		}

		wk, err := strconv.Atoi(workers)
		if err != nil {
			fmt.Printf("\nNúmero inválido de workers\n\n")
			continue
		}

		if wk == actualCfg.Workers {
			fmt.Printf("\nNúmero de workers preservados\n\n")
			break
		}

		actualCfg.Workers = wk
		fmt.Printf("\nNúmero de workers alterados\n\n")
		break
	}

	for {
		fmt.Printf("----- Tamanho do buffer de transferência - bytes: %d -----\n", actualCfg.BufferSize)
		fmt.Println("Deseje o tamanho do buffer desejado (pressione Return para pular esse item): ")
		fmt.Print("> ")

		scanner.Scan()
		buffer := strings.TrimSpace(scanner.Text())
		if buffer == "" {
			fmt.Printf("\nBuffer preservado\n\n")
			break
		}

		bf, err := strconv.Atoi(buffer)
		if err != nil {
			fmt.Printf("\nErro ao fazer a alteração, garanta que o número de bytes foi digitado corretamente\n\n")
			continue
		}

		if bf == int(actualCfg.BufferSize) {
			fmt.Printf("\nBuffer preservado\n\n")
			break
		}

		actualCfg.Workers = bf
		fmt.Printf("\nBuffer de transferência alterado\n\n")
		break
	}

	if err = actualCfg.SaveConfig(); err != nil {
		fmt.Printf("\n\nErro ao salvar as alterações: %v", err)
	}

	return nil
}
