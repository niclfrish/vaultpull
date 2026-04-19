package sync

import (
	"testing"
)

func TestNamespacedPath_NoNamespace(t *testing.T) {
	got, err := NamespacedPath("", "secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/data/app" {
		t.Errorf("got %q, want %q", got, "secret/data/app")
	}
}

func TestNamespacedPath_WithNamespace(t *testing.T) {
	got, err := NamespacedPath("staging", "secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "staging/secret/data/app"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNamespacedPath_EmptyPath(t *testing.T) {
	_, err := NamespacedPath("staging", "")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestPrefixKeys_NoNamespace(t *testing.T) {
	input := map[string]string{"DB_PASS": "secret"}
	result := PrefixKeys("", input)
	if result["DB_PASS"] != "secret" {
		t.Error("expected unchanged keys when namespace is empty")
	}
}

func TestPrefixKeys_WithNamespace(t *testing.T) {
	input := map[string]string{"DB_PASS": "secret"}
	result := PrefixKeys("my-app", input)
	if v, ok := result["MY_APP_DB_PASS"]; !ok || v != "secret" {
		t.Errorf("expected MY_APP_DB_PASS=secret, got %v", result)
	}
}
