package sync

import "context"

// HookEvent represents the lifecycle stage at which a hook fires.
type HookEvent string

const (
	HookPreFetch  HookEvent = "pre_fetch"
	HookPostFetch HookEvent = "post_fetch"
	HookPreApply  HookEvent = "pre_apply"
	HookPostApply HookEvent = "post_apply"
)

// HookFunc is a function invoked at a specific lifecycle event.
// It receives the event name and the current secrets map (may be nil for pre-fetch).
// Returning an error aborts the sync pipeline.
type HookFunc func(ctx context.Context, event HookEvent, secrets map[string]string) error

// HookRunner manages and executes ordered lifecycle hooks.
type HookRunner struct {
	hooks []hookEntry
}

type hookEntry struct {
	event HookEvent
	fn    HookFunc
}

// NewHookRunner returns an initialised HookRunner.
func NewHookRunner() *HookRunner {
	return &HookRunner{}
}

// Register appends a hook for the given event.
func (r *HookRunner) Register(event HookEvent, fn HookFunc) {
	r.hooks = append(r.hooks, hookEntry{event: event, fn: fn})
}

// Run executes all hooks registered for the given event in order.
// The first error short-circuits execution.
func (r *HookRunner) Run(ctx context.Context, event HookEvent, secrets map[string]string) error {
	for _, h := range r.hooks {
		if h.event != event {
			continue
		}
		if err := h.fn(ctx, event, secrets); err != nil {
			return err
		}
	}
	return nil
}
