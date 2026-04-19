package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddr  string
	vaultToken string
	namespace  string
	outputFile string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync HashiCorp Vault secrets to local .env files",
	Long: `vaultpull fetches secrets from HashiCorp Vault and writes them
to a local .env file. Supports namespaces for multi-tenant Vault setups.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&vaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Vault namespace")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", ".env", "Output .env file path")
}
