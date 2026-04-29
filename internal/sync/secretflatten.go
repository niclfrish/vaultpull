package sync

import (
	"fmt"
	"strings"
)

// DefaultFlattenConfig returns a FlattenConfig with sensible defaults.
func DefaultFlattenConfig() FlattenConfig {
	return FlattenConfig{
		Separator: "_",
		MaxDepth:  10,
		UpperCase: true,
	}
}

// FlattenConfig controls how nested key paths are flattened.
type FlattenConfig struct {
	// Separator is placed between path segments (default "_").
	Separator string
	// MaxDepth limits how many segments are joined (0 = unlimited).
	MaxDepth int
	// UpperCase converts the resulting key to upper-case.
	UpperCase bool
}

// FlattenSecrets takes a map whose keys may contain dot-separated segments
// (e.g. "db.host") and returns a new map with those segments joined by the
// configured separator (e.g. "DB_HOST").
func FlattenSecrets(secrets map[string]string, cfg FlattenConfig) (map[string]string, error) {
	if secrets == nil {
		return nil, fmt.Errorf("secretflatten: secrets map is nil")
	}
	if cfg.Separator == "" {
		return nil, fmt.Errorf("secretflatten: separator must not be empty")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		flat := flattenKey(k, cfg)
		if _, exists := out[flat]; exists {
			return nil, fmt.Errorf("secretflatten: key collision after flattening: %q", flat)
		}
		out[flat] = v
	}
	return out, nil
}

// FlattenSummary returns a human-readable description of the operation.
func FlattenSummary(before, after map[string]string) string {
	return fmt.Sprintf("flattened %d keys into %d keys", len(before), len(after))
}

func flattenKey(key string, cfg FlattenConfig) string {
	segments := strings.Split(key, ".")
	if cfg.MaxDepth > 0 && len(segments) > cfg.MaxDepth {
		segments = segments[:cfg.MaxDepth]
	}
	result := strings.Join(segments, cfg.Separator)
	if cfg.UpperCase {
		result = strings.ToUpper(result)
	}
	return result
}
