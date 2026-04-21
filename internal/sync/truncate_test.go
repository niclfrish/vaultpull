package sync

import (
	"strings"
	"testing"
)

func TestTruncateValue_ShortValue(t *testing.T) {
	cfg := DefaultTruncateConfig()
	input := "short"
	got := TruncateValue(input, cfg)
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestTruncateValue_ExactLength(t *testing.T) {
	cfg := TruncateConfig{MaxLength: 5, Suffix: "..."}
	input := "hello"
	got := TruncateValue(input, cfg)
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestTruncateValue_LongValue(t *testing.T) {
	cfg := TruncateConfig{MaxLength: 10, Suffix: "..."}
	input := "this is a very long secret value"
	got := TruncateValue(input, cfg)
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
	if len([]rune(got)) != 13 { // 10 + len("...")
		t.Errorf("unexpected length %d for %q", len(got), got)
	}
}

func TestTruncateValue_ZeroMaxLength(t *testing.T) {
	cfg := TruncateConfig{MaxLength: 0, Suffix: "..."}
	input := "should not be truncated"
	got := TruncateValue(input, cfg)
	if got != input {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestTruncateValue_CustomSuffix(t *testing.T) {
	cfg := TruncateConfig{MaxLength: 4, Suffix: "[…]"}
	input := "supersecret"
	got := TruncateValue(input, cfg)
	if !strings.HasSuffix(got, "[…]") {
		t.Errorf("expected custom suffix, got %q", got)
	}
}

func TestTruncateSecrets_AllValues(t *testing.T) {
	cfg := TruncateConfig{MaxLength: 5, Suffix: "..."}
	secrets := map[string]string{
		"SHORT": "hi",
		"LONG":  "this is definitely too long",
	}
	out := TruncateSecrets(secrets, cfg)
	if out["SHORT"] != "hi" {
		t.Errorf("short value should be unchanged, got %q", out["SHORT"])
	}
	if !strings.HasSuffix(out["LONG"], "...") {
		t.Errorf("long value should be truncated, got %q", out["LONG"])
	}
}

func TestTruncateSummary_Truncated(t *testing.T) {
	original := "supersecretvalue"
	truncated := "super..."
	summary := TruncateSummary(original, truncated)
	if summary == "" {
		t.Error("expected non-empty summary for truncated value")
	}
}

func TestTruncateSummary_NotTruncated(t *testing.T) {
	original := "hi"
	truncated := "hi"
	summary := TruncateSummary(original, truncated)
	if summary != "" {
		t.Errorf("expected empty summary, got %q", summary)
	}
}
