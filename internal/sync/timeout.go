package sync

import (
	"context"
	"fmt"
	"time"
)

// TimeoutConfig holds configuration for operation timeouts.
type TimeoutConfig struct {
	// FetchTimeout is the maximum duration allowed for fetching secrets from Vault.
	FetchTimeout time.Duration
	// WriteTimeout is the maximum duration allowed for writing the .env file.
	WriteTimeout time.Duration
}

// DefaultTimeoutConfig returns a TimeoutConfig with sensible defaults.
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		FetchTimeout: 15 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

// WithTimeout executes fn within the given timeout duration.
// It returns an error if the context is cancelled or the deadline is exceeded.
func WithTimeout(timeout time.Duration, fn func(ctx context.Context) error) error {
	if timeout <= 0 {
		return fmt.Errorf("timeout must be greater than zero")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type result struct {
		err error
	}

	ch := make(chan result, 1)
	go func() {
		ch <- result{err: fn(ctx)}
	}()

	select {
	case res := <-ch:
		return res.err
	case <-ctx.Done():
		return fmt.Errorf("operation timed out after %s: %w", timeout, ctx.Err())
	}
}
