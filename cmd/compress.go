package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compress long secret values using gzip+base64 encoding",
	Long: `Reads secrets from Vault and compresses values that exceed a
configurable length threshold using gzip compression encoded as base64.
Useful for storing large certificate or key material in .env files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		minLen, _ := cmd.Flags().GetInt("min-length")
		decompress, _ := cmd.Flags().GetBool("decompress")
		key, _ := cmd.Flags().GetString("key")

		if key == "" {
			return fmt.Errorf("--key is required")
		}

		value, _ := cmd.Flags().GetString("value")
		if value == "" {
			return fmt.Errorf("--value is required")
		}

		if decompress {
			out, err := sync.DecompressValue(value)
			if err != nil {
				return fmt.Errorf("decompress %q: %w", key, err)
			}
			fmt.Fprintf(os.Stdout, "%s=%s\n", key, out)
			return nil
		}

		cfg := sync.CompressConfig{MinLength: minLen}
		out, err := sync.CompressValue(value, cfg)
		if err != nil {
			return fmt.Errorf("compress %q: %w", key, err)
		}
		fmt.Fprintf(os.Stdout, "%s=%s\n", key, out)
		return nil
	},
}

func init() {
	compressCmd.Flags().Int("min-length", sync.DefaultCompressConfig().MinLength,
		"minimum value length to trigger compression")
	compressCmd.Flags().Bool("decompress", false,
		"decompress a previously compressed value")
	compressCmd.Flags().String("key", "", "secret key name (required)")
	compressCmd.Flags().String("value", "", "secret value to compress/decompress (required)")

	rootCmd.AddCommand(compressCmd)
}
