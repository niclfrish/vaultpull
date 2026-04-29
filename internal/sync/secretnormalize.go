package sync

import (
	"fmt"
	"strings"
	"unicode"
)

// DefaultNormalizeConfig returns a NormalizeConfig with sensible defaults.
func DefaultNormalizeConfig() NormalizeConfig {
	return NormalizeConfig{
		UppercaseKeys:   true,
		TrimValues:      true,
		ReplaceHyphens:  true,
		ReplaceDots:     true,
		ReplacementChar: "_",
	}
}

// NormalizeConfig controls how secrets are normalized.
type NormalizeConfig struct {
	UppercaseKeys   bool
	TrimValues      bool
	ReplaceHyphens  bool
	ReplaceDots     bool
	ReplacementChar string
}

// NormalizeSummary holds statistics from a normalization run.
type NormalizeSummary struct {
	Total     int
	Modified  int
	Skipped   int
}

// NormalizeSecrets applies normalization rules to secret keys and values.
func NormalizeSecrets(secrets map[string]string, cfg NormalizeConfig) (map[string]string, NormalizeSummary, error) {
	if secrets == nil {
		return nil, NormalizeSummary{}, fmt.Errorf("secrets map is nil")
	}
	if cfg.ReplacementChar == "" {
		cfg.ReplacementChar = "_"
	}

	out := make(map[string]string, len(secrets))
	summary := NormalizeSummary{Total: len(secrets)}

	for k, v := range secrets {
		newKey := k
		newVal := v
		changed := false

		if cfg.ReplaceHyphens {
			replaced := strings.ReplaceAll(newKey, "-", cfg.ReplacementChar)
			if replaced != newKey {
				newKey = replaced
				changed = true
			}
		}
		if cfg.ReplaceDots {
			replaced := strings.ReplaceAll(newKey, ".", cfg.ReplacementChar)
			if replaced != newKey {
				newKey = replaced
				changed = true
			}
		}
		if cfg.UppercaseKeys {
			upper := strings.ToUpper(newKey)
			if upper != newKey {
				newKey = upper
				changed = true
			}
		}
		if cfg.TrimValues {
			trimmed := strings.TrimFunc(newVal, unicode.IsSpace)
			if trimmed != newVal {
				newVal = trimmed
				changed = true
			}
		}

		if changed {
			summary.Modified++
		} else {
			summary.Skipped++
		}
		out[newKey] = newVal
	}

	return out, summary, nil
}
