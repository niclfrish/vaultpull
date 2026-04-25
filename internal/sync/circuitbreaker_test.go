package sync

import (
	"errors"
	"testing"
	"time"
)

var errFake = errors.New("fake error")

func TestNewCircuitBreaker_InvalidMaxFailures(t *testing.T) {
	_, err := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 0, OpenDuration: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxFailures=0")
	}
}

func TestNewCircuitBreaker_InvalidOpenDuration(t *testing.T) {
	_, err := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 1, OpenDuration: 0})
	if err == nil {
		t.Fatal("expected error for OpenDuration=0")
	}
}

func TestNewCircuitBreaker_Success(t *testing.T) {
	cb, err := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cb.State() != "closed" {
		t.Fatalf("expected closed, got %s", cb.State())
	}
}

func TestCircuitBreaker_OpensAfterMaxFailures(t *testing.T) {
	cb, _ := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 2, OpenDuration: time.Minute})

	for i := 0; i < 2; i++ {
		_ = cb.Do(func() error { return errFake })
	}

	if cb.State() != "open" {
		t.Fatalf("expected open, got %s", cb.State())
	}
}

func TestCircuitBreaker_ReturnsErrCircuitOpen(t *testing.T) {
	cb, _ := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 1, OpenDuration: time.Minute})
	_ = cb.Do(func() error { return errFake })

	err := cb.Do(func() error { return nil })
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_ResetsOnSuccess(t *testing.T) {
	cb, _ := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 3, OpenDuration: time.Minute})

	_ = cb.Do(func() error { return errFake })
	_ = cb.Do(func() error { return errFake })
	_ = cb.Do(func() error { return nil }) // success resets

	if cb.State() != "closed" {
		t.Fatalf("expected closed after success, got %s", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenAfterDuration(t *testing.T) {
	cb, _ := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 1, OpenDuration: 10 * time.Millisecond})
	_ = cb.Do(func() error { return errFake })

	time.Sleep(20 * time.Millisecond)

	// Next call should be allowed (half-open probe)
	err := cb.Do(func() error { return nil })
	if err != nil {
		t.Fatalf("expected nil after open duration elapsed, got %v", err)
	}
	if cb.State() != "closed" {
		t.Fatalf("expected closed after successful probe, got %s", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	cb, _ := NewCircuitBreaker(CircuitBreakerConfig{MaxFailures: 1, OpenDuration: 10 * time.Millisecond})
	_ = cb.Do(func() error { return errFake })

	time.Sleep(20 * time.Millisecond)

	_ = cb.Do(func() error { return errFake })
	if cb.State() != "open" {
		t.Fatalf("expected open after half-open failure, got %s", cb.State())
	}
}
