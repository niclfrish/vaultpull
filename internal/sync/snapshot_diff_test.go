package sync

import (
	"strings"
	"testing"
	"time"
)

func TestSnapshotDiff_NilPrevious_AllAdded(t *testing.T) {
	current := map[string]string{"A": "1", "B": "2"}
	result := SnapshotDiff(nil, current)
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Removed) != 0 || len(result.Changed) != 0 {
		t.Error("expected no removed or changed")
	}
}

func TestSnapshotDiff_WithPrevious_DetectsChanges(t *testing.T) {
	prev := &Snapshot{
		Timestamp: time.Now(),
		Secrets:   map[string]string{"A": "old", "B": "same", "C": "gone"},
	}
	current := map[string]string{"A": "new", "B": "same", "D": "fresh"}
	result := SnapshotDiff(prev, current)

	if _, ok := result.Changed["A"]; !ok {
		t.Error("expected A to be changed")
	}
	if _, ok := result.Unchanged["B"]; !ok {
		t.Error("expected B to be unchanged")
	}
	if _, ok := result.Removed["C"]; !ok {
		t.Error("expected C to be removed")
	}
	if _, ok := result.Added["D"]; !ok {
		t.Error("expected D to be added")
	}
}

func TestSnapshotDiff_EmptyCurrentAndPrev(t *testing.T) {
	prev := &Snapshot{Secrets: map[string]string{}}
	result := SnapshotDiff(prev, map[string]string{})
	if len(result.Added)+len(result.Removed)+len(result.Changed) != 0 {
		t.Error("expected no changes for empty inputs")
	}
}

func TestSnapshotSummary_NoChanges(t *testing.T) {
	prev := &Snapshot{Secrets: map[string]string{"A": "1"}}
	summary := SnapshotSummary(prev, map[string]string{"A": "1"})
	if !strings.Contains(summary, "no changes") {
		t.Errorf("expected 'no changes' in summary, got: %s", summary)
	}
}

func TestSnapshotSummary_WithChanges(t *testing.T) {
	prev := &Snapshot{Secrets: map[string]string{"A": "old"}}
	summary := SnapshotSummary(prev, map[string]string{"A": "new", "B": "added"})
	if !strings.Contains(summary, "1 added") {
		t.Errorf("expected '1 added' in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "1 changed") {
		t.Errorf("expected '1 changed' in summary, got: %s", summary)
	}
}
