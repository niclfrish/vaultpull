package sync

import (
	"fmt"
	"strings"
)

// DefaultLookupConfig returns a LookupConfig with sensible defaults.
func DefaultLookupConfig() LookupConfig {
	return LookupConfig{
		CaseSensitive: false,
		PartialMatch:  false,
	}
}

// LookupConfig controls how LookupSecrets searches for keys.
type LookupConfig struct {
	CaseSensitive bool
	PartialMatch  bool
}

// LookupResult holds a single lookup result.
type LookupResult struct {
	Key   string
	Value string
}

// LookupSecrets searches secrets for keys matching any of the provided queries.
// Returns an ordered slice of LookupResult.
func LookupSecrets(secrets map[string]string, queries []string, cfg LookupConfig) ([]LookupResult, error) {
	if secrets == nil {
		return nil, fmt.Errorf("lookup: secrets map is nil")
	}
	if len(queries) == 0 {
		return nil, fmt.Errorf("lookup: no queries provided")
	}

	seen := make(map[string]bool)
	var results []LookupResult

	for _, q := range queries {
		norm := q
		if !cfg.CaseSensitive {
			norm = strings.ToLower(q)
		}
		for k, v := range secrets {
			compare := k
			if !cfg.CaseSensitive {
				compare = strings.ToLower(k)
			}
			matched := false
			if cfg.PartialMatch {
				matched = strings.Contains(compare, norm)
			} else {
				matched = compare == norm
			}
			if matched && !seen[k] {
				seen[k] = true
				results = append(results, LookupResult{Key: k, Value: v})
			}
		}
	}
	return results, nil
}

// LookupSummary returns a human-readable summary of lookup results.
func LookupSummary(results []LookupResult) string {
	if len(results) == 0 {
		return "lookup: no results found"
	}
	return fmt.Sprintf("lookup: found %d result(s)", len(results))
}
