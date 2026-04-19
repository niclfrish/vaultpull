package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".env")
}

func TestWrite_CreatesFile(t *testing.T) {
	path := tempFile(t)
	w := New(path)
	err := w.Write(map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to exist")
	}
}

func TestWrite_ContentSortedAndSanitized(t *testing.T) {
	path := tempFile(t)
	w := New(path)
	secrets := map[string]string{
		"db-host": "localhost",
		"api-key": "secret123",
	}
	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "API_KEY=secret123" {
		t.Errorf("unexpected line: %q", lines[0])
	}
	if lines[1] != "DB_HOST=localhost" {
		t.Errorf("unexpected line: %q", lines[1])
	}
}

func TestWrite_EscapesValuesWithSpaces(t *testing.T) {
	path := tempFile(t)
	w := New(path)
	if err := w.Write(map[string]string{"MSG": "hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(path)
	if !strings.Contains(string(content), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", content)
	}
}

func TestWrite_InvalidPath(t *testing.T) {
	w := New("/nonexistent/dir/.env")
	err := w.Write(map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
