package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintPlan_NoChanges(t *testing.T) {
	plan := &Plan{
		Entries: []PlanEntry{{Key: "A", Change: ChangeNone}},
	}
	var buf bytes.Buffer
	PrintPlan(&buf, plan)
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestPrintPlan_WithChanges(t *testing.T) {
	plan := &Plan{
		Entries: []PlanEntry{
			{Key: "DB_PASS", Change: ChangeAdd, New: "secret123"},
			{Key: "OLD_KEY", Change: ChangeRemove, Old: "val"},
			{Key: "API_KEY", Change: ChangeUpdate, Old: "abc", New: "xyz"},
		},
	}
	var buf bytes.Buffer
	PrintPlan(&buf, plan)
	out := buf.String()
	if !strings.Contains(out, "+ DB_PASS") {
		t.Errorf("expected add line, got: %s", out)
	}
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected remove line, got: %s", out)
	}
	if !strings.Contains(out, "~ API_KEY") {
		t.Errorf("expected update line, got: %s", out)
	}
	if !strings.Contains(out, "Plan:") {
		t.Errorf("expected summary line, got: %s", out)
	}
}

func TestMaskValue(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"secret", "se****"},
		{"ab", "**"},
		{"a", "*"},
		{"", ""},
	}
	for _, c := range cases {
		got := maskValue(c.input)
		if got != c.expected {
			t.Errorf("maskValue(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}
