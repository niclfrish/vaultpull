package sync

import (
	"strings"
	"testing"
)

func TestMergeByPriority_EmptySources(t *testing.T) {
	result, err := MergeByPriority(DefaultPriorityConfig(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestMergeByPriority_InvalidPriority(t *testing.T) {
	sources := []PrioritySource{
		{Name: "a", Priority: 0, Secrets: map[string]string{"K": "v"}},
	}
	_, err := MergeByPriority(DefaultPriorityConfig(), sources)
	if err == nil {
		t.Fatal("expected error for priority 0")
	}
}

func TestMergeByPriority_DuplicatePriority(t *testing.T) {
	sources := []PrioritySource{
		{Name: "a", Priority: 1, Secrets: map[string]string{"K": "v"}},
		{Name: "b", Priority: 1, Secrets: map[string]string{"K": "v2"}},
	}
	_, err := MergeByPriority(DefaultPriorityConfig(), sources)
	if err == nil {
		t.Fatal("expected error for duplicate priority")
	}
}

func TestMergeByPriority_HigherPriorityWins(t *testing.T) {
	sources := []PrioritySource{
		{Name: "low", Priority: 2, Secrets: map[string]string{"DB_PASS": "low-value"}},
		{Name: "high", Priority: 1, Secrets: map[string]string{"DB_PASS": "high-value"}},
	}
	result, err := MergeByPriority(DefaultPriorityConfig(), sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_PASS"] != "high-value" {
		t.Errorf("expected high-value, got %q", result["DB_PASS"])
	}
}

func TestMergeByPriority_ConflictAnnotation(t *testing.T) {
	cfg := PriorityConfig{ConflictPrefix: "__conflict_"}
	sources := []PrioritySource{
		{Name: "vault", Priority: 1, Secrets: map[string]string{"TOKEN": "vault-token"}},
		{Name: "local", Priority: 2, Secrets: map[string]string{"TOKEN": "local-token"}},
	}
	result, err := MergeByPriority(cfg, sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["TOKEN"] != "vault-token" {
		t.Errorf("expected vault-token, got %q", result["TOKEN"])
	}
	conflictKey := "__conflict_local_TOKEN"
	if result[conflictKey] != "local-token" {
		t.Errorf("expected conflict annotation %q = local-token, got %q", conflictKey, result[conflictKey])
	}
}

func TestMergeByPriority_NoConflictAnnotationWhenPrefixEmpty(t *testing.T) {
	cfg := PriorityConfig{ConflictPrefix: ""}
	sources := []PrioritySource{
		{Name: "a", Priority: 1, Secrets: map[string]string{"K": "winner"}},
		{Name: "b", Priority: 2, Secrets: map[string]string{"K": "loser"}},
	}
	result, err := MergeByPriority(cfg, sources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
}

func TestPrioritySummary_Output(t *testing.T) {
	sources := []PrioritySource{
		{Name: "vault", Priority: 1, Secrets: map[string]string{"A": "1", "B": "2"}},
		{Name: "local", Priority: 2, Secrets: map[string]string{"C": "3"}},
	}
	merged := map[string]string{"A": "1", "B": "2", "C": "3"}
	summary := PrioritySummary(sources, merged)
	if !strings.Contains(summary, "3 keys") {
		t.Errorf("expected merged key count in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "vault") {
		t.Errorf("expected source name in summary, got: %s", summary)
	}
}
