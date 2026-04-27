package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

// priorityCmd merges secrets from multiple named sources with explicit priorities.
var priorityCmd = &cobra.Command{
	Use:   "priority",
	Short: "Merge secrets from multiple sources respecting priority order",
	Long: `Merge secrets from multiple named sources.
Each source is specified as NAME:PRIORITY:key=val,key=val.
Priority 1 is highest. Conflicts are annotated with --conflict-prefix.`,
	Example: `  vaultpull priority \
    --source "vault:1:DB_PASS=secret,API_KEY=abc" \
    --source "local:2:DB_PASS=override" \
    --conflict-prefix __conflict_`,
	RunE: runPriority,
}

func runPriority(cmd *cobra.Command, _ []string) error {
	rawSources, _ := cmd.Flags().GetStringArray("source")
	conflictPrefix, _ := cmd.Flags().GetString("conflict-prefix")

	if len(rawSources) == 0 {
		return fmt.Errorf("at least one --source is required")
	}

	sources, err := parsePrioritySources(rawSources)
	if err != nil {
		return err
	}

	cfg := sync.PriorityConfig{ConflictPrefix: conflictPrefix}
	merged, err := sync.MergeByPriority(cfg, sources)
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	for k, v := range merged {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	fmt.Fprint(w, sync.PrioritySummary(sources, merged))
	return nil
}

// parsePrioritySources parses "NAME:PRIORITY:key=val,key=val" entries.
func parsePrioritySources(raw []string) ([]sync.PrioritySource, error) {
	sources := make([]sync.PrioritySource, 0, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid source format %q: want NAME:PRIORITY:key=val,...", r)
		}
		prio, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid priority %q in source %q: %w", parts[1], parts[0], err)
		}
		secrets := map[string]string{}
		for _, pair := range strings.Split(parts[2], ",") {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 || kv[0] == "" {
				continue
			}
			secrets[kv[0]] = kv[1]
		}
		sources = append(sources, sync.PrioritySource{
			Name:     parts[0],
			Priority: prio,
			Secrets:  secrets,
		})
	}
	return sources, nil
}

func init() {
	_ = os.Stderr // ensure os import used
	priorityCmd.Flags().StringArray("source", nil, "Source in NAME:PRIORITY:key=val format (repeatable)")
	priorityCmd.Flags().String("conflict-prefix", "__conflict_", "Prefix for conflict annotation keys (empty to disable)")
	RootCmd.AddCommand(priorityCmd)
}
