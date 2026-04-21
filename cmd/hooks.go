package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var (
	hookRequireKeys  []string
	hookCountLimit   int
	hookEnableLog    bool
)

// buildHookRunner constructs a HookRunner from the CLI flags registered by
// this file. It is called by the root sync command before running the pipeline.
func buildHookRunner() *sync.HookRunner {
	r := sync.NewHookRunner()

	if hookEnableLog {
		r.Register(sync.HookPreFetch, sync.LoggingHook(os.Stdout))
		r.Register(sync.HookPostFetch, sync.LoggingHook(os.Stdout))
		r.Register(sync.HookPreApply, sync.LoggingHook(os.Stdout))
		r.Register(sync.HookPostApply, sync.LoggingHook(os.Stdout))
	}

	if len(hookRequireKeys) > 0 {
		keys := splitTrimmed(hookRequireKeys)
		r.Register(sync.HookPostFetch, sync.RequireKeysHook(keys))
	}

	if hookCountLimit > 0 {
		r.Register(sync.HookPostFetch, sync.CountLimitHook(hookCountLimit))
	}

	return r
}

// splitTrimmed flattens a slice of potentially comma-separated strings into
// individual trimmed tokens (supports both --flag a,b and --flag a --flag b).
func splitTrimmed(raw []string) []string {
	var out []string
	for _, s := range raw {
		for _, part := range strings.Split(s, ",") {
			if t := strings.TrimSpace(part); t != "" {
				out = append(out, t)
			}
		}
	}
	return out
}

func init() {
	hooksCmd := &cobra.Command{
		Use:   "hooks",
		Short: "List configured lifecycle hooks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Configured hooks:")
			if hookEnableLog {
				fmt.Fprintln(cmd.OutOrStdout(), "  logging: enabled on all events")
			}
			if len(hookRequireKeys) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  require-keys: %s\n", strings.Join(splitTrimmed(hookRequireKeys), ", "))
			}
			if hookCountLimit > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  count-limit: %d\n", hookCountLimit)
			}
			return nil
		},
	}

	hooksCmd.Flags().StringArrayVar(&hookRequireKeys, "require-key", nil, "require key(s) to be present after fetch (comma-separated or repeated)")
	hooksCmd.Flags().IntVar(&hookCountLimit, "count-limit", 0, "fail if fetched secrets exceed this count (0 = disabled)")
	hooksCmd.Flags().BoolVar(&hookEnableLog, "log", false, "enable lifecycle event logging to stdout")

	rootCmd.AddCommand(hooksCmd)
}
