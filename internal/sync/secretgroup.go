package sync

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultGroupConfig returns a GroupConfig with sensible defaults.
func DefaultGroupConfig() GroupConfig {
	return GroupConfig{
		Separator: "_",
		MaxDepth:  2,
	}
}

// GroupConfig controls how secrets are grouped by key prefix.
type GroupConfig struct {
	// Separator is the delimiter used to split key segments (default: "_").
	Separator string
	// MaxDepth is the maximum number of prefix segments used for grouping.
	MaxDepth int
}

// SecretGroup holds secrets that share a common prefix.
type SecretGroup struct {
	Prefix  string
	Secrets map[string]string
}

// GroupSecrets partitions secrets into groups based on key prefix segments.
// Keys that do not contain the separator are placed in a group with an empty prefix.
func GroupSecrets(secrets map[string]string, cfg GroupConfig) ([]SecretGroup, error) {
	if secrets == nil {
		return nil, fmt.Errorf("secretgroup: secrets map is nil")
	}
	if cfg.Separator == "" {
		return nil, fmt.Errorf("secretgroup: separator must not be empty")
	}
	if cfg.MaxDepth < 1 {
		return nil, fmt.Errorf("secretgroup: MaxDepth must be >= 1")
	}

	groupMap := make(map[string]map[string]string)
	for k, v := range secrets {
		prefix := extractPrefix(k, cfg.Separator, cfg.MaxDepth)
		if groupMap[prefix] == nil {
			groupMap[prefix] = make(map[string]string)
		}
		groupMap[prefix][k] = v
	}

	prefixes := make([]string, 0, len(groupMap))
	for p := range groupMap {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	groups := make([]SecretGroup, 0, len(prefixes))
	for _, p := range prefixes {
		groups = append(groups, SecretGroup{Prefix: p, Secrets: groupMap[p]})
	}
	return groups, nil
}

// GroupSummary returns a human-readable summary of the grouped secrets.
func GroupSummary(groups []SecretGroup) string {
	if len(groups) == 0 {
		return "no groups"
	}
	var sb strings.Builder
	for _, g := range groups {
		label := g.Prefix
		if label == "" {
			label = "(ungrouped)"
		}
		fmt.Fprintf(&sb, "group=%s keys=%d\n", label, len(g.Secrets))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func extractPrefix(key, sep string, maxDepth int) string {
	parts := strings.Split(key, sep)
	if len(parts) <= 1 {
		return ""
	}
	depth := maxDepth
	if depth > len(parts)-1 {
		depth = len(parts) - 1
	}
	return strings.Join(parts[:depth], sep)
}
