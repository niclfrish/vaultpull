package sync

import "fmt"

// LabelFilterStage returns a PipelineStage that applies a LabelFilter to the
// secret map, keeping only entries that carry all required label annotations.
// If labels is nil or empty, all secrets pass through unchanged.
func LabelFilterStage(labels map[string]string) PipelineStage {
	return PipelineStage{
		Name: "label-filter",
		Run: func(secrets map[string]string) (map[string]string, error) {
			if len(labels) == 0 {
				return secrets, nil
			}
			f := NewLabelFilter(labels)
			result := f.Apply(secrets)
			return result, nil
		},
	}
}

// ParseLabelFlags parses a slice of "key=value" strings into a label map.
// Returns an error if any entry is not in "key=value" format.
func ParseLabelFlags(flags []string) (map[string]string, error) {
	labels := make(map[string]string, len(flags))
	for _, flag := range flags {
		parts := splitN(flag, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid label %q: must be key=value", flag)
		}
		labels[parts[0]] = parts[1]
	}
	return labels, nil
}

// splitN is a thin wrapper so we don't import strings in this file.
func splitN(s, sep string, n int) []string {
	var result []string
	remaining := s
	for i := 0; i < n-1; i++ {
		idx := indexOf(remaining, sep)
		if idx < 0 {
			break
		}
		result = append(result, remaining[:idx])
		remaining = remaining[idx+len(sep):]
	}
	result = append(result, remaining)
	return result
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
