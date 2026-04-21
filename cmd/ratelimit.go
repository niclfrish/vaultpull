package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	defaultRPS   = 10
	defaultBurst = 5
)

var (
	rateLimitRPS   int
	rateLimitBurst int
)

var rateLimitCmd = &cobra.Command{
	Use:   "ratelimit",
	Short: "Show or validate the current rate limit configuration",
	Long: `Display the effective rate limit settings that vaultpull uses
when making requests to HashiCorp Vault. Values can be overridden
via flags or the VAULTPULL_RPS / VAULTPULL_BURST environment variables.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rps := resolveInt("VAULTPULL_RPS", rateLimitRPS, defaultRPS)
		burst := resolveInt("VAULTPULL_BURST", rateLimitBurst, defaultBurst)

		if rps <= 0 {
			return fmt.Errorf("ratelimit: requests-per-second must be > 0")
		}
		if burst <= 0 {
			return fmt.Errorf("ratelimit: burst must be > 0")
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Rate limit config:\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  requests-per-second: %d\n", rps)
		fmt.Fprintf(cmd.OutOrStdout(), "  burst:               %d\n", burst)
		return nil
	},
}

// resolveInt returns the flag value if non-zero, else parses the env var,
// else falls back to the default.
func resolveInt(envKey string, flagVal, defaultVal int) int {
	if flagVal != 0 {
		return flagVal
	}
	if v := os.Getenv(envKey); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}

func init() {
	rateLimitCmd.Flags().IntVar(&rateLimitRPS, "rps", 0, "max requests per second to Vault (default 10)")
	rateLimitCmd.Flags().IntVar(&rateLimitBurst, "burst", 0, "max burst size above rate limit (default 5)")
	rootCmd.AddCommand(rateLimitCmd)
}
