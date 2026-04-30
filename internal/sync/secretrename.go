package sync

import (
	"fmt"
	"strings"
)

// DefaultRenameConfig returns a RenameConfig with sensible defaults.
func DefaultRenameConfig() RenameConfig {
	return RenameConfig{
		CaseSensitive: true,
	}
}

// RenameConfig controls how keys are renamed in a secrets map.
type RenameConfig struct {
	// Rules maps old key names to new key names.
	Rules map[string]string
	// CaseSensitive controls whether key matching is case-sensitive.
	CaseSensitive bool
	// FailOnMissing causes RenameSecrets to return an error if a rule's source key is absent.
	FailOnMissing bool
}

// RenameSummary holds statistics from a RenameSecrets call.
type RenameSummary struct {
	Renamed int
	Missed  int
}

// RenameSecrets applies rename rules to secrets, returning the updated map and a summary.
func RenameSecrets(secrets map[string]string, cfg RenameConfig) (map[string]string, RenameSummary, error) {
	if secrets == nil {
		return nil, RenameSummary{}, fmt.Errorf("secrets map is nil")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var summary RenameSummary

	for oldKey, newKey := range cfg.Rules {
		matched := ""
		if cfg.CaseSensitive {
			if _, ok := out[oldKey]; ok {
				matched = oldKey
			}
		} else {
			for k := range out {
				if strings.EqualFold(k, oldKey) {
					matched = k
					break
				}
			}
		}

		if matched == "" {
			summary.Missed++
			if cfg.FailOnMissing {
				return nil, summary, fmt.Errorf("rename: source key %q not found", oldKey)
			}
			continue
		}

		out[newKey] = out[matched]
		delete(out, matched)
		summary.Renamed++
	}

	return out, summary, nil
}
