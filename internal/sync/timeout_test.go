package sync

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithTimeout_SucceedsWithinDeadline(t *testing.T) {
	err := WithTimeout(1*time.Second, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWithTimeout_PropagatesFnError(t *testing.T) {
	expected := errors.New("vault unavailable")
	err := WithTimeout(1*time.Second, func(ctx context.Context) error {
		return expected
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestWithTimeout_ExpiresWhenFnTooSlow(t *testing.T) {
	err := WithTimeout(50*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded wrapped in error, got %v", err)
	}
}

func TestWithTimeout_ZeroDurationReturnsError(t *testing.T) {
	err := WithTimeout(0, func(ctx context.Context) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error for zero timeout, got nil")
	}
}

func TestDefaultTimeoutConfig(t *testing.T) {
	cfg := DefaultTimeoutConfig()
	if cfg.FetchTimeout != 15*time.Second {
		t.Errorf("expected FetchTimeout 15s, got %v", cfg.FetchTimeout)
	}
	if cfg.WriteTimeout != 5*time.Second {
		t.Errorf("expected WriteTimeout 5s, got %v", cfg.WriteTimeout)
	}
}

func TestWithTimeout_ContextPassedToFn(t *testing.T) {
	var received context.Context
	_ = WithTimeout(1*time.Second, func(ctx context.Context) error {
		received = ctx
		return nil
	})
	if received == nil {
		t.Fatal("expected context to be passed to fn")
	}
}
