package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestPinStage_StageName(t *testing.T) {
	stage := PinStage(DefaultPinConfig())
	if stage.Name != "pin" {
		t.Errorf("expected stage name 'pin', got %q", stage.Name)
	}
}

func TestPinStage_PassesThroughNoPins(t *testing.T) {
	stage := PinStage(DefaultPinConfig())
	secrets := map[string]string{"FOO": "bar"}
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
}

func TestPinStage_PinsKey(t *testing.T) {
	cfg := DefaultPinConfig()
	cfg.Pins = map[string]string{"API_KEY": "v5"}
	stage := PinStage(cfg)
	secrets := map[string]string{"API_KEY": "mytoken"}
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out[cfg.AnnotationKey], "API_KEY@v5") {
		t.Errorf("expected annotation, got %q", out[cfg.AnnotationKey])
	}
}

func TestPinAndReport_NilSecrets(t *testing.T) {
	_, err := PinAndReport(nil, DefaultPinConfig(), nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestPinAndReport_NilWriter_UsesStdout(t *testing.T) {
	cfg := DefaultPinConfig()
	cfg.Pins = map[string]string{"X": "v1"}
	secrets := map[string]string{"X": "value"}
	// Should not panic when w is nil.
	_, err := PinAndReport(secrets, cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPinAndReport_WritesOutput(t *testing.T) {
	cfg := DefaultPinConfig()
	cfg.Pins = map[string]string{"SECRET": "v2", "MISSING": "v1"}
	secrets := map[string]string{"SECRET": "s3cr3t"}

	var buf bytes.Buffer
	out, err := PinAndReport(secrets, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected non-nil output")
	}
	output := buf.String()
	if !strings.Contains(output, "1 pinned") {
		t.Errorf("expected '1 pinned' in output, got: %s", output)
	}
	if !strings.Contains(output, "1 missing") {
		t.Errorf("expected '1 missing' in output, got: %s", output)
	}
	if !strings.Contains(output, "[pinned]") {
		t.Errorf("expected '[pinned]' in output, got: %s", output)
	}
	if !strings.Contains(output, "[missing]") {
		t.Errorf("expected '[missing]' in output, got: %s", output)
	}
}
