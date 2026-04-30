package sync

import (
	"fmt"
	"io"
	"os"
)

// RenameStage returns a pipeline stage that renames secret keys according to cfg.
func RenameStage(cfg RenameConfig) PipelineStage {
	return PipelineStage{
		Name: "rename",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			if secrets == nil {
				return secrets, nil
			}
			out, _, err := RenameSecrets(secrets, cfg)
			return out, err
		},
	}
}

// RenameAndReport renames secrets and writes a human-readable summary to w.
// If w is nil, os.Stdout is used.
func RenameAndReport(secrets map[string]string, cfg RenameConfig, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		fmt.Fprintln(w, "rename: no secrets provided")
		return nil, fmt.Errorf("secrets map is nil")
	}

	out, summary, err := RenameSecrets(secrets, cfg)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "rename: %d key(s) renamed, %d rule(s) had no match\n", summary.Renamed, summary.Missed)
	return out, nil
}
