package sync

import (
	"context"
	"fmt"
	"time"
)

// ThrottleConfig holds configuration for secret fetch throttling.
type ThrottleConfig struct {
	// MinInterval is the minimum time between successive fetch operations.
	MinInterval time.Duration
	// MaxBatchSize limits how many secrets are processed per batch (0 = unlimited).
	MaxBatchSize int
}

// DefaultThrottleConfig returns a sensible default throttle configuration.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		MinInterval:  200 * time.Millisecond,
		MaxBatchSize: 50,
	}
}

// Throttler enforces a minimum interval between calls and optional batch size limits.
type Throttler struct {
	cfg      ThrottleConfig
	lastCall time.Time
	sleepFn  func(time.Duration)
}

// NewThrottler creates a Throttler with the given config.
func NewThrottler(cfg ThrottleConfig) (*Throttler, error) {
	if cfg.MinInterval < 0 {
		return nil, fmt.Errorf("throttle: MinInterval must be non-negative, got %s", cfg.MinInterval)
	}
	if cfg.MaxBatchSize < 0 {
		return nil, fmt.Errorf("throttle: MaxBatchSize must be non-negative, got %d", cfg.MaxBatchSize)
	}
	return &Throttler{
		cfg:     cfg,
		sleepFn: time.Sleep,
	}, nil
}

// Wait blocks until the minimum interval since the last call has elapsed,
// or until ctx is cancelled.
func (t *Throttler) Wait(ctx context.Context) error {
	if t.cfg.MinInterval == 0 {
		t.lastCall = time.Now()
		return nil
	}
	elapsed := time.Since(t.lastCall)
	remaining := t.cfg.MinInterval - elapsed
	if remaining > 0 {
		select {
		case <-ctx.Done():
			return fmt.Errorf("throttle: context cancelled while waiting: %w", ctx.Err())
		case <-time.After(remaining):
		}
	}
	t.lastCall = time.Now()
	return nil
}

// Batch splits secrets into chunks no larger than MaxBatchSize.
// If MaxBatchSize is 0, the original map is returned as a single batch.
func (t *Throttler) Batch(secrets map[string]string) []map[string]string {
	if t.cfg.MaxBatchSize == 0 || len(secrets) <= t.cfg.MaxBatchSize {
		return []map[string]string{secrets}
	}
	var batches []map[string]string
	current := make(map[string]string, t.cfg.MaxBatchSize)
	for k, v := range secrets {
		current[k] = v
		if len(current) == t.cfg.MaxBatchSize {
			batches = append(batches, current)
			current = make(map[string]string, t.cfg.MaxBatchSize)
		}
	}
	if len(current) > 0 {
		batches = append(batches, current)
	}
	return batches
}
