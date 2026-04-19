package sync

import "fmt"

// ChangeType represents the type of change for a secret key.
type ChangeType string

const (
	ChangeAdd    ChangeType = "add"
	ChangeRemove ChangeType = "remove"
	ChangeUpdate ChangeType = "update"
	ChangeNone   ChangeType = "none"
)

// PlanEntry describes a single planned change.
type PlanEntry struct {
	Key    string
	Change ChangeType
	Old    string
	New    string
}

// Plan holds all planned changes before applying them.
type Plan struct {
	Entries []PlanEntry
}

// HasChanges returns true if there are any non-trivial changes.
func (p *Plan) HasChanges() bool {
	for _, e := range p.Entries {
		if e.Change != ChangeNone {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of the plan.
func (p *Plan) Summary() string {
	add, remove, update := 0, 0, 0
	for _, e := range p.Entries {
		switch e.Change {
		case ChangeAdd:
			add++
		case ChangeRemove:
			remove++
		case ChangeUpdate:
			update++
		}
	}
	return fmt.Sprintf("Plan: +%d added, -%d removed, ~%d updated", add, remove, update)
}

// BuildPlan constructs a Plan from a Diff result.
func BuildPlan(d *DiffResult) *Plan {
	plan := &Plan{}
	for _, e := range d.Entries {
		ct := ChangeNone
		switch e.Status {
		case StatusAdded:
			ct = ChangeAdd
		case StatusRemoved:
			ct = ChangeRemove
		case StatusChanged:
			ct = ChangeUpdate
		}
		plan.Entries = append(plan.Entries, PlanEntry{
			Key:    e.Key,
			Change: ct,
			Old:    e.OldValue,
			New:    e.NewValue,
		})
	}
	return plan
}
