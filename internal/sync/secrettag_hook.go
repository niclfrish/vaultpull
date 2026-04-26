package sync

import (
	"fmt"
	"io"
	"os"
)

// TagAndReport injects metadata tags into secrets and optionally writes
// a summary of injected keys to w. If w is nil, os.Stdout is used.
func TagAndReport(secrets map[string]string, cfg SecretTagConfig, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}

	tagged, err := TagSecrets(secrets, cfg)
	if err != nil {
		return nil, fmt.Errorf("TagAndReport: %w", err)
	}

	injected := len(tagged) - len(secrets)
	if injected > 0 {
		fmt.Fprintf(w, "[tag] injected %d metadata key(s) with prefix %q\n", injected, cfg.Prefix)
	}

	return tagged, nil
}

// TagStage returns a pipeline Stage that injects metadata tags into secrets.
func TagStage(cfg SecretTagConfig) Stage {
	return Stage{
		Name: "tag",
		Run: func(secrets map[string]string) (map[string]string, error) {
			return TagSecrets(secrets, cfg)
		},
	}
}

// StripTagStage returns a pipeline Stage that removes metadata tag keys.
func StripTagStage(prefix string) Stage {
	return Stage{
		Name: "strip-tags",
		Run: func(secrets map[string]string) (map[string]string, error) {
			return StripTagKeys(secrets, prefix), nil
		},
	}
}
