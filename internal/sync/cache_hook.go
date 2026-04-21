package sync

import (
	"fmt"
	"io"
	"os"
	"time"
)

// CachedFetcher wraps a secret-fetching function with cache read/write logic.
type CachedFetcher struct {
	cache *SecretCache
	ttl   time.Duration
	out   io.Writer
}

// NewCachedFetcher returns a CachedFetcher using the provided cache and TTL.
func NewCachedFetcher(cache *SecretCache, ttl time.Duration) *CachedFetcher {
	return &CachedFetcher{cache: cache, ttl: ttl, out: os.Stdout}
}

// Fetch returns cached secrets if fresh, otherwise calls fetchFn and caches the result.
func (cf *CachedFetcher) Fetch(
	path, namespace string,
	fetchFn func() (map[string]string, error),
) (map[string]string, error) {
	entry, err := cf.cache.Get(path, namespace)
	if err != nil {
		return nil, fmt.Errorf("cache lookup: %w", err)
	}

	if entry != nil && time.Since(entry.FetchedAt) < cf.ttl {
		fmt.Fprintf(cf.out, "cache: using cached secrets for %q (age: %s)\n",
			path, time.Since(entry.FetchedAt).Round(time.Second))
		return entry.Secrets, nil
	}

	secrets, err := fetchFn()
	if err != nil {
		if entry != nil {
			fmt.Fprintf(cf.out, "cache: fetch failed, falling back to stale cache for %q\n", path)
			return entry.Secrets, nil
		}
		return nil, err
	}

	if putErr := cf.cache.Put(path, namespace, secrets); putErr != nil {
		fmt.Fprintf(cf.out, "cache: warning: failed to store cache: %v\n", putErr)
	}
	return secrets, nil
}
