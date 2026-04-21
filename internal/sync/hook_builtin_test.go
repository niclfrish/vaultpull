package sync

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestLoggingHook_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	hook := LoggingHook(&buf)
	secrets := map[string]string{"A": "1", "B": "2"}

	if err := hook(context.Background(), HookPostFetch, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, string(HookPostFetch)) {
		t.Errorf("expected event name in output, got: %s", out)
	}
	if !strings.Contains(out, "secrets=2") {
		t.Errorf("expected secret count in output, got: %s", out)
	}
}

func TestLoggingHook_NilWriterUsesStdout(t *testing.T) {
	// Should not panic.
	hook := LoggingHook(nil)
	if err := hook(context.Background(), HookPreFetch, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireKeysHook_AllPresent(t *testing.T) {
	hook := RequireKeysHook([]string{"DB_URL", "API_KEY"})
	secrets := map[string]string{"DB_URL": "postgres://", "API_KEY": "secret"}
	if err := hook(context.Background(), HookPostFetch, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireKeysHook_MissingKey(t *testing.T) {
	hook := RequireKeysHook([]string{"DB_URL", "MISSING"})
	secrets := map[string]string{"DB_URL": "postgres://"}
	err := hook(context.Background(), HookPostFetch, secrets)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "MISSING") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestCountLimitHook_UnderLimit(t *testing.T) {
	hook := CountLimitHook(5)
	secrets := map[string]string{"A": "1", "B": "2"}
	if err := hook(context.Background(), HookPostFetch, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCountLimitHook_ExceedsLimit(t *testing.T) {
	hook := CountLimitHook(2)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	err := hook(context.Background(), HookPostFetch, secrets)
	if err == nil {
		t.Fatal("expected error when count exceeds limit")
	}
}
