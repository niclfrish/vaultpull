package sync

import (
	"testing"
)

func TestMaskSecrets_NilSecrets(t *testing.T) {
	cfg := DefaultSecretMaskConfig()
	res, err := MaskSecrets(nil, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Secrets) != 0 {
		t.Errorf("expected empty map, got %v", res.Secrets)
	}
}

func TestMaskSecrets_MasksSensitiveKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abcd1234efgh",
		"APP_NAME":    "myapp",
	}
	cfg := DefaultSecretMaskConfig()
	res, err := MaskSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["APP_NAME"] != "myapp" {
		t.Errorf("non-sensitive key should be unchanged, got %q", res.Secrets["APP_NAME"])
	}
	if res.Secrets["DB_PASSWORD"] == "supersecret" {
		t.Error("DB_PASSWORD should be masked")
	}
	if res.Secrets["API_KEY"] == "abcd1234efgh" {
		t.Error("API_KEY should be masked")
	}
	if len(res.MaskedKeys) != 2 {
		t.Errorf("expected 2 masked keys, got %d", len(res.MaskedKeys))
	}
}

func TestMaskSecrets_VisibleChars(t *testing.T) {
	secrets := map[string]string{
		"MY_TOKEN": "abcdefgh",
	}
	cfg := DefaultSecretMaskConfig()
	cfg.VisibleChars = 3
	res, err := MaskSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := res.Secrets["MY_TOKEN"]
	if got[len(got)-3:] != "fgh" {
		t.Errorf("expected last 3 chars visible, got %q", got)
	}
}

func TestMaskSecrets_ZeroVisibleChars(t *testing.T) {
	secrets := map[string]string{
		"SECRET_KEY": "hunter2",
	}
	cfg := DefaultSecretMaskConfig()
	cfg.VisibleChars = 0
	res, err := MaskSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["SECRET_KEY"] != "*******" {
		t.Errorf("expected full mask, got %q", res.Secrets["SECRET_KEY"])
	}
}

func TestMaskSecrets_InvalidPattern(t *testing.T) {
	secrets := map[string]string{"key": "val"}
	cfg := DefaultSecretMaskConfig()
	cfg.KeyPatterns = []string{`[invalid`}
	_, err := MaskSecrets(secrets, cfg)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestApplyMask_EmptyValue(t *testing.T) {
	result := applyMask("", "*", 4)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestApplyMask_VisibleCharsExceedsLength(t *testing.T) {
	result := applyMask("abc", "*", 10)
	if result != "***" {
		t.Errorf("expected full mask when visibleChars >= len, got %q", result)
	}
}
