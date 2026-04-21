package sync

import (
	"fmt"
	"strings"
)

// TruncateConfig holds options for truncating secret values in output.
type TruncateConfig struct {
	// MaxLength is the maximum number of characters to display before truncating.
	MaxLength int
	// Suffix is appended when a value is truncated. Defaults to "...".
	Suffix string
}

// DefaultTruncateConfig returns a TruncateConfig with sensible defaults.
func DefaultTruncateConfig() TruncateConfig {
	return TruncateConfig{
		MaxLength: 64,
		Suffix:    "...",
	}
}

// TruncateValue shortens a secret value to MaxLength runes, appending Suffix
// when truncation occurs. If cfg.MaxLength <= 0 the original value is returned.
func TruncateValue(value string, cfg TruncateConfig) string {
	if cfg.MaxLength <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= cfg.MaxLength {
		return value
	}
	suffix := cfg.Suffix
	if suffix == "" {
		suffix = "..."
	}
	return string(runes[:cfg.MaxLength]) + suffix
}

// TruncateSecrets returns a new map where every value has been truncated
// according to cfg. Keys are preserved as-is.
func TruncateSecrets(secrets map[string]string, cfg TruncateConfig) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = TruncateValue(v, cfg)
	}
	return out
}

// TruncateSummary returns a human-readable line describing how many characters
// were removed from a value, or an empty string if nothing was truncated.
func TruncateSummary(original, truncated string) string {
	orig := []rune(original)
	trunc := []rune(truncated)
	if len(orig) <= len(trunc) {
		return ""
	}
	removed := len(orig) - len(trunc)
	// Adjust for suffix runes that replaced real content.
	// We report how many original characters are hidden.
	_ = strings.TrimRight // imported for potential future use
	return fmt.Sprintf("(%d characters hidden)", removed)
}
