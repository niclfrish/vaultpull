package sync

import (
	"errors"
	"sync"
	"time"
)

// TokenRotateConfig holds configuration for token rotation.
type TokenRotateConfig struct {
	// RotateAfter is the duration after which the token should be rotated.
	RotateAfter time.Duration
	// GracePeriod is the overlap window where both old and new tokens are valid.
	GracePeriod time.Duration
}

// DefaultTokenRotateConfig returns sensible defaults.
func DefaultTokenRotateConfig() TokenRotateConfig {
	return TokenRotateConfig{
		RotateAfter: 24 * time.Hour,
		GracePeriod: 5 * time.Minute,
	}
}

// TokenFetcher retrieves a fresh token from an external source.
type TokenFetcher func() (string, error)

// TokenRotator manages automatic rotation of a Vault token.
type TokenRotator struct {
	mu        sync.RWMutex
	current   string
	fetchedAt time.Time
	cfg       TokenRotateConfig
	fetcher   TokenFetcher
}

// NewTokenRotator creates a TokenRotator with an initial token.
func NewTokenRotator(initial string, fetcher TokenFetcher, cfg TokenRotateConfig) (*TokenRotator, error) {
	if initial == "" {
		return nil, errors.New("tokenrotate: initial token must not be empty")
	}
	if fetcher == nil {
		return nil, errors.New("tokenrotate: fetcher must not be nil")
	}
	if cfg.RotateAfter <= 0 {
		return nil, errors.New("tokenrotate: RotateAfter must be positive")
	}
	return &TokenRotator{
		current:   initial,
		fetchedAt: time.Now(),
		cfg:       cfg,
		fetcher:   fetcher,
	}, nil
}

// Token returns the current token, rotating it if expired.
func (r *TokenRotator) Token() (string, error) {
	r.mu.RLock()
	expired := time.Since(r.fetchedAt) >= r.cfg.RotateAfter
	r.mu.RUnlock()

	if !expired {
		r.mu.RLock()
		defer r.mu.RUnlock()
		return r.current, nil
	}

	return r.rotate()
}

func (r *TokenRotator) rotate() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock.
	if time.Since(r.fetchedAt) < r.cfg.RotateAfter {
		return r.current, nil
	}

	newToken, err := r.fetcher()
	if err != nil {
		// Return stale token within grace period.
		if time.Since(r.fetchedAt) < r.cfg.RotateAfter+r.cfg.GracePeriod {
			return r.current, nil
		}
		return "", errors.New("tokenrotate: failed to rotate token: " + err.Error())
	}

	r.current = newToken
	r.fetchedAt = time.Now()
	return r.current, nil
}

// Age returns how long the current token has been held.
func (r *TokenRotator) Age() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return time.Since(r.fetchedAt)
}
