package sync

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// SchemaCheckResult holds the outcome of a schema validation pass.
type SchemaCheckResult struct {
	Violations []SchemaViolation
	Passed     bool
}

// ValidateSchema runs schema validation against secrets and writes a
// human-readable report to w. If w is nil, os.Stdout is used.
// Returns an error if any violations are found.
func ValidateSchema(schema *Schema, secrets map[string]string, w io.Writer) (*SchemaCheckResult, error) {
	if w == nil {
		w = os.Stdout
	}
	if schema == nil {
		return &SchemaCheckResult{Passed: true}, nil
	}

	violations := schema.Validate(secrets)
	result := &SchemaCheckResult{
		Violations: violations,
		Passed:     len(violations) == 0,
	}

	if result.Passed {
		fmt.Fprintln(w, "schema validation passed: all rules satisfied")
		return result, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("schema validation failed: %d violation(s)\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  - [%s] %s\n", v.Key, v.Message))
	}
	fmt.Fprint(w, sb.String())

	return result, fmt.Errorf("schema validation failed with %d violation(s)", len(violations))
}
