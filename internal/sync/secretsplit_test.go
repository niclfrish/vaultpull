package sync

import (
	"testing"
)

func TestSplitSecrets_NilSecrets(t *testing.T) {
	_, _, _, err := SplitSecrets(nil, DefaultSplitConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestSplitSecrets_EmptyDelimiter(t *testing.T) {
	cfg := DefaultSplitConfig()
	cfg.Delimiter = ""
	_, _, _, err := SplitSecrets(map[string]string{"K": "v"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestSplitSecrets_BasicSplit(t *testing.T) {
	secrets := map[string]string{
		"DB": "host:localhost",
	}
	cfg := DefaultSplitConfig()
	out, results, summary, err := SplitSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Split != 1 {
		t.Errorf("expected 1 split, got %d", summary.Split)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].NewKey != "DB_host" {
		t.Errorf("unexpected new key: %s", results[0].NewKey)
	}
	if out["DB_host"] != "localhost" {
		t.Errorf("unexpected value: %s", out["DB_host"])
	}
}

func TestSplitSecrets_SkipsWhenIndexOutOfRange(t *testing.T) {
	secrets := map[string]string{
		"PLAIN": "nodelimiterhere",
	}
	cfg := DefaultSplitConfig()
	out, _, summary, err := SplitSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", summary.Skipped)
	}
	if out["PLAIN"] != "nodelimiterhere" {
		t.Error("original key/value should be preserved when skipped")
	}
}

func TestSplitSecrets_OnlyKeysFilters(t *testing.T) {
	secrets := map[string]string{
		"TARGET": "k:v",
		"OTHER":  "x:y",
	}
	cfg := DefaultSplitConfig()
	cfg.OnlyKeys = []string{"TARGET"}
	_, results, summary, err := SplitSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Split != 1 {
		t.Errorf("expected 1 split, got %d", summary.Split)
	}
	if summary.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", summary.Skipped)
	}
	if len(results) != 1 || results[0].OriginalKey != "TARGET" {
		t.Error("only TARGET should have been split")
	}
}

func TestSplitSecrets_TrimsWhitespace(t *testing.T) {
	secrets := map[string]string{
		"KEY": " name : value with spaces ",
	}
	cfg := DefaultSplitConfig()
	out, results, _, err := SplitSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].NewKey != "KEY_name" {
		t.Errorf("expected trimmed key, got %s", results[0].NewKey)
	}
	if out["KEY_name"] != "value with spaces" {
		t.Errorf("expected trimmed value, got %s", out["KEY_name"])
	}
}
