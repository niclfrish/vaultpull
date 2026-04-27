package sync

import (
	"fmt"
	"sort"
	"strings"
)

// PrioritySource represents a named secret source with a priority level.
// Lower numbers indicate higher priority (1 = highest).
type PrioritySource struct {
	Name     string
	Priority int
	Secrets  map[string]string
}

// DefaultPriorityConfig returns a PriorityConfig with sensible defaults.
func DefaultPriorityConfig() PriorityConfig {
	return PriorityConfig{
		ConflictPrefix: "__conflict_",
	}
}

// PriorityConfig controls behaviour of MergeByPriority.
type PriorityConfig struct {
	// ConflictPrefix is prepended to keys that were overridden by a higher-priority source.
	// Set to "" to suppress conflict annotations.
	ConflictPrefix string
}

// MergeByPriority merges multiple PrioritySources into a single map.
// Sources must be provided in any order; priority field determines precedence.
// When two sources define the same key the lower-priority value is discarded.
// If cfg.ConflictPrefix is non-empty, the discarded value is preserved under
// "<prefix><source>_<key>" for audit purposes.
func MergeByPriority(cfg PriorityConfig, sources []PrioritySource) (map[string]string, error) {
	if len(sources) == 0 {
		return map[string]string{}, nil
	}

	// Validate priorities are unique and positive.
	seen := map[int]string{}
	for _, s := range sources {
		if s.Priority < 1 {
			return nil, fmt.Errorf("source %q has invalid priority %d: must be >= 1", s.Name, s.Priority)
		}
		if prev, ok := seen[s.Priority]; ok {
			return nil, fmt.Errorf("sources %q and %q share priority %d", prev, s.Name, s.Priority)
		}
		seen[s.Priority] = s.Name
	}

	// Sort ascending so highest priority (lowest number) is applied last and wins.
	sorted := make([]PrioritySource, len(sources))
	copy(sorted, sources)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})

	result := map[string]string{}
	origin := map[string]string{} // key -> source name that currently owns it

	for _, src := range sorted {
		for k, v := range src.Secrets {
			if existing, conflict := result[k]; conflict && cfg.ConflictPrefix != "" {
				conflictKey := cfg.ConflictPrefix + strings.ToLower(origin[k]) + "_" + k
				result[conflictKey] = existing
			}
			result[k] = v
			origin[k] = src.Name
		}
	}

	return result, nil
}

// PrioritySummary returns a human-readable summary of how many keys each source contributed.
func PrioritySummary(sources []PrioritySource, merged map[string]string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("merged %d keys from %d sources\n", len(merged), len(sources)))
	for _, s := range sources {
		sb.WriteString(fmt.Sprintf("  [%d] %s: %d keys\n", s.Priority, s.Name, len(s.Secrets)))
	}
	return sb.String()
}
