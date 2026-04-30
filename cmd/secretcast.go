package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

// parseCastRules parses a slice of "KEY:TYPE" strings into CastRule values.
func parseCastRules(flags []string) ([]sync.CastRule, error) {
	rules := make([]sync.CastRule, 0, len(flags))
	for _, f := range flags {
		parts := strings.SplitN(f, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid cast rule %q: expected KEY:TYPE", f)
		}
		rules = append(rules, sync.CastRule{
			Key:    strings.TrimSpace(parts[0]),
			CastTo: sync.CastType(strings.TrimSpace(parts[1])),
		})
	}
	return rules, nil
}

func init() {
	var castFlags []string

	castCmd := &cobra.Command{
		Use:   "cast",
		Short: "Fetch secrets and cast specified keys to target types",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
			if err != nil {
				return fmt.Errorf("vault client: %w", err)
			}

			secrets, err := client.GetSecrets(cfg.SecretPath)
			if err != nil {
				return fmt.Errorf("fetch secrets: %w", err)
			}

			rules, err := parseCastRules(castFlags)
			if err != nil {
				return err
			}

			casted, err := sync.CastSecrets(secrets, rules)
			if err != nil {
				return fmt.Errorf("cast secrets: %w", err)
			}

			cmd.PrintErrln(sync.CastSummary(rules))
			for k, v := range casted {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}

	castCmd.Flags().StringArrayVar(&castFlags, "rule", nil,
		"cast rule in KEY:TYPE format (types: string, int, float, bool); repeatable")

	rootCmd.AddCommand(castCmd)
}
