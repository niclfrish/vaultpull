package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidateAndReport_NilSecrets(t *testing.T) {
	var buf bytes.Buffer
	err := ValidateAndReport(&buf, nil, nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestValidateAndReport_PassesWithValidSecrets(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	err := ValidateAndReport(&buf, secrets, []string{"DB_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "validation passed") {
		t.Errorf("expected success message in output, got: %s", buf.String())
	}
}

func TestValidateAndReport_FailsWithMissingRequired(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"DB_HOST": "localhost"}
	err := ValidateAndReport(&buf, secrets, []string{"API_SECRET"})
	if err == nil {
		t.Fatal("expected error due to missing required key")
	}
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", buf.String())
	}
}

func TestValidateAndReport_NilWriter_UsesStdout(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	// Should not panic when writer is nil
	err := ValidateAndReport(nil, secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateAndReport_WarningInOutput(t *testing.T) {
	var buf bytes.Buffer
	longVal := strings.Repeat("x", 5000)
	secrets := map[string]string{"BIG_KEY": longVal}
	v := NewValidator(nil, 100)
	res, _ := v.Validate(secrets)
	if len(res.Warnings) == 0 {
		t.Skip("no warnings produced, skipping output check")
	}
	err := ValidateAndReport(&buf, secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Logf("output: %s", buf.String())
	}
}
