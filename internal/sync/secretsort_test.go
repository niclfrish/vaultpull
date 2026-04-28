package sync

import (
	"testing"
)

func TestSortSecrets_NilSecrets(t *testing.T) {
	_, err := SortSecrets(nil, DefaultSortConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestSortSecrets_AlphaOrder(t *testing.T) {
	secrets := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	pairs, err := SortSecrets(secrets, SortConfig{Order: SortOrderAlpha})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
	if pairs[0][0] != "APPLE" || pairs[1][0] != "MANGO" || pairs[2][0] != "ZEBRA" {
		t.Errorf("unexpected order: %v", pairs)
	}
}

func TestSortSecrets_AlphaDescOrder(t *testing.T) {
	secrets := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	pairs, err := SortSecrets(secrets, SortConfig{Order: SortOrderAlphaDesc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pairs[0][0] != "ZEBRA" || pairs[1][0] != "MANGO" || pairs[2][0] != "APPLE" {
		t.Errorf("unexpected order: %v", pairs)
	}
}

func TestSortSecrets_KeyLengthOrder(t *testing.T) {
	secrets := map[string]string{"AB": "x", "A": "y", "ABC": "z"}
	pairs, err := SortSecrets(secrets, SortConfig{Order: SortOrderKeyLength})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs[0][0]) > len(pairs[1][0]) || len(pairs[1][0]) > len(pairs[2][0]) {
		t.Errorf("expected ascending key length order, got: %v", pairs)
	}
}

func TestSortSecrets_ValueLengthOrder(t *testing.T) {
	secrets := map[string]string{"A": "hello world", "B": "hi", "C": "hey"}
	pairs, err := SortSecrets(secrets, SortConfig{Order: SortOrderValueLength})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs[0][1]) > len(pairs[1][1]) {
		t.Errorf("expected ascending value length, got: %v", pairs)
	}
}

func TestSortSecrets_UnknownOrder(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	_, err := SortSecrets(secrets, SortConfig{Order: "random"})
	if err == nil {
		t.Fatal("expected error for unknown order")
	}
}

func TestSortSecrets_PrefixFirst(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "h", "APP_NAME": "n", "DB_PORT": "p"}
	pairs, err := SortSecrets(secrets, SortConfig{Order: SortOrderAlpha, Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pairs[0][0] != "DB_HOST" || pairs[1][0] != "DB_PORT" || pairs[2][0] != "APP_NAME" {
		t.Errorf("expected DB_ keys first, got: %v", pairs)
	}
}

func TestSortSummary(t *testing.T) {
	pairs := [][2]string{{"A", "1"}, {"B", "2"}}
	summary := SortSummary(pairs, SortConfig{Order: SortOrderAlpha})
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSortSecrets_EmptyMap(t *testing.T) {
	pairs, err := SortSecrets(map[string]string{}, DefaultSortConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected empty result, got %d", len(pairs))
	}
}
