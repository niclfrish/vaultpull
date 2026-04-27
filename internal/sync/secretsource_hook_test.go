package sync

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAnnotateStage_NilSecrets(t *testing.T) {
	stage := AnnotateStage(DefaultSecretSourceConfig(SourceTypeVault, "kv/app"))
	_, err := stage(nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestAnnotateStage_InjectsSourceType(t *testing.T) {
	src := SecretSource{
		Type:      SourceTypeVault,
		Location:  "kv/myapp",
		FetchedAt: time.Now().UTC(),
	}
	stage := AnnotateStage(src)
	out, err := stage(map[string]string{"KEY": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["__source_type"] != "vault" {
		t.Errorf("expected vault, got %s", out["__source_type"])
	}
	if out["KEY"] != "val" {
		t.Error("original key should be preserved")
	}
}

func TestStripAnnotationsStage_NilSecrets(t *testing.T) {
	stage := StripAnnotationsStage()
	_, err := stage(nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestStripAnnotationsStage_RemovesAnnotations(t *testing.T) {
	stage := StripAnnotationsStage()
	input := map[string]string{
		"DB_URL":             "postgres://",
		"__source_type":     "vault",
		"__source_location": "kv/app",
	}
	out, err := stage(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["__source_type"]; ok {
		t.Error("annotation should be stripped")
	}
	if out["DB_URL"] != "postgres://" {
		t.Error("DB_URL should remain")
	}
}

func TestReportSource_NilSecrets(t *testing.T) {
	var buf bytes.Buffer
	_, err := ReportSource(nil, SecretSource{}, &buf)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestReportSource_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	src := SecretSource{
		Type:      SourceTypeEnv,
		Location:  ".env",
		FetchedAt: time.Now().UTC(),
		Namespace: "dev",
	}
	secrets := map[string]string{"X": "1"}
	out, err := ReportSource(secrets, src, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "1" {
		t.Error("secrets should pass through")
	}
	if !strings.Contains(buf.String(), "env") {
		t.Errorf("expected source type in output: %s", buf.String())
	}
}

func TestReportSource_NilWriter_UsesStdout(t *testing.T) {
	secrets := map[string]string{"A": "b"}
	out, err := ReportSource(secrets, DefaultSecretSourceConfig(SourceTypeFile, "secrets.env"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "b" {
		t.Error("expected pass-through")
	}
}
