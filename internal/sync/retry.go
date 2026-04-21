package sync

import (
	"errors"
	"fmt"
	"time"
)

// RetryConfig holds configuration for retry behaviour.
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Retrier executes a function with exponential back-off retry logic.
type Retrier struct {
	cfg   RetryConfig
	sleep func(time.Duration)
}

// NewRetrier creates a new Retrier with the given config.
func NewRetrier(cfg RetryConfig) *Retrier {
	return &Retrier{cfg: cfg, sleep: time.Sleep}
}

// Run executes fn up to MaxAttempts times, sleeping between failures.
// It returns the last error if all attempts are exhausted.
func (r *Retrier) Run(fn func() error) error {
	var err error
	delay := r.cfg.Delay

	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if attempt < r.cfg.MaxAttempts {
			r.sleep(delay)
			delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		}
	}

	return fmt.Errorf("all %d attempts failed: %w", r.cfg.MaxAttempts, err)
}

// IsRetryable returns true for errors that should trigger a retry.
// Callers may wrap errors with a sentinel to opt out of retrying.
func IsRetryable(err error) bool {
	return !errors.Is(err, ErrNonRetryable)
}

// ErrNonRetryable signals that an error should not be retried.
var ErrNonRetryable = errors.New("non-retryable error")

// NonRetryable wraps err so that the retrier will not retry it.
func NonRetryable(err error) error {
	return fmt.Errorf("%w: %w", ErrNonRetryable, err)
}
