package sync

import (
	"regexp"
	"strings"
)

// SecretMaskConfig holds configuration for secret masking rules.
type SecretMaskConfig struct {
	// KeyPatterns are regex patterns; matching keys will have values masked.
	KeyPatterns []string
	// MaskChar is the character used to replace masked values.
	MaskChar string
	// VisibleChars is the number of trailing characters to leave visible (0 = full mask).
	VisibleChars int
}

// DefaultSecretMaskConfig returns a sensible default config.
func DefaultSecretMaskConfig() SecretMaskConfig {
	return SecretMaskConfig{
		KeyPatterns:  []string{`(?i)(password|secret|token|key|credential|passwd|apikey|api_key)`},
		MaskChar:     "*",
		VisibleChars: 4,
	}
}

// SecretMaskResult holds the masked secrets and a summary.
type SecretMaskResult struct {
	Secrets map[string]string
	MaskedKeys []string
}

// MaskSecrets returns a copy of secrets with sensitive values masked.
// Keys matching any pattern have their values replaced.
func MaskSecrets(secrets map[string]string, cfg SecretMaskConfig) (SecretMaskResult, error) {
	if secrets == nil {
		return SecretMaskResult{Secrets: map[string]string{}}, nil
	}

	var compiled []*regexp.Regexp
	for _, p := range cfg.KeyPatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return SecretMaskResult{}, err
		}
		compiled = append(compiled, re)
	}

	out := make(map[string]string, len(secrets))
	var masked []string

	for k, v := range secrets {
		if matchesAny(k, compiled) {
			out[k] = applyMask(v, cfg.MaskChar, cfg.VisibleChars)
			masked = append(masked, k)
		} else {
			out[k] = v
		}
	}

	return SecretMaskResult{Secrets: out, MaskedKeys: masked}, nil
}

func matchesAny(key string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

func applyMask(value, maskChar string, visibleChars int) string {
	if value == "" {
		return value
	}
	if visibleChars <= 0 || visibleChars >= len(value) {
		return strings.Repeat(maskChar, len(value))
	}
	tail := value[len(value)-visibleChars:]
	return strings.Repeat(maskChar, len(value)-visibleChars) + tail
}
