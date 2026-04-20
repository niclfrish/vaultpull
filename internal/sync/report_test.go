package sync

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReport_Print_NoChanges(t *testing.T) {
	r := &Report{
		EnvFile:    ".env",
		SecretPath: "secret/app",
		Duration:   50 * time.Millisecond,
		Plan:       Plan{Changes: []Change{}},
	}

	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "Sync Report") {
		t.Error("expected 'Sync Report' header")
	}
	if !strings.Contains(out, ".env") {
		t.Error("expected env file name in output")
	}
	if !strings.Contains(out, "Added:") {
		t.Error("expected 'Added:' in output")
	}
}

func TestReport_Print_WithChanges(t *testing.T) {
	plan := Plan{
		Changes: []Change{
			{Key: "FOO", Action: ActionAdd, NewValue: "bar"},
			{Key: "OLD", Action: ActionRemove, OldValue: "x"},
			{Key: "DB", Action: ActionUpdate, OldValue: "a", NewValue: "b"},
			{Key: "SAME", Action: ActionUnchanged, OldValue: "v", NewValue: "v"},
		},
	}
	r := &Report{
		EnvFile:    ".env",
		SecretPath: "secret/app",
		Duration:   120 * time.Millisecond,
		Plan:       plan,
	}

	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "[ADD]") {
		t.Error("expected [ADD] in detail section")
	}
	if !strings.Contains(out, "[REMOVE]") {
		t.Error("expected [REMOVE] in detail section")
	}
	if !strings.Contains(out, "[UPDATE]") {
		t.Error("expected [UPDATE] in detail section")
	}
	if strings.Contains(out, "SAME") {
		t.Error("unchanged keys should not appear in detail")
	}
}

func TestReport_Print_WithNamespace(t *testing.T) {
	r := &Report{
		EnvFile:    ".env",
		SecretPath: "secret/app",
		Namespace:  "staging",
		Duration:   10 * time.Millisecond,
		Plan:       Plan{Changes: []Change{}},
	}

	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "staging") {
		t.Error("expected namespace in output")
	}
}

func TestReport_Print_DryRun(t *testing.T) {
	r := &Report{
		EnvFile:    ".env",
		SecretPath: "secret/app",
		DryRun:     true,
		Duration:   5 * time.Millisecond,
		Plan:       Plan{Changes: []Change{}},
	}

	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "true") {
		t.Error("expected dry_run=true in output")
	}
}
