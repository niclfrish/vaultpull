package sync

import "sort"

// Status represents the diff state of a key.
type Status string

const (
	StatusAdded     Status = "added"
	StatusRemoved   Status = "removed"
	StatusChanged   Status = "changed"
	StatusUnchanged Status = "unchanged"
)

// DiffEntry holds information about a single key difference.
type DiffEntry struct {
	Key      string
	Status   Status
	OldValue string
	NewValue string
}

// DiffResult holds all diff entries.
type DiffResult struct {
	Entries []DiffEntry
}

// Diff computes the difference between current (local) and incoming (vault) secrets.
func Diff(current, incoming map[string]string) *DiffResult {
	result := &DiffResult{}
	seen := map[string]bool{}

	for k, newVal := range incoming {
		seen[k] = true
		oldVal, exists := current[k]
		if !exists {
			result.Entries = append(result.Entries, DiffEntry{Key: k, Status: StatusAdded, NewValue: newVal})
		} else if oldVal != newVal {
			result.Entries = append(result.Entries, DiffEntry{Key: k, Status: StatusChanged, OldValue: oldVal, NewValue: newVal})
		} else {
			result.Entries = append(result.Entries, DiffEntry{Key: k, Status: StatusUnchanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range current {
		if !seen[k] {
			result.Entries = append(result.Entries, DiffEntry{Key: k, Status: StatusRemoved, OldValue: oldVal})
		}
	}

	sort.Slice(result.Entries, func(i, j int) bool {
		return result.Entries[i].Key < result.Entries[j].Key
	})
	return result
}
