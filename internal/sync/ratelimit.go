package sync

import (
	"context"
	"fmt"
	"time"
)

// RateLimitConfig holds configuration for rate limiting Vault API calls.
type RateLimitConfig struct {
	// RequestsPerSecond is the maximum number of requests allowed per second.
	RequestsPerSecond int
	// Burst is the maximum number of requests allowed to burst above the rate.
	Burst int
}

// DefaultRateLimitConfig returns a sensible default rate limit configuration.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond: 10,
		Burst:             5,
	}
}

// RateLimiter controls the rate of outgoing requests using a token bucket.
type RateLimiter struct {
	tokens   chan struct{}
	interval time.Duration
	stop     chan struct{}
}

// NewRateLimiter creates a RateLimiter from the given config and starts
// the background token refill goroutine.
func NewRateLimiter(cfg RateLimitConfig) (*RateLimiter, error) {
	if cfg.RequestsPerSecond <= 0 {
		return nil, fmt.Errorf("ratelimit: RequestsPerSecond must be > 0, got %d", cfg.RequestsPerSecond)
	}
	if cfg.Burst <= 0 {
		return nil, fmt.Errorf("ratelimit: Burst must be > 0, got %d", cfg.Burst)
	}

	tokens := make(chan struct{}, cfg.Burst)
	for i := 0; i < cfg.Burst; i++ {
		tokens <- struct{}{}
	}

	rl := &RateLimiter{
		tokens:   tokens,
		interval: time.Second / time.Duration(cfg.RequestsPerSecond),
		stop:     make(chan struct{}),
	}
	go rl.refill()
	return rl, nil
}

// Wait blocks until a token is available or the context is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.tokens:
		return nil
	}
}

// Stop shuts down the background refill goroutine.
func (r *RateLimiter) Stop() {
	close(r.stop)
}

func (r *RateLimiter) refill() {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-r.stop:
			return
		case <-ticker.C:
			select {
			case r.tokens <- struct{}{}:
			default:
				// bucket full, discard token
			}
		}
	}
}
