package sync

import (
	"errors"
	"testing"
	"time"
)

func noSleep(_ time.Duration) {}

func newFastRetrier(max int) *Retrier {
	r := NewRetrier(RetryConfig{
		MaxAttempts: max,
		Delay:       0,
		Multiplier:  1.0,
	})
	r.sleep = noSleep
	return r
}

func TestRetrier_SucceedsFirstAttempt(t *testing.T) {
	r := newFastRetrier(3)
	calls := 0
	err := r.Run(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetrier_RetriesOnFailure(t *testing.T) {
	r := newFastRetrier(3)
	calls := 0
	sentinel := errors.New("transient")

	err := r.Run(func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetrier_ExhaustsAttempts(t *testing.T) {
	r := newFastRetrier(3)
	sentinel := errors.New("persistent")

	err := r.Run(func() error { return sentinel })
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel in error chain, got %v", err)
	}
}

func TestNonRetryable_Wrapping(t *testing.T) {
	base := errors.New("base error")
	wrapped := NonRetryable(base)

	if !errors.Is(wrapped, ErrNonRetryable) {
		t.Error("expected ErrNonRetryable in chain")
	}
	if !errors.Is(wrapped, base) {
		t.Error("expected base error in chain")
	}
}

func TestIsRetryable(t *testing.T) {
	if !IsRetryable(errors.New("ordinary")) {
		t.Error("ordinary error should be retryable")
	}
	if IsRetryable(NonRetryable(errors.New("fatal"))) {
		t.Error("NonRetryable error should not be retryable")
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}
