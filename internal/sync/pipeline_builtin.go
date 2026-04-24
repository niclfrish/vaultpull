package sync

import (
	"fmt"
	"strings"
)

// FilterStage returns a Stage function that removes keys not matching the
// provided Filter criteria. It wraps the existing NewFilter helper so pipeline
// users don't need to construct a Filter manually.
func FilterStage(f *Filter) func(map[string]string) (map[string]string, error) {
	return func(secrets map[string]string) (map[string]string, error) {
		return f.Apply(secrets), nil
	}
}

// TransformStage returns a Stage function that applies a Transformer to every
// value in the secrets map.
func TransformStage(tr *Transformer) func(map[string]string) (map[string]string, error) {
	return func(secrets map[string]string) (map[string]string, error) {
		return tr.Apply(secrets)
	}
}

// DedupeStage returns a Stage function that deduplicates keys according to the
// provided strategy ("first" or "last").
func DedupeStage(strategy string) func(map[string]string) (map[string]string, error) {
	return func(secrets map[string]string) (map[string]string, error) {
		// Convert map to ordered pairs for Dedupe (preserve insertion order via
		// sorted keys so behaviour is deterministic in tests).
		pairs := make([]string, 0, len(secrets)*2)
		for k, v := range secrets {
			pairs = append(pairs, k, v)
		}
		return Dedupe(pairs, strategy)
	}
}

// RequiredKeysStage returns a Stage function that errors if any of the
// required keys are absent from the secrets map.
func RequiredKeysStage(keys ...string) func(map[string]string) (map[string]string, error) {
	return func(secrets map[string]string) (map[string]string, error) {
		var missing []string
		for _, k := range keys {
			if _, ok := secrets[k]; !ok {
				missing = append(missing, k)
			}
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("required keys missing: %s", strings.Join(missing, ", "))
		}
		return secrets, nil
	}
}

// TruncateStage returns a Stage function that truncates long values using the
// default truncation configuration.
func TruncateStage() func(map[string]string) (map[string]string, error) {
	return func(secrets map[string]string) (map[string]string, error) {
		return TruncateSecrets(secrets, DefaultTruncateConfig()), nil
	}
}
