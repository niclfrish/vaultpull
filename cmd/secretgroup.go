package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

var (
	groupSeparator string
	groupMaxDepth  int
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Display Vault secrets grouped by key prefix",
	Long: `Fetches secrets from Vault and displays them partitioned
into groups based on a key prefix separator (default: "_").`,
	RunE: runGroup,
}

func runGroup(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("group: config: %w", err)
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("group: vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("group: fetch secrets: %w", err)
	}

	groupCfg := sync.GroupConfig{
		Separator: groupSeparator,
		MaxDepth:  groupMaxDepth,
	}

	groups, err := sync.GroupSecrets(secrets, groupCfg)
	if err != nil {
		return fmt.Errorf("group: %w", err)
	}

	fmt.Fprintln(os.Stdout, sync.GroupSummary(groups))
	for _, g := range groups {
		label := g.Prefix
		if label == "" {
			label = "(ungrouped)"
		}
		fmt.Fprintf(os.Stdout, "\n[%s]\n", label)
		for k, v := range g.Secrets {
			fmt.Fprintf(os.Stdout, "  %s = %s\n", k, v)
		}
	}
	return nil
}

func init() {
	groupCmd.Flags().StringVar(&groupSeparator, "separator", "_", "Key segment separator used for grouping")
	groupCmd.Flags().IntVar(&groupMaxDepth, "max-depth", 1, "Number of prefix segments to use for group name")
	rootCmd.AddCommand(groupCmd)
}
