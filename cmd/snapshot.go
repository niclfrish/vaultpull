package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Show diff between last snapshot and current Vault secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		namespace, _ := cmd.Flags().GetString("namespace")
		snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")

		effectivePath := sync.NamespacedPath(cfg.SecretPath, namespace)
		current, err := client.GetSecrets(effectivePath)
		if err != nil {
			return fmt.Errorf("fetch secrets: %w", err)
		}

		store, err := sync.NewSnapshotStore(snapshotDir)
		if err != nil {
			return fmt.Errorf("snapshot store: %w", err)
		}

		prev, err := store.Load(cfg.SecretPath, namespace)
		if err != nil {
			return fmt.Errorf("load snapshot: %w", err)
		}

		saveFlag, _ := cmd.Flags().GetBool("save")
		if saveFlag {
			snap := sync.Snapshot{
				Path:      cfg.SecretPath,
				Namespace: namespace,
				Secrets:   current,
			}
			if err := store.Save(snap); err != nil {
				return fmt.Errorf("save snapshot: %w", err)
			}
			fmt.Fprintln(os.Stdout, "Snapshot saved.")
			return nil
		}

		summary := sync.SnapshotSummary(prev, current)
		fmt.Fprintln(os.Stdout, summary)
		return nil
	},
}

func init() {
	snapshotCmd.Flags().String("namespace", "", "Vault namespace prefix")
	snapshotCmd.Flags().String("snapshot-dir", ".vaultpull/snapshots", "Directory to store snapshots")
	snapshotCmd.Flags().Bool("save", false, "Save current secrets as new snapshot baseline")
	rootCmd.AddCommand(snapshotCmd)
}
