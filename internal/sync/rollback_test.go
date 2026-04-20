package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateBackup_NoOriginalFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	bak, err := CreateBackup(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bak != nil {
		t.Fatalf("expected nil backup when file does not exist, got %+v", bak)
	}
}

func TestCreateBackup_CreatesBackupFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	origContent := []byte("KEY=value\n")
	if err := os.WriteFile(envPath, origContent, 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := CreateBackup(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bak == nil {
		t.Fatal("expected backup, got nil")
	}
	if !strings.HasSuffix(bak.BackupPath, ".bak") {
		t.Errorf("backup path should end with .bak, got %s", bak.BackupPath)
	}

	data, err := os.ReadFile(bak.BackupPath)
	if err != nil {
		t.Fatalf("cannot read backup file: %v", err)
	}
	if string(data) != string(origContent) {
		t.Errorf("backup content mismatch: got %q, want %q", data, origContent)
	}

	_ = bak.Discard()
}

func TestBackup_Restore(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	origContent := []byte("KEY=original\n")
	if err := os.WriteFile(envPath, origContent, 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := CreateBackup(envPath)
	if err != nil || bak == nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Overwrite the original to simulate a failed apply
	if err := os.WriteFile(envPath, []byte("KEY=changed\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := bak.Restore(); err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(origContent) {
		t.Errorf("restored content mismatch: got %q, want %q", data, origContent)
	}

	// Backup file should be removed after restore
	if _, err := os.Stat(bak.BackupPath); !os.IsNotExist(err) {
		t.Error("backup file should be removed after restore")
	}
}

func TestBackup_Discard(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	if err := os.WriteFile(envPath, []byte("X=1\n"), 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := CreateBackup(envPath)
	if err != nil || bak == nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	if err := bak.Discard(); err != nil {
		t.Fatalf("Discard failed: %v", err)
	}

	if _, err := os.Stat(bak.BackupPath); !os.IsNotExist(err) {
		t.Error("backup file should not exist after discard")
	}
}
