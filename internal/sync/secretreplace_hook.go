package sync

import (
	"fmt"
	"io"
	"os"
)

// ReplaceStage returns a pipeline stage that applies value replacements.
func ReplaceStage(cfg ReplaceConfig) PipelineStage {
	return PipelineStage{
		Name: "replace",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			out, _, err := ReplaceSecrets(secrets, cfg)
			return out, err
		},
	}
}

// ReplaceAndReport applies replacements to secrets and writes a summary to w.
// If w is nil, os.Stdout is used.
func ReplaceAndReport(secrets map[string]string, cfg ReplaceConfig, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		return nil, fmt.Errorf("ReplaceAndReport: secrets is nil")
	}
	out, summary, err := ReplaceSecrets(secrets, cfg)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(w, "replace: modified=%d skipped=%d\n", summary.Modified, summary.Skipped)
	return out, nil
}
