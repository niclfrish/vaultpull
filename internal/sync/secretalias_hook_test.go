package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestAliasStage_NoAliases_PassThrough(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	stage := AliasStage(DefaultAliasConfig())
	result, err := stage.Run(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", result["KEY"])
	}
}

func TestAliasStage_RenamesKey(t *testing.T) {
	secrets := map[string]string{"OLD": "secret"}
	cfg := AliasConfig{Aliases: map[string]string{"OLD": "NEW"}, KeepOriginal: false}
	stage := AliasStage(cfg)
	result, err := stage.Run(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["OLD"]; ok {
		t.Error("OLD key should be removed")
	}
	if result["NEW"] != "secret" {
		t.Errorf("expected NEW=secret, got %q", result["NEW"])
	}
}

func TestAliasStage_StageName(t *testing.T) {
	stage := AliasStage(DefaultAliasConfig())
	if stage.Name != "alias" {
		t.Errorf("expected stage name 'alias', got %q", stage.Name)
	}
}

func TestAliasAndReport_NilSecrets(t *testing.T) {
	_, err := AliasAndReport(nil, DefaultAliasConfig(), nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestAliasAndReport_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"A": "1", "B": "2"}
	cfg := AliasConfig{Aliases: map[string]string{"A": "ALPHA"}, KeepOriginal: false}
	_, err := AliasAndReport(secrets, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alias:") {
		t.Errorf("expected 'alias:' in output, got: %q", out)
	}
	if !strings.Contains(out, "1 aliases applied") {
		t.Errorf("expected '1 aliases applied' in output, got: %q", out)
	}
}

func TestAliasAndReport_NilWriter_UsesStdout(t *testing.T) {
	secrets := map[string]string{"X": "y"}
	// Should not panic; output goes to stdout.
	_, err := AliasAndReport(secrets, DefaultAliasConfig(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAliasAndReport_ConflictReturnsError(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"A": "1", "B": "2"}
	cfg := AliasConfig{Aliases: map[string]string{"A": "B"}, KeepOriginal: false}
	_, err := AliasAndReport(secrets, cfg, &buf)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}
