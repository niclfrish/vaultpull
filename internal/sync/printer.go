package sync

import (
	"fmt"
	"io"
	"strings"
)

// PrintPlan writes a human-readable diff plan to the given writer.
func PrintPlan(w io.Writer, plan *Plan) {
	if !plan.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return
	}
	for _, e := range plan.Entries {
		switch e.Change {
		case ChangeAdd:
			fmt.Fprintf(w, "  + %s=%s\n", e.Key, maskValue(e.New))
		case ChangeRemove:
			fmt.Fprintf(w, "  - %s\n", e.Key)
		case ChangeUpdate:
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", e.Key, maskValue(e.Old), maskValue(e.New))
		}
	}
	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintln(w, plan.Summary())
}

// maskValue replaces all but the first 2 chars with asterisks for display.
func maskValue(v string) string {
	if len(v) <= 2 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-2)
}
