package sync

import (
	"testing"
)

func TestLabelFilterStage_EmptyLabels_PassThrough(t *testing.T) {
	stage := LabelFilterStage(nil)
	input := map[string]string{"A": "1", "B": "2"}
	out, err := stage.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestLabelFilterStage_FiltersCorrectly(t *testing.T) {
	stage := LabelFilterStage(map[string]string{"env": "prod"})
	input := map[string]string{
		"DB__label__env=prod": "db-val",
		"CACHE__label__env=dev": "cache-val",
	}
	out, err := stage.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["DB"]; !ok {
		t.Error("expected DB in output")
	}
}

func TestParseLabelFlags_Valid(t *testing.T) {
	labels, err := ParseLabelFlags([]string{"env=prod", "region=us-east"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if labels["env"] != "prod" || labels["region"] != "us-east" {
		t.Errorf("unexpected labels: %v", labels)
	}
}

func TestParseLabelFlags_Invalid(t *testing.T) {
	_, err := ParseLabelFlags([]string{"no-equals-sign"})
	if err == nil {
		t.Error("expected error for invalid label flag")
	}
}

func TestParseLabelFlags_EmptyKey(t *testing.T) {
	_, err := ParseLabelFlags([]string{"=value"})
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestLabelFilterStage_Name(t *testing.T) {
	stage := LabelFilterStage(map[string]string{"x": "y"})
	if stage.Name != "label-filter" {
		t.Errorf("expected stage name 'label-filter', got %q", stage.Name)
	}
}
