package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/sync"
)

var secretsourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Annotate or inspect secret source provenance",
	Long:  `Annotate fetched secrets with provenance metadata (type, location, namespace, timestamp).`,
	RunE:  runSecretSource,
}

var (
	sourceType      string
	sourceLocation  string
	sourceNamespace string
	sourceStrip     bool
	sourceReport    bool
)

func runSecretSource(cmd *cobra.Command, args []string) error {
	srcType := sync.SourceType(sourceType)
	switch srcType {
	case sync.SourceTypeVault, sync.SourceTypeEnv, sync.SourceTypeFile, sync.SourceTypeUnknown:
		// valid
	default:
		return fmt.Errorf("unknown source type %q; use vault, env, file, or unknown", sourceType)
	}

	src := sync.SecretSource{
		Type:      srcType,
		Location:  sourceLocation,
		FetchedAt: time.Now().UTC(),
		Namespace: sourceNamespace,
	}

	if sourceReport {
		fmt.Fprintf(os.Stdout, "%s\n", sync.SourceSummary(src))
		return nil
	}

	if sourceStrip {
		fmt.Fprintln(os.Stdout, "[source] strip mode: annotations would be removed from secrets")
		return nil
	}

	fmt.Fprintf(os.Stdout, "[source] annotate mode: %s\n", sync.SourceSummary(src))
	return nil
}

func init() {
	secretSourceFlags := secretsourceCmd.Flags()
	secretSourceFlags.StringVar(&sourceType, "type", "vault", "Source type: vault, env, file, unknown")
	secretSourceFlags.StringVar(&sourceLocation, "location", "", "Source location (e.g. secret/data/app)")
	secretSourceFlags.StringVar(&sourceNamespace, "namespace", "", "Optional namespace for the source")
	secretSourceFlags.BoolVar(&sourceStrip, "strip", false, "Strip source annotations from secrets")
	secretSourceFlags.BoolVar(&sourceReport, "report", false, "Print source summary and exit")

	rootCmd.AddCommand(secretsourceCmd)
}
