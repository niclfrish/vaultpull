package sync

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ExpiryInfo holds parsed expiry metadata for a secret.
type ExpiryInfo struct {
	Key       string
	ExpiresAt time.Time
	TTL       time.Duration
	Expired   bool
}

// ParseExpiryHeader parses a "expires_at" value from secret metadata.
// Expected format: Unix timestamp (seconds) or RFC3339 string.
func ParseExpiryHeader(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty expiry header")
	}
	if ts, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return time.Unix(ts, 0).UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("unrecognised expiry format %q: %w", raw, err)
	}
	return t.UTC(), nil
}

// ClassifyExpiry returns ExpiryInfo for each secret key that carries an
// "__expires_at__" suffix encoding its expiry timestamp.
func ClassifyExpiry(secrets map[string]string, now time.Time) []ExpiryInfo {
	const suffix = "__expires_at__"
	var result []ExpiryInfo
	for k, v := range secrets {
		if !strings.HasSuffix(k, suffix) {
			continue
		}
		baseKey := strings.TrimSuffix(k, suffix)
		at, err := ParseExpiryHeader(v)
		if err != nil {
			continue
		}
		ttl := at.Sub(now)
		result = append(result, ExpiryInfo{
			Key:       baseKey,
			ExpiresAt: at,
			TTL:       ttl,
			Expired:   ttl <= 0,
		})
	}
	return result
}

// ExpirySummary returns a human-readable summary line.
func ExpirySummary(infos []ExpiryInfo) string {
	if len(infos) == 0 {
		return "no expiry metadata found"
	}
	expired := 0
	for _, i := range infos {
		if i.Expired {
			expired++
		}
	}
	return fmt.Sprintf("%d secret(s) with expiry metadata; %d expired", len(infos), expired)
}
