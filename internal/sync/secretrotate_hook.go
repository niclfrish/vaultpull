package sync

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// CheckRotationAndReport evaluates rotation status for all secrets and writes
// a summary to w (defaults to os.Stdout when nil). It returns an error if any
// secret is stale and failOnStale is true.
func CheckRotationAndReport(secrets map[string]string, cfg RotateConfig, failOnStale bool, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		fmt.Fprintln(w, "rotation: no secrets to check")
		return nil
	}

	statuses := CheckRotation(secrets, cfg)
	fmt.Fprintln(w, RotateSummary(statuses))

	// Collect stale keys in deterministic order for reporting.
	var staleKeys []string
	for k, s := range statuses {
		if s == RotationStale {
			staleKeys = append(staleKeys, k)
		}
	}
	sort.Strings(staleKeys)

	for _, k := range staleKeys {
		fmt.Fprintf(w, "  stale: %s\n", k)
	}

	if failOnStale && len(staleKeys) > 0 {
		return fmt.Errorf("rotation: %d stale secret(s) detected", len(staleKeys))
	}
	return nil
}

// RotationStage returns a pipeline Stage that injects rotation status as
// metadata keys ("__rotation_status.<key>") into the secret map.
func RotationStage(cfg RotateConfig) Stage {
	return Stage{
		Name: "rotation",
		Fn: func(secrets map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(secrets))
			for k, v := range secrets {
				out[k] = v
			}
			statuses := CheckRotation(secrets, cfg)
			for k, s := range statuses {
				out["__rotation_status."+k] = s.String()
			}
			return out, nil
		},
	}
}
