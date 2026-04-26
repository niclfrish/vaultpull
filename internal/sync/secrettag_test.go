package sync

import (
	"strings"
	"testing"
)

func TestTagSecrets_NilSecrets(t *testing.T) {
	_, err := TagSecrets(nil, DefaultSecretTagConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestTagSecrets_EmptyPrefix(t *testing.T) {
	cfg := DefaultSecretTagConfig()
	cfg.Prefix = ""
	_, err := TagSecrets(map[string]string{"A": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestTagSecrets_InjectsSource(t *testing.T) {
	cfg := DefaultSecretTagConfig()
	cfg.Timestamp = false
	out, err := TagSecrets(map[string]string{"KEY": "val"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["__meta_source"]; got != "vault" {
		t.Errorf("expected source=vault, got %q", got)
	}
}

func TestTagSecrets_InjectsCount(t *testing.T) {
	cfg := DefaultSecretTagConfig()
	cfg.Timestamp = false
	input := map[string]string{"A": "1", "B": "2"}
	out, err := TagSecrets(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["__meta_count"]; got != "2" {
		t.Errorf("expected count=2, got %q", got)
	}
}

func TestTagSecrets_InjectsTimestamp(t *testing.T) {
	cfg := DefaultSecretTagConfig()
	out, err := TagSecrets(map[string]string{"X": "y"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts, ok := out["__meta_synced_at"]
	if !ok || ts == "" {
		t.Error("expected __meta_synced_at to be set")
	}
}

func TestTagSecrets_OriginalKeysPreserved(t *testing.T) {
	cfg := DefaultSecretTagConfig()
	cfg.Timestamp = false
	input := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, err := TagSecrets(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range input {
		if out[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestStripTagKeys_RemovesMetaKeys(t *testing.T) {
	secrets := map[string]string{
		"APP_KEY":         "abc",
		"__meta_source":   "vault",
		"__meta_count":    "1",
		"__meta_synced_at": "2024-01-01T00:00:00Z",
	}
	out := StripTagKeys(secrets, "__meta")
	for k := range out {
		if strings.HasPrefix(k, "__meta_") {
			t.Errorf("unexpected tag key in output: %q", k)
		}
	}
	if out["APP_KEY"] != "abc" {
		t.Error("expected APP_KEY to be preserved")
	}
}

func TestStripTagKeys_NilInput(t *testing.T) {
	out := StripTagKeys(nil, "__meta")
	if out != nil {
		t.Error("expected nil output for nil input")
	}
}
