package cmd

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/sync"
)

func TestParseTargets_Valid(t *testing.T) {
	raw := []string{"prod:.env.prod", "staging:.env.staging:STAGING"}
	targets, err := parseTargets(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}
	if targets[0].Name != "prod" || targets[0].Path != ".env.prod" || targets[0].Namespace != "" {
		t.Errorf("unexpected prod target: %+v", targets[0])
	}
	if targets[1].Name != "staging" || targets[1].Namespace != "STAGING" {
		t.Errorf("unexpected staging target: %+v", targets[1])
	}
}

func TestParseTargets_MissingPath(t *testing.T) {
	_, err := parseTargets([]string{"onlyname"})
	if err == nil {
		t.Fatal("expected error for target without path")
	}
}

func TestParseTargets_Empty(t *testing.T) {
	targets, err := parseTargets([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(targets) != 0 {
		t.Errorf("expected empty targets")
	}
}

func TestParseTargets_NoNamespace(t *testing.T) {
	targets, err := parseTargets([]string{"dev:.env.dev"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if targets[0].Namespace != "" {
		t.Errorf("expected empty namespace, got %q", targets[0].Namespace)
	}
}

func TestParseTargets_IntegrationWithMultiEnvWriter(t *testing.T) {
	targets, err := parseTargets([]string{"a:/tmp/a.env:NS"})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	written := map[string]map[string]string{}
	w, err := sync.NewMultiEnvWriter(targets, func(path string, s map[string]string) error {
		written[path] = s
		return nil
	})
	if err != nil {
		t.Fatalf("writer error: %v", err)
	}
	w.WriteAll(map[string]string{"KEY": "value"})
	if written["/tmp/a.env"]["NS_KEY"] != "value" {
		t.Errorf("expected NS_KEY=value in output, got %v", written["/tmp/a.env"])
	}
}
