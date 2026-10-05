package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"Ivan-Martins-DevProjects/localsend/cmd/internal/config"
)

var list bool

func listConfig() {
	if err := config.ListConfig(); err != nil {
		fmt.Println(err.Error())
	}
}

var update bool

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure os parâmetros de transferência para sua aplicação",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case list:
			listConfig()
			return nil

		case update:
			if err := config.CreateConfig(); err != nil {
				return err
			}
			return nil

		default:
			return cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(
		&list,
		"list",
		"l",
		false,
		"Liste as configurações já registradas.",
	)

	configCmd.Flags().BoolVarP(
		&update,
		"update",
		"u",
		false,
		"Atualize as configurações definidas",
	)
}
