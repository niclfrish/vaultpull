package sync

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LeaseInfo holds parsed TTL and renewable metadata for a Vault secret lease.
type LeaseInfo struct {
	LeaseID   string
	Duration  time.Duration
	Renewable bool
	ExpiresAt time.Time
}

// DefaultLeaseTTL is the fallback TTL when none is provided.
const DefaultLeaseTTL = 30 * time.Minute

// ParseLeaseHeader parses a lease descriptor string of the form:
//
//	"lease_id=<id>,ttl=<seconds>,renewable=<bool>"
//
// All fields are optional; missing fields fall back to defaults.
func ParseLeaseHeader(raw string) (*LeaseInfo, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("secretlease: empty lease header")
	}

	info := &LeaseInfo{
		Duration:  DefaultLeaseTTL,
		Renewable: false,
	}

	for _, part := range strings.Split(raw, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		switch key {
		case "lease_id":
			info.LeaseID = val
		case "ttl":
			secs, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("secretlease: invalid ttl %q: %w", val, err)
			}
			if secs <= 0 {
				return nil, fmt.Errorf("secretlease: ttl must be positive, got %d", secs)
			}
			info.Duration = time.Duration(secs) * time.Second
		case "renewable":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("secretlease: invalid renewable %q: %w", val, err)
			}
			info.Renewable = b
		}
	}

	info.ExpiresAt = time.Now().Add(info.Duration)
	return info, nil
}

// IsExpired reports whether the lease has passed its expiry time.
func (l *LeaseInfo) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

// LeaseSummary returns a human-readable one-liner for the lease.
func LeaseSummary(l *LeaseInfo) string {
	if l == nil {
		return "no lease"
	}
	renewStr := "non-renewable"
	if l.Renewable {
		renewStr = "renewable"
	}
	return fmt.Sprintf("lease=%s ttl=%s %s expires=%s",
		l.LeaseID, l.Duration, renewStr, l.ExpiresAt.Format(time.RFC3339))
}
