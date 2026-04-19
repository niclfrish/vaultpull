package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

func newVaultServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestRun_Success(t *testing.T) {
	server := newVaultServer(t, map[string]interface{}{"KEY": "value"})
	defer server.Close()

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, ".env")

	cfg := &config.Config{
		VaultAddr:  server.URL,
		VaultToken: "test-token",
		SecretPath: "secret/data/app",
		OutputFile: outFile,
	}

	s, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestRun_EmptySecrets(t *testing.T) {
	server := newVaultServer(t, map[string]interface{}{})
	defer server.Close()

	tmpDir := t.TempDir()
	cfg := &config.Config{
		VaultAddr:  server.URL,
		VaultToken: "test-token",
		SecretPath: "secret/data/empty",
		OutputFile: filepath.Join(tmpDir, ".env"),
	}

	s, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	if err := s.Run(); err == nil {
		t.Error("expected error for empty secrets")
	}
}
