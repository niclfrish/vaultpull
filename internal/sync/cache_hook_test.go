package sync

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

func newTestCachedFetcher(t *testing.T, ttl time.Duration) (*CachedFetcher, *SecretCache) {
	t.Helper()
	cache, _ := NewSecretCache(t.TempDir())
	cf := NewCachedFetcher(cache, ttl)
	cf.out = &bytes.Buffer{}
	return cf, cache
}

func TestCachedFetcher_CallsFetchFnOnMiss(t *testing.T) {
	cf, _ := newTestCachedFetcher(t, time.Minute)
	called := 0
	fetchFn := func() (map[string]string, error) {
		called++
		return map[string]string{"KEY": "val"}, nil
	}
	secrets, err := cf.Fetch("secret/app", "", fetchFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Errorf("expected fetchFn called once, got %d", called)
	}
	if secrets["KEY"] != "val" {
		t.Errorf("expected KEY=val")
	}
}

func TestCachedFetcher_UsesCacheOnHit(t *testing.T) {
	cf, cache := newTestCachedFetcher(t, time.Minute)
	_ = cache.Put("secret/app", "", map[string]string{"CACHED": "yes"})

	called := 0
	secrets, err := cf.Fetch("secret/app", "", func() (map[string]string, error) {
		called++
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 0 {
		t.Error("expected fetchFn not to be called on cache hit")
	}
	if secrets["CACHED"] != "yes" {
		t.Errorf("expected CACHED=yes from cache")
	}
}

func TestCachedFetcher_ExpiredCacheRefetches(t *testing.T) {
	cf, cache := newTestCachedFetcher(t, time.Nanosecond)
	_ = cache.Put("secret/app", "", map[string]string{"OLD": "data"})
	time.Sleep(2 * time.Millisecond)

	called := 0
	_, err := cf.Fetch("secret/app", "", func() (map[string]string, error) {
		called++
		return map[string]string{"NEW": "data"}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Error("expected fetchFn to be called after TTL expiry")
	}
}

func TestCachedFetcher_FallsBackToStaleOnError(t *testing.T) {
	cf, cache := newTestCachedFetcher(t, time.Nanosecond)
	_ = cache.Put("secret/app", "", map[string]string{"STALE": "ok"})
	time.Sleep(2 * time.Millisecond)

	secrets, err := cf.Fetch("secret/app", "", func() (map[string]string, error) {
		return nil, errors.New("vault unreachable")
	})
	if err != nil {
		t.Fatalf("expected stale fallback, got error: %v", err)
	}
	if secrets["STALE"] != "ok" {
		t.Error("expected stale cache value")
	}
}

func TestCachedFetcher_NoFallback_ReturnsError(t *testing.T) {
	cf, _ := newTestCachedFetcher(t, time.Minute)
	_, err := cf.Fetch("secret/app", "", func() (map[string]string, error) {
		return nil, errors.New("vault down")
	})
	if err == nil {
		t.Fatal("expected error when no cache and fetch fails")
	}
}
