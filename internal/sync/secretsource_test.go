package sync

import (
	"strings"
	"testing"
	"time"
)

func TestAnnotateWithSource_NilSecrets(t *testing.T) {
	result := AnnotateWithSource(nil, SecretSource{Type: SourceTypeVault})
	if result != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestAnnotateWithSource_InjectsKeys(t *testing.T) {
	src := SecretSource{
		Type:      SourceTypeVault,
		Location:  "secret/data/app",
		FetchedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Namespace: "prod",
	}
	secrets := map[string]string{"FOO": "bar"}
	out := AnnotateWithSource(secrets, src)

	if out["FOO"] != "bar" {
		t.Errorf("original key missing")
	}
	if out["__source_type"] != "vault" {
		t.Errorf("expected __source_type=vault, got %s", out["__source_type"])
	}
	if out["__source_location"] != "secret/data/app" {
		t.Errorf("unexpected location: %s", out["__source_location"])
	}
	if out["__source_namespace"] != "prod" {
		t.Errorf("expected namespace=prod, got %s", out["__source_namespace"])
	}
	if out["__source_fetched_at"] == "" {
		t.Error("expected fetched_at to be set")
	}
}

func TestAnnotateWithSource_NoNamespace(t *testing.T) {
	src := DefaultSecretSourceConfig(SourceTypeEnv, ".env")
	out := AnnotateWithSource(map[string]string{"A": "1"}, src)
	if _, ok := out["__source_namespace"]; ok {
		t.Error("namespace key should not be injected when empty")
	}
}

func TestStripSourceAnnotations_NilSecrets(t *testing.T) {
	result := StripSourceAnnotations(nil)
	if result != nil {
		t.Fatal("expected nil")
	}
}

func TestStripSourceAnnotations_RemovesSourceKeys(t *testing.T) {
	secrets := map[string]string{
		"FOO":               "bar",
		"__source_type":     "vault",
		"__source_location": "secret/app",
	}
	out := StripSourceAnnotations(secrets)
	if _, ok := out["__source_type"]; ok {
		t.Error("__source_type should be stripped")
	}
	if out["FOO"] != "bar" {
		t.Error("FOO should be preserved")
	}
}

func TestSourceSummary_ContainsFields(t *testing.T) {
	src := SecretSource{
		Type:      SourceTypeFile,
		Location:  "/etc/secrets",
		FetchedAt: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		Namespace: "staging",
	}
	summary := SourceSummary(src)
	for _, part := range []string{"file", "/etc/secrets", "staging"} {
		if !strings.Contains(summary, part) {
			t.Errorf("summary missing %q: %s", part, summary)
		}
	}
}

func TestSourceSummary_NoNamespace(t *testing.T) {
	src := SecretSource{Type: SourceTypeVault, Location: "kv/app", FetchedAt: time.Now()}
	summary := SourceSummary(src)
	if !strings.Contains(summary, "(none)") {
		t.Errorf("expected (none) for empty namespace: %s", summary)
	}
}
