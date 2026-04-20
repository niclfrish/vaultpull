package sync

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

// Report summarises a completed sync operation for human-readable output.
type Report struct {
	EnvFile    string
	SecretPath string
	Namespace  string
	DryRun     bool
	Duration   time.Duration
	Plan       Plan
}

// Print writes a formatted sync report to w.
func (r *Report) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Sync Report\n")
	fmt.Fprintf(tw, "-----------\n")
	fmt.Fprintf(tw, "File:\t%s\n", r.EnvFile)
	fmt.Fprintf(tw, "Secret Path:\t%s\n", r.SecretPath)
	if r.Namespace != "" {
		fmt.Fprintf(tw, "Namespace:\t%s\n", r.Namespace)
	}
	fmt.Fprintf(tw, "Dry Run:\t%v\n", r.DryRun)
	fmt.Fprintf(tw, "Duration:\t%s\n", r.Duration.Round(time.Millisecond))
	fmt.Fprintf(tw, "\n")

	summary := r.Plan.Summary()
	fmt.Fprintf(tw, "Changes:\n")
	fmt.Fprintf(tw, "  Added:\t%d\n", summary[ActionAdd])
	fmt.Fprintf(tw, "  Updated:\t%d\n", summary[ActionUpdate])
	fmt.Fprintf(tw, "  Removed:\t%d\n", summary[ActionRemove])
	fmt.Fprintf(tw, "  Unchanged:\t%d\n", summary[ActionUnchanged])

	if len(r.Plan.Changes) > 0 {
		fmt.Fprintf(tw, "\nDetail:\n")
		for _, c := range r.Plan.Changes {
			if c.Action == ActionUnchanged {
				continue
			}
			action := strings.ToUpper(string(c.Action))
			fmt.Fprintf(tw, "  [%s]\t%s\n", action, c.Key)
		}
	}

	tw.Flush()
}
