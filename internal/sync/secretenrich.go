package sync

import (
	"fmt"
	"sort"
	"strings"
)

// DefaultEnrichConfig returns a default EnrichConfig.
func DefaultEnrichConfig() EnrichConfig {
	return EnrichConfig{
		Prefix:    "",
		Suffix:    "",
		StaticKeys: map[string]string{},
	}
}

// EnrichConfig controls how secrets are enriched with additional metadata.
type EnrichConfig struct {
	// Prefix prepended to every value (e.g. "vault://").
	Prefix string
	// Suffix appended to every value.
	Suffix string
	// StaticKeys are injected verbatim into the output map.
	StaticKeys map[string]string
	// OnlyKeys restricts enrichment to specific keys. Empty means all.
	OnlyKeys []string
}

// EnrichSecrets applies prefix/suffix decoration and injects static keys.
// It returns a new map; the original is not modified.
func EnrichSecrets(secrets map[string]string, cfg EnrichConfig) (map[string]string, error) {
	if secrets == nil {
		return nil, fmt.Errorf("enrichsecrets: secrets map is nil")
	}

	onlySet := make(map[string]struct{}, len(cfg.OnlyKeys))
	for _, k := range cfg.OnlyKeys {
		onlySet[k] = struct{}{}
	}

	out := make(map[string]string, len(secrets)+len(cfg.StaticKeys))

	for k, v := range secrets {
		if len(onlySet) > 0 {
			if _, ok := onlySet[k]; !ok {
				out[k] = v
				continue
			}
		}
		out[k] = cfg.Prefix + v + cfg.Suffix
	}

	for k, v := range cfg.StaticKeys {
		if _, exists := out[k]; exists {
			return nil, fmt.Errorf("enrichsecrets: static key %q conflicts with existing secret", k)
		}
		out[k] = v
	}

	return out, nil
}

// EnrichSummary returns a human-readable summary of enrichment results.
func EnrichSummary(original, enriched map[string]string, cfg EnrichConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("enriched %d secret(s)", len(original)))
	if cfg.Prefix != "" || cfg.Suffix != "" {
		sb.WriteString(fmt.Sprintf(", decorated with prefix=%q suffix=%q", cfg.Prefix, cfg.Suffix))
	}
	if len(cfg.StaticKeys) > 0 {
		keys := make([]string, 0, len(cfg.StaticKeys))
		for k := range cfg.StaticKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf(", injected static keys: [%s]", strings.Join(keys, ", ")))
	}
	return sb.String()
}
