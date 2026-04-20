package sync

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationResult holds the outcome of a secrets validation.
type ValidationResult struct {
	Errors   []string
	Warnings []string
}

// IsValid returns true if there are no validation errors.
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// Summary returns a human-readable summary of the validation result.
func (r *ValidationResult) Summary() string {
	if r.IsValid() && len(r.Warnings) == 0 {
		return "validation passed: no issues found"
	}
	parts := []string{}
	if len(r.Errors) > 0 {
		parts = append(parts, fmt.Sprintf("%d error(s)", len(r.Errors)))
	}
	if len(r.Warnings) > 0 {
		parts = append(parts, fmt.Sprintf("%d warning(s)", len(r.Warnings)))
	}
	return "validation finished: " + strings.Join(parts, ", ")
}

// Validator checks secrets maps for common issues.
type Validator struct {
	requiredKeys []string
	maxValueLen  int
}

// NewValidator creates a Validator with optional required keys and a max value length.
func NewValidator(requiredKeys []string, maxValueLen int) *Validator {
	if maxValueLen <= 0 {
		maxValueLen = 4096
	}
	return &Validator{requiredKeys: requiredKeys, maxValueLen: maxValueLen}
}

// Validate inspects the provided secrets map and returns a ValidationResult.
func (v *Validator) Validate(secrets map[string]string) (*ValidationResult, error) {
	if secrets == nil {
		return nil, errors.New("secrets map must not be nil")
	}
	result := &ValidationResult{}

	for _, req := range v.requiredKeys {
		if _, ok := secrets[req]; !ok {
			result.Errors = append(result.Errors, fmt.Sprintf("required key missing: %s", req))
		}
	}

	for k, val := range secrets {
		if k == "" {
			result.Errors = append(result.Errors, "empty key found in secrets")
			continue
		}
		if len(val) > v.maxValueLen {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("value for key %q exceeds %d characters", k, v.maxValueLen))
		}
		if strings.ContainsAny(k, " \t") {
			result.Errors = append(result.Errors, fmt.Sprintf("key contains whitespace: %q", k))
		}
	}

	return result, nil
}
