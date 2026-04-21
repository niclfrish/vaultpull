package sync

import (
	"context"
	"errors"
	"testing"
)

func TestHookRunner_NoHooks(t *testing.T) {
	r := NewHookRunner()
	if err := r.Run(context.Background(), HookPreFetch, nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestHookRunner_RunsMatchingHooks(t *testing.T) {
	r := NewHookRunner()
	called := false
	r.Register(HookPostFetch, func(_ context.Context, _ HookEvent, _ map[string]string) error {
		called = true
		return nil
	})
	if err := r.Run(context.Background(), HookPostFetch, map[string]string{"K": "V"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected hook to be called")
	}
}

func TestHookRunner_SkipsNonMatchingEvent(t *testing.T) {
	r := NewHookRunner()
	called := false
	r.Register(HookPreApply, func(_ context.Context, _ HookEvent, _ map[string]string) error {
		called = true
		return nil
	})
	_ = r.Run(context.Background(), HookPostApply, nil)
	if called {
		t.Fatal("hook should not have been called for a different event")
	}
}

func TestHookRunner_ShortCircuitsOnError(t *testing.T) {
	r := NewHookRunner()
	sentinel := errors.New("hook error")
	secondCalled := false

	r.Register(HookPreApply, func(_ context.Context, _ HookEvent, _ map[string]string) error {
		return sentinel
	})
	r.Register(HookPreApply, func(_ context.Context, _ HookEvent, _ map[string]string) error {
		secondCalled = true
		return nil
	})

	err := r.Run(context.Background(), HookPreApply, nil)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if secondCalled {
		t.Fatal("second hook must not run after first returns error")
	}
}

func TestHookRunner_MultipleEvents(t *testing.T) {
	r := NewHookRunner()
	var log []HookEvent

	for _, ev := range []HookEvent{HookPreFetch, HookPostFetch, HookPreApply, HookPostApply} {
		ev := ev
		r.Register(ev, func(_ context.Context, e HookEvent, _ map[string]string) error {
			log = append(log, e)
			return nil
		})
	}

	for _, ev := range []HookEvent{HookPreFetch, HookPostFetch, HookPreApply, HookPostApply} {
		if err := r.Run(context.Background(), ev, nil); err != nil {
			t.Fatalf("unexpected error on %s: %v", ev, err)
		}
	}

	if len(log) != 4 {
		t.Fatalf("expected 4 hook calls, got %d", len(log))
	}
}
