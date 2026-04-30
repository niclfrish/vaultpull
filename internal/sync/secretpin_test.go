package sync

import (
	"strings"
	"testing"
)

func TestPinSecrets_NilSecrets(t *testing.T) {
	_, _, err := PinSecrets(nil, DefaultPinConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestPinSecrets_NoPins_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	out, summary, err := PinSecrets(secrets, DefaultPinConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
	if summary.Pinned != 0 || summary.Missing != 0 {
		t.Errorf("expected empty summary, got %+v", summary)
	}
}

func TestPinSecrets_AnnotatesMatchingKey(t *testing.T) {
	cfg := DefaultPinConfig()
	cfg.Pins = map[string]string{"DB_PASS": "v3"}
	secrets := map[string]string{"DB_PASS": "secret"}

	out, summary, err := PinSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Pinned != 1 {
		t.Errorf("expected 1 pinned, got %d", summary.Pinned)
	}
	annotation := out[cfg.AnnotationKey]
	if !strings.Contains(annotation, "DB_PASS@v3") {
		t.Errorf("expected annotation to contain DB_PASS@v3, got %q", annotation)
	}
}

func TestPinSecrets_MissingKeyNonStrict(t *testing.T) {
	cfg := DefaultPinConfig()
	cfg.Pins = map[string]string{"MISSING_KEY": "v1"}
	secrets := map[string]string{"OTHER": "val"}

	_, summary, err := PinSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error in non-strict mode: %v", err)
	}
	if summary.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", summary.Missing)
	}
}

func TestPinSecrets_MissingKeyStrict(t *testing.T) {
	cfg := DefaultPinConfig()
	cfg.StrictMode = true
	cfg.Pins = map[string]string{"MISSING_KEY": "v1"}
	secrets := map[string]string{"OTHER": "val"}

	_, _, err := PinSecrets(secrets, cfg)
	if err == nil {
		t.Fatal("expected error in strict mode for missing key")
	}
}

func TestPinSecrets_MultipleAnnotationsMerged(t *testing.T) {
	cfg := DefaultPinConfig()
	cfg.Pins = map[string]string{"KEY_A": "v1", "KEY_B": "v2"}
	secrets := map[string]string{"KEY_A": "alpha", "KEY_B": "beta"}

	out, summary, err := PinSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Pinned != 2 {
		t.Errorf("expected 2 pinned, got %d", summary.Pinned)
	}
	annotation := out[cfg.AnnotationKey]
	if !strings.Contains(annotation, "KEY_A@v1") || !strings.Contains(annotation, "KEY_B@v2") {
		t.Errorf("expected both annotations, got %q", annotation)
	}
}
