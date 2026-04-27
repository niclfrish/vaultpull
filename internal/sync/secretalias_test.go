package sync

import (
	"testing"
)

func TestApplyAliases_NilSecrets(t *testing.T) {
	_, err := ApplyAliases(nil, DefaultAliasConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestApplyAliases_NoAliases_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	result, err := ApplyAliases(secrets, DefaultAliasConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", result["FOO"])
	}
}

func TestApplyAliases_RenamesKey(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "value1"}
	cfg := AliasConfig{Aliases: map[string]string{"OLD_KEY": "NEW_KEY"}, KeepOriginal: false}
	result, err := ApplyAliases(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["OLD_KEY"]; ok {
		t.Error("original key should have been removed")
	}
	if result["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", result["NEW_KEY"])
	}
}

func TestApplyAliases_KeepOriginal(t *testing.T) {
	secrets := map[string]string{"SRC": "hello"}
	cfg := AliasConfig{Aliases: map[string]string{"SRC": "DST"}, KeepOriginal: true}
	result, err := ApplyAliases(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["SRC"] != "hello" {
		t.Error("original key should be retained")
	}
	if result["DST"] != "hello" {
		t.Error("alias key should be present")
	}
}

func TestApplyAliases_ConflictReturnsError(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	cfg := AliasConfig{Aliases: map[string]string{"A": "B"}, KeepOriginal: false}
	_, err := ApplyAliases(secrets, cfg)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestApplyAliases_EmptyTargetReturnsError(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	cfg := AliasConfig{Aliases: map[string]string{"FOO": "   "}, KeepOriginal: false}
	_, err := ApplyAliases(secrets, cfg)
	if err == nil {
		t.Fatal("expected error for empty alias target")
	}
}

func TestApplyAliases_MissingSourceKey_Skipped(t *testing.T) {
	secrets := map[string]string{"OTHER": "val"}
	cfg := AliasConfig{Aliases: map[string]string{"MISSING": "NEW"}, KeepOriginal: false}
	result, err := ApplyAliases(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["NEW"]; ok {
		t.Error("alias should not be created for missing source key")
	}
}

func TestAliasSummary_NoAliases(t *testing.T) {
	s := AliasSummary(DefaultAliasConfig())
	if s != "no aliases configured" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestAliasSummary_WithAliases(t *testing.T) {
	cfg := AliasConfig{Aliases: map[string]string{"A": "B"}}
	s := AliasSummary(cfg)
	if s == "no aliases configured" {
		t.Error("expected non-empty summary")
	}
}
