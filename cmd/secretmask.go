package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var secretmaskCmd = &cobra.Command{
	Use:   "secretmask",
	Short: "Preview secrets with sensitive values masked",
	Long:  `Fetches secrets and prints them with sensitive key values masked based on key name patterns.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		patterns, _ := cmd.Flags().GetStringSlice("pattern")
		visible, _ := cmd.Flags().GetInt("visible-chars")
		maskChar, _ := cmd.Flags().GetString("mask-char")

		cfg := sync.DefaultSecretMaskConfig()
		if len(patterns) > 0 {
			cfg.KeyPatterns = patterns
		}
		if visible >= 0 {
			cfg.VisibleChars = visible
		}
		if maskChar != "" {
			cfg.MaskChar = maskChar
		}

		// Example secrets sourced from environment for preview purposes.
		secrets := map[string]string{}
		for _, e := range os.Environ() {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				secrets[parts[0]] = parts[1]
			}
		}

		res, err := sync.MaskSecrets(secrets, cfg)
		if err != nil {
			return fmt.Errorf("secretmask: %w", err)
		}

		w := cmd.OutOrStdout()
		fmt.Fprintf(w, "Masked %d key(s):\n", len(res.MaskedKeys))
		for _, k := range res.MaskedKeys {
			fmt.Fprintf(w, "  %s=%s\n", k, res.Secrets[k])
		}
		return nil
	},
}

func init() {
	secretmaskCmd.Flags().StringSlice("pattern", nil, "Regex patterns to match sensitive key names (overrides default)")
	secretmaskCmd.Flags().Int("visible-chars", 4, "Number of trailing characters to leave unmasked (0 = full mask)")
	secretmaskCmd.Flags().String("mask-char", "*", "Character to use for masking")
	rootCmd.AddCommand(secretmaskCmd)
}
