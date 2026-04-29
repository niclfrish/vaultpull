package sync

import (
	"testing"
)

func TestPromoteSecrets_NilSecrets(t *testing.T) {
	_, _, err := PromoteSecrets(nil, DefaultPromoteConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestPromoteSecrets_EmptyToPrefix(t *testing.T) {
	cfg := DefaultPromoteConfig()
	cfg.ToPrefix = ""
	_, _, err := PromoteSecrets(map[string]string{"A": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty ToPrefix")
	}
}

func TestPromoteSecrets_CopiesMatchingKeys(t *testing.T) {
	secrets := map[string]string{
		"dev_DB_HOST": "localhost",
		"dev_DB_PASS": "secret",
		"OTHER":       "keep",
	}
	cfg := PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_", Overwrite: false}
	out, result, err := PromoteSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["prod_DB_HOST"] != "localhost" {
		t.Errorf("expected prod_DB_HOST=localhost, got %q", out["prod_DB_HOST"])
	}
	if out["prod_DB_PASS"] != "secret" {
		t.Errorf("expected prod_DB_PASS=secret, got %q", out["prod_DB_PASS"])
	}
	if len(result.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(result.Promoted))
	}
	if out["OTHER"] != "keep" {
		t.Error("original keys should be preserved")
	}
}

func TestPromoteSecrets_SkipsExistingWithoutOverwrite(t *testing.T) {
	secrets := map[string]string{
		"dev_KEY": "new",
		"prod_KEY": "existing",
	}
	cfg := PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_", Overwrite: false}
	out, result, err := PromoteSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["prod_KEY"] != "existing" {
		t.Errorf("expected existing value preserved, got %q", out["prod_KEY"])
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
}

func TestPromoteSecrets_OverwritesWhenEnabled(t *testing.T) {
	secrets := map[string]string{
		"dev_KEY":  "new",
		"prod_KEY": "old",
	}
	cfg := PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_", Overwrite: true}
	out, result, err := PromoteSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["prod_KEY"] != "new" {
		t.Errorf("expected overwritten value, got %q", out["prod_KEY"])
	}
	if len(result.Overwrote) != 1 {
		t.Errorf("expected 1 overwrote, got %d", len(result.Overwrote))
	}
}

func TestPromoteSecrets_DryRun_DoesNotMutate(t *testing.T) {
	secrets := map[string]string{"dev_KEY": "val"}
	cfg := PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_", DryRun: true}
	out, result, err := PromoteSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := out["prod_KEY"]; exists {
		t.Error("dry run should not write destination key")
	}
	if len(result.Promoted) != 1 {
		t.Errorf("expected 1 in promoted list for dry run, got %d", len(result.Promoted))
	}
}

func TestPromoteSummary_Counts(t *testing.T) {
	r := PromoteResult{
		Promoted:  []string{"a", "b"},
		Skipped:   []string{"c"},
		Overwrote: []string{"d"},
	}
	got := PromoteSummary(r)
	want := "promoted=2 skipped=1 overwrote=1"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
