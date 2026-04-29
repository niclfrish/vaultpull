package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var flattenCmd = &cobra.Command{
	Use:   "flatten",
	Short: "Flatten dot-separated secret keys into env-style names",
	Long: `Fetch secrets from Vault and flatten any dot-separated key segments
into a single key joined by a separator (default "_"), then write to
the configured .env file.`,
	RunE: runFlatten,
}

func runFlatten(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("flatten: config: %w", err)
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("flatten: vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cmd.Context(), cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("flatten: fetch secrets: %w", err)
	}

	flatCfg := sync.DefaultFlattenConfig()

	sep, _ := cmd.Flags().GetString("separator")
	if sep != "" {
		flatCfg.Separator = sep
	}

	maxDepth, _ := cmd.Flags().GetInt("max-depth")
	if maxDepth > 0 {
		flatCfg.MaxDepth = maxDepth
	}

	noUpper, _ := cmd.Flags().GetBool("no-upper")
	if noUpper {
		flatCfg.UpperCase = false
	}

	flattened, err := sync.FlattenSecrets(secrets, flatCfg)
	if err != nil {
		return fmt.Errorf("flatten: %w", err)
	}

	fmt.Fprintln(os.Stdout, sync.FlattenSummary(secrets, flattened))
	return nil
}

func init() {
	flattenCmd.Flags().String("separator", "_", "Separator to join key segments")
	flattenCmd.Flags().Int("max-depth", 0, "Maximum number of segments to join (0 = unlimited)")
	flattenCmd.Flags().Bool("no-upper", false, "Disable automatic upper-casing of keys")
	rootCmd.AddCommand(flattenCmd)
}
