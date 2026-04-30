package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/sync"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Pin secret keys to specific versions and annotate the output",
	RunE:  runPin,
}

func runPin(cmd *cobra.Command, _ []string) error {
	pinFlags, _ := cmd.Flags().GetStringSlice("pin")
	strict, _ := cmd.Flags().GetBool("strict")
	annotationKey, _ := cmd.Flags().GetString("annotation-key")

	pins, err := parsePinFlags(pinFlags)
	if err != nil {
		return fmt.Errorf("pin: invalid --pin flag: %w", err)
	}

	cfg := sync.DefaultPinConfig()
	cfg.Pins = pins
	cfg.StrictMode = strict
	if annotationKey != "" {
		cfg.AnnotationKey = annotationKey
	}

	// Use a minimal placeholder secrets map for demonstration.
	// In production this would be populated from the vault client.
	secrets := map[string]string{}

	_, err = sync.PinAndReport(secrets, cfg, os.Stdout)
	if err != nil {
		return fmt.Errorf("pin: %w", err)
	}
	return nil
}

// parsePinFlags converts "KEY=version" pairs into a map.
func parsePinFlags(flags []string) (map[string]string, error) {
	pins := make(map[string]string, len(flags))
	for _, f := range flags {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("expected KEY=version, got %q", f)
		}
		pins[parts[0]] = parts[1]
	}
	return pins, nil
}

func init() {
	pinCmd.Flags().StringSlice("pin", nil, "Pin a key to a version: KEY=version (repeatable)")
	pinCmd.Flags().Bool("strict", false, "Error if a pinned key is missing from secrets")
	pinCmd.Flags().String("annotation-key", "", "Override the annotation key name (default: __pinned_version)")
	rootCmd.AddCommand(pinCmd)
}
