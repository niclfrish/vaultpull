package sync

import (
	"errors"
	"strings"
	"testing"
)

func TestPipeline_Empty_ReturnsInput(t *testing.T) {
	p := NewPipeline()
	input := map[string]string{"KEY": "value"}
	out, err := p.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value %q, got %q", "value", out["KEY"])
	}
}

func TestPipeline_SingleStage_Transforms(t *testing.T) {
	p := NewPipeline()
	p.AddStage("uppercase-values", func(s map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(s))
		for k, v := range s {
			out[k] = strings.ToUpper(v)
		}
		return out, nil
	})
	out, err := p.Run(map[string]string{"KEY": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["KEY"])
	}
}

func TestPipeline_MultipleStages_ChainedCorrectly(t *testing.T) {
	p := NewPipeline()
	p.AddStage("add-prefix", func(s map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(s))
		for k, v := range s {
			out[k] = "pre_" + v
		}
		return out, nil
	})
	p.AddStage("add-suffix", func(s map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(s))
		for k, v := range s {
			out[k] = v + "_suf"
		}
		return out, nil
	})
	out, err := p.Run(map[string]string{"K": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "pre_val_suf" {
		t.Errorf("expected pre_val_suf, got %q", out["K"])
	}
}

func TestPipeline_StageError_WrapsName(t *testing.T) {
	p := NewPipeline()
	p.AddStage("failing-stage", func(s map[string]string) (map[string]string, error) {
		return nil, errors.New("boom")
	})
	_, err := p.Run(map[string]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failing-stage") {
		t.Errorf("expected stage name in error, got: %v", err)
	}
}

// TestPipeline_StageError_StopsExecution verifies that when a stage returns an
// error, subsequent stages are not executed.
func TestPipeline_StageError_StopsExecution(t *testing.T) {
	p := NewPipeline()
	executed := false
	p.AddStage("failing-stage", func(s map[string]string) (map[string]string, error) {
		return nil, errors.New("boom")
	})
	p.AddStage("should-not-run", func(s map[string]string) (map[string]string, error) {
		executed = true
		return s, nil
	})
	_, err := p.Run(map[string]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if executed {
		t.Error("expected subsequent stage to not execute after error")
	}
}

func TestPipeline_StageNames_ReturnsOrder(t *testing.T) {
	p := NewPipeline()
	p.AddStage("a", func(s map[string]string) (map[string]string, error) { return s, nil })
	p.AddStage("b", func(s map[string]string) (map[string]string, error) { return s, nil })
	names := p.StageNames()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("unexpected stage names: %v", names)
	}
}

func TestPipeline_StageCount(t *testing.T) {
	p := NewPipeline()
	if p.StageCount() != 0 {
		t.Errorf("expected 0, got %d", p.StageCount())
	}
	p.AddStage("x", func(s map[string]string) (map[string]string, error) { return s, nil })
	if p.StageCount() != 1 {
		t.Errorf("expected 1, got %d", p.StageCount())
	}
}
