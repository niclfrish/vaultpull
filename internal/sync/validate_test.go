package sync

import (
	"testing"
)

func TestValidate_NilSecrets(t *testing.T) {
	v := NewValidator(nil, 0)
	_, err := v.Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil secrets map")
	}
}

func TestValidate_EmptyMap_NoRequired(t *testing.T) {
	v := NewValidator(nil, 0)
	res, err := v.Validate(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsValid() {
		t.Errorf("expected valid result, got errors: %v", res.Errors)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	v := NewValidator([]string{"DB_PASSWORD", "API_KEY"}, 0)
	res, err := v.Validate(map[string]string{"DB_PASSWORD": "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsValid() {
		t.Error("expected validation to fail due to missing required key")
	}
	if len(res.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(res.Errors))
	}
}

func TestValidate_ValueExceedsMaxLength(t *testing.T) {
	v := NewValidator(nil, 10)
	secrets := map[string]string{"SHORT": "ok", "LONG": "this value is definitely too long"}
	res, err := v.Validate(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsValid() {
		t.Errorf("expected no errors, got: %v", res.Errors)
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestValidate_KeyWithWhitespace(t *testing.T) {
	v := NewValidator(nil, 0)
	res, err := v.Validate(map[string]string{"BAD KEY": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsValid() {
		t.Error("expected validation error for key with whitespace")
	}
}

func TestValidationResult_Summary(t *testing.T) {
	res := &ValidationResult{}
	if res.Summary() != "validation passed: no issues found" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
	res.Errors = append(res.Errors, "some error")
	res.Warnings = append(res.Warnings, "some warning")
	summary := res.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
