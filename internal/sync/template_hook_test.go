package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderTemplate_NilRenderer_ReturnsSecrets(t *testing.T) {
	fn := RenderTemplate(nil, "OUT", &bytes.Buffer{})
	secrets := map[string]string{"A": "1"}
	out, err := fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected secrets to be returned unchanged")
	}
}

func TestRenderTemplate_WritesRenderedOutput(t *testing.T) {
	r, _ := NewTemplateRenderer(`HOST={{ index . "HOST" }}`)
	var buf bytes.Buffer
	fn := RenderTemplate(r, "OUT", &buf)
	_, err := fn(map[string]string{"HOST": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "localhost") {
		t.Errorf("expected buffer to contain 'localhost', got: %s", buf.String())
	}
}

func TestRenderTemplate_NilWriter_UsesStdout(t *testing.T) {
	r, _ := NewTemplateRenderer(`static`)
	fn := RenderTemplate(r, "OUT", nil)
	_, err := fn(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderTemplate_RenderError_ReturnsError(t *testing.T) {
	r, _ := NewTemplateRenderer(`{{ index . "MISSING" }}`)
	fn := RenderTemplate(r, "OUT", &bytes.Buffer{})
	_, err := fn(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestTemplateStage_AddsRenderedKey(t *testing.T) {
	r, _ := NewTemplateRenderer(`{{ index . "FOO" }}-rendered`)
	stage := TemplateStage(r, "FOO_RENDERED")
	secrets := map[string]string{"FOO": "bar"}
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO_RENDERED"] != "bar-rendered" {
		t.Errorf("expected 'bar-rendered', got: %s", out["FOO_RENDERED"])
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected original key FOO to be preserved")
	}
}

func TestTemplateStage_NilRenderer_PassesThrough(t *testing.T) {
	stage := TemplateStage(nil, "OUT")
	secrets := map[string]string{"X": "y"}
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "y" {
		t.Errorf("expected secrets to pass through unchanged")
	}
}
