package sync

import (
	"bytes"
	"strings"
	"testing"
)

func newTestCheckpointStore(t *testing.T) *CheckpointStore {
	t.Helper()
	store, err := NewCheckpointStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewCheckpointStore: %v", err)
	}
	return store
}

func TestCheckpointAfterSync_SavesCheckpoint(t *testing.T) {
	store := newTestCheckpointStore(t)
	var buf bytes.Buffer
	secrets := map[string]string{"KEY_A": "val1", "KEY_B": "val2"}

	err := CheckpointAfterSync(store, "secret/app", "staging", secrets, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[checkpoint] saved") {
		t.Errorf("expected log output, got: %q", buf.String())
	}

	cp, err := store.Load("secret/app", "staging")
	if err != nil || cp == nil {
		t.Fatalf("expected checkpoint to be saved: %v", err)
	}
	if cp.KeyCount != 2 {
		t.Errorf("KeyCount: want 2, got %d", cp.KeyCount)
	}
}

func TestCheckpointAfterSync_NilStore_NoOp(t *testing.T) {
	err := CheckpointAfterSync(nil, "secret/app", "", map[string]string{"A": "1"}, nil)
	if err != nil {
		t.Fatalf("expected no error with nil store, got %v", err)
	}
}

func TestLoadCheckpoint_NoExisting(t *testing.T) {
	store := newTestCheckpointStore(t)
	var buf bytes.Buffer
	cp, err := LoadCheckpoint(store, "secret/new", "", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp != nil {
		t.Fatal("expected nil checkpoint")
	}
	if !strings.Contains(buf.String(), "no previous checkpoint") {
		t.Errorf("expected 'no previous checkpoint' in output, got: %q", buf.String())
	}
}

func TestLoadCheckpoint_ReturnsExisting(t *testing.T) {
	store := newTestCheckpointStore(t)
	secrets := map[string]string{"X": "y"}
	_ = CheckpointAfterSync(store, "secret/svc", "", secrets, nil)

	var buf bytes.Buffer
	cp, err := LoadCheckpoint(store, "secret/svc", "", &buf)
	if err != nil || cp == nil {
		t.Fatalf("expected checkpoint: %v", err)
	}
	if !strings.Contains(buf.String(), "last sync") {
		t.Errorf("expected 'last sync' in output, got: %q", buf.String())
	}
}

func TestSecretsChangedSinceCheckpoint_NilCheckpoint(t *testing.T) {
	if !SecretsChangedSinceCheckpoint(nil, map[string]string{"A": "1"}) {
		t.Error("expected true when checkpoint is nil")
	}
}

func TestSecretsChangedSinceCheckpoint_Unchanged(t *testing.T) {
	store := newTestCheckpointStore(t)
	secrets := map[string]string{"K": "v", "K2": "v2"}
	_ = CheckpointAfterSync(store, "secret/app", "", secrets, nil)
	cp, _ := store.Load("secret/app", "")

	if SecretsChangedSinceCheckpoint(cp, secrets) {
		t.Error("expected no change detected")
	}
}

func TestSecretsChangedSinceCheckpoint_Changed(t *testing.T) {
	store := newTestCheckpointStore(t)
	original := map[string]string{"K": "v1"}
	_ = CheckpointAfterSync(store, "secret/app", "", original, nil)
	cp, _ := store.Load("secret/app", "")

	updated := map[string]string{"K": "v2"}
	if !SecretsChangedSinceCheckpoint(cp, updated) {
		t.Error("expected change to be detected")
	}
}
