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
	encryptPassphrase string
	encryptOutput     string
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Fetch secrets from Vault, encrypt them, and write to a .env file",
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

		outPath := encryptOutput
		if outPath == "" {
			outPath = ".env.enc"
		}

		writeFn := func(m map[string]string) error {
			f, err := os.Create(outPath)
			if err != nil {
				return err
			}
			defer f.Close()
			for k, v := range m {
				fmt.Fprintf(f, "%s=%s\n", k, v)
			}
			return nil
		}

		return sync.EncryptAndWrite(secrets, encryptPassphrase, writeFn, os.Stdout)
	},
}

func init() {
	encryptCmd.Flags().StringVarP(&encryptPassphrase, "passphrase", "p", "", "Passphrase used to encrypt secrets (required)")
	encryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", ".env.enc", "Output file path for encrypted secrets")
	_ = encryptCmd.MarkFlagRequired("passphrase")
	rootCmd.AddCommand(encryptCmd)
}
