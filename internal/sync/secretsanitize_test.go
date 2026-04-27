package sync

import (
	"strings"
	"testing"
)

func TestSanitizeSecrets_NilSecrets(t *testing.T) {
	result := SanitizeSecrets(nil, DefaultSanitizeConfig())
	if result.Secrets == nil {
		t.Fatal("expected non-nil Secrets map")
	}
	if len(result.Secrets) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result.Secrets))
	}
}

func TestSanitizeSecrets_NoViolations(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	result := SanitizeSecrets(secrets, DefaultSanitizeConfig())
	if len(result.Violations) != 0 {
		t.Errorf("expected no violations, got %d", len(result.Violations))
	}
	if result.Secrets["KEY"] != "value" {
		t.Errorf("unexpected value: %s", result.Secrets["KEY"])
	}
}

func TestSanitizeSecrets_StripControlChars(t *testing.T) {
	secrets := map[string]string{"KEY": "val\x01ue\x07"}
	cfg := DefaultSanitizeConfig()
	cfg.StripControlChars = true
	result := SanitizeSecrets(secrets, cfg)
	if result.Secrets["KEY"] != "value" {
		t.Errorf("expected control chars stripped, got %q", result.Secrets["KEY"])
	}
}

func TestSanitizeSecrets_NormalizeWhitespace(t *testing.T) {
	secrets := map[string]string{"KEY": "  hello world  "}
	cfg := DefaultSanitizeConfig()
	cfg.NormalizeWhitespace = true
	result := SanitizeSecrets(secrets, cfg)
	if result.Secrets["KEY"] != "hello world" {
		t.Errorf("expected trimmed value, got %q", result.Secrets["KEY"])
	}
}

func TestSanitizeSecrets_TruncatesLongKey(t *testing.T) {
	longKey := strings.Repeat("K", 200)
	secrets := map[string]string{longKey: "v"}
	cfg := DefaultSanitizeConfig()
	cfg.MaxKeyLength = 128
	result := SanitizeSecrets(secrets, cfg)
	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}
	for k := range result.Secrets {
		if len(k) != 128 {
			t.Errorf("expected truncated key length 128, got %d", len(k))
		}
	}
}

func TestSanitizeSecrets_TruncatesLongValue(t *testing.T) {
	secrets := map[string]string{"KEY": strings.Repeat("x", 5000)}
	cfg := DefaultSanitizeConfig()
	cfg.MaxValueLength = 4096
	result := SanitizeSecrets(secrets, cfg)
	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}
	if len(result.Secrets["KEY"]) != 4096 {
		t.Errorf("expected value truncated to 4096, got %d", len(result.Secrets["KEY"]))
	}
}

func TestSanitizeSummary_NoViolations(t *testing.T) {
	summary := SanitizeSummary(nil)
	if summary != "sanitize: no violations" {
		t.Errorf("unexpected summary: %s", summary)
	}
}

func TestSanitizeSummary_WithViolations(t *testing.T) {
	violations := []SanitizeViolation{
		{Key: "FOO", Message: "key truncated from 200 to 128 chars"},
		{Key: "BAR", Message: "value truncated from 5000 to 4096 chars"},
	}
	summary := SanitizeSummary(violations)
	if !strings.Contains(summary, "2 violation") {
		t.Errorf("expected violation count in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "FOO") {
		t.Errorf("expected FOO in summary")
	}
}
