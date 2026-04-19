package sync

// DiffResult holds the changes between existing and new secrets.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]string
	Unchanged map[string]string
}

// Diff compares existing env values against incoming secrets and
// returns a categorised DiffResult.
func Diff(existing, incoming map[string]string) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, newVal := range incoming {
		oldVal, exists := existing[k]
		if !exists {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = newVal
		} else {
			result.Unchanged[k] = newVal
		}
	}

	for k, oldVal := range existing {
		if _, exists := incoming[k]; !exists {
			result.Removed[k] = oldVal
		}
	}

	return result
}

// HasChanges returns true if there are any added, removed, or changed keys.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
