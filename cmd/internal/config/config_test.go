package config

import (
	"testing"

	apperror "Ivan-Martins-DevProjects/localsend/cmd/internal/app_errors"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name     string
		data     *Config
		wantErr  bool
		whichErr apperror.AppError
	}{
		{
			name: "Arquivo encontrado com sucesso",
		},
	}
}
