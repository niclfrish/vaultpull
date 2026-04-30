package sync

import (
	"fmt"
	"io"
	"os"
)

// PinStage returns a pipeline stage that pins secrets according to cfg.
func PinStage(cfg PinConfig) Stage {
	return Stage{
		Name: "pin",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			out, _, err := PinSecrets(secrets, cfg)
			return out, err
		},
	}
}

// PinAndReport pins secrets and writes a human-readable summary to w.
// If w is nil, os.Stdout is used.
func PinAndReport(secrets map[string]string, cfg PinConfig, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		return nil, fmt.Errorf("PinAndReport: secrets map is nil")
	}

	out, summary, err := PinSecrets(secrets, cfg)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "secret pin: %d pinned, %d missing\n", summary.Pinned, summary.Missing)
	for _, r := range summary.Results {
		if r.Missing {
			fmt.Fprintf(w, "  [missing] %s (wanted version %s)\n", r.Key, r.Version)
		} else {
			fmt.Fprintf(w, "  [pinned]  %s @ %s\n", r.Key, r.Version)
		}
	}

	return out, nil
}
