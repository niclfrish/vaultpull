package sync

// MergeStrategy defines how existing local values are handled during sync.
type MergeStrategy int

const (
	// StrategyOverwrite replaces all local values with Vault values.
	StrategyOverwrite MergeStrategy = iota
	// StrategyKeepLocal preserves local values when a key already exists.
	StrategyKeepLocal
)

// Merge combines vault secrets with existing local env vars according to the
// given strategy. The vault map is never mutated.
func Merge(local, vault map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(vault))

	for k, v := range vault {
		result[k] = v
	}

	if strategy == StrategyKeepLocal {
		for k, v := range local {
			if _, exists := result[k]; exists {
				result[k] = v
			}
		}
	}

	return result
}
