package sync

import (
	"testing"
)

func TestNewMasker_InvalidPattern(t *testing.T) {
	cfg := MaskConfig{
		Patterns: []string{`[invalid`},
	}
	_, err := NewMasker(cfg)
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestNewMasker_InvalidPartialPattern(t *testing.T) {
	cfg := MaskConfig{
		PartialPatterns: []string{`[bad`},
	}
	_, err := NewMasker(cfg)
	if err == nil {
		t.Fatal("expected error for invalid partial pattern, got nil")
	}
}

func TestMasker_Apply_NoPatterns_ReturnsUnchanged(t *testing.T) {
	m, err := NewMasker(MaskConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := m.Apply(secrets)
	for k, v := range secrets {
		if out[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestMasker_Apply_FullMask(t *testing.T) {
	m, err := NewMasker(DefaultMaskConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"AUTH_TOKEN":  "tok123",
		"PLAIN_KEY":   "visible",
	}
	out := m.Apply(secrets)
	if out["DB_PASSWORD"] != "***********" {
		t.Errorf("DB_PASSWORD: expected fully masked, got %q", out["DB_PASSWORD"])
	}
	if out["AUTH_TOKEN"] != "******" {
		t.Errorf("AUTH_TOKEN: expected fully masked, got %q", out["AUTH_TOKEN"])
	}
}

func TestMasker_Apply_PartialMask(t *testing.T) {
	m, err := NewMasker(DefaultMaskConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{
		"API_KEY": "abcdefgh",
	}
	out := m.Apply(secrets)
	// VisibleChars=4, value len=8 => 4 masked + last 4 visible
	expected := "****efgh"
	if out["API_KEY"] != expected {
		t.Errorf("API_KEY: expected %q, got %q", expected, out["API_KEY"])
	}
}

func TestMasker_Apply_PartialMask_ShortValue(t *testing.T) {
	m, err := NewMasker(DefaultMaskConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{
		"ACCESS_KEY": "ab",
	}
	out := m.Apply(secrets)
	// value shorter than visibleChars => fully masked
	if out["ACCESS_KEY"] != "**" {
		t.Errorf("ACCESS_KEY: expected %q, got %q", "**", out["ACCESS_KEY"])
	}
}

func TestMasker_DefaultMaskChar(t *testing.T) {
	cfg := MaskConfig{
		Patterns:     []string{`(?i)secret`},
		MaskChar:     0, // should default to '*'
	}
	m, err := NewMasker(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := m.Apply(map[string]string{"MY_SECRET": "val"})
	if out["MY_SECRET"] != "***" {
		t.Errorf("expected '***', got %q", out["MY_SECRET"])
	}
}

func TestPartialMask_ZeroVisible(t *testing.T) {
	result := partialMask("hello", 0, '*')
	if result != "*****" {
		t.Errorf("expected fully masked, got %q", result)
	}
}
