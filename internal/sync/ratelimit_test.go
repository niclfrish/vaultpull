package sync

import (
	"context"
	"testing"
	"time"
)

func TestNewRateLimiter_InvalidRPS(t *testing.T) {
	_, err := NewRateLimiter(RateLimitConfig{RequestsPerSecond: 0, Burst: 5})
	if err == nil {
		t.Fatal("expected error for RequestsPerSecond=0")
	}
}

func TestNewRateLimiter_InvalidBurst(t *testing.T) {
	_, err := NewRateLimiter(RateLimitConfig{RequestsPerSecond: 10, Burst: 0})
	if err == nil {
		t.Fatal("expected error for Burst=0")
	}
}

func TestNewRateLimiter_Success(t *testing.T) {
	rl, err := NewRateLimiter(DefaultRateLimitConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rl.Stop()

	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
}

func TestRateLimiter_WaitConsumesToken(t *testing.T) {
	rl, err := NewRateLimiter(RateLimitConfig{RequestsPerSecond: 100, Burst: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rl.Stop()

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("Wait() error on attempt %d: %v", i, err)
		}
	}
}

func TestRateLimiter_WaitRespectsContextCancel(t *testing.T) {
	// Burst of 1: consume the only token first, then cancel should fire.
	rl, err := NewRateLimiter(RateLimitConfig{RequestsPerSecond: 1, Burst: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rl.Stop()

	ctx := context.Background()
	// Drain the single burst token.
	_ = rl.Wait(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = rl.Wait(ctx2)
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
}

func TestDefaultRateLimitConfig(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	if cfg.RequestsPerSecond <= 0 {
		t.Errorf("expected positive RequestsPerSecond, got %d", cfg.RequestsPerSecond)
	}
	if cfg.Burst <= 0 {
		t.Errorf("expected positive Burst, got %d", cfg.Burst)
	}
}
