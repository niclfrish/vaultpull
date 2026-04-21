package sync

import (
	"os"
	"testing"
)

func TestNewSecretCache_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/cache"
	c, err := NewSecretCache(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected cache dir to be created")
	}
}

func TestCache_PutAndGet(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir())
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}

	if err := c.Put("secret/app", "prod", secrets); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	entry, err := c.Get("secret/app", "prod")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if entry == nil {
		t.Fatal("expected cache entry, got nil")
	}
	if entry.Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", entry.Secrets["FOO"])
	}
	if entry.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
}

func TestCache_Get_MissReturnsNil(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir())
	entry, err := c.Get("secret/missing", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil for cache miss")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir())
	secrets := map[string]string{"KEY": "val"}
	_ = c.Put("secret/app", "staging", secrets)

	if err := c.Invalidate("secret/app", "staging"); err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	entry, _ := c.Get("secret/app", "staging")
	if entry != nil {
		t.Fatal("expected nil after invalidation")
	}
}

func TestCache_Invalidate_NonExistent(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir())
	if err := c.Invalidate("secret/nope", ""); err != nil {
		t.Fatalf("expected no error for missing entry, got: %v", err)
	}
}

func TestChecksumSecrets_Deterministic(t *testing.T) {
	s := map[string]string{"A": "1", "B": "2"}
	if checksumSecrets(s) != checksumSecrets(s) {
		t.Error("checksum should be deterministic")
	}
}
