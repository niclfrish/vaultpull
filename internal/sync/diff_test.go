package sync

import (
	"testing"
)

func TestDiff_AllAdded(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(existing, incoming)

	if len(result.Added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(result.Added))
	}
	if result.HasChanges() == false {
		t.Fatal("expected HasChanges to be true")
	}
}

func TestDiff_AllRemoved(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{}

	result := Diff(existing, incoming)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(result.Removed))
	}
}

func TestDiff_Changed(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	result := Diff(existing, incoming)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(result.Changed))
	}
	if result.Changed["FOO"] != "new" {
		t.Errorf("expected changed value 'new', got %s", result.Changed["FOO"])
	}
}

func TestDiff_Unchanged(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(existing, incoming)

	if len(result.Unchanged) != 1 {
		t.Fatalf("expected 1 unchanged, got %d", len(result.Unchanged))
	}
	if result.HasChanges() {
		t.Fatal("expected HasChanges to be false")
	}
}

func TestDiff_Mixed(t *testing.T) {
	existing := map[string]string{"KEEP": "same", "CHANGE": "old", "REMOVE": "gone"}
	incoming := map[string]string{"KEEP": "same", "CHANGE": "new", "ADD": "fresh"}

	result := Diff(existing, incoming)

	if len(result.Added) != 1 || result.Added["ADD"] != "fresh" {
		t.Errorf("unexpected Added: %v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed["REMOVE"] != "gone" {
		t.Errorf("unexpected Removed: %v", result.Removed)
	}
	if len(result.Changed) != 1 || result.Changed["CHANGE"] != "new" {
		t.Errorf("unexpected Changed: %v", result.Changed)
	}
	if len(result.Unchanged) != 1 {
		t.Errorf("unexpected Unchanged: %v", result.Unchanged)
	}
}
