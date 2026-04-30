package sync

import (
	"testing"
)

func TestSampleSecrets_NilSecrets(t *testing.T) {
	_, err := SampleSecrets(nil, DefaultSampleConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestSampleSecrets_InvalidMaxSamples(t *testing.T) {
	cfg := DefaultSampleConfig()
	cfg.MaxSamples = 0
	_, err := SampleSecrets(map[string]string{"A": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for MaxSamples=0")
	}
}

func TestSampleSecrets_ReturnsAllWhenUnderLimit(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	cfg := DefaultSampleConfig()
	cfg.MaxSamples = 10
	out, err := SampleSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}

func TestSampleSecrets_LimitsToMaxSamples(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5"}
	cfg := DefaultSampleConfig()
	cfg.MaxSamples = 3
	out, err := SampleSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 sampled keys, got %d", len(out))
	}
}

func TestSampleSecrets_DeterministicWithSameSeed(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5"}
	cfg := DefaultSampleConfig()
	cfg.MaxSamples = 2
	out1, _ := SampleSecrets(secrets, cfg)
	out2, _ := SampleSecrets(secrets, cfg)
	for k := range out1 {
		if _, ok := out2[k]; !ok {
			t.Errorf("key %q present in first sample but not second", k)
		}
	}
}

func TestSampleSecrets_OnlyKeysFilter(t *testing.T) {
	secrets := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	cfg := DefaultSampleConfig()
	cfg.MaxSamples = 10
	cfg.OnlyKeys = []string{"APP_"}
	out, err := SampleSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys with APP_ prefix, got %d", len(out))
	}
	if _, ok := out["DB_URL"]; ok {
		t.Error("DB_URL should have been excluded")
	}
}

func TestSampleSummary(t *testing.T) {
	s := SampleSummary(100, 10)
	if s == "" {
		t.Error("expected non-empty summary")
	}
	expected := "sampled 10 of 100 secrets"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
