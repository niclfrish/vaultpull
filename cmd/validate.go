package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/env"
	"vaultpull/internal/sync"
)

var (
	validateEnvFile   string
	validateRequired  string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a local .env file against required keys",
	Long:  "Reads a local .env file and checks for required keys, empty keys, and oversized values.",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		reader := env.NewReader(validateEnvFile)
		secrets, err := reader.Read()
		if err != nil {
			return fmt.Errorf("reading env file: %w", err)
		}

		var required []string
		if validateRequired != "" {
			for _, k := range strings.Split(validateRequired, ",") {
				k = strings.TrimSpace(k)
				if k != "" {
					required = append(required, k)
				}
			}
		}

		return sync.ValidateAndReport(os.Stdout, secrets, required)
	},
}

func init() {
	validateCmd.Flags().StringVarP(&validateEnvFile, "file", "f", ".env", "path to the .env file to validate")
	validateCmd.Flags().StringVar(&validateRequired, "require", "", "comma-separated list of required keys")
	rootCmd.AddCommand(validateCmd)
}
