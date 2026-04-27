package sync

import (
	"strings"
	"testing"
)

func TestRedactSecrets_NilSecrets(t *testing.T) {
	out, summary, err := RedactSecrets(nil, DefaultRedactConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Errorf("expected nil output, got %v", out)
	}
	if len(summary.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys, got %v", summary.RedactedKeys)
	}
}

func TestRedactSecrets_RedactsSensitiveKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"APP_NAME":    "vaultpull",
	}
	out, summary, err := RedactSecrets(secrets, DefaultRedactConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY to be redacted, got %q", out["API_KEY"])
	}
	if out["APP_NAME"] != "vaultpull" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out["APP_NAME"])
	}
	if summary.Total != 3 {
		t.Errorf("expected total 3, got %d", summary.Total)
	}
	if len(summary.RedactedKeys) != 2 {
		t.Errorf("expected 2 redacted keys, got %d", len(summary.RedactedKeys))
	}
}

func TestRedactSecrets_CustomReplacement(t *testing.T) {
	secrets := map[string]string{"MY_TOKEN": "tok-xyz"}
	cfg := RedactConfig{
		Patterns:    []string{`(?i)token`},
		Replacement: "***",
	}
	out, _, err := RedactSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MY_TOKEN"] != "***" {
		t.Errorf("expected '***', got %q", out["MY_TOKEN"])
	}
}

func TestRedactSecrets_InvalidPattern(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	cfg := RedactConfig{Patterns: []string{`[invalid`}}
	_, _, err := RedactSecrets(secrets, cfg)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRedactSecrets_EmptyPatterns_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"SECRET_KEY": "value"}
	cfg := RedactConfig{Patterns: []string{}, Replacement: "[REDACTED]"}
	out, summary, err := RedactSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET_KEY"] != "value" {
		t.Errorf("expected unchanged value, got %q", out["SECRET_KEY"])
	}
	if len(summary.RedactedKeys) != 0 {
		t.Errorf("expected no redacted keys")
	}
}

func TestRedactSummary_String(t *testing.T) {
	s := RedactSummary{RedactedKeys: []string{"DB_PASSWORD", "API_KEY"}, Total: 5}
	result := s.String()
	if !strings.Contains(result, "2/5") {
		t.Errorf("expected '2/5' in summary string, got %q", result)
	}
	if !strings.Contains(result, "DB_PASSWORD") {
		t.Errorf("expected key name in summary string")
	}
}
