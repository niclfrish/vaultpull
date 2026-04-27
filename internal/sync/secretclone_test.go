package sync

import (
	"strings"
	"testing"
)

func TestCloneSecrets_NilSecrets(t *testing.T) {
	_, _, err := CloneSecrets(nil, DefaultCloneConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestCloneSecrets_EmptyDestPrefix(t *testing.T) {
	cfg := DefaultCloneConfig()
	cfg.DestPrefix = ""
	_, _, err := CloneSecrets(map[string]string{"KEY": "val"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty DestPrefix")
	}
}

func TestCloneSecrets_ClonesAllKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "1", "BAR": "2"}
	cfg := DefaultCloneConfig()
	cfg.DestPrefix = "COPY"

	out, res, err := CloneSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 2 {
		t.Errorf("expected 2 cloned, got %d", res.Cloned)
	}
	if out["COPY_FOO"] != "1" || out["COPY_BAR"] != "2" {
		t.Errorf("cloned keys missing or wrong values: %v", out)
	}
	// originals preserved
	if out["FOO"] != "1" || out["BAR"] != "2" {
		t.Error("original keys should be preserved")
	}
}

func TestCloneSecrets_SourcePrefixFilters(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "APP_NAME": "myapp"}
	cfg := DefaultCloneConfig()
	cfg.SourcePrefix = "DB_"
	cfg.DestPrefix = "BACKUP"

	out, res, err := CloneSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 1 {
		t.Errorf("expected 1 cloned, got %d", res.Cloned)
	}
	if _, ok := out["BACKUP_DB_HOST"]; !ok {
		t.Error("expected BACKUP_DB_HOST in output")
	}
	if _, ok := out["BACKUP_APP_NAME"]; ok {
		t.Error("APP_NAME should not be cloned")
	}
}

func TestCloneSecrets_SkipsExistingWithoutOverwrite(t *testing.T) {
	secrets := map[string]string{"KEY": "original", "CLONE_KEY": "existing"}
	cfg := DefaultCloneConfig()
	cfg.DestPrefix = "CLONE"
	cfg.Overwrite = false

	out, res, err := CloneSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
	if out["CLONE_KEY"] != "existing" {
		t.Error("existing key should not be overwritten")
	}
}

func TestCloneSecrets_OverwritesWhenEnabled(t *testing.T) {
	secrets := map[string]string{"KEY": "new", "CLONE_KEY": "old"}
	cfg := DefaultCloneConfig()
	cfg.DestPrefix = "CLONE"
	cfg.Overwrite = true

	out, res, err := CloneSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Overwritten != 1 {
		t.Errorf("expected 1 overwritten, got %d", res.Overwritten)
	}
	if out["CLONE_KEY"] != "new" {
		t.Errorf("expected overwritten value 'new', got %q", out["CLONE_KEY"])
	}
}

func TestCloneSummary(t *testing.T) {
	r := CloneResult{Cloned: 3, Skipped: 1, Overwritten: 2}
	s := CloneSummary(r)
	if !strings.Contains(s, "cloned=3") || !strings.Contains(s, "skipped=1") || !strings.Contains(s, "overwritten=2") {
		t.Errorf("unexpected summary: %s", s)
	}
}
