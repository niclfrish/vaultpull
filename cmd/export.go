package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Vault secrets to stdout in a chosen format",
	Long: `Fetch secrets from Vault and print them to stdout.

Supported formats:
  dotenv  — KEY="value" lines (default)
  export  — export KEY="value" lines (shell-sourceable)
  json    — pretty-printed JSON object`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		vc, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		formatFlag, _ := cmd.Flags().GetString("format")
		namespace, _ := cmd.Flags().GetString("namespace")

		path := sync.NamespacedPath(cfg.SecretPath, namespace)
		secrets, err := vc.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("fetch secrets: %w", err)
		}

		if namespace != "" {
			secrets = sync.PrefixKeys(secrets, namespace)
		}

		exporter, err := sync.NewExporter(sync.ExportFormat(formatFlag), os.Stdout)
		if err != nil {
			return fmt.Errorf("exporter: %w", err)
		}

		if err := exporter.Export(secrets); err != nil {
			return fmt.Errorf("export: %w", err)
		}
		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("format", "f", "dotenv", "Output format: dotenv, export, json")
	exportCmd.Flags().StringP("namespace", "n", "", "Namespace prefix for secret path and keys")
	rootCmd.AddCommand(exportCmd)
}
