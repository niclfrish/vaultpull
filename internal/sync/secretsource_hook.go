package sync

import (
	"fmt"
	"io"
	"os"
)

// AnnotateStage returns a pipeline stage that injects source annotations.
func AnnotateStage(src SecretSource) func(map[string]string) (map[string]string, error) {
	return func(secrets map[string]string) (map[string]string, error) {
		if secrets == nil {
			return nil, fmt.Errorf("annotate stage: nil secrets")
		}
		return AnnotateWithSource(secrets, src), nil
	}
}

// StripAnnotationsStage returns a pipeline stage that removes source annotations.
func StripAnnotationsStage() func(map[string]string) (map[string]string, error) {
	return func(secrets map[string]string) (map[string]string, error) {
		if secrets == nil {
			return nil, fmt.Errorf("strip annotations stage: nil secrets")
		}
		return StripSourceAnnotations(secrets), nil
	}
}

// ReportSource writes a human-readable source summary to w (defaults to stdout).
// It returns the secrets map unchanged.
func ReportSource(secrets map[string]string, src SecretSource, w io.Writer) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	if secrets == nil {
		return nil, fmt.Errorf("report source: nil secrets")
	}
	_, err := fmt.Fprintf(w, "[source] %s\n", SourceSummary(src))
	if err != nil {
		return nil, fmt.Errorf("report source: write failed: %w", err)
	}
	return secrets, nil
}
