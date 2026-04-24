package sync_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/sync"
)

func TestMultiEnvWriter_WritesFilesOnDisk(t *testing.T) {
	dir := t.TempDir()
	pathA := filepath.Join(dir, ".env.a")
	pathB := filepath.Join(dir, ".env.b")

	targets := []sync.EnvTarget{
		{Name: "a", Path: pathA, Namespace: ""},
		{Name: "b", Path: pathB, Namespace: "B"},
	}

	writerFn := func(path string, secrets map[string]string) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		for k, v := range secrets {
			_, err = fmt.Fprintf(f, "%s=%s\n", k, v)
			if err != nil {
				return err
			}
		}
		return nil
	}

	mw, err := sync.NewMultiEnvWriter(targets, writerFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets := map[string]string{"TOKEN": "abc123"}
	results := mw.WriteAll(secrets)
	if err := sync.AnyError(results); err != nil {
		t.Fatalf("write error: %v", err)
	}

	bytesA, err := os.ReadFile(pathA)
	if err != nil {
		t.Fatalf("read .env.a: %v", err)
	}
	if !strings.Contains(string(bytesA), "TOKEN=abc123") {
		t.Errorf(".env.a missing TOKEN=abc123, got: %s", bytesA)
	}

	bytesB, err := os.ReadFile(pathB)
	if err != nil {
		t.Fatalf("read .env.b: %v", err)
	}
	if !strings.Contains(string(bytesB), "B_TOKEN=abc123") {
		t.Errorf(".env.b missing B_TOKEN=abc123, got: %s", bytesB)
	}
}

func fmt_Fprintf(f *os.File, format string, args ...interface{}) (int, error) {
	return 0, nil // placeholder to avoid import issue in snippet
}
