package sync

import (
	"testing"
)

func TestMergeWithStrategy_NoSources(t *testing.T) {
	out, err := MergeWithStrategy(DefaultMergeStrategyConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestMergeWithStrategy_UnknownStrategy(t *testing.T) {
	cfg := MergeStrategyConfig{Strategy: "bogus"}
	_, err := MergeWithStrategy(cfg, map[string]string{"A": "1"})
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestMergeWithStrategy_Last_OverridesEarlier(t *testing.T) {
	cfg := MergeStrategyConfig{Strategy: MergeStrategyLast}
	src1 := map[string]string{"KEY": "first", "ONLY_IN_1": "yes"}
	src2 := map[string]string{"KEY": "second", "ONLY_IN_2": "yes"}
	out, err := MergeWithStrategy(cfg, src1, src2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", out["KEY"])
	}
	if out["ONLY_IN_1"] != "yes" || out["ONLY_IN_2"] != "yes" {
		t.Errorf("missing keys from sources: %v", out)
	}
}

func TestMergeWithStrategy_First_KeepsEarlier(t *testing.T) {
	cfg := MergeStrategyConfig{Strategy: MergeStrategyFirst}
	src1 := map[string]string{"KEY": "first"}
	src2 := map[string]string{"KEY": "second"}
	out, err := MergeWithStrategy(cfg, src1, src2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", out["KEY"])
	}
}

func TestMergeWithStrategy_Error_ConflictReturnsError(t *testing.T) {
	cfg := MergeStrategyConfig{Strategy: MergeStrategyError}
	src1 := map[string]string{"KEY": "v1"}
	src2 := map[string]string{"KEY": "v2"}
	_, err := MergeWithStrategy(cfg, src1, src2)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestMergeWithStrategy_Error_SameValueNoConflict(t *testing.T) {
	cfg := MergeStrategyConfig{Strategy: MergeStrategyError}
	src1 := map[string]string{"KEY": "same"}
	src2 := map[string]string{"KEY": "same"}
	out, err := MergeWithStrategy(cfg, src1, src2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "same" {
		t.Errorf("expected 'same', got %q", out["KEY"])
	}
}

func TestMergeWithStrategy_Prefix_AddsPrefix(t *testing.T) {
	cfg := MergeStrategyConfig{Strategy: MergeStrategyPrefix}
	src1 := map[string]string{"key": "v1"}
	src2 := map[string]string{"key": "v2"}
	out, err := MergeWithStrategy(cfg, src1, src2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "v1" {
		t.Errorf("expected original key 'KEY'='v1', got %q", out["KEY"])
	}
	if out["SRC1_KEY"] != "v2" {
		t.Errorf("expected prefixed key 'SRC1_KEY'='v2', got %q", out["SRC1_KEY"])
	}
}

func TestDefaultMergeStrategyConfig(t *testing.T) {
	cfg := DefaultMergeStrategyConfig()
	if cfg.Strategy != MergeStrategyLast {
		t.Errorf("expected default strategy 'last', got %q", cfg.Strategy)
	}
}
