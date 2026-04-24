package sync

import (
	"strings"
	"testing"
)

func TestFilterStage_RemovesNonMatchingKeys(t *testing.T) {
	f := NewFilter(FilterCriteria{Prefix: "APP_"})
	stage := FilterStage(f)
	out, err := stage(map[string]string{"APP_KEY": "v1", "OTHER": "v2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_KEY"]; !ok {
		t.Error("expected APP_KEY to be present")
	}
	if _, ok := out["OTHER"]; ok {
		t.Error("expected OTHER to be removed")
	}
}

func TestTransformStage_AppliesTransforms(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform)
	stage := TransformStage(tr)
	out, err := stage(map[string]string{"KEY": "  hello  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "hello" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestRequiredKeysStage_AllPresent(t *testing.T) {
	stage := RequiredKeysStage("A", "B")
	_, err := stage(map[string]string{"A": "1", "B": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequiredKeysStage_MissingKey(t *testing.T) {
	stage := RequiredKeysStage("A", "MISSING")
	_, err := stage(map[string]string{"A": "1"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "MISSING") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestTruncateStage_TruncatesLongValues(t *testing.T) {
	stage := TruncateStage()
	cfg := DefaultTruncateConfig()
	long := strings.Repeat("x", cfg.MaxLength+10)
	out, err := stage(map[string]string{"KEY": long})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out["KEY"]) > cfg.MaxLength {
		t.Errorf("value not truncated: len=%d", len(out["KEY"]))
	}
}

func TestPipeline_WithBuiltins_EndToEnd(t *testing.T) {
	p := NewPipeline()
	p.AddStage("filter", FilterStage(NewFilter(FilterCriteria{Prefix: "APP_"})))
	p.AddStage("require", RequiredKeysStage("APP_TOKEN"))
	p.AddStage("transform", TransformStage(NewTransformer(TrimSpaceTransform)))

	input := map[string]string{
		"APP_TOKEN": "  secret  ",
		"IGNORED":   "noise",
	}
	out, err := p.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_TOKEN"] != "secret" {
		t.Errorf("expected trimmed token, got %q", out["APP_TOKEN"])
	}
	if _, ok := out["IGNORED"]; ok {
		t.Error("IGNORED should have been filtered out")
	}
}
