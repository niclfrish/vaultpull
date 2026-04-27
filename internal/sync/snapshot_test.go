package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSnapshotStore_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snapshots")
	_, err := NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSnapshot_SaveAndLoad(t *testing.T) {
	store, _ := NewSnapshotStore(t.TempDir())
	snap := Snapshot{
		Path:      "secret/app",
		Namespace: "prod",
		Secrets:   map[string]string{"KEY": "value", "DB_PASS": "s3cr3t"},
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := store.Load("secret/app", "prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if loaded.Secrets["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Secrets["KEY"])
	}
	if loaded.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSnapshot_Load_NonExistent(t *testing.T) {
	store, _ := NewSnapshotStore(t.TempDir())
	snap, err := store.Load("missing/path", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil for missing snapshot")
	}
}

func TestSnapshotFileName_NoNamespace(t *testing.T) {
	name := snapshotFileName("secret/app", "")
	expected := "secret_app.snapshot.json"
	if name != expected {
		t.Errorf("expected %q, got %q", expected, name)
	}
}

func TestSnapshotFileName_WithNamespace(t *testing.T) {
	name := snapshotFileName("secret/app", "prod")
	expected := "prod_secret_app.snapshot.json"
	if name != expected {
		t.Errorf("expected %q, got %q", expected, name)
	}
}

func TestSnapshot_Save_OverwritesPrevious(t *testing.T) {
	store, _ := NewSnapshotStore(t.TempDir())
	first := Snapshot{Path: "secret/app", Secrets: map[string]string{"A": "1"}}
	second := Snapshot{Path: "secret/app", Secrets: map[string]string{"A": "2", "B": "3"}}
	_ = store.Save(first)
	_ = store.Save(second)
	loaded, _ := store.Load("secret/app", "")
	if loaded.Secrets["A"] != "2" {
		t.Errorf("expected overwritten value A=2, got %q", loaded.Secrets["A"])
	}
	if len(loaded.Secrets) != 2 {
		t.Errorf("expected 2 keys, got %d", len(loaded.Secrets))
	}
}

func TestSnapshot_SaveAndLoad_PreservesNamespace(t *testing.T) {
	store, _ := NewSnapshotStore(t.TempDir())
	snap := Snapshot{
		Path:      "secret/app",
		Namespace: "staging",
		Secrets:   map[string]string{"ENV": "staging"},
	}
	_ = store.Save(snap)
	loaded, err := store.Load("secret/app", "staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if loaded.Namespace != "staging" {
		t.Errorf("expected namespace %q, got %q", "staging", loaded.Namespace)
	}
}
