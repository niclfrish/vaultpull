package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var (
	transformRedactKeys []string
	transformTrimSpace  bool
)

var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Preview value transformations applied to secrets before sync",
	Long: `Applies configured transformations (trim whitespace, redact sensitive keys)
to secrets fetched from Vault and prints the resulting key/value pairs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fns := buildTransformFuncs()
		if len(fns) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No transformations configured. Use --trim or --redact flags.")
			return nil
		}

		// Example static secrets for preview; real implementation would pull from Vault.
		sample := map[string]string{
			"DB_PASSWORD": "  secret123  ",
			"API_TOKEN":   "tok-abc",
			"HOST":        "  localhost  ",
		}

		tr := sync.NewTransformer(fns...)
		result, err := tr.Apply(sample)
		if err != nil {
			fmt.Fprintf(os.Stderr, "transform error: %v\n", err)
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Transformed secrets preview:")
		for k, v := range result {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s=%s\n", k, v)
		}
		return nil
	},
}

func buildTransformFuncs() []sync.TransformFunc {
	var fns []sync.TransformFunc
	if transformTrimSpace {
		fns = append(fns, sync.TrimSpaceTransform())
	}
	if len(transformRedactKeys) > 0 {
		fns = append(fns, sync.RedactTransform(transformRedactKeys, "[REDACTED]"))
	}
	return fns
}

func init() {
	transformCmd.Flags().BoolVar(&transformTrimSpace, "trim", false, "Trim whitespace from secret values")
	transformCmd.Flags().StringSliceVar(&transformRedactKeys, "redact", nil, "Key substrings whose values should be redacted (e.g. password,token)")
	rootCmd.AddCommand(transformCmd)
}
