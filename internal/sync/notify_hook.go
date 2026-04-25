package sync

import (
	"fmt"
	"io"
)

// NotifyOnSync emits an INFO event after a successful secret fetch.
// It is intended to be used as a post-fetch hook in the sync pipeline.
func NotifyOnSync(notifier *Notifier, namespace string, w io.Writer) func(secrets map[string]string) error {
	return func(secrets map[string]string) error {
		if notifier == nil {
			return nil
		}
		meta := map[string]string{
			"keys": fmt.Sprintf("%d", len(secrets)),
		}
		errs := notifier.Emit(NotifyInfo, namespace, "secrets synced successfully", meta)
		if len(errs) > 0 {
			// Non-fatal: log to writer but do not abort the sync.
			if w != nil {
				for _, e := range errs {
					fmt.Fprintf(w, "notify warning: %v\n", e)
				}
			}
		}
		return nil
	}
}

// NotifyOnError emits an ERROR event and returns the original error unchanged.
func NotifyOnError(notifier *Notifier, namespace string, w io.Writer) func(err error) error {
	return func(err error) error {
		if notifier == nil || err == nil {
			return err
		}
		meta := map[string]string{
			"error": err.Error(),
		}
		notifyErrs := notifier.Emit(NotifyError, namespace, "sync failed", meta)
		if len(notifyErrs) > 0 && w != nil {
			for _, ne := range notifyErrs {
				fmt.Fprintf(w, "notify warning: %v\n", ne)
			}
		}
		return err
	}
}
