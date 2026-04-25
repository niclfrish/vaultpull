package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	syncp "vaultpull/internal/sync"
)

var tokenRotateCmd = &cobra.Command{
	Use:   "token-status",
	Short: "Show the current Vault token rotation status",
	Long:  `Displays the age of the currently cached Vault token and whether it is due for rotation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := os.Getenv("VAULT_TOKEN")
		if token == "" {
			return fmt.Errorf("VAULT_TOKEN environment variable is not set")
		}

		rotateAfterStr, _ := cmd.Flags().GetString("rotate-after")
		gracePeriodStr, _ := cmd.Flags().GetString("grace-period")

		rotateAfter, err := time.ParseDuration(rotateAfterStr)
		if err != nil {
			return fmt.Errorf("invalid --rotate-after: %w", err)
		}
		gracePeriod, err := time.ParseDuration(gracePeriodStr)
		if err != nil {
			return fmt.Errorf("invalid --grace-period: %w", err)
		}

		cfg := syncp.TokenRotateConfig{
			RotateAfter: rotateAfter,
			GracePeriod: gracePeriod,
		}

		// Static fetcher: in a real integration this would call Vault's token renewal.
		fetcher := func() (string, error) {
			return os.Getenv("VAULT_TOKEN"), nil
		}

		rotator, err := syncp.NewTokenRotator(token, fetcher, cfg)
		if err != nil {
			return fmt.Errorf("failed to initialise token rotator: %w", err)
		}

		syncp.LogTokenAge(rotator, cmd.OutOrStdout())

		tok, err := rotator.Token()
		if err != nil {
			return fmt.Errorf("token rotation failed: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "active token (last 8 chars): ...%s\n", safeSuffix(tok, 8))
		return nil
	},
}

func safeSuffix(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

func init() {
	tokenRotateCmd.Flags().String("rotate-after", "24h", "Duration after which the token should be rotated (e.g. 12h, 30m)")
	tokenRotateCmd.Flags().String("grace-period", "5m", "Grace period to use stale token when rotation fails")
	rootCmd.AddCommand(tokenRotateCmd)
}
