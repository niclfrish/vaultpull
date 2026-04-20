package sync

import (
	"fmt"
	"io"
	"os"
)

// ValidateAndReport runs validation on the given secrets and prints the result
// to the provided writer. It returns an error if validation itself fails or if
// the result contains errors.
func ValidateAndReport(w io.Writer, secrets map[string]string, requiredKeys []string) error {
	if w == nil {
		w = os.Stdout
	}

	v := NewValidator(requiredKeys, 0)
	result, err := v.Validate(secrets)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	for _, warn := range result.Warnings {
		fmt.Fprintf(w, "  [WARN]  %s\n", warn)
	}
	for _, e := range result.Errors {
		fmt.Fprintf(w, "  [ERROR] %s\n", e)
	}

	fmt.Fprintf(w, "  %s\n", result.Summary())

	if !result.IsValid() {
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}
	return nil
}
