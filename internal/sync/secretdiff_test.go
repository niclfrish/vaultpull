package sync

import (
	"bytes"
	"testing"
)

func TestDiffSecrets_AllAdded(t *testing.T) {
	next := map[string]string{"A": "1", "B": "2"}
	entries := DiffSecrets(nil, next)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Op != "added" {
			t.Errorf("expected op=added, got %s", e.Op)
		}
	}
}

func TestDiffSecrets_AllRemoved(t *testing.T) {
	prev := map[string]string{"X": "val"}
	entries := DiffSecrets(prev, nil)
	if len(entries) != 1 || entries[0].Op != "removed" {
		t.Fatalf("expected 1 removed entry, got %+v", entries)
	}
}

func TestDiffSecrets_Changed(t *testing.T) {
	prev := map[string]string{"KEY": "old"}
	next := map[string]string{"KEY": "new"}
	entries := DiffSecrets(prev, next)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Op != "changed" {
		t.Errorf("expected changed, got %s", entries[0].Op)
	}
	if entries[0].OldVal != "old" || entries[0].NewVal != "new" {
		t.Errorf("unexpected values: %+v", entries[0])
	}
}

func TestDiffSecrets_Unchanged(t *testing.T) {
	m := map[string]string{"K": "v"}
	entries := DiffSecrets(m, m)
	if len(entries) != 0 {
		t.Errorf("expected no diff, got %+v", entries)
	}
}

func TestDiffSecrets_SortedByKey(t *testing.T) {
	next := map[string]string{"Z": "1", "A": "2", "M": "3"}
	entries := DiffSecrets(nil, next)
	keys := []string{entries[0].Key, entries[1].Key, entries[2].Key}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}

func TestPrintSecretDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	PrintSecretDiff(nil, DefaultSecretDiffConfig(), &buf)
	if buf.String() != "no secret changes detected\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestPrintSecretDiff_MasksValues(t *testing.T) {
	entries := []SecretDiffEntry{
		{Key: "TOKEN", OldVal: "secret", NewVal: "newsecret", Op: "changed"},
	}
	var buf bytes.Buffer
	cfg := DefaultSecretDiffConfig()
	PrintSecretDiff(entries, cfg, &buf)
	out := buf.String()
	if contains(out, "secret") {
		t.Errorf("expected values to be masked, got: %s", out)
	}
	if !contains(out, "[redacted]") {
		t.Errorf("expected redacted marker in output, got: %s", out)
	}
}

func TestPrintSecretDiff_ShowsValues(t *testing.T) {
	entries := []SecretDiffEntry{
		{Key: "DB_PASS", OldVal: "", NewVal: "hunter2", Op: "added"},
	}
	var buf bytes.Buffer
	cfg := SecretDiffConfig{MaskValues: false}
	PrintSecretDiff(entries, cfg, &buf)
	if !contains(buf.String(), "hunter2") {
		t.Errorf("expected plain value in output, got: %s", buf.String())
	}
}

func TestSecretDiffSummary(t *testing.T) {
	entries := []SecretDiffEntry{
		{Op: "added"}, {Op: "added"},
		{Op: "removed"},
		{Op: "changed"}, {Op: "changed"}, {Op: "changed"},
	}
	a, r, c := SecretDiffSummary(entries)
	if a != 2 || r != 1 || c != 3 {
		t.Errorf("unexpected summary: added=%d removed=%d changed=%d", a, r, c)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
