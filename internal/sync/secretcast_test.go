package sync

import (
	"testing"
)

func TestCastSecrets_NilSecrets(t *testing.T) {
	out, err := CastSecrets(nil, DefaultCastConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestCastSecrets_NoRules_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"PORT": "8080", "DEBUG": "true"}
	out, err := CastSecrets(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" || out["DEBUG"] != "true" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestCastSecrets_CastsInt(t *testing.T) {
	secrets := map[string]string{"PORT": " 8080 "}
	rules := []CastRule{{Key: "PORT", CastTo: CastInt}}
	out, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", out["PORT"])
	}
}

func TestCastSecrets_CastsBool(t *testing.T) {
	secrets := map[string]string{"DEBUG": "1"}
	rules := []CastRule{{Key: "DEBUG", CastTo: CastBool}}
	out, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DEBUG"] != "true" {
		t.Errorf("expected 'true', got %q", out["DEBUG"])
	}
}

func TestCastSecrets_CastsFloat(t *testing.T) {
	secrets := map[string]string{"RATE": "3.14"}
	rules := []CastRule{{Key: "RATE", CastTo: CastFloat}}
	out, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RATE"] != "3.14" {
		t.Errorf("expected '3.14', got %q", out["RATE"])
	}
}

func TestCastSecrets_InvalidInt_ReturnsError(t *testing.T) {
	secrets := map[string]string{"PORT": "not-a-number"}
	rules := []CastRule{{Key: "PORT", CastTo: CastInt}}
	_, err := CastSecrets(secrets, rules)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCastSecrets_UnknownCastType_ReturnsError(t *testing.T) {
	secrets := map[string]string{"X": "val"}
	rules := []CastRule{{Key: "X", CastTo: CastType("json")}}
	_, err := CastSecrets(secrets, rules)
	if err == nil {
		t.Fatal("expected error for unknown cast type")
	}
}

func TestCastSecrets_MissingKey_Skipped(t *testing.T) {
	secrets := map[string]string{"OTHER": "val"}
	rules := []CastRule{{Key: "MISSING", CastTo: CastInt}}
	out, err := CastSecrets(secrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MISSING"]; ok {
		t.Error("missing key should not appear in output")
	}
}

func TestCastSummary_Empty(t *testing.T) {
	s := CastSummary(nil)
	if s != "no cast rules applied" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestCastSummary_WithRules(t *testing.T) {
	rules := []CastRule{
		{Key: "A", CastTo: CastInt},
		{Key: "B", CastTo: CastBool},
	}
	s := CastSummary(rules)
	if s != "2 cast rule(s) applied" {
		t.Errorf("unexpected summary: %q", s)
	}
}
