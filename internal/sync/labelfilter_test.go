package sync

import (
	"testing"
)

func TestLabelFilter_NoRequired_ReturnsAll(t *testing.T) {
	f := NewLabelFilter(nil)
	input := map[string]string{
		"KEY_A__label__env=prod": "val1",
		"KEY_B": "val2",
	}
	out := f.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestLabelFilter_MatchesSingleLabel(t *testing.T) {
	f := NewLabelFilter(map[string]string{"env": "prod"})
	input := map[string]string{
		"DB_HOST__label__env=prod": "localhost",
		"API_KEY__label__env=staging": "secret",
		"PLAIN_KEY": "plain",
	}
	out := f.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if v, ok := out["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
}

func TestLabelFilter_StripsLabelsFromOutputKeys(t *testing.T) {
	f := NewLabelFilter(map[string]string{"tier": "backend"})
	input := map[string]string{
		"SECRET__label__tier=backend__label__region=us": "value",
	}
	out := f.Apply(input)
	if _, ok := out["SECRET"]; !ok {
		t.Errorf("expected stripped key 'SECRET', got %v", out)
	}
}

func TestLabelFilter_RequiresAllLabels(t *testing.T) {
	f := NewLabelFilter(map[string]string{"env": "prod", "tier": "backend"})
	input := map[string]string{
		"KEY_A__label__env=prod__label__tier=backend": "yes",
		"KEY_B__label__env=prod": "no",
	}
	out := f.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d: %v", len(out), out)
	}
	if _, ok := out["KEY_A"]; !ok {
		t.Error("expected KEY_A in output")
	}
}

func TestLabelFilter_EmptyInput_ReturnsEmpty(t *testing.T) {
	f := NewLabelFilter(map[string]string{"env": "prod"})
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestSplitKeyLabels_NoLabels(t *testing.T) {
	base, labels := splitKeyLabels("MY_KEY")
	if base != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %s", base)
	}
	if len(labels) != 0 {
		t.Errorf("expected no labels, got %v", labels)
	}
}

func TestSplitKeyLabels_MultipleLabels(t *testing.T) {
	base, labels := splitKeyLabels("MY_KEY__label__env=prod__label__region=eu")
	if base != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %s", base)
	}
	if labels["env"] != "prod" || labels["region"] != "eu" {
		t.Errorf("unexpected labels: %v", labels)
	}
}
