package sync

import (
	"strings"
	"testing"
)

func TestNewTemplateRenderer_EmptyText(t *testing.T) {
	_, err := NewTemplateRenderer("")
	if err == nil {
		t.Fatal("expected error for empty template text")
	}
}

func TestNewTemplateRenderer_InvalidSyntax(t *testing.T) {
	_, err := NewTemplateRenderer("{{ .Unclosed")
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestNewTemplateRenderer_Success(t *testing.T) {
	_, err := NewTemplateRenderer("KEY={{ index . \"KEY\" }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRender_BasicInterpolation(t *testing.T) {
	r, _ := NewTemplateRenderer(`DB_URL={{ index . "DB_URL" }}`)
	secrets := map[string]string{"DB_URL": "postgres://localhost/db"}
	out, err := r.Render(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "postgres://localhost/db") {
		t.Errorf("expected rendered output to contain DB_URL value, got: %s", out)
	}
}

func TestRender_MissingKeyReturnsError(t *testing.T) {
	r, _ := NewTemplateRenderer(`{{ index . "MISSING" }}`)
	_, err := r.Render(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing key with missingkey=error")
	}
}

func TestRender_NilSecretsUsesEmptyMap(t *testing.T) {
	r, _ := NewTemplateRenderer(`static-value`)
	out, err := r.Render(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "static-value" {
		t.Errorf("expected 'static-value', got: %s", out)
	}
}

func TestRenderToMap_Success(t *testing.T) {
	r, _ := NewTemplateRenderer(`{{ index . "TOKEN" }}`)
	secrets := map[string]string{"TOKEN": "abc123"}
	m, err := r.RenderToMap(secrets, "RENDERED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["RENDERED"] != "abc123" {
		t.Errorf("expected 'abc123', got: %s", m["RENDERED"])
	}
}

func TestRenderToMap_EmptyOutputKey(t *testing.T) {
	r, _ := NewTemplateRenderer(`hello`)
	_, err := r.RenderToMap(map[string]string{}, "")
	if err == nil {
		t.Fatal("expected error for empty output key")
	}
}
