package sync

import (
	"fmt"
	"io"
	"os"
)

// PromoteStage returns a pipeline stage that promotes secrets between prefixes.
func PromoteStage(cfg PromoteConfig) Stage {
	return Stage{
		Name: "promote",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			if len(secrets) == 0 {
				return secrets, nil
			}
			out, _, err := PromoteSecrets(secrets, cfg)
			return out, err
		},
	}
}

// PromoteAndReport promotes secrets and writes a summary to w.
// If w is nil, os.Stdout is used.
func PromoteAndReport(secrets map[string]string, cfg PromoteConfig, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		fmt.Fprintln(w, "promote: no secrets provided")
		return nil, nil
	}
	out, result, err := PromoteSecrets(secrets, cfg)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(w, "promote: %s\n", PromoteSummary(result))
	if cfg.DryRun {
		fmt.Fprintln(w, "promote: dry-run mode — no changes written")
	}
	return out, nil
}
