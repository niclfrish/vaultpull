package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Search for secrets by key name",
	Long:  `Lookup searches Vault secrets by exact or partial key name and prints matching results.`,
	RunE:  runLookup,
}

func runLookup(cmd *cobra.Command, args []string) error {
	queries, _ := cmd.Flags().GetStringSlice("query")
	partial, _ := cmd.Flags().GetBool("partial")
	caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")

	if len(queries) == 0 {
		return fmt.Errorf("at least one --query is required")
	}

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

	lookupCfg := sync.DefaultLookupConfig()
	lookupCfg.PartialMatch = partial
	lookupCfg.CaseSensitive = caseSensitive

	normalized := make([]string, len(queries))
	for i, q := range queries {
		normalized[i] = strings.TrimSpace(q)
	}

	_, err = sync.LookupAndReport(secrets, normalized, lookupCfg, os.Stdout)
	return err
}

func init() {
	lookupCmd.Flags().StringSlice("query", nil, "Key names to search for (repeatable)")
	lookupCmd.Flags().Bool("partial", false, "Enable partial/substring matching")
	lookupCmd.Flags().Bool("case-sensitive", false, "Use case-sensitive key matching")
	rootCmd.AddCommand(lookupCmd)
}
