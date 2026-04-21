package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of secrets.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
	Namespace string            `json:"namespace,omitempty"`
	Secrets   map[string]string `json:"secrets"`
}

// SnapshotStore persists and retrieves snapshots from disk.
type SnapshotStore struct {
	dir string
}

// NewSnapshotStore creates a SnapshotStore backed by the given directory.
func NewSnapshotStore(dir string) (*SnapshotStore, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &SnapshotStore{dir: dir}, nil
}

// Save writes a snapshot to disk, keyed by path and namespace.
func (s *SnapshotStore) Save(snap Snapshot) error {
	snap.Timestamp = time.Now().UTC()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	fileName := snapshotFileName(snap.Path, snap.Namespace)
	return os.WriteFile(filepath.Join(s.dir, fileName), data, 0600)
}

// Load retrieves the most recent snapshot for a given path and namespace.
func (s *SnapshotStore) Load(path, namespace string) (*Snapshot, error) {
	fileName := snapshotFileName(path, namespace)
	data, err := os.ReadFile(filepath.Join(s.dir, fileName))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

// snapshotFileName produces a stable filename from path and namespace.
func snapshotFileName(path, namespace string) string {
	safe := func(s string) string {
		out := make([]byte, len(s))
		for i := range s {
			if s[i] == '/' || s[i] == '\\' || s[i] == ':' {
				out[i] = '_'
			} else {
				out[i] = s[i]
			}
		}
		return string(out)
	}
	if namespace == "" {
		return safe(path) + ".snapshot.json"
	}
	return safe(namespace) + "_" + safe(path) + ".snapshot.json"
}
