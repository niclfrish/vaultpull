package sync

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// DefaultSecretDiffConfig returns a DiffConfig with sensible defaults.
func DefaultSecretDiffConfig() SecretDiffConfig {
	return SecretDiffConfig{
		MaskValues: true,
		MaskReplacement: "[redacted]",
	}
}

// SecretDiffConfig controls how secret diffs are rendered.
type SecretDiffConfig struct {
	MaskValues      bool
	MaskReplacement string
}

// SecretDiffEntry represents a single changed secret between two sets.
type SecretDiffEntry struct {
	Key    string
	OldVal string
	NewVal string
	Op     string // "added", "removed", "changed"
}

// DiffSecrets compares two secret maps and returns ordered diff entries.
func DiffSecrets(prev, next map[string]string) []SecretDiffEntry {
	if prev == nil {
		prev = map[string]string{}
	}
	if next == nil {
		next = map[string]string{}
	}

	seen := map[string]bool{}
	var entries []SecretDiffEntry

	for k, newVal := range next {
		seen[k] = true
		if oldVal, ok := prev[k]; !ok {
			entries = append(entries, SecretDiffEntry{Key: k, OldVal: "", NewVal: newVal, Op: "added"})
		} else if oldVal != newVal {
			entries = append(entries, SecretDiffEntry{Key: k, OldVal: oldVal, NewVal: newVal, Op: "changed"})
		}
	}

	for k, oldVal := range prev {
		if !seen[k] {
			entries = append(entries, SecretDiffEntry{Key: k, OldVal: oldVal, NewVal: "", Op: "removed"})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// PrintSecretDiff writes a human-readable diff to w using cfg.
func PrintSecretDiff(entries []SecretDiffEntry, cfg SecretDiffConfig, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if len(entries) == 0 {
		fmt.Fprintln(w, "no secret changes detected")
		return
	}
	for _, e := range entries {
		old, nw := e.OldVal, e.NewVal
		if cfg.MaskValues {
			if old != "" {
				old = cfg.MaskReplacement
			}
			if nw != "" {
				nw = cfg.MaskReplacement
			}
		}
		switch e.Op {
		case "added":
			fmt.Fprintf(w, "+ %s = %s\n", e.Key, nw)
		case "removed":
			fmt.Fprintf(w, "- %s = %s\n", e.Key, old)
		case "changed":
			fmt.Fprintf(w, "~ %s: %s -> %s\n", e.Key, old, nw)
		}
	}
}

// SecretDiffSummary returns counts of added, removed, and changed entries.
func SecretDiffSummary(entries []SecretDiffEntry) (added, removed, changed int) {
	for _, e := range entries {
		switch e.Op {
		case "added":
			added++
		case "removed":
			removed++
		case "changed":
			changed++
		}
	}
	return
}
