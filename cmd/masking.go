package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var maskingCmd = &cobra.Command{
	Use:   "mask",
	Short: "Fetch secrets and print them with sensitive values masked",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		maskCfg := sync.DefaultMaskConfig()

		if extra, _ := cmd.Flags().GetStringSlice("mask-pattern"); len(extra) > 0 {
			maskCfg.Patterns = append(maskCfg.Patterns, extra...)
		}
		if extra, _ := cmd.Flags().GetStringSlice("partial-pattern"); len(extra) > 0 {
			maskCfg.PartialPatterns = append(maskCfg.PartialPatterns, extra...)
		}
		if v, _ := cmd.Flags().GetInt("visible-chars"); v > 0 {
			maskCfg.VisibleChars = v
		}

		masker, err := sync.NewMasker(maskCfg)
		if err != nil {
			return fmt.Errorf("masker: %w", err)
		}

		masked := masker.Apply(secrets)

		w := cmd.OutOrStdout()
		for k, v := range masked {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	maskedFlags := maskedCmd()
	if maskedFlags == nil {
		os.Exit(1)
	}
	rootCmd.AddCommand(maskingCmd)
	maskingCmd.Flags().StringSlice("mask-pattern", nil, "Additional regex patterns for full masking")
	maskingCmd.Flags().StringSlice("partial-pattern", nil, "Additional regex patterns for partial masking")
	maskingCmd.Flags().Int("visible-chars", 4, "Number of trailing characters to reveal in partial masking")
}

func maskedCmd() *cobra.Command { return maskingCmd }
