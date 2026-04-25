package sync

import (
	"regexp"
	"strings"
)

// MaskConfig controls how secret values are masked.
type MaskConfig struct {
	// Patterns is a list of key patterns (regex) whose values should be fully masked.
	Patterns []string
	// PartialPatterns is a list of key patterns whose values are partially revealed.
	PartialPatterns []string
	// VisibleChars is the number of trailing characters to reveal for partial masking.
	VisibleChars int
	// MaskChar is the character used for masking (default '*').
	MaskChar rune
}

// DefaultMaskConfig returns a sensible default masking configuration.
func DefaultMaskConfig() MaskConfig {
	return MaskConfig{
		Patterns:        []string{`(?i)(password|secret|token|key|passphrase)`},
		PartialPatterns: []string{`(?i)(api_?key|access_?key)`},
		VisibleChars:    4,
		MaskChar:        '*',
	}
}

// Masker applies masking rules to a map of secrets.
type Masker struct {
	cfg      MaskConfig
	full     []*regexp.Regexp
	partial  []*regexp.Regexp
}

// NewMasker compiles the patterns in cfg and returns a Masker.
func NewMasker(cfg MaskConfig) (*Masker, error) {
	if cfg.MaskChar == 0 {
		cfg.MaskChar = '*'
	}
	m := &Masker{cfg: cfg}
	for _, p := range cfg.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		m.full = append(m.full, re)
	}
	for _, p := range cfg.PartialPatterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		m.partial = append(m.partial, re)
	}
	return m, nil
}

// Apply returns a new map with values masked according to the configured rules.
func (m *Masker) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.maskValue(k, v)
	}
	return out
}

func (m *Masker) maskValue(key, value string) string {
	for _, re := range m.full {
		if re.MatchString(key) {
			return strings.Repeat(string(m.cfg.MaskChar), len(value))
		}
	}
	for _, re := range m.partial {
		if re.MatchString(key) {
			return partialMask(value, m.cfg.VisibleChars, m.cfg.MaskChar)
		}
	}
	return value
}

func partialMask(value string, visible int, ch rune) string {
	if visible <= 0 || len(value) <= visible {
		return strings.Repeat(string(ch), len(value))
	}
	maskLen := len(value) - visible
	return strings.Repeat(string(ch), maskLen) + value[maskLen:]
}
