package config

import (
	"encoding/json"
	"errors"
	"os"

	apperror "Ivan-Martins-DevProjects/localsend/cmd/internal/errors"
)

type Config struct {
	Port       string `json:"port"`
	Workers    int    `json:"workers"`
	BufferSize int64  `json:"buffer_size"`
}

const path = "config.json"

var defaultConfig = &Config{
	Port:       ":8080",
	Workers:    4,
	BufferSize: 1024,
}

func GetConfig() (*Config, error) {
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		if err := defaultConfig.SaveConfig(); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperror.InternalServerError("Erro ao coletar dados do arquivo de configuração", err)
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, apperror.InternalServerError("Erro ao realizar o parsing dos dados do arquivo de configuração", err)
	}

	return &cfg, nil
}

func (c *Config) SaveConfig() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return apperror.InternalServerError("Erro ao inserir dados do arquivo de configuração.", err)
	}

	return os.WriteFile(path, data, 0o644)
}
