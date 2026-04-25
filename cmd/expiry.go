package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var (
	expiryFailOnExpired bool
)

var expiryCmd = &cobra.Command{
	Use:   "expiry",
	Short: "Check secret expiry metadata from Vault",
	Long:  "Fetches secrets and inspects __expires_at__ companion keys, reporting which secrets have expired.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		path := sync.NamespacedPath(cfg.SecretPath, cfg.Namespace)
		secrets, err := client.GetSecrets(path)
		if err != nil {
			return fmt.Errorf("fetch secrets: %w", err)
		}

		result, err := sync.CheckExpiry(secrets, expiryFailOnExpired, os.Stdout)
		if err != nil {
			return err
		}

		if len(result.Infos) == 0 {
			fmt.Println("No expiry metadata found in secrets.")
		}
		return nil
	},
}

func init() {
	expiryCmd.Flags().BoolVar(&expiryFailOnExpired, "fail-on-expired", false,
		"Exit with non-zero status if any secrets are expired")
	rootCmd.AddCommand(expiryCmd)
}
