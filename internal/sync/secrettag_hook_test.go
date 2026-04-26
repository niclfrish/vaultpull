package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestTagAndReport_NilSecrets(t *testing.T) {
	_, err := TagAndReport(nil, DefaultSecretTagConfig(), nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestTagAndReport_WritesInjectedCount(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultSecretTagConfig()
	cfg.Timestamp = false
	input := map[string]string{"KEY": "val"}
	out, err := TagAndReport(input, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) <= len(input) {
		t.Error("expected tagged output to have more keys than input")
	}
	if !strings.Contains(buf.String(), "injected") {
		t.Errorf("expected report to mention injected keys, got: %q", buf.String())
	}
}

func TestTagAndReport_NilWriter_UsesStdout(t *testing.T) {
	cfg := DefaultSecretTagConfig()
	cfg.Timestamp = false
	_, err := TagAndReport(map[string]string{"A": "1"}, cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagStage_RunsSuccessfully(t *testing.T) {
	cfg := DefaultSecretTagConfig()
	cfg.Timestamp = false
	stage := TagStage(cfg)
	if stage.Name != "tag" {
		t.Errorf("expected stage name 'tag', got %q", stage.Name)
	}
	input := map[string]string{"X": "y"}
	out, err := stage.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["__meta_source"]; !ok {
		t.Error("expected __meta_source key in output")
	}
}

func TestStripTagStage_RemovesTagKeys(t *testing.T) {
	stage := StripTagStage("__meta")
	if stage.Name != "strip-tags" {
		t.Errorf("expected stage name 'strip-tags', got %q", stage.Name)
	}
	input := map[string]string{
		"APP": "value",
		"__meta_source": "vault",
	}
	out, err := stage.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["__meta_source"]; ok {
		t.Error("expected __meta_source to be stripped")
	}
	if out["APP"] != "value" {
		t.Error("expected APP key to be preserved")
	}
}
