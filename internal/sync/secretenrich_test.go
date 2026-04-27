package sync

import (
	"testing"
)

func TestEnrichSecrets_NilSecrets(t *testing.T) {
	_, err := EnrichSecrets(nil, DefaultEnrichConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestEnrichSecrets_NoConfig_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	out, err := EnrichSecrets(secrets, DefaultEnrichConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value unchanged, got %q", out["KEY"])
	}
}

func TestEnrichSecrets_AppliesPrefixAndSuffix(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc123", "HOST": "localhost"}
	cfg := EnrichConfig{Prefix: "vault://", Suffix: "!"}
	out, err := EnrichSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "vault://abc123!" {
		t.Errorf("expected vault://abc123!, got %q", out["TOKEN"])
	}
	if out["HOST"] != "vault://localhost!" {
		t.Errorf("expected vault://localhost!, got %q", out["HOST"])
	}
}

func TestEnrichSecrets_OnlyKeys_LimitsDecoration(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	cfg := EnrichConfig{Prefix: ">>>", OnlyKeys: []string{"B"}}
	out, err := EnrichSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("A should be unchanged, got %q", out["A"])
	}
	if out["B"] != ">>>2" {
		t.Errorf("B should be decorated, got %q", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("C should be unchanged, got %q", out["C"])
	}
}

func TestEnrichSecrets_InjectsStaticKeys(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	cfg := EnrichConfig{StaticKeys: map[string]string{"STATIC": "injected"}}
	out, err := EnrichSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["STATIC"] != "injected" {
		t.Errorf("expected injected, got %q", out["STATIC"])
	}
}

func TestEnrichSecrets_StaticKeyConflict_ReturnsError(t *testing.T) {
	secrets := map[string]string{"CONFLICT": "original"}
	cfg := EnrichConfig{StaticKeys: map[string]string{"CONFLICT": "override"}}
	_, err := EnrichSecrets(secrets, cfg)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestEnrichSummary_WithPrefixAndStaticKeys(t *testing.T) {
	original := map[string]string{"A": "1", "B": "2"}
	enriched := map[string]string{"A": ">>1", "B": ">>2", "STATIC": "x"}
	cfg := EnrichConfig{
		Prefix:     ">>",
		StaticKeys: map[string]string{"STATIC": "x"},
	}
	summary := EnrichSummary(original, enriched, cfg)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestEnrichSummary_NoDecorations(t *testing.T) {
	original := map[string]string{"K": "v"}
	enriched := map[string]string{"K": "v"}
	summary := EnrichSummary(original, enriched, DefaultEnrichConfig())
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
