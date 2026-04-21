package sync

// SnapshotDiff compares two snapshots and returns a DiffResult.
// If previous is nil, all current keys are treated as Added.
func SnapshotDiff(previous *Snapshot, current map[string]string) DiffResult {
	var prev map[string]string
	if previous != nil {
		prev = previous.Secrets
	} else {
		prev = map[string]string{}
	}
	return Diff(prev, current)
}

// SnapshotSummary returns a human-readable summary of changes between snapshots.
func SnapshotSummary(previous *Snapshot, current map[string]string) string {
	result := SnapshotDiff(previous, current)
	plan := BuildPlan(result)
	return plan.Summary()
}
