package sync

import (
	"errors"
	"testing"
)

func TestTransformer_Apply_NoFunctions(t *testing.T) {
	tr := NewTransformer()
	secrets := map[string]string{"KEY": "value"}
	out, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestTransformer_Apply_TrimSpace(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	secrets := map[string]string{
		"KEY": "  hello world  ",
		"OTHER": "clean",
	}
	out, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", out["KEY"])
	}
	if out["OTHER"] != "clean" {
		t.Errorf("expected 'clean', got %q", out["OTHER"])
	}
}

func TestTransformer_Apply_Redact(t *testing.T) {
	tr := NewTransformer(RedactTransform([]string{"password", "secret"}, "***"))
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_SECRET":  "abc123",
		"HOST":        "localhost",
	}
	out, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected '***', got %q", out["DB_PASSWORD"])
	}
	if out["API_SECRET"] != "***" {
		t.Errorf("expected '***', got %q", out["API_SECRET"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["HOST"])
	}
}

func TestTransformer_Apply_ErrorPropagates(t *testing.T) {
	failing := TransformFunc(func(key, value string) (string, error) {
		if key == "BAD" {
			return "", errors.New("invalid value")
		}
		return value, nil
	})
	tr := NewTransformer(failing)
	secrets := map[string]string{"BAD": "x", "GOOD": "y"}
	_, err := tr.Apply(secrets)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTransformer_Apply_ChainedFunctions(t *testing.T) {
	tr := NewTransformer(
		TrimSpaceTransform(),
		RedactTransform([]string{"token"}, "[REDACTED]"),
	)
	secrets := map[string]string{
		"API_TOKEN": "  mytoken123  ",
		"NAME":      "  alice  ",
	}
	out, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_TOKEN"] != "[REDACTED]" {
		t.Errorf("expected '[REDACTED]', got %q", out["API_TOKEN"])
	}
	if out["NAME"] != "alice" {
		t.Errorf("expected 'alice', got %q", out["NAME"])
	}
}

func TestTransformer_Apply_DoesNotMutateInput(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	secrets := map[string]string{"KEY": "  value  "}
	_, err := tr.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["KEY"] != "  value  " {
		t.Error("input map was mutated")
	}
}

func TestTransformer_Apply_EmptyMap(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	out, err := tr.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
