package sync

import (
	"strings"
	"testing"
)

func TestJoinSecrets_NilSecrets(t *testing.T) {
	_, _, err := JoinSecrets(nil, DefaultJoinConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestJoinSecrets_EmptyOutputKey(t *testing.T) {
	cfg := DefaultJoinConfig()
	cfg.Keys = []string{"A"}
	cfg.OutputKey = ""
	_, _, err := JoinSecrets(map[string]string{"A": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty OutputKey")
	}
}

func TestJoinSecrets_EmptyKeys(t *testing.T) {
	cfg := DefaultJoinConfig()
	cfg.OutputKey = "JOINED"
	_, _, err := JoinSecrets(map[string]string{"A": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty Keys slice")
	}
}

func TestJoinSecrets_BasicJoin(t *testing.T) {
	secrets := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	cfg := DefaultJoinConfig()
	cfg.Keys = []string{"HOST", "PORT"}
	cfg.OutputKey = "DSN"
	cfg.Separator = ":"

	out, summary, err := JoinSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "localhost:5432" {
		t.Errorf("expected 'localhost:5432', got %q", out["DSN"])
	}
	if summary.Joined != 2 {
		t.Errorf("expected Joined=2, got %d", summary.Joined)
	}
	if summary.Skipped != 0 {
		t.Errorf("expected Skipped=0, got %d", summary.Skipped)
	}
}

func TestJoinSecrets_SkipsMissingKeys(t *testing.T) {
	secrets := map[string]string{"HOST": "db"}
	cfg := DefaultJoinConfig()
	cfg.Keys = []string{"HOST", "PORT"}
	cfg.OutputKey = "ADDR"
	cfg.Separator = ":"

	out, summary, err := JoinSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ADDR"] != "db" {
		t.Errorf("expected 'db', got %q", out["ADDR"])
	}
	if summary.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", summary.Skipped)
	}
}

func TestJoinSecrets_StripParts(t *testing.T) {
	secrets := map[string]string{
		"USER": "admin",
		"PASS": "secret",
		"OTHER": "keep",
	}
	cfg := DefaultJoinConfig()
	cfg.Keys = []string{"USER", "PASS"}
	cfg.OutputKey = "CREDENTIALS"
	cfg.Separator = ":"
	cfg.StripParts = true

	out, summary, err := JoinSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["USER"]; ok {
		t.Error("expected USER to be stripped")
	}
	if _, ok := out["PASS"]; ok {
		t.Error("expected PASS to be stripped")
	}
	if out["OTHER"] != "keep" {
		t.Errorf("expected OTHER to be preserved")
	}
	if summary.Stripped != 2 {
		t.Errorf("expected Stripped=2, got %d", summary.Stripped)
	}
}

func TestJoinSummaryString_ContainsFields(t *testing.T) {
	s := JoinSummary{Joined: 3, Skipped: 1, Stripped: 2}
	result := JoinSummaryString(s)
	for _, want := range []string{"joined=3", "skipped=1", "stripped=2"} {
		if !strings.Contains(result, want) {
			t.Errorf("summary string %q missing %q", result, want)
		}
	}
}
