package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// DefaultRedactConfig returns a sensible default configuration for redaction.
func DefaultRedactConfig() RedactConfig {
	return RedactConfig{
		Patterns: []string{
			`(?i)password`,
			`(?i)secret`,
			`(?i)token`,
			`(?i)api[_-]?key`,
			`(?i)private[_-]?key`,
		},
		Replacement: "[REDACTED]",
	}
}

// RedactConfig controls which keys are redacted and how.
type RedactConfig struct {
	Patterns    []string
	Replacement string
}

// RedactSecrets replaces the values of matching keys with the configured replacement string.
// It returns a new map and a summary of what was redacted.
func RedactSecrets(secrets map[string]string, cfg RedactConfig) (map[string]string, RedactSummary, error) {
	if secrets == nil {
		return nil, RedactSummary{}, nil
	}

	compiled := make([]*regexp.Regexp, 0, len(cfg.Patterns))
	for _, p := range cfg.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, RedactSummary{}, fmt.Errorf("invalid redact pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	replacement := cfg.Replacement
	if replacement == "" {
		replacement = "[REDACTED]"
	}

	out := make(map[string]string, len(secrets))
	var redacted []string

	for k, v := range secrets {
		if matchesAnyPattern(k, compiled) {
			out[k] = replacement
			redacted = append(redacted, k)
		} else {
			out[k] = v
		}
	}

	return out, RedactSummary{RedactedKeys: redacted, Total: len(secrets)}, nil
}

// RedactSummary holds the result of a redaction pass.
type RedactSummary struct {
	RedactedKeys []string
	Total        int
}

// String returns a human-readable summary.
func (s RedactSummary) String() string {
	return fmt.Sprintf("redacted %d/%d keys: %s", len(s.RedactedKeys), s.Total, strings.Join(s.RedactedKeys, ", "))
}

func matchesAnyPattern(key string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
