package sync

import (
	"testing"
)

func TestNormalizeSecrets_NilSecrets(t *testing.T) {
	_, _, err := NormalizeSecrets(nil, DefaultNormalizeConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestNormalizeSecrets_UppercaseKeys(t *testing.T) {
	cfg := DefaultNormalizeConfig()
	cfg.TrimValues = false

	out, summary, err := NormalizeSecrets(map[string]string{"db_host": "localhost", "DB_PORT": "5432"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
	if _, ok := out["DB_PORT"]; !ok {
		t.Error("expected DB_PORT in output")
	}
	if summary.Total != 2 {
		t.Errorf("expected total=2, got %d", summary.Total)
	}
}

func TestNormalizeSecrets_ReplaceHyphens(t *testing.T) {
	cfg := DefaultNormalizeConfig()
	cfg.UppercaseKeys = false
	cfg.TrimValues = false

	out, _, err := NormalizeSecrets(map[string]string{"my-key": "val"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["my_key"]; !ok {
		t.Error("expected my_key after hyphen replacement")
	}
}

func TestNormalizeSecrets_ReplaceDots(t *testing.T) {
	cfg := DefaultNormalizeConfig()
	cfg.UppercaseKeys = false
	cfg.TrimValues = false
	cfg.ReplaceHyphens = false

	out, _, err := NormalizeSecrets(map[string]string{"app.config": "v"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["app_config"]; !ok {
		t.Error("expected app_config after dot replacement")
	}
}

func TestNormalizeSecrets_TrimValues(t *testing.T) {
	cfg := DefaultNormalizeConfig()
	cfg.UppercaseKeys = false
	cfg.ReplaceHyphens = false
	cfg.ReplaceDots = false

	out, summary, err := NormalizeSecrets(map[string]string{"KEY": "  hello  "}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello" {
		t.Errorf("expected trimmed value 'hello', got %q", out["KEY"])
	}
	if summary.Modified != 1 {
		t.Errorf("expected modified=1, got %d", summary.Modified)
	}
}

func TestNormalizeSecrets_NoChanges(t *testing.T) {
	cfg := DefaultNormalizeConfig()

	out, summary, err := NormalizeSecrets(map[string]string{"MY_KEY": "value"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MY_KEY"] != "value" {
		t.Errorf("expected unchanged value")
	}
	if summary.Skipped != 1 {
		t.Errorf("expected skipped=1, got %d", summary.Skipped)
	}
}

func TestNormalizeSecrets_EmptySecrets(t *testing.T) {
	out, summary, err := NormalizeSecrets(map[string]string{}, DefaultNormalizeConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Error("expected empty output")
	}
	if summary.Total != 0 {
		t.Errorf("expected total=0, got %d", summary.Total)
	}
}
