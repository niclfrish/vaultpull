package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestRead_NonExistentFile(t *testing.T) {
	r := NewReader("/nonexistent/.env")
	got, err := r.Read()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRead_ParsesKeyValues(t *testing.T) {
	p := writeEnvFile(t, "FOO=bar\nBAZ=\"hello world\"\n# comment\n\nQUX=123\n")
	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]string{"FOO": "bar", "BAZ": "hello world", "QUX": "123"}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("key %s: want %q, got %q", k, v, got[k])
		}
	}
	if len(got) != len(expected) {
		t.Errorf("expected %d keys, got %d", len(expected), len(got))
	}
}

func TestRead_IgnoresInvalidLines(t *testing.T) {
	p := writeEnvFile(t, "VALID=yes\nNOEQUALSIGN\n")
	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d: %v", len(got), got)
	}
}
