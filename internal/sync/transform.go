package sync

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a secret value.
type TransformFunc func(key, value string) (string, error)

// Transformer applies a chain of TransformFuncs to secrets.
type Transformer struct {
	fns []TransformFunc
}

// NewTransformer creates a Transformer with the given transform functions.
func NewTransformer(fns ...TransformFunc) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs all transform functions over each key/value pair in secrets.
// Returns a new map with transformed values. Stops on first error.
func (t *Transformer) Apply(secrets map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}
	for _, fn := range t.fns {
		for k, v := range result {
			transformed, err := fn(k, v)
			if err != nil {
				return nil, fmt.Errorf("transform error on key %q: %w", k, err)
			}
			result[k] = transformed
		}
	}
	return result, nil
}

// TrimSpaceTransform trims leading and trailing whitespace from values.
func TrimSpaceTransform() TransformFunc {
	return func(key, value string) (string, error) {
		return strings.TrimSpace(value), nil
	}
}

// UpperKeyTransform converts all keys to uppercase (returns value unchanged).
// Note: caller must rebuild map if key casing matters; this only signals intent.
func UpperKeyTransform() TransformFunc {
	return func(key, value string) (string, error) {
		_ = strings.ToUpper(key)
		return value, nil
	}
}

// RedactTransform replaces values for keys matching any of the given substrings with a redacted placeholder.
func RedactTransform(sensitiveSubstrings []string, placeholder string) TransformFunc {
	return func(key, value string) (string, error) {
		lower := strings.ToLower(key)
		for _, sub := range sensitiveSubstrings {
			if strings.Contains(lower, strings.ToLower(sub)) {
				return placeholder, nil
			}
		}
		return value, nil
	}
}
