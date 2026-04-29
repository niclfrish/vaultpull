package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromoteStage_PassesThroughEmpty(t *testing.T) {
	stage := PromoteStage(PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_"})
	out, err := stage.Fn(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}

func TestPromoteStage_PromotesKeys(t *testing.T) {
	cfg := PromoteConfig{FromPrefix: "staging_", ToPrefix: "prod_", Overwrite: true}
	stage := PromoteStage(cfg)
	in := map[string]string{"staging_TOKEN": "abc123"}
	out, err := stage.Fn(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["prod_TOKEN"] != "abc123" {
		t.Errorf("expected prod_TOKEN=abc123, got %q", out["prod_TOKEN"])
	}
}

func TestPromoteStage_StageName(t *testing.T) {
	stage := PromoteStage(DefaultPromoteConfig())
	if stage.Name != "promote" {
		t.Errorf("expected stage name 'promote', got %q", stage.Name)
	}
}

func TestPromoteAndReport_NilSecrets(t *testing.T) {
	var buf bytes.Buffer
	out, err := PromoteAndReport(nil, DefaultPromoteConfig(), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Error("expected nil output for nil secrets")
	}
	if !strings.Contains(buf.String(), "no secrets") {
		t.Errorf("expected 'no secrets' in output, got %q", buf.String())
	}
}

func TestPromoteAndReport_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"dev_KEY": "value"}
	cfg := PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_"}
	_, err := PromoteAndReport(secrets, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "promote:") {
		t.Errorf("expected summary in output, got %q", buf.String())
	}
}

func TestPromoteAndReport_NilWriter_UsesStdout(t *testing.T) {
	secrets := map[string]string{"dev_X": "1"}
	cfg := PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_"}
	_, err := PromoteAndReport(secrets, cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error using nil writer: %v", err)
	}
}

func TestPromoteAndReport_DryRunNote(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"dev_KEY": "val"}
	cfg := PromoteConfig{FromPrefix: "dev_", ToPrefix: "prod_", DryRun: true}
	_, err := PromoteAndReport(secrets, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "dry-run") {
		t.Errorf("expected dry-run note in output, got %q", buf.String())
	}
}
