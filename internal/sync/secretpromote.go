package sync

import (
	"fmt"
	"strings"
)

// DefaultPromoteConfig returns a sensible default promotion configuration.
func DefaultPromoteConfig() PromoteConfig {
	return PromoteConfig{
		FromPrefix: "",
		ToPrefix:   "",
		Overwrite:  false,
		DryRun:     false,
	}
}

// PromoteConfig controls how secrets are promoted between namespaces/prefixes.
type PromoteConfig struct {
	FromPrefix string
	ToPrefix   string
	Overwrite  bool
	DryRun     bool
}

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	Promoted  []string
	Skipped   []string
	Overwrote []string
}

// PromoteSummary returns a human-readable summary of the promotion result.
func PromoteSummary(r PromoteResult) string {
	return fmt.Sprintf("promoted=%d skipped=%d overwrote=%d",
		len(r.Promoted), len(r.Skipped), len(r.Overwrote))
}

// PromoteSecrets copies secrets matching FromPrefix into ToPrefix.
// If Overwrite is false, existing keys in the destination are skipped.
func PromoteSecrets(secrets map[string]string, cfg PromoteConfig) (map[string]string, PromoteResult, error) {
	if secrets == nil {
		return nil, PromoteResult{}, fmt.Errorf("promote: secrets map is nil")
	}
	if cfg.ToPrefix == "" {
		return nil, PromoteResult{}, fmt.Errorf("promote: ToPrefix must not be empty")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var result PromoteResult

	for k, v := range secrets {
		if cfg.FromPrefix != "" && !strings.HasPrefix(k, cfg.FromPrefix) {
			continue
		}
		base := strings.TrimPrefix(k, cfg.FromPrefix)
		destKey := cfg.ToPrefix + base

		if _, exists := out[destKey]; exists && !cfg.Overwrite {
			result.Skipped = append(result.Skipped, destKey)
			continue
		}

		if cfg.DryRun {
			result.Promoted = append(result.Promoted, destKey)
			continue
		}

		if _, exists := out[destKey]; exists {
			result.Overwrote = append(result.Overwrote, destKey)
		} else {
			result.Promoted = append(result.Promoted, destKey)
		}
		out[destKey] = v
	}

	return out, result, nil
}
