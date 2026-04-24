package sync

import "strings"

// TagFilter filters secrets based on key tag annotations.
// Tags are expressed as key suffixes in the form "KEY#tag1,tag2".
// TagFilter extracts the base key and keeps only entries whose tags
// match at least one of the required tags (OR semantics).
type TagFilter struct {
	requiredTags []string
}

// NewTagFilter creates a TagFilter that retains secrets whose key
// carries at least one of the given tags. If no tags are provided,
// all secrets pass through unchanged.
func NewTagFilter(tags []string) *TagFilter {
	return &TagFilter{requiredTags: tags}
}

// Apply returns a new map containing only the secrets whose keys
// match the required tags. Tag annotations are stripped from the
// returned keys so downstream consumers receive clean key names.
func (tf *TagFilter) Apply(secrets map[string]string) map[string]string {
	if len(tf.requiredTags) == 0 {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for rawKey, val := range secrets {
		baseKey, keyTags := splitKeyTags(rawKey)
		if tf.matchesAny(keyTags) {
			out[baseKey] = val
		}
	}
	return out
}

// matchesAny returns true if at least one of the key's tags is in
// the filter's required tag list.
func (tf *TagFilter) matchesAny(keyTags []string) bool {
	for _, kt := range keyTags {
		for _, rt := range tf.requiredTags {
			if strings.EqualFold(kt, rt) {
				return true
			}
		}
	}
	return false
}

// splitKeyTags splits a raw key of the form "KEY#tag1,tag2" into the
// base key and a slice of tag strings. If no "#" separator is present
// the key is returned as-is with an empty tag slice.
func splitKeyTags(raw string) (string, []string) {
	parts := strings.SplitN(raw, "#", 2)
	if len(parts) == 1 {
		return raw, nil
	}
	tags := strings.Split(parts[1], ",")
	clean := make([]string, 0, len(tags))
	for _, t := range tags {
		if trimmed := strings.TrimSpace(t); trimmed != "" {
			clean = append(clean, trimmed)
		}
	}
	return parts[0], clean
}
