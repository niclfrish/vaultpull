package sync

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewAuditLogger_LogCreatesFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	logger := NewAuditLogger(logPath)
	entry := AuditEntry{
		Timestamp:  time.Now().UTC(),
		EnvFile:    ".env",
		SecretPath: "secret/app",
		Changes:    map[string]string{"FOO": "add"},
		DryRun:     false,
	}

	if err := logger.Log(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("expected audit log file to be created")
	}
}

func TestAuditLogger_AppendMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	logger := NewAuditLogger(logPath)

	for i := 0; i < 3; i++ {
		err := logger.Log(AuditEntry{
			Timestamp:  time.Now().UTC(),
			EnvFile:    ".env",
			SecretPath: "secret/app",
			Changes:    map[string]string{},
		})
		if err != nil {
			t.Fatalf("log entry %d failed: %v", i, err)
		}
	}

	f, _ := os.Open(logPath)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 log lines, got %d", count)
	}
}

func TestBuildAuditEntry(t *testing.T) {
	plan := Plan{
		Changes: []Change{
			{Key: "FOO", Action: ActionAdd},
			{Key: "BAR", Action: ActionRemove},
		},
	}

	entry := BuildAuditEntry(plan, ".env", "secret/app", "prod", true)

	if entry.EnvFile != ".env" {
		t.Errorf("expected .env, got %s", entry.EnvFile)
	}
	if entry.Namespace != "prod" {
		t.Errorf("expected prod namespace")
	}
	if !entry.DryRun {
		t.Error("expected dry_run to be true")
	}
	if entry.Changes["FOO"] != "add" {
		t.Errorf("expected FOO=add, got %s", entry.Changes["FOO"])
	}
	if entry.Changes["BAR"] != "remove" {
		t.Errorf("expected BAR=remove, got %s", entry.Changes["BAR"])
	}

	// Ensure JSON serialisable
	if _, err := json.Marshal(entry); err != nil {
		t.Errorf("entry not JSON serialisable: %v", err)
	}
}

func TestAuditLogger_InvalidPath(t *testing.T) {
	logger := NewAuditLogger("/nonexistent/dir/audit.log")
	err := logger.Log(AuditEntry{})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
