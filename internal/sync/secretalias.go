package sync

import (
	"fmt"
	"strings"
)

// AliasConfig defines a mapping from source key to one or more alias keys.
type AliasConfig struct {
	// Aliases maps original key names to their desired alias names.
	Aliases map[string]string
	// KeepOriginal controls whether the original key is retained alongside the alias.
	KeepOriginal bool
}

// DefaultAliasConfig returns a sensible default AliasConfig.
func DefaultAliasConfig() AliasConfig {
	return AliasConfig{
		Aliases:      make(map[string]string),
		KeepOriginal: false,
	}
}

// ApplyAliases renames keys in secrets according to the provided AliasConfig.
// If KeepOriginal is true, both the original and aliased key are present.
// Returns an error if an alias target already exists and would cause a conflict.
func ApplyAliases(secrets map[string]string, cfg AliasConfig) (map[string]string, error) {
	if secrets == nil {
		return nil, fmt.Errorf("secretalias: secrets map is nil")
	}
	if len(cfg.Aliases) == 0 {
		return secrets, nil
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	for src, dst := range cfg.Aliases {
		dst = strings.TrimSpace(dst)
		if dst == "" {
			return nil, fmt.Errorf("secretalias: alias target for %q is empty", src)
		}
		val, ok := result[src]
		if !ok {
			continue
		}
		if _, conflict := result[dst]; conflict && dst != src {
			return nil, fmt.Errorf("secretalias: alias target %q already exists", dst)
		}
		result[dst] = val
		if !cfg.KeepOriginal && dst != src {
			delete(result, src)
		}
	}
	return result, nil
}

// AliasSummary returns a human-readable summary of applied aliases.
func AliasSummary(cfg AliasConfig) string {
	if len(cfg.Aliases) == 0 {
		return "no aliases configured"
	}
	parts := make([]string, 0, len(cfg.Aliases))
	for src, dst := range cfg.Aliases {
		parts = append(parts, fmt.Sprintf("%s -> %s", src, dst))
	}
	return fmt.Sprintf("aliases: %s", strings.Join(parts, ", "))
}
