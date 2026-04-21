package sync

import "strings"

// DedupeStrategy defines how duplicate keys are resolved.
type DedupeStrategy int

const (
	// DedupeKeepFirst retains the first occurrence of a duplicate key.
	DedupeKeepFirst DedupeStrategy = iota
	// DedupeKeepLast retains the last occurrence of a duplicate key.
	DedupeKeepLast
)

// DedupeResult holds the output of a deduplication pass.
type DedupeResult struct {
	Secrets    map[string]string
	Duplicates []string
}

// Dedupe removes duplicate keys from a slice of key=value pairs according to
// the chosen strategy. It returns a clean map and the list of duplicate keys
// that were discarded.
func Dedupe(pairs []string, strategy DedupeStrategy) DedupeResult {
	seen := make(map[string]bool)
	duplicates := []string{}
	result := make(map[string]string)

	type entry struct {
		key string
		val string
	}

	var entries []entry
	for _, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			continue
		}
		entries = append(entries, entry{key: pair[:idx], val: pair[idx+1:]})
	}

	if strategy == DedupeKeepLast {
		// Reverse so that the first write wins after reversal (= keep last).
		for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}
	}

	for _, e := range entries {
		if seen[e.key] {
			duplicates = append(duplicates, e.key)
			continue
		}
		seen[e.key] = true
		result[e.key] = e.val
	}

	return DedupeResult{
		Secrets:    result,
		Duplicates: duplicates,
	}
}
