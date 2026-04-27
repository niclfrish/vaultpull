package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// DefaultSanitizeConfig returns a SanitizeConfig with sensible defaults.
func DefaultSanitizeConfig() SanitizeConfig {
	return SanitizeConfig{
		StripControlChars: true,
		NormalizeWhitespace: true,
		MaxKeyLength: 128,
		MaxValueLength: 4096,
	}
}

// SanitizeConfig controls how secrets are sanitized.
type SanitizeConfig struct {
	StripControlChars   bool
	NormalizeWhitespace bool
	MaxKeyLength        int
	MaxValueLength      int
}

// SanitizeViolation records a single sanitization issue.
type SanitizeViolation struct {
	Key     string
	Message string
}

// SanitizeResult holds the sanitized map and any violations found.
type SanitizeResult struct {
	Secrets    map[string]string
	Violations []SanitizeViolation
}

var controlCharRe = regexp.MustCompile(`[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]`)

// SanitizeSecrets cleans secret keys and values according to cfg.
// Keys or values that exceed max lengths are truncated and a violation is recorded.
func SanitizeSecrets(secrets map[string]string, cfg SanitizeConfig) SanitizeResult {
	if secrets == nil {
		return SanitizeResult{Secrets: map[string]string{}}
	}

	result := SanitizeResult{
		Secrets: make(map[string]string, len(secrets)),
	}

	for k, v := range secrets {
		origKey := k

		if cfg.MaxKeyLength > 0 && len(k) > cfg.MaxKeyLength {
			result.Violations = append(result.Violations, SanitizeViolation{
				Key:     origKey,
				Message: fmt.Sprintf("key truncated from %d to %d chars", len(k), cfg.MaxKeyLength),
			})
			k = k[:cfg.MaxKeyLength]
		}

		if cfg.StripControlChars {
			v = controlCharRe.ReplaceAllString(v, "")
		}

		if cfg.NormalizeWhitespace {
			v = strings.TrimSpace(v)
		}

		if cfg.MaxValueLength > 0 && len(v) > cfg.MaxValueLength {
			result.Violations = append(result.Violations, SanitizeViolation{
				Key:     origKey,
				Message: fmt.Sprintf("value truncated from %d to %d chars", len(v), cfg.MaxValueLength),
			})
			v = v[:cfg.MaxValueLength]
		}

		result.Secrets[k] = v
	}

	return result
}

// SanitizeSummary returns a human-readable summary of violations.
func SanitizeSummary(violations []SanitizeViolation) string {
	if len(violations) == 0 {
		return "sanitize: no violations"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "sanitize: %d violation(s)\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&sb, "  [%s] %s\n", v.Key, v.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
