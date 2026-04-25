package sync

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRotateTokenAndInject_NilRotator(t *testing.T) {
	_, err := RotateTokenAndInject(nil, "KEY", map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil rotator")
	}
}

func TestRotateTokenAndInject_NilSecrets(t *testing.T) {
	r, _ := NewTokenRotator("tok", func() (string, error) { return "tok", nil }, DefaultTokenRotateConfig())
	_, err := RotateTokenAndInject(r, "KEY", nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestRotateTokenAndInject_DefaultKey(t *testing.T) {
	r, _ := NewTokenRotator("mytoken", func() (string, error) { return "mytoken", nil }, DefaultTokenRotateConfig())
	out, err := RotateTokenAndInject(r, "", map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAULT_TOKEN"] != "mytoken" {
		t.Errorf("expected VAULT_TOKEN=mytoken, got %q", out["VAULT_TOKEN"])
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
}

func TestRotateTokenAndInject_CustomKey(t *testing.T) {
	r, _ := NewTokenRotator("secret", func() (string, error) { return "secret", nil }, DefaultTokenRotateConfig())
	out, err := RotateTokenAndInject(r, "APP_TOKEN", map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_TOKEN"] != "secret" {
		t.Errorf("expected APP_TOKEN=secret, got %q", out["APP_TOKEN"])
	}
}

func TestRotateTokenAndInject_FetchError(t *testing.T) {
	cfg := TokenRotateConfig{RotateAfter: 1 * time.Millisecond, GracePeriod: 1 * time.Millisecond}
	r, _ := NewTokenRotator("old", func() (string, error) {
		return "", errors.New("fetch failed")
	}, cfg)
	time.Sleep(10 * time.Millisecond)

	_, err := RotateTokenAndInject(r, "", map[string]string{})
	if err == nil {
		t.Fatal("expected error when token fetch fails beyond grace period")
	}
}

func TestLogTokenAge_NilRotator(t *testing.T) {
	var buf bytes.Buffer
	LogTokenAge(nil, &buf)
	if !strings.Contains(buf.String(), "nil") {
		t.Errorf("expected nil mention in output, got: %q", buf.String())
	}
}

func TestLogTokenAge_WritesToBuffer(t *testing.T) {
	r, _ := NewTokenRotator("tok", func() (string, error) { return "tok", nil }, DefaultTokenRotateConfig())
	var buf bytes.Buffer
	LogTokenAge(r, &buf)
	if !strings.Contains(buf.String(), "age") {
		t.Errorf("expected age in output, got: %q", buf.String())
	}
}

func TestLogTokenAge_NilWriter(t *testing.T) {
	r, _ := NewTokenRotator("tok", func() (string, error) { return "tok", nil }, DefaultTokenRotateConfig())
	// Should not panic with nil writer (falls back to stdout).
	LogTokenAge(r, nil)
}
