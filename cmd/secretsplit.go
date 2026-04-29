package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/sync"
)

func init() {
	var (
		delimiter   string
		keyIndex    int
		valueIndex  int
		onlyKeys    []string
	)

	cmd := &cobra.Command{
		Use:   "secretsplit",
		Short: "Split compound secret values into separate key/value pairs",
		Long: `Reads secrets from stdin-style flags and splits compound values
using a delimiter. The extracted key part becomes a suffix on the
original key, and the extracted value part becomes the new value.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := sync.SplitConfig{
				Delimiter:  delimiter,
				KeyIndex:   keyIndex,
				ValueIndex: valueIndex,
				OnlyKeys:   onlyKeys,
			}

			// Build a small demo map from positional KEY=VALUE args.
			secrets := make(map[string]string)
			for _, arg := range args {
				parts := splitN(arg, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid argument %q: expected KEY=VALUE", arg)
				}
				secrets[parts[0]] = parts[1]
			}

			out, results, summary, err := sync.SplitSecrets(secrets, cfg)
			if err != nil {
				return fmt.Errorf("secretsplit: %w", err)
			}

			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "Split: %d  Skipped: %d\n", summary.Split, summary.Skipped)
			for _, r := range results {
				fmt.Fprintf(w, "  %s -> %s = %s\n", r.OriginalKey, r.NewKey, r.NewValue)
			}
			if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
				for k, v := range out {
					fmt.Fprintf(w, "  [out] %s=%s\n", k, v)
				}
			}
			_ = out
			return nil
		},
	}

	cmd.Flags().StringVar(&delimiter, "delimiter", ":", "delimiter used to split compound values")
	cmd.Flags().IntVar(&keyIndex, "key-index", 0, "part index to use as the new key suffix")
	cmd.Flags().IntVar(&valueIndex, "value-index", 1, "part index to use as the new value")
	cmd.Flags().StringSliceVar(&onlyKeys, "only", nil, "restrict splitting to these keys")
	cmd.Flags().Bool("verbose", false, "print resulting map")

	if err := cmd.MarkFlagRequired("delimiter"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rootCmd.AddCommand(cmd)
}
