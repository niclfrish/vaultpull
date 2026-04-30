package sync

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultJoinConfig returns a JoinConfig with sensible defaults.
func DefaultJoinConfig() JoinConfig {
	return JoinConfig{
		Separator:  "_",
		OutputKey:  "",
		StripParts: false,
	}
}

// JoinConfig controls how multiple secret values are joined into one.
type JoinConfig struct {
	// Keys is the ordered list of secret keys whose values will be joined.
	Keys []string
	// Separator is placed between each value.
	Separator string
	// OutputKey is the key under which the joined value is stored.
	OutputKey string
	// StripParts removes the source keys from the result map after joining.
	StripParts bool
}

// JoinSummary holds statistics about a JoinSecrets operation.
type JoinSummary struct {
	Joined  int
	Stripped int
	Skipped int
}

// JoinSecrets concatenates the values of cfg.Keys (in order) with cfg.Separator
// and stores the result under cfg.OutputKey. Missing keys are skipped unless
// strict mode is desired by the caller checking the summary.
func JoinSecrets(secrets map[string]string, cfg JoinConfig) (map[string]string, JoinSummary, error) {
	if secrets == nil {
		return nil, JoinSummary{}, fmt.Errorf("joinSecrets: secrets map is nil")
	}
	if cfg.OutputKey == "" {
		return nil, JoinSummary{}, fmt.Errorf("joinSecrets: OutputKey must not be empty")
	}
	if len(cfg.Keys) == 0 {
		return nil, JoinSummary{}, fmt.Errorf("joinSecrets: Keys must not be empty")
	}

	var parts []string
	var summary JoinSummary

	for _, k := range cfg.Keys {
		v, ok := secrets[k]
		if !ok {
			summary.Skipped++
			continue
		}
		parts = append(parts, v)
		summary.Joined++
	}

	out := make(map[string]string, len(secrets)+1)
	for k, v := range secrets {
		out[k] = v
	}
	out[cfg.OutputKey] = strings.Join(parts, cfg.Separator)

	if cfg.StripParts {
		for _, k := range cfg.Keys {
			if k == cfg.OutputKey {
				continue
			}
			if _, exists := out[k]; exists {
				delete(out, k)
				summary.Stripped++
			}
		}
	}

	return out, summary, nil
}

// JoinSummaryString returns a human-readable summary line.
func JoinSummaryString(s JoinSummary) string {
	parts := []string{
		fmt.Sprintf("joined=%d", s.Joined),
		fmt.Sprintf("skipped=%d", s.Skipped),
		fmt.Sprintf("stripped=%d", s.Stripped),
	}
	sort.Strings(parts)
	return strings.Join(parts, " ")
}
