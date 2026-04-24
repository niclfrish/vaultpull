package sync

import (
	"strings"
)

// LabelFilter filters secrets based on key=value label annotations embedded
// in the secret key using a "__label__KEY=VALUE" suffix convention.
type LabelFilter struct {
	required map[string]string
}

// NewLabelFilter creates a LabelFilter that keeps only secrets whose keys
// carry all of the specified label key=value pairs.
func NewLabelFilter(labels map[string]string) *LabelFilter {
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	return &LabelFilter{required: copy}
}

// Apply returns a new map containing only the entries that satisfy all
// required labels. Labels are stripped from the output keys.
func (f *LabelFilter) Apply(secrets map[string]string) map[string]string {
	result := make(map[string]string)
	for rawKey, value := range secrets {
		baseKey, labels := splitKeyLabels(rawKey)
		if f.matchesAll(labels) {
			result[baseKey] = value
		}
	}
	return result
}

func (f *LabelFilter) matchesAll(labels map[string]string) bool {
	if len(f.required) == 0 {
		return true
	}
	for wantK, wantV := range f.required {
		if got, ok := labels[wantK]; !ok || got != wantV {
			return false
		}
	}
	return true
}

// splitKeyLabels parses a key of the form "BASE__label__K1=V1__label__K2=V2"
// into the base key and a map of label key→value pairs.
func splitKeyLabels(raw string) (string, map[string]string) {
	const marker = "__label__"
	parts := strings.Split(raw, marker)
	base := parts[0]
	labels := make(map[string]string)
	for _, part := range parts[1:] {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			labels[kv[0]] = kv[1]
		}
	}
	return base, labels
}
