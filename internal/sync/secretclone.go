package sync

import (
	"fmt"
	"strings"
)

// DefaultCloneConfig returns a CloneConfig with sensible defaults.
func DefaultCloneConfig() CloneConfig {
	return CloneConfig{
		Overwrite: false,
		Separator: "_",
	}
}

// CloneConfig controls how secrets are cloned under a new prefix.
type CloneConfig struct {
	// SourcePrefix filters which keys to clone. Empty means all keys.
	SourcePrefix string
	// DestPrefix is prepended to each cloned key.
	DestPrefix string
	// Overwrite allows the clone to overwrite existing destination keys.
	Overwrite bool
	// Separator is placed between DestPrefix and the original key.
	Separator string
}

// CloneResult holds the outcome of a CloneSecrets call.
type CloneResult struct {
	Cloned    int
	Skipped   int
	Overwritten int
}

// CloneSummary returns a human-readable summary of the clone result.
func CloneSummary(r CloneResult) string {
	return fmt.Sprintf("cloned=%d skipped=%d overwritten=%d", r.Cloned, r.Skipped, r.Overwritten)
}

// CloneSecrets copies secrets whose keys match SourcePrefix into new keys
// under DestPrefix. The original keys are preserved in the output map.
func CloneSecrets(secrets map[string]string, cfg CloneConfig) (map[string]string, CloneResult, error) {
	if secrets == nil {
		return nil, CloneResult{}, fmt.Errorf("secretclone: secrets map is nil")
	}
	if cfg.DestPrefix == "" {
		return nil, CloneResult{}, fmt.Errorf("secretclone: DestPrefix must not be empty")
	}
	sep := cfg.Separator

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var res CloneResult
	for k, v := range secrets {
		if cfg.SourcePrefix != "" && !strings.HasPrefix(k, cfg.SourcePrefix) {
			continue
		}
		newKey := cfg.DestPrefix + sep + k
		if _, exists := out[newKey]; exists {
			if !cfg.Overwrite {
				res.Skipped++
				continue
			}
			res.Overwritten++
		} else {
			res.Cloned++
		}
		out[newKey] = v
	}
	return out, res, nil
}
