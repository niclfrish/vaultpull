package sync

import (
	"fmt"
	"strings"
)

// DefaultSplitConfig returns a SplitConfig with sensible defaults.
func DefaultSplitConfig() SplitConfig {
	return SplitConfig{
		Delimiter:  ":",
		KeyIndex:   0,
		ValueIndex: 1,
	}
}

// SplitConfig controls how SplitSecrets parses compound values.
type SplitConfig struct {
	// Delimiter separates the parts of a compound value.
	Delimiter string
	// KeyIndex is the part index used as the new key suffix.
	KeyIndex int
	// ValueIndex is the part index used as the new value.
	ValueIndex int
	// OnlyKeys, when non-empty, restricts splitting to those keys.
	OnlyKeys []string
}

// SplitResult holds a single split outcome.
type SplitResult struct {
	OriginalKey string
	NewKey      string
	NewValue    string
}

// SplitSummary describes the outcome of a SplitSecrets call.
type SplitSummary struct {
	Split   int
	Skipped int
}

// SplitSecrets splits compound secret values into separate key/value pairs.
// For each matching key whose value contains the delimiter, the value is split
// and a new entry is added using the original key plus the extracted key part.
func SplitSecrets(secrets map[string]string, cfg SplitConfig) (map[string]string, []SplitResult, SplitSummary, error) {
	if secrets == nil {
		return nil, nil, SplitSummary{}, fmt.Errorf("splitSecrets: secrets map is nil")
	}
	if cfg.Delimiter == "" {
		return nil, nil, SplitSummary{}, fmt.Errorf("splitSecrets: delimiter must not be empty")
	}

	allowed := make(map[string]bool, len(cfg.OnlyKeys))
	for _, k := range cfg.OnlyKeys {
		allowed[k] = true
	}

	out := make(map[string]string, len(secrets))
	var results []SplitResult
	var summary SplitSummary

	for k, v := range secrets {
		if len(allowed) > 0 && !allowed[k] {
			out[k] = v
			summary.Skipped++
			continue
		}
		parts := strings.Split(v, cfg.Delimiter)
		if cfg.KeyIndex >= len(parts) || cfg.ValueIndex >= len(parts) {
			out[k] = v
			summary.Skipped++
			continue
		}
		newKey := k + "_" + strings.TrimSpace(parts[cfg.KeyIndex])
		newVal := strings.TrimSpace(parts[cfg.ValueIndex])
		out[newKey] = newVal
		results = append(results, SplitResult{OriginalKey: k, NewKey: newKey, NewValue: newVal})
		summary.Split++
	}

	return out, results, summary, nil
}
