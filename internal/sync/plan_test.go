package sync

import (
	"testing"
)

func TestBuildPlan_AllTypes(t *testing.T) {
	d := &DiffResult{
		Entries: []DiffEntry{
			{Key: "A", Status: StatusAdded, NewValue: "1"},
			{Key: "B", Status: StatusRemoved, OldValue: "2"},
			{Key: "C", Status: StatusChanged, OldValue: "3", NewValue: "4"},
			{Key: "D", Status: StatusUnchanged, OldValue: "5", NewValue: "5"},
		},
	}
	plan := BuildPlan(d)
	if len(plan.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(plan.Entries))
	}
	if plan.Entries[0].Change != ChangeAdd {
		t.Errorf("expected ChangeAdd for A")
	}
	if plan.Entries[1].Change != ChangeRemove {
		t.Errorf("expected ChangeRemove for B")
	}
	if plan.Entries[2].Change != ChangeUpdate {
		t.Errorf("expected ChangeUpdate for C")
	}
	if plan.Entries[3].Change != ChangeNone {
		t.Errorf("expected ChangeNone for D")
	}
}

func TestPlan_HasChanges(t *testing.T) {
	p := &Plan{Entries: []PlanEntry{{Key: "X", Change: ChangeNone}}}
	if p.HasChanges() {
		t.Error("expected no changes")
	}
	p.Entries = append(p.Entries, PlanEntry{Key: "Y", Change: ChangeAdd})
	if !p.HasChanges() {
		t.Error("expected changes")
	}
}

func TestPlan_Summary(t *testing.T) {
	d := &DiffResult{
		Entries: []DiffEntry{
			{Key: "A", Status: StatusAdded},
			{Key: "B", Status: StatusAdded},
			{Key: "C", Status: StatusRemoved},
			{Key: "D", Status: StatusChanged},
		},
	}
	plan := BuildPlan(d)
	got := plan.Summary()
	expected := "Plan: +2 added, -1 removed, ~1 updated"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
