package sync

import (
	"fmt"
	"regexp"
)

// SchemaRule defines a validation rule for a secret key.
type SchemaRule struct {
	Key      string
	Pattern  string
	Required bool
}

// Schema holds a set of rules used to validate secret maps.
type Schema struct {
	rules    []SchemaRule
	compiled map[string]*regexp.Regexp
}

// NewSchema creates a Schema from the provided rules.
// Returns an error if any pattern fails to compile.
func NewSchema(rules []SchemaRule) (*Schema, error) {
	s := &Schema{
		rules:    rules,
		compiled: make(map[string]*regexp.Regexp),
	}
	for _, r := range rules {
		if r.Pattern == "" {
			continue
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("schema: invalid pattern for key %q: %w", r.Key, err)
		}
		s.compiled[r.Key] = re
	}
	return s, nil
}

// SchemaViolation describes a single rule violation.
type SchemaViolation struct {
	Key     string
	Message string
}

// Validate checks secrets against the schema rules.
// It returns all violations found (not just the first).
func (s *Schema) Validate(secrets map[string]string) []SchemaViolation {
	var violations []SchemaViolation
	for _, rule := range s.rules {
		val, exists := secrets[rule.Key]
		if rule.Required && !exists {
			violations = append(violations, SchemaViolation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}
		if !exists {
			continue
		}
		if re, ok := s.compiled[rule.Key]; ok {
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern),
				})
			}
		}
	}
	return violations
}
