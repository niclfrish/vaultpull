package sync

import (
	"fmt"
	"math/rand"
	"sort"
)

// DefaultSampleConfig returns a SampleConfig with sensible defaults.
func DefaultSampleConfig() SampleConfig {
	return SampleConfig{
		Seed:       42,
		MaxSamples: 10,
	}
}

// SampleConfig controls how secrets are sampled.
type SampleConfig struct {
	// Seed is used for deterministic sampling.
	Seed int64
	// MaxSamples is the maximum number of secrets to include.
	MaxSamples int
	// OnlyKeys restricts sampling to keys matching these prefixes.
	OnlyKeys []string
}

// SampleSecrets returns a deterministic random subset of secrets up to MaxSamples.
// If MaxSamples <= 0 or >= len(secrets), all secrets are returned unchanged.
func SampleSecrets(secrets map[string]string, cfg SampleConfig) (map[string]string, error) {
	if secrets == nil {
		return nil, fmt.Errorf("secretsample: secrets map is nil")
	}
	if cfg.MaxSamples <= 0 {
		return nil, fmt.Errorf("secretsample: MaxSamples must be greater than zero")
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		if len(cfg.OnlyKeys) == 0 || hasAnyPrefix(k, cfg.OnlyKeys) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	if cfg.MaxSamples >= len(keys) {
		out := make(map[string]string, len(keys))
		for _, k := range keys {
			out[k] = secrets[k]
		}
		return out, nil
	}

	r := rand.New(rand.NewSource(cfg.Seed)) //nolint:gosec
	r.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	sampled := keys[:cfg.MaxSamples]
	sort.Strings(sampled)

	out := make(map[string]string, len(sampled))
	for _, k := range sampled {
		out[k] = secrets[k]
	}
	return out, nil
}

// SampleSummary returns a human-readable summary of the sampling operation.
func SampleSummary(total, sampled int) string {
	return fmt.Sprintf("sampled %d of %d secrets", sampled, total)
}

func hasAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return true
		}
	}
	return false
}
