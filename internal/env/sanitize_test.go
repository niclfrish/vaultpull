package env

import "testing"

func TestSanitizeKey(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"foo", "FOO"},
		{"foo-bar", "FOO_BAR"},
		{"FOO_BAR", "FOO_BAR"},
		{"mixed-Case-Key", "MIXED_CASE_KEY"},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := sanitizeKey(c.input)
			if got != c.expected {
				t.Errorf("sanitizeKey(%q) = %q, want %q", c.input, got, c.expected)
			}
		})
	}
}

func TestEscapeValue(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with space", `"with space"`},
		{"with#hash", `"with#hash"`},
		{"with\ttab", `"with\ttab"`},
		{`quo"te`, `"quo\"te"`},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := escapeValue(c.input)
			if got != c.expected {
				t.Errorf("escapeValue(%q) = %q, want %q", c.input, got, c.expected)
			}
		})
	}
}
