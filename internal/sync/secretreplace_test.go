package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestReplaceSecrets_NilSecrets(t *testing.T) {
	_, _, err := ReplaceSecrets(nil, DefaultReplaceConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestReplaceSecrets_NoReplacements_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"KEY": "hello world"}
	cfg := DefaultReplaceConfig()
	out, summary, err := ReplaceSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello world" {
		t.Errorf("expected unchanged value, got %q", out["KEY"])
	}
	if summary.Skipped != 1 || summary.Modified != 0 {
		t.Errorf("unexpected summary: %+v", summary)
	}
}

func TestReplaceSecrets_ReplacesSubstring(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://localhost:5432/mydb"}
	cfg := DefaultReplaceConfig()
	cfg.Replacements = map[string]string{"localhost": "db.internal"}
	out, summary, err := ReplaceSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out["DB_URL"], "db.internal") {
		t.Errorf("expected replacement, got %q", out["DB_URL"])
	}
	if summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", summary.Modified)
	}
}

func TestReplaceSecrets_OnlyKeys_LimitsScope(t *testing.T) {
	secrets := map[string]string{
		"A": "foo bar",
		"B": "foo baz",
	}
	cfg := DefaultReplaceConfig()
	cfg.Replacements = map[string]string{"foo": "qux"}
	cfg.OnlyKeys = []string{"A"}
	out, summary, err := ReplaceSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "qux bar" {
		t.Errorf("expected A replaced, got %q", out["A"])
	}
	if out["B"] != "foo baz" {
		t.Errorf("expected B unchanged, got %q", out["B"])
	}
	if summary.Modified != 1 || summary.Skipped != 1 {
		t.Errorf("unexpected summary: %+v", summary)
	}
}

func TestReplaceSecrets_CaseInsensitive(t *testing.T) {
	secrets := map[string]string{"MSG": "Hello WORLD"}
	cfg := DefaultReplaceConfig()
	cfg.CaseSensitive = false
	cfg.Replacements = map[string]string{"world": "Go"}
	out, _, err := ReplaceSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out["MSG"], "Go") {
		t.Errorf("expected case-insensitive replacement, got %q", out["MSG"])
	}
}

func TestReplaceStage_Name(t *testing.T) {
	stage := ReplaceStage(DefaultReplaceConfig())
	if stage.Name != "replace" {
		t.Errorf("expected stage name 'replace', got %q", stage.Name)
	}
}

func TestReplaceAndReport_NilSecrets(t *testing.T) {
	_, err := ReplaceAndReport(nil, DefaultReplaceConfig(), nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestReplaceAndReport_WritesOutput(t *testing.T) {
	secrets := map[string]string{"KEY": "old_value"}
	cfg := DefaultReplaceConfig()
	cfg.Replacements = map[string]string{"old": "new"}
	var buf bytes.Buffer
	out, err := ReplaceAndReport(secrets, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "new_value" {
		t.Errorf("expected replacement, got %q", out["KEY"])
	}
	if !strings.Contains(buf.String(), "modified=1") {
		t.Errorf("expected summary in output, got %q", buf.String())
	}
}
