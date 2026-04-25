package sync

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// MetadataConfig holds configuration for metadata injection.
type MetadataConfig struct {
	// KeyPrefix is prepended to every injected metadata key.
	KeyPrefix string
	// IncludeTimestamp injects a sync timestamp key.
	IncludeTimestamp bool
	// IncludeCount injects the total number of secrets synced.
	IncludeCount bool
	// IncludeKeys injects a comma-separated list of synced key names.
	IncludeKeys bool
	// TimestampFormat is the time layout used when IncludeTimestamp is true.
	// Defaults to time.RFC3339 if empty.
	TimestampFormat string
}

// DefaultMetadataConfig returns sensible defaults.
func DefaultMetadataConfig() MetadataConfig {
	return MetadataConfig{
		KeyPrefix:        "VAULTPULL_META_",
		IncludeTimestamp: true,
		IncludeCount:     true,
		IncludeKeys:      false,
		TimestampFormat:  time.RFC3339,
	}
}

// InjectMetadata adds metadata entries to secrets according to cfg.
// It never overwrites existing keys.
func InjectMetadata(secrets map[string]string, cfg MetadataConfig, now time.Time) (map[string]string, error) {
	if secrets == nil {
		return nil, fmt.Errorf("metadata: secrets map must not be nil")
	}

	prefix := cfg.KeyPrefix
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	if cfg.IncludeTimestamp {
		fmt := cfg.TimestampFormat
		if fmt == "" {
			fmt = time.RFC3339
		}
		k := prefix + "SYNCED_AT"
		if _, exists := result[k]; !exists {
			result[k] = now.Format(fmt)
		}
	}

	if cfg.IncludeCount {
		k := prefix + "SECRET_COUNT"
		if _, exists := result[k]; !exists {
			result[k] = fmt.Sprintf("%d", len(secrets))
		}
	}

	if cfg.IncludeKeys {
		keys := make([]string, 0, len(secrets))
		for k := range secrets {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		k := prefix + "KEYS"
		if _, exists := result[k]; !exists {
			result[k] = strings.Join(keys, ",")
		}
	}

	return result, nil
}
