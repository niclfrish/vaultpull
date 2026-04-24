package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

func init() {
	var (
		requiredKeys  string
		prefixFilter  string
		excludeKeys   string
		enableTrim    bool
		enableTruncate bool
	)

	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Run a configurable processing pipeline on fetched secrets",
		Long: `Fetch secrets from Vault and pass them through a series of
built-in pipeline stages (filter, transform, validate) before writing
them to the target .env file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := sync.NewPipeline()

			// Filter stage
			if prefixFilter != "" || excludeKeys != "" {
				criteria := sync.FilterCriteria{Prefix: prefixFilter}
				if excludeKeys != "" {
					criteria.Exclude = splitTrimmed(excludeKeys, ",")
				}
				p.AddStage("filter", sync.FilterStage(sync.NewFilter(criteria)))
			}

			// Transform stage
			if enableTrim {
				p.AddStage("trim", sync.TransformStage(sync.NewTransformer(sync.TrimSpaceTransform)))
			}

			// Truncate stage
			if enableTruncate {
				p.AddStage("truncate", sync.TruncateStage())
			}

			// Required-keys validation stage
			if requiredKeys != "" {
				keys := splitTrimmed(requiredKeys, ",")
				p.AddStage("require", sync.RequiredKeysStage(keys...))
			}

			fmt.Fprintf(os.Stdout, "Pipeline configured with %d stage(s): %s\n",
				p.StageCount(), strings.Join(p.StageNames(), " → "))
			return nil
		},
	}

	cmd.Flags().StringVar(&requiredKeys, "require", "", "Comma-separated list of required secret keys")
	cmd.Flags().StringVar(&prefixFilter, "prefix", "", "Only keep keys with this prefix")
	cmd.Flags().StringVar(&excludeKeys, "exclude", "", "Comma-separated list of keys to exclude")
	cmd.Flags().BoolVar(&enableTrim, "trim", false, "Trim whitespace from all values")
	cmd.Flags().BoolVar(&enableTruncate, "truncate", false, "Truncate values exceeding the default max length")

	rootCmd.AddCommand(cmd)
}
