package sync

import (
	"errors"
	"testing"
	"time"
)

func TestNewTokenRotator_EmptyInitial(t *testing.T) {
	_, err := NewTokenRotator("", func() (string, error) { return "x", nil }, DefaultTokenRotateConfig())
	if err == nil {
		t.Fatal("expected error for empty initial token")
	}
}

func TestNewTokenRotator_NilFetcher(t *testing.T) {
	_, err := NewTokenRotator("tok", nil, DefaultTokenRotateConfig())
	if err == nil {
		t.Fatal("expected error for nil fetcher")
	}
}

func TestNewTokenRotator_InvalidRotateAfter(t *testing.T) {
	cfg := DefaultTokenRotateConfig()
	cfg.RotateAfter = 0
	_, err := NewTokenRotator("tok", func() (string, error) { return "x", nil }, cfg)
	if err == nil {
		t.Fatal("expected error for zero RotateAfter")
	}
}

func TestTokenRotator_ReturnsInitialToken(t *testing.T) {
	r, err := NewTokenRotator("initial-token", func() (string, error) {
		return "new-token", nil
	}, DefaultTokenRotateConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tok, err := r.Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "initial-token" {
		t.Errorf("expected initial-token, got %q", tok)
	}
}

func TestTokenRotator_RotatesAfterExpiry(t *testing.T) {
	cfg := TokenRotateConfig{RotateAfter: 1 * time.Millisecond, GracePeriod: time.Minute}
	r, _ := NewTokenRotator("old", func() (string, error) {
		return "rotated", nil
	}, cfg)

	time.Sleep(5 * time.Millisecond)

	tok, err := r.Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "rotated" {
		t.Errorf("expected rotated, got %q", tok)
	}
}

func TestTokenRotator_GracePeriodOnFetchError(t *testing.T) {
	cfg := TokenRotateConfig{RotateAfter: 1 * time.Millisecond, GracePeriod: time.Minute}
	r, _ := NewTokenRotator("stale", func() (string, error) {
		return "", errors.New("vault unavailable")
	}, cfg)

	time.Sleep(5 * time.Millisecond)

	tok, err := r.Token()
	if err != nil {
		t.Fatalf("expected stale token within grace period, got error: %v", err)
	}
	if tok != "stale" {
		t.Errorf("expected stale, got %q", tok)
	}
}

func TestTokenRotator_ErrorAfterGracePeriodExpires(t *testing.T) {
	cfg := TokenRotateConfig{RotateAfter: 1 * time.Millisecond, GracePeriod: 1 * time.Millisecond}
	r, _ := NewTokenRotator("stale", func() (string, error) {
		return "", errors.New("vault down")
	}, cfg)

	time.Sleep(10 * time.Millisecond)

	_, err := r.Token()
	if err == nil {
		t.Fatal("expected error after grace period expired")
	}
}

func TestTokenRotator_Age(t *testing.T) {
	r, _ := NewTokenRotator("tok", func() (string, error) { return "x", nil }, DefaultTokenRotateConfig())
	time.Sleep(5 * time.Millisecond)
	if r.Age() < 5*time.Millisecond {
		t.Error("expected age to be at least 5ms")
	}
}
