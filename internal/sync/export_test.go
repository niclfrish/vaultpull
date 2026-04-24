package sync

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewExporter_InvalidFormat(t *testing.T) {
	_, err := NewExporter("xml", nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestNewExporter_NilWriterDefaultsToStdout(t *testing.T) {
	e, err := NewExporter(FormatJSON, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestExport_JSON(t *testing.T) {
	var buf bytes.Buffer
	e, _ := NewExporter(FormatJSON, &buf)

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := e.Export(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected JSON output: %v", got)
	}
}

func TestExport_Dotenv(t *testing.T) {
	var buf bytes.Buffer
	e, _ := NewExporter(FormatDotenv, &buf)

	secrets := map[string]string{"KEY": "value with spaces"}
	if err := e.Export(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "KEY=") {
		t.Errorf("expected KEY= in output, got: %s", output)
	}
	if strings.HasPrefix(output, "export ") {
		t.Errorf("dotenv format should not include export prefix")
	}
}

func TestExport_ExportFormat(t *testing.T) {
	var buf bytes.Buffer
	e, _ := NewExporter(FormatExport, &buf)

	secrets := map[string]string{"TOKEN": "abc123"}
	if err := e.Export(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "export TOKEN=") {
		t.Errorf("expected 'export TOKEN=' prefix, got: %s", output)
	}
}

func TestExport_NilSecrets(t *testing.T) {
	var buf bytes.Buffer
	e, _ := NewExporter(FormatDotenv, &buf)
	if err := e.Export(nil); err != nil {
		t.Fatalf("unexpected error for nil secrets: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil secrets")
	}
}

func TestExport_Dotenv_SortedKeys(t *testing.T) {
	var buf bytes.Buffer
	e, _ := NewExporter(FormatDotenv, &buf)

	secrets := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	_ = e.Export(secrets)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA=") {
		t.Errorf("expected first line to start with ALPHA=, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA=") {
		t.Errorf("expected last line to start with ZEBRA=, got: %s", lines[2])
	}
}
