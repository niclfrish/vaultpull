package sync

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// LoggingHook returns a HookFunc that writes a timestamped message to w
// whenever the hook fires. If w is nil, os.Stdout is used.
func LoggingHook(w io.Writer) HookFunc {
	if w == nil {
		w = os.Stdout
	}
	return func(_ context.Context, event HookEvent, secrets map[string]string) error {
		count := len(secrets)
		fmt.Fprintf(w, "[%s] hook fired: event=%s secrets=%d\n",
			time.Now().UTC().Format(time.RFC3339), event, count)
		return nil
	}
}

// RequireKeysHook returns a HookFunc that fails if any of the required keys
// are absent from the secrets map. Intended for use on HookPostFetch.
func RequireKeysHook(required []string) HookFunc {
	return func(_ context.Context, _ HookEvent, secrets map[string]string) error {
		for _, k := range required {
			if _, ok := secrets[k]; !ok {
				return fmt.Errorf("required secret key %q not found in fetched secrets", k)
			}
		}
		return nil
	}
}

// CountLimitHook returns a HookFunc that fails when the number of secrets
// exceeds max. Useful as a sanity-check guard.
func CountLimitHook(max int) HookFunc {
	return func(_ context.Context, _ HookEvent, secrets map[string]string) error {
		if len(secrets) > max {
			return fmt.Errorf("secrets count %d exceeds limit %d", len(secrets), max)
		}
		return nil
	}
}

// FilterEventHook wraps a HookFunc so that it only executes for the specified
// event. For all other events the hook is a no-op. This is useful when
// composing hooks that should only run at a particular lifecycle stage.
func FilterEventHook(event HookEvent, h HookFunc) HookFunc {
	return func(ctx context.Context, e HookEvent, secrets map[string]string) error {
		if e != event {
			return nil
		}
		return h(ctx, e, secrets)
	}
}
