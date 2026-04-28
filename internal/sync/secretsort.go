package sync

import (
	"fmt"
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for secrets.
type SortOrder string

const (
	SortOrderAlpha      SortOrder = "alpha"
	SortOrderAlphaDesc  SortOrder = "alpha-desc"
	SortOrderKeyLength  SortOrder = "key-length"
	SortOrderValueLength SortOrder = "value-length"
)

// DefaultSortConfig returns a sensible default sort configuration.
func DefaultSortConfig() SortConfig {
	return SortConfig{
		Order: SortOrderAlpha,
	}
}

// SortConfig controls how secrets are sorted before output.
type SortConfig struct {
	Order  SortOrder
	Prefix string // optional: sort keys with this prefix first
}

// SortSecrets returns a new map with keys sorted into a slice of key-value pairs.
// Since maps are unordered, this returns an ordered slice of [2]string pairs.
func SortSecrets(secrets map[string]string, cfg SortConfig) ([]([2]string), error) {
	if secrets == nil {
		return nil, fmt.Errorf("secretsort: secrets map is nil")
	}
	if cfg.Order == "" {
		cfg.Order = SortOrderAlpha
	}

	pairs := make([][2]string, 0, len(secrets))
	for k, v := range secrets {
		pairs = append(pairs, [2]string{k, v})
	}

	switch cfg.Order {
	case SortOrderAlpha:
		sort.Slice(pairs, func(i, j int) bool {
			return prefixFirst(pairs[i][0], pairs[j][0], cfg.Prefix)
		})
	case SortOrderAlphaDesc:
		sort.Slice(pairs, func(i, j int) bool {
			return prefixFirst(pairs[j][0], pairs[i][0], cfg.Prefix)
		})
	case SortOrderKeyLength:
		sort.Slice(pairs, func(i, j int) bool {
			return len(pairs[i][0]) < len(pairs[j][0])
		})
	case SortOrderValueLength:
		sort.Slice(pairs, func(i, j int) bool {
			return len(pairs[i][1]) < len(pairs[j][1])
		})
	default:
		return nil, fmt.Errorf("secretsort: unknown order %q", cfg.Order)
	}

	return pairs, nil
}

// SortSummary returns a human-readable description of the sort result.
func SortSummary(pairs [][2]string, cfg SortConfig) string {
	return fmt.Sprintf("sorted %d secrets by %s", len(pairs), cfg.Order)
}

func prefixFirst(a, b, prefix string) bool {
	if prefix == "" {
		return a < b
	}
	aHas := strings.HasPrefix(a, prefix)
	bHas := strings.HasPrefix(b, prefix)
	if aHas != bHas {
		return aHas
	}
	return a < b
}
