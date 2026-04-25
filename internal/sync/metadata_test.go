package sync

import (
	"strings"
	"testing"
	"time"
)

var fixedTime = time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

func TestInjectMetadata_NilSecrets(t *testing.T) {
	_, err := InjectMetadata(nil, DefaultMetadataConfig(), fixedTime)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestInjectMetadata_Timestamp(t *testing.T) {
	cfg := DefaultMetadataConfig()
	cfg.IncludeCount = false
	cfg.IncludeKeys = false

	out, err := InjectMetadata(map[string]string{"FOO": "bar"}, cfg, fixedTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out["VAULTPULL_META_SYNCED_AT"]
	want := "2024-06-15T12:00:00Z"
	if got != want {
		t.Errorf("timestamp: got %q, want %q", got, want)
	}
}

func TestInjectMetadata_Count(t *testing.T) {
	cfg := DefaultMetadataConfig()
	cfg.IncludeTimestamp = false
	cfg.IncludeKeys = false

	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	out, err := InjectMetadata(secrets, cfg, fixedTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAULTPULL_META_SECRET_COUNT"] != "3" {
		t.Errorf("count: got %q, want \"3\"", out["VAULTPULL_META_SECRET_COUNT"])
	}
}

func TestInjectMetadata_Keys(t *testing.T) {
	cfg := DefaultMetadataConfig()
	cfg.IncludeTimestamp = false
	cfg.IncludeCount = false
	cfg.IncludeKeys = true

	secrets := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	out, err := InjectMetadata(secrets, cfg, fixedTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out["VAULTPULL_META_KEYS"]
	if !strings.Contains(got, "ALPHA") || !strings.Contains(got, "ZEBRA") {
		t.Errorf("keys: got %q, expected sorted key list", got)
	}
	// Verify sorted order
	parts := strings.Split(got, ",")
	if parts[0] != "ALPHA" {
		t.Errorf("keys not sorted: first element = %q, want ALPHA", parts[0])
	}
}

func TestInjectMetadata_DoesNotOverwriteExisting(t *testing.T) {
	cfg := DefaultMetadataConfig()
	cfg.IncludeKeys = false

	existing := "already-set"
	secrets := map[string]string{
		"VAULTPULL_META_SYNCED_AT":    existing,
		"VAULTPULL_META_SECRET_COUNT": existing,
	}
	out, err := InjectMetadata(secrets, cfg, fixedTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAULTPULL_META_SYNCED_AT"] != existing {
		t.Errorf("should not overwrite SYNCED_AT")
	}
	if out["VAULTPULL_META_SECRET_COUNT"] != existing {
		t.Errorf("should not overwrite SECRET_COUNT")
	}
}

func TestInjectMetadata_CustomPrefix(t *testing.T) {
	cfg := DefaultMetadataConfig()
	cfg.KeyPrefix = "META_"
	cfg.IncludeKeys = false

	out, err := InjectMetadata(map[string]string{"X": "y"}, cfg, fixedTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["META_SYNCED_AT"]; !ok {
		t.Error("expected META_SYNCED_AT with custom prefix")
	}
}

func TestInjectMetadata_EmptySecrets_CountZero(t *testing.T) {
	cfg := DefaultMetadataConfig()
	cfg.IncludeTimestamp = false
	cfg.IncludeKeys = false

	out, err := InjectMetadata(map[string]string{}, cfg, fixedTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAULTPULL_META_SECRET_COUNT"] != "0" {
		t.Errorf("empty secrets count: got %q, want \"0\"", out["VAULTPULL_META_SECRET_COUNT"])
	}
}
