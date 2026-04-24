package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidateSchema_NilSchema(t *testing.T) {
	result, err := ValidateSchema(nil, map[string]string{"KEY": "val"}, nil)
	if err != nil {
		t.Fatalf("expected no error for nil schema, got %v", err)
	}
	if !result.Passed {
		t.Error("expected Passed=true for nil schema")
	}
}

func TestValidateSchema_PassesWithValidSecrets(t *testing.T) {
	schema, _ := NewSchema([]SchemaRule{
		{Key: "PORT", Pattern: `^\d+$`, Required: true},
	})
	var buf bytes.Buffer
	result, err := ValidateSchema(schema, map[string]string{"PORT": "9090"}, &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Passed {
		t.Error("expected Passed=true")
	}
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected 'passed' in output, got: %q", buf.String())
	}
}

func TestValidateSchema_FailsWithViolations(t *testing.T) {
	schema, _ := NewSchema([]SchemaRule{
		{Key: "HOST", Required: true},
	})
	var buf bytes.Buffer
	result, err := ValidateSchema(schema, map[string]string{}, &buf)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	if result.Passed {
		t.Error("expected Passed=false")
	}
	if len(result.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result.Violations))
	}
	if !strings.Contains(buf.String(), "HOST") {
		t.Errorf("expected key name in output, got: %q", buf.String())
	}
}

func TestValidateSchema_NilWriter_UsesStdout(t *testing.T) {
	schema, _ := NewSchema([]SchemaRule{
		{Key: "X", Required: false},
	})
	// Should not panic when writer is nil
	_, err := ValidateSchema(schema, map[string]string{"X": "ok"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSchema_MultipleViolationsReported(t *testing.T) {
	schema, _ := NewSchema([]SchemaRule{
		{Key: "A", Required: true},
		{Key: "B", Required: true},
	})
	var buf bytes.Buffer
	result, err := ValidateSchema(schema, map[string]string{}, &buf)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(result.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(result.Violations))
	}
	output := buf.String()
	if !strings.Contains(output, "[A]") || !strings.Contains(output, "[B]") {
		t.Errorf("expected both keys in output, got: %q", output)
	}
}
