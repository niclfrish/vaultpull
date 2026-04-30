package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestLookupAndReport_NilSecrets(t *testing.T) {
	_, err := LookupAndReport(nil, []string{"KEY"}, DefaultLookupConfig(), nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestLookupAndReport_WritesOutput(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "token"}
	var buf bytes.Buffer
	results, err := LookupAndReport(secrets, []string{"db_host"}, DefaultLookupConfig(), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(buf.String(), "DB_HOST") {
		t.Errorf("expected output to contain DB_HOST, got: %s", buf.String())
	}
}

func TestLookupAndReport_NilWriter_UsesStdout(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	// Should not panic with nil writer
	_, err := LookupAndReport(secrets, []string{"foo"}, DefaultLookupConfig(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLookupAndReport_SortedOutput(t *testing.T) {
	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	cfg := DefaultLookupConfig()
	cfg.PartialMatch = true
	var buf bytes.Buffer
	_, err := LookupAndReport(secrets, []string{"_key"}, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	aIdx := strings.Index(out, "A_KEY")
	mIdx := strings.Index(out, "M_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Errorf("expected sorted output A < M < Z, got:\n%s", out)
	}
}

func TestLookupStage_StageName(t *testing.T) {
	stage := LookupStage([]string{"FOO"}, DefaultLookupConfig())
	if stage.Name != "lookup" {
		t.Errorf("expected stage name 'lookup', got %s", stage.Name)
	}
}

func TestLookupStage_InjectsAnnotation(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	stage := LookupStage([]string{"FOO"}, DefaultLookupConfig())
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["__lookup_FOO"]; !ok || v != "bar" {
		t.Errorf("expected __lookup_FOO=bar, got %q", v)
	}
}

func TestLookupStage_MissingKeyInjectsEmpty(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	stage := LookupStage([]string{"MISSING"}, DefaultLookupConfig())
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["__lookup_MISSING"]; !ok || v != "" {
		t.Errorf("expected __lookup_MISSING='', got %q (ok=%v)", v, ok)
	}
}
