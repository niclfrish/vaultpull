package sync

import (
	"testing"
)

func TestNewSchema_InvalidPattern(t *testing.T) {
	_, err := NewSchema([]SchemaRule{
		{Key: "FOO", Pattern: "["},
	})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestNewSchema_Success(t *testing.T) {
	_, err := NewSchema([]SchemaRule{
		{Key: "FOO", Pattern: `^[a-z]+$`, Required: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchema_Validate_RequiredMissing(t *testing.T) {
	s, _ := NewSchema([]SchemaRule{
		{Key: "DB_URL", Required: true},
	})
	violations := s.Validate(map[string]string{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_URL" {
		t.Errorf("expected key DB_URL, got %q", violations[0].Key)
	}
}

func TestSchema_Validate_PatternMismatch(t *testing.T) {
	s, _ := NewSchema([]SchemaRule{
		{Key: "PORT", Pattern: `^\d+$`},
	})
	violations := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestSchema_Validate_PatternMatch(t *testing.T) {
	s, _ := NewSchema([]SchemaRule{
		{Key: "PORT", Pattern: `^\d+$`},
	})
	violations := s.Validate(map[string]string{"PORT": "8080"})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestSchema_Validate_NoRules(t *testing.T) {
	s, _ := NewSchema(nil)
	violations := s.Validate(map[string]string{"ANY": "value"})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestSchema_Validate_MultipleViolations(t *testing.T) {
	s, _ := NewSchema([]SchemaRule{
		{Key: "HOST", Required: true},
		{Key: "PORT", Pattern: `^\d+$`, Required: true},
	})
	violations := s.Validate(map[string]string{})
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestSchema_Validate_OptionalKeyAbsent(t *testing.T) {
	s, _ := NewSchema([]SchemaRule{
		{Key: "OPTIONAL", Pattern: `^yes|no$`, Required: false},
	})
	violations := s.Validate(map[string]string{})
	if len(violations) != 0 {
		t.Errorf("expected no violations for absent optional key, got %v", violations)
	}
}
