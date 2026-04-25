package sync

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is in the open state.
var ErrCircuitOpen = errors.New("circuit breaker is open")

type circuitState int

const (
	stateClosed circuitState = iota
	stateOpen
	stateHalfOpen
)

// CircuitBreakerConfig holds configuration for the circuit breaker.
type CircuitBreakerConfig struct {
	MaxFailures  int
	OpenDuration time.Duration
}

// DefaultCircuitBreakerConfig returns sensible defaults.
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxFailures:  3,
		OpenDuration: 30 * time.Second,
	}
}

// CircuitBreaker protects a resource from repeated failures.
type CircuitBreaker struct {
	cfg      CircuitBreakerConfig
	mu       sync.Mutex
	state    circuitState
	failures int
	openedAt time.Time
}

// NewCircuitBreaker creates a new CircuitBreaker with the given config.
func NewCircuitBreaker(cfg CircuitBreakerConfig) (*CircuitBreaker, error) {
	if cfg.MaxFailures <= 0 {
		return nil, fmt.Errorf("circuit breaker: MaxFailures must be > 0, got %d", cfg.MaxFailures)
	}
	if cfg.OpenDuration <= 0 {
		return nil, fmt.Errorf("circuit breaker: OpenDuration must be > 0, got %v", cfg.OpenDuration)
	}
	return &CircuitBreaker{cfg: cfg, state: stateClosed}, nil
}

// Do executes fn if the circuit is closed or half-open, recording success or failure.
func (cb *CircuitBreaker) Do(fn func() error) error {
	cb.mu.Lock()
	switch cb.state {
	case stateOpen:
		if time.Since(cb.openedAt) >= cb.cfg.OpenDuration {
			cb.state = stateHalfOpen
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	case stateClosed, stateHalfOpen:
		// allowed to proceed
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()
	if err != nil {
		cb.failures++
		if cb.failures >= cb.cfg.MaxFailures || cb.state == stateHalfOpen {
			cb.state = stateOpen
			cb.openedAt = time.Now()
		}
		return err
	}
	cb.failures = 0
	cb.state = stateClosed
	return nil
}

// State returns the current circuit state as a string.
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case stateOpen:
		return "open"
	case stateHalfOpen:
		return "half-open"
	default:
		return "closed"
	}
}
