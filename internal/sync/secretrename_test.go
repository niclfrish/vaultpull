package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenameSecrets_NilSecrets(t *testing.T) {
	_, _, err := RenameSecrets(nil, DefaultRenameConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestRenameSecrets_NoRules_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, summary, err := RenameSecrets(secrets, DefaultRenameConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 || out["FOO"] != "bar" {
		t.Errorf("expected unchanged map, got %v", out)
	}
	if summary.Renamed != 0 || summary.Missed != 0 {
		t.Errorf("expected zero counts, got %+v", summary)
	}
}

func TestRenameSecrets_RenamesKey(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "value"}
	cfg := RenameConfig{Rules: map[string]string{"OLD_KEY": "NEW_KEY"}, CaseSensitive: true}
	out, summary, err := RenameSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("old key should be removed")
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected new key to hold value, got %v", out)
	}
	if summary.Renamed != 1 {
		t.Errorf("expected 1 renamed, got %d", summary.Renamed)
	}
}

func TestRenameSecrets_MissingKeyNonStrict(t *testing.T) {
	secrets := map[string]string{"PRESENT": "v"}
	cfg := RenameConfig{Rules: map[string]string{"ABSENT": "NEW"}, CaseSensitive: true}
	out, summary, err := RenameSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["PRESENT"]; !ok {
		t.Error("unrelated key should remain")
	}
	if summary.Missed != 1 {
		t.Errorf("expected 1 missed, got %d", summary.Missed)
	}
}

func TestRenameSecrets_MissingKeyStrict(t *testing.T) {
	secrets := map[string]string{"PRESENT": "v"}
	cfg := RenameConfig{Rules: map[string]string{"ABSENT": "NEW"}, CaseSensitive: true, FailOnMissing: true}
	_, _, err := RenameSecrets(secrets, cfg)
	if err == nil {
		t.Fatal("expected error for missing key in strict mode")
	}
}

func TestRenameSecrets_CaseInsensitive(t *testing.T) {
	secrets := map[string]string{"old_key": "val"}
	cfg := RenameConfig{Rules: map[string]string{"OLD_KEY": "NEW_KEY"}, CaseSensitive: false}
	out, summary, err := RenameSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "val" {
		t.Errorf("expected case-insensitive rename, got %v", out)
	}
	if summary.Renamed != 1 {
		t.Errorf("expected 1 renamed, got %d", summary.Renamed)
	}
}

func TestRenameAndReport_WritesOutput(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	cfg := RenameConfig{Rules: map[string]string{"A": "ALPHA"}, CaseSensitive: true}
	var buf bytes.Buffer
	out, err := RenameAndReport(secrets, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ALPHA"] != "1" {
		t.Errorf("expected ALPHA=1, got %v", out)
	}
	if !strings.Contains(buf.String(), "1 key(s) renamed") {
		t.Errorf("expected summary in output, got: %s", buf.String())
	}
}

func TestRenameAndReport_NilSecrets(t *testing.T) {
	var buf bytes.Buffer
	_, err := RenameAndReport(nil, DefaultRenameConfig(), &buf)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}
