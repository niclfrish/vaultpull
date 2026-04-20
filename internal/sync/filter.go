package sync

import "strings"

// Filter holds criteria for including or excluding secret keys.
type Filter struct {
	Prefixes []string
	Excludes []string
}

// NewFilter creates a Filter from include-prefix and exclude-key slices.
func NewFilter(prefixes, excludes []string) *Filter {
	return &Filter{
		Prefixes: prefixes,
		Excludes: excludes,
	}
}

// Apply returns a new map containing only the entries that pass the filter.
// If no prefixes are specified, all keys are included unless excluded.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	result := make(map[string]string, len(secrets))

	excludeSet := make(map[string]struct{}, len(f.Excludes))
	for _, e := range f.Excludes {
		excludeSet[strings.ToUpper(e)] = struct{}{}
	}

	for k, v := range secrets {
		upper := strings.ToUpper(k)

		if _, excluded := excludeSet[upper]; excluded {
			continue
		}

		if len(f.Prefixes) == 0 {
			result[k] = v
			continue
		}

		for _, p := range f.Prefixes {
			if strings.HasPrefix(upper, strings.ToUpper(p)) {
				result[k] = v
				break
			}
		}
	}

	return result
}
