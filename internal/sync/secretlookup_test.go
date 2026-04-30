package sync

import (
	"testing"
)

func TestLookupSecrets_NilSecrets(t *testing.T) {
	_, err := LookupSecrets(nil, []string{"KEY"}, DefaultLookupConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestLookupSecrets_NoQueries(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	_, err := LookupSecrets(secrets, nil, DefaultLookupConfig())
	if err == nil {
		t.Fatal("expected error for empty queries")
	}
}

func TestLookupSecrets_ExactMatch(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	cfg := DefaultLookupConfig()
	results, err := LookupSecrets(secrets, []string{"db_host"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", results[0].Key)
	}
}

func TestLookupSecrets_CaseSensitiveNoMatch(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	cfg := DefaultLookupConfig()
	cfg.CaseSensitive = true
	results, err := LookupSecrets(secrets, []string{"db_host"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for case-sensitive mismatch")
	}
}

func TestLookupSecrets_PartialMatch(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "secret"}
	cfg := DefaultLookupConfig()
	cfg.PartialMatch = true
	results, err := LookupSecrets(secrets, []string{"db_"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestLookupSecrets_DeduplicatesResults(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	cfg := DefaultLookupConfig()
	results, err := LookupSecrets(secrets, []string{"foo", "FOO"}, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 deduplicated result, got %d", len(results))
	}
}

func TestLookupSummary_NoResults(t *testing.T) {
	s := LookupSummary(nil)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestLookupSummary_WithResults(t *testing.T) {
	results := []LookupResult{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	s := LookupSummary(results)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
