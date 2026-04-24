package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var labelFilterCmd = &cobra.Command{
	Use:   "label-filter",
	Short: "Fetch secrets and keep only those matching specified labels",
	Long: `Fetch secrets from Vault and filter them by label annotations.

Labels are embedded in secret keys using the convention:
  BASE_KEY__label__KEY=VALUE

Only secrets carrying ALL specified --label flags are written to the output file.`,
	RunE: runLabelFilter,
}

var labelFlags []string

func runLabelFilter(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	labels, err := sync.ParseLabelFlags(labelFlags)
	if err != nil {
		return fmt.Errorf("parse labels: %w", err)
	}

	pipeline := sync.NewPipeline()
	pipeline.Add(sync.LabelFilterStage(labels))

	filtered, err := pipeline.Run(secrets)
	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}

	fmt.Fprintf(os.Stdout, "label-filter: %d secret(s) matched\n", len(filtered))
	return nil
}

func init() {
	labelFilterCmd.Flags().StringArrayVar(&labelFlags, "label", nil,
		"Label selector in key=value format (repeatable; all must match)")
	rootCmd.AddCommand(labelFilterCmd)
}
