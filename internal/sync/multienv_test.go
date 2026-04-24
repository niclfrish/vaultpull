package sync

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewMultiEnvWriter_NoTargets(t *testing.T) {
	_, err := NewMultiEnvWriter(nil, func(string, map[string]string) error { return nil })
	if err == nil {
		t.Fatal("expected error for empty targets")
	}
}

func TestNewMultiEnvWriter_NilWriter(t *testing.T) {
	targets := []EnvTarget{{Name: "a", Path: ".env"}}
	_, err := NewMultiEnvWriter(targets, nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestMultiEnvWriter_WriteAll_Success(t *testing.T) {
	written := map[string]map[string]string{}
	writer := func(path string, secrets map[string]string) error {
		written[path] = secrets
		return nil
	}
	targets := []EnvTarget{
		{Name: "prod", Path: ".env.prod", Namespace: ""},
		{Name: "staging", Path: ".env.staging", Namespace: "STAGING"},
	}
	m, err := NewMultiEnvWriter(targets, writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{"KEY": "val"}
	results := m.WriteAll(secrets)
	if err := AnyError(results); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if _, ok := written[".env.prod"]["KEY"]; !ok {
		t.Error("expected KEY in prod target")
	}
	if _, ok := written[".env.staging"]["STAGING_KEY"]; !ok {
		t.Error("expected STAGING_KEY in staging target")
	}
}

func TestMultiEnvWriter_WriteAll_PartialError(t *testing.T) {
	writer := func(path string, secrets map[string]string) error {
		if path == ".env.bad" {
			return errors.New("disk full")
		}
		return nil
	}
	targets := []EnvTarget{
		{Name: "good", Path: ".env.good"},
		{Name: "bad", Path: ".env.bad"},
	}
	m, _ := NewMultiEnvWriter(targets, writer)
	results := m.WriteAll(map[string]string{"X": "1"})
	if results["good"] != nil {
		t.Error("expected no error for good target")
	}
	if results["bad"] == nil {
		t.Error("expected error for bad target")
	}
}

func TestMultiEnvWriter_TargetNames_Sorted(t *testing.T) {
	targets := []EnvTarget{
		{Name: "z", Path: ".env.z"},
		{Name: "a", Path: ".env.a"},
		{Name: "m", Path: ".env.m"},
	}
	m, _ := NewMultiEnvWriter(targets, func(string, map[string]string) error { return nil })
	names := m.TargetNames()
	expected := []string{"a", "m", "z"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("pos %d: got %q want %q", i, n, expected[i])
		}
	}
}

func TestAnyError_NilMap(t *testing.T) {
	if err := AnyError(map[string]error{"a": nil, "b": nil}); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestAnyError_ReturnsFirst(t *testing.T) {
	results := map[string]error{
		"alpha": nil,
		"beta":  fmt.Errorf("oops"),
	}
	err := AnyError(results)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, results["beta"]) {
		t.Errorf("unexpected error: %v", err)
	}
}
