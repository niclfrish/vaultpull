package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint records the last successful sync state for a given path/namespace.
type Checkpoint struct {
	Path      string            `json:"path"`
	Namespace string            `json:"namespace"`
	SyncedAt  time.Time         `json:"synced_at"`
	KeyCount  int               `json:"key_count"`
	Checksum  string            `json:"checksum"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// CheckpointStore persists and retrieves sync checkpoints.
type CheckpointStore struct {
	dir string
}

// NewCheckpointStore creates a CheckpointStore backed by dir, creating it if needed.
func NewCheckpointStore(dir string) (*CheckpointStore, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("checkpoint: create dir: %w", err)
	}
	return &CheckpointStore{dir: dir}, nil
}

// Save writes a checkpoint to disk.
func (s *CheckpointStore) Save(cp Checkpoint) error {
	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	path := filepath.Join(s.dir, s.fileName(cp.Path, cp.Namespace))
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("checkpoint: write: %w", err)
	}
	return nil
}

// Load retrieves the last checkpoint for the given vault path and namespace.
// Returns nil, nil if no checkpoint exists yet.
func (s *CheckpointStore) Load(vaultPath, namespace string) (*Checkpoint, error) {
	path := filepath.Join(s.dir, s.fileName(vaultPath, namespace))
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("checkpoint: read: %w", err)
	}
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}
	return &cp, nil
}

// Delete removes the checkpoint for the given vault path and namespace.
func (s *CheckpointStore) Delete(vaultPath, namespace string) error {
	path := filepath.Join(s.dir, s.fileName(vaultPath, namespace))
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("checkpoint: delete: %w", err)
	}
	return nil
}

func (s *CheckpointStore) fileName(vaultPath, namespace string) string {
	if namespace == "" {
		return fmt.Sprintf("checkpoint_%s.json", sanitizeKey(vaultPath))
	}
	return fmt.Sprintf("checkpoint_%s_%s.json", sanitizeKey(namespace), sanitizeKey(vaultPath))
}
