package sync

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func rotatedAtKey(key string) string {
	return DefaultRotateConfig().RotatedAtKey + "." + key
}

func TestCheckRotation_NoTimestamp_Unknown(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret"}
	statuses := CheckRotation(secrets, DefaultRotateConfig())
	if statuses["DB_PASS"] != RotationUnknown {
		t.Fatalf("expected unknown, got %v", statuses["DB_PASS"])
	}
}

func TestCheckRotation_RecentTimestamp_OK(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	secrets := map[string]string{
		"API_KEY":             "val",
		rotatedAtKey("API_KEY"): recent,
	}
	statuses := CheckRotation(secrets, DefaultRotateConfig())
	if statuses["API_KEY"] != RotationOK {
		t.Fatalf("expected ok, got %v", statuses["API_KEY"])
	}
}

func TestCheckRotation_OldTimestamp_Stale(t *testing.T) {
	old := time.Now().Add(-48 * time.Hour).Format(time.RFC3339)
	secrets := map[string]string{
		"TOKEN":             "val",
		rotatedAtKey("TOKEN"): old,
	}
	statuses := CheckRotation(secrets, DefaultRotateConfig())
	if statuses["TOKEN"] != RotationStale {
		t.Fatalf("expected stale, got %v", statuses["TOKEN"])
	}
}

func TestRotateSummary_Counts(t *testing.T) {
	statuses := map[string]RotationStatus{
		"A": RotationOK,
		"B": RotationStale,
		"C": RotationUnknown,
	}
	summary := RotateSummary(statuses)
	for _, want := range []string{"ok=1", "stale=1", "unknown=1"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary %q missing %q", summary, want)
		}
	}
}

func TestCheckRotationAndReport_FailOnStale(t *testing.T) {
	old := time.Now().Add(-48 * time.Hour).Format(time.RFC3339)
	secrets := map[string]string{
		"SECRET":             "val",
		rotatedAtKey("SECRET"): old,
	}
	var buf bytes.Buffer
	err := CheckRotationAndReport(secrets, DefaultRotateConfig(), true, &buf)
	if err == nil {
		t.Fatal("expected error for stale secret")
	}
	if !strings.Contains(buf.String(), "stale") {
		t.Errorf("expected stale in output, got: %s", buf.String())
	}
}

func TestCheckRotationAndReport_NilSecrets_NoOp(t *testing.T) {
	var buf bytes.Buffer
	if err := CheckRotationAndReport(nil, DefaultRotateConfig(), true, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRotationStage_InjectsStatusKeys(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	secrets := map[string]string{
		"DB_PASS":               "val",
		rotatedAtKey("DB_PASS"): recent,
	}
	stage := RotationStage(DefaultRotateConfig())
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	statusKey := "__rotation_status.DB_PASS"
	if out[statusKey] != "ok" {
		t.Errorf("expected ok, got %q", out[statusKey])
	}
}
