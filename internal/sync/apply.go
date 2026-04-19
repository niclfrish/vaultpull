package sync

import (
	"fmt"
	"io"
	"os"
)

// Applier applies a Plan to a .env file.
type Applier struct {
	writer EnvWriter
	out    io.Writer
}

// EnvWriter is the interface used to write env files.
type EnvWriter interface {
	Write(path string, data map[string]string) error
}

// NewApplier creates a new Applier.
func NewApplier(writer EnvWriter) *Applier {
	return &Applier{writer: writer, out: os.Stdout}
}

// Apply writes the merged secrets from the plan to the given file path.
// It returns an error if no changes are present or writing fails.
func (a *Applier) Apply(path string, plan *Plan) error {
	if !plan.HasChanges() {
		fmt.Fprintln(a.out, "nothing to apply: no changes detected")
		return nil
	}

	merged := make(map[string]string)
	for _, e := range plan.Unchanged {
		merged[e.Key] = e.OldValue
	}
	for _, e := range plan.Added {
		merged[e.Key] = e.NewValue
	}
	for _, e := range plan.Changed {
		merged[e.Key] = e.NewValue
	}
	// Removed entries are intentionally omitted.

	if err := a.writer.Write(path, merged); err != nil {
		return fmt.Errorf("apply: write %s: %w", path, err)
	}

	fmt.Fprintf(a.out, "applied %d changes to %s\n", plan.Summary().Total, path)
	return nil
}
