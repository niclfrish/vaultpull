package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/vault"
)

var (
	notifyNamespace string
	notifyQuiet     bool
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Sync secrets and emit lifecycle notifications",
	Long:  `Pull secrets from Vault and emit INFO/ERROR notifications to stdout.`,
	RunE:  runNotify,
}

func runNotify(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	client, err := vault.New(cfg.VaultAddress, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	var w = os.Stdout
	if notifyQuiet {
		w = nil
	}

	writerSink := sync.NewWriterSink(w)
	notifier := sync.NewNotifier(writerSink)

	secretPath := sync.NamespacedPath(cfg.SecretPath, notifyNamespace)
	secrets, fetchErr := client.GetSecrets(secretPath)
	if fetchErr != nil {
		hook := sync.NotifyOnError(notifier, notifyNamespace, os.Stderr)
		return hook(fetchErr)
	}

	hook := sync.NotifyOnSync(notifier, notifyNamespace, os.Stderr)
	if err := hook(secrets); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "synced %d secret(s) from %s\n", len(secrets), secretPath)
	return nil
}

func init() {
	notifyCmd.Flags().StringVar(&notifyNamespace, "namespace", "", "Vault namespace prefix")
	notifyCmd.Flags().BoolVar(&notifyQuiet, "quiet", false, "Suppress notification output")
	RootCmd.AddCommand(notifyCmd)
}
