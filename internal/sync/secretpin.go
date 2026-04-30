package sync

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultPinConfig returns a PinConfig with sensible defaults.
func DefaultPinConfig() PinConfig {
	return PinConfig{
		AnnotationKey: "__pinned_version",
	}
}

// PinConfig controls how secret version pinning is applied.
type PinConfig struct {
	// Pins maps secret key names to required version strings.
	Pins map[string]string
	// AnnotationKey is the key used to store the pinned version annotation.
	AnnotationKey string
	// StrictMode causes an error when a key in Pins is not found in secrets.
	StrictMode bool
}

// PinResult holds the outcome of a single pin operation.
type PinResult struct {
	Key     string
	Version string
	Missing bool
}

// PinSummary describes the overall result of PinSecrets.
type PinSummary struct {
	Pinned  int
	Missing int
	Results []PinResult
}

// PinSecrets annotates secrets with pinned version information and
// optionally errors when a pinned key is absent.
func PinSecrets(secrets map[string]string, cfg PinConfig) (map[string]string, PinSummary, error) {
	if secrets == nil {
		return nil, PinSummary{}, fmt.Errorf("pinSecrets: secrets map is nil")
	}
	if len(cfg.Pins) == 0 {
		return secrets, PinSummary{}, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var summary PinSummary
	keys := make([]string, 0, len(cfg.Pins))
	for k := range cfg.Pins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		version := cfg.Pins[key]
		if _, ok := out[key]; !ok {
			summary.Missing++
			summary.Results = append(summary.Results, PinResult{Key: key, Version: version, Missing: true})
			if cfg.StrictMode {
				return nil, summary, fmt.Errorf("pinSecrets: pinned key %q not found in secrets", key)
			}
			continue
		}
		annotation := fmt.Sprintf("%s@%s", key, version)
		existing := out[cfg.AnnotationKey]
		if existing == "" {
			out[cfg.AnnotationKey] = annotation
		} else if !strings.Contains(existing, annotation) {
			out[cfg.AnnotationKey] = existing + "," + annotation
		}
		summary.Pinned++
		summary.Results = append(summary.Results, PinResult{Key: key, Version: version})
	}

	return out, summary, nil
}
