package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry holds a cached snapshot of secrets with metadata.
type CacheEntry struct {
	Path      string            `json:"path"`
	Namespace string            `json:"namespace"`
	Secrets   map[string]string `json:"secrets"`
	Checksum  string            `json:"checksum"`
	FetchedAt time.Time         `json:"fetched_at"`
}

// SecretCache manages on-disk caching of Vault secrets.
type SecretCache struct {
	cacheDir string
}

// NewSecretCache creates a SecretCache backed by the given directory.
func NewSecretCache(cacheDir string) (*SecretCache, error) {
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("cache: create dir: %w", err)
	}
	return &SecretCache{cacheDir: cacheDir}, nil
}

// Put stores a secrets map under a key derived from path and namespace.
func (c *SecretCache) Put(path, namespace string, secrets map[string]string) error {
	entry := CacheEntry{
		Path:      path,
		Namespace: namespace,
		Secrets:   secrets,
		Checksum:  checksumSecrets(secrets),
		FetchedAt: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	return os.WriteFile(c.cacheFile(path, namespace), data, 0600)
}

// Get retrieves a cached entry. Returns nil, nil when no cache exists.
func (c *SecretCache) Get(path, namespace string) (*CacheEntry, error) {
	data, err := os.ReadFile(c.cacheFile(path, namespace))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read: %w", err)
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	return &entry, nil
}

// Invalidate removes the cached entry for the given path and namespace.
func (c *SecretCache) Invalidate(path, namespace string) error {
	err := os.Remove(c.cacheFile(path, namespace))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *SecretCache) cacheFile(path, namespace string) string {
	key := fmt.Sprintf("%s::%s", namespace, path)
	h := sha256.Sum256([]byte(key))
	return filepath.Join(c.cacheDir, hex.EncodeToString(h[:8])+".json")
}

func checksumSecrets(secrets map[string]string) string {
	data, _ := json.Marshal(secrets)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
