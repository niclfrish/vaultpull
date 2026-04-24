package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

var tagFilterCmd = &cobra.Command{
	Use:   "tagfilter",
	Short: "Sync only secrets whose keys carry specified tags (KEY#tag1,tag2)",
	RunE:  runTagFilter,
}

func runTagFilter(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	rawTags, _ := cmd.Flags().GetString("tags")
	outPath, _ := cmd.Flags().GetString("output")

	tags := splitTrimmed(rawTags, ",")
	if len(tags) == 0 {
		return fmt.Errorf("--tags must specify at least one tag")
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.GetSecrets(cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}

	filtered := sync.NewTagFilter(tags).Apply(secrets)
	if len(filtered) == 0 {
		fmt.Fprintln(os.Stderr, "warning: no secrets matched the provided tags")
	}

	fmt.Fprintf(cmd.OutOrStdout(), "# filtered by tags: %s\n", strings.Join(tags, ", "))
	for k, v := range filtered {
		fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
	}

	if outPath != "" {
		_ = outPath // wired to env writer in full integration
	}

	return nil
}

func init() {
	tagFilterCmd.Flags().String("tags", "", "Comma-separated list of tags to keep (e.g. prod,staging)")
	tagFilterCmd.Flags().String("output", "", "Path to write filtered .env file (optional)")
	_ = tagFilterCmd.MarkFlagRequired("tags")
	rootCmd.AddCommand(tagFilterCmd)
}
