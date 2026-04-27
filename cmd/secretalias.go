package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Rename secret keys using alias mappings before writing",
	Long: `Apply key aliases to fetched secrets.
Aliases are specified as KEY=ALIAS pairs. By default the original key is removed.
Use --keep-original to retain both the source and alias keys.`,
	RunE: runAlias,
}

func runAlias(cmd *cobra.Command, _ []string) error {
	aliasFlags, _ := cmd.Flags().GetStringArray("alias")
	keepOriginal, _ := cmd.Flags().GetBool("keep-original")

	aliasMap, err := parseAliasFlags(aliasFlags)
	if err != nil {
		return err
	}

	cfg := sync.AliasConfig{
		Aliases:      aliasMap,
		KeepOriginal: keepOriginal,
	}

	// Demo: apply to a static set; real integration would wire into the pipeline.
	secrets := map[string]string{"EXAMPLE_KEY": "example_value"}
	result, err := sync.AliasAndReport(secrets, cfg, os.Stdout)
	if err != nil {
		return fmt.Errorf("alias: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", sync.AliasSummary(cfg))
	for k, v := range result {
		fmt.Fprintf(os.Stdout, "  %s=%s\n", k, v)
	}
	return nil
}

func parseAliasFlags(flags []string) (map[string]string, error) {
	result := make(map[string]string, len(flags))
	for _, f := range flags {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("alias: invalid format %q, expected KEY=ALIAS", f)
		}
		src := strings.TrimSpace(parts[0])
		dst := strings.TrimSpace(parts[1])
		if src == "" || dst == "" {
			return nil, fmt.Errorf("alias: key and alias must be non-empty in %q", f)
		}
		result[src] = dst
	}
	return result, nil
}

func init() {
	aliasCmd.Flags().StringArray("alias", nil, "Key alias mapping as KEY=ALIAS (repeatable)")
	aliasCmd.Flags().Bool("keep-original", false, "Retain the original key alongside the alias")
	rootCmd.AddCommand(aliasCmd)
}
