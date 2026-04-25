package sync

import (
	"os"
	"testing"
	"time"
)

func TestNewCheckpointStore_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/checkpoints"
	_, err := NewCheckpointStore(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestCheckpoint_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewCheckpointStore(dir)

	cp := Checkpoint{
		Path:      "secret/myapp",
		Namespace: "prod",
		SyncedAt:  time.Now().UTC().Truncate(time.Second),
		KeyCount:  5,
		Checksum:  "abc123",
		Meta:      map[string]string{"env": "production"},
	}

	if err := store.Save(cp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load("secret/myapp", "prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected checkpoint, got nil")
	}
	if loaded.KeyCount != cp.KeyCount {
		t.Errorf("KeyCount: want %d, got %d", cp.KeyCount, loaded.KeyCount)
	}
	if loaded.Checksum != cp.Checksum {
		t.Errorf("Checksum: want %q, got %q", cp.Checksum, loaded.Checksum)
	}
	if loaded.Meta["env"] != "production" {
		t.Errorf("Meta: want 'production', got %q", loaded.Meta["env"])
	}
}

func TestCheckpoint_Load_NonExistent(t *testing.T) {
	store, _ := NewCheckpointStore(t.TempDir())
	cp, err := store.Load("secret/missing", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp != nil {
		t.Fatalf("expected nil checkpoint, got %+v", cp)
	}
}

func TestCheckpoint_Delete(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewCheckpointStore(dir)

	cp := Checkpoint{Path: "secret/app", Namespace: "", SyncedAt: time.Now(), KeyCount: 2}
	_ = store.Save(cp)

	if err := store.Delete("secret/app", ""); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	loaded, err := store.Load("secret/app", "")
	if err != nil {
		t.Fatalf("Load after delete: %v", err)
	}
	if loaded != nil {
		t.Fatal("expected nil after delete")
	}
}

func TestCheckpoint_Delete_NonExistent(t *testing.T) {
	store, _ := NewCheckpointStore(t.TempDir())
	if err := store.Delete("secret/nope", "staging"); err != nil {
		t.Fatalf("expected no error deleting non-existent checkpoint, got %v", err)
	}
}

func TestCheckpoint_NoNamespace_FileName(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewCheckpointStore(dir)
	cp := Checkpoint{Path: "secret/svc", Namespace: "", SyncedAt: time.Now(), KeyCount: 1}
	_ = store.Save(cp)

	loaded, err := store.Load("secret/svc", "")
	if err != nil || loaded == nil {
		t.Fatalf("Load without namespace failed: %v", err)
	}
}
