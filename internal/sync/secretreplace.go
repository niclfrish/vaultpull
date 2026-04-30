package sync

import (
	"fmt"
	"strings"
)

// DefaultReplaceConfig returns a ReplaceConfig with sensible defaults.
func DefaultReplaceConfig() ReplaceConfig {
	return ReplaceConfig{
		CaseSensitive: true,
	}
}

// ReplaceConfig controls how value replacements are applied.
type ReplaceConfig struct {
	// Replacements maps old substrings to new substrings.
	Replacements map[string]string
	// OnlyKeys restricts replacement to the given keys. Empty means all keys.
	OnlyKeys []string
	// CaseSensitive determines whether matching is case-sensitive.
	CaseSensitive bool
}

// ReplaceSummary holds statistics from a ReplaceSecrets call.
type ReplaceSummary struct {
	Modified int
	Skipped  int
}

// ReplaceSecrets applies substring replacements to secret values.
// If cfg.OnlyKeys is non-empty, only those keys are modified.
func ReplaceSecrets(secrets map[string]string, cfg ReplaceConfig) (map[string]string, ReplaceSummary, error) {
	if secrets == nil {
		return nil, ReplaceSummary{}, fmt.Errorf("replaceSecrets: secrets map is nil")
	}
	if len(cfg.Replacements) == 0 {
		return secrets, ReplaceSummary{Skipped: len(secrets)}, nil
	}

	filter := make(map[string]bool, len(cfg.OnlyKeys))
	for _, k := range cfg.OnlyKeys {
		filter[k] = true
	}

	out := make(map[string]string, len(secrets))
	var summary ReplaceSummary

	for k, v := range secrets {
		if len(filter) > 0 && !filter[k] {
			out[k] = v
			summary.Skipped++
			continue
		}
		original := v
		for old, newVal := range cfg.Replacements {
			if cfg.CaseSensitive {
				v = strings.ReplaceAll(v, old, newVal)
			} else {
				v = replaceAllInsensitive(v, old, newVal)
			}
		}
		out[k] = v
		if v != original {
			summary.Modified++
		} else {
			summary.Skipped++
		}
	}
	return out, summary, nil
}

// replaceAllInsensitive replaces all case-insensitive occurrences of old with newVal in s.
func replaceAllInsensitive(s, old, newVal string) string {
	if old == "" {
		return s
	}
	lower := strings.ToLower(s)
	lowerOld := strings.ToLower(old)
	var result strings.Builder
	for {
		idx := strings.Index(lower, lowerOld)
		if idx < 0 {
			result.WriteString(s)
			break
		}
		result.WriteString(s[:idx])
		result.WriteString(newVal)
		s = s[idx+len(old):]
		lower = lower[idx+len(old):]
	}
	return result.String()
}
