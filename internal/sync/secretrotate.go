package sync

import (
	"fmt"
	"time"
)

// RotateConfig controls secret rotation behaviour.
type RotateConfig struct {
	// RotateAfter is the age threshold after which a secret is considered stale.
	RotateAfter time.Duration
	// RotatedAtKey is the metadata key that stores the last rotation timestamp.
	RotatedAtKey string
}

// DefaultRotateConfig returns sensible defaults.
func DefaultRotateConfig() RotateConfig {
	return RotateConfig{
		RotateAfter:  24 * time.Hour,
		RotatedAtKey: "__rotated_at",
	}
}

// RotationStatus describes whether a secret needs rotation.
type RotationStatus int

const (
	RotationOK      RotationStatus = iota // within threshold
	RotationStale                         // past threshold, should rotate
	RotationUnknown                       // no timestamp present
)

// String returns a human-readable label for the status.
func (s RotationStatus) String() string {
	switch s {
	case RotationOK:
		return "ok"
	case RotationStale:
		return "stale"
	default:
		return "unknown"
	}
}

// CheckRotation inspects secrets for a rotation timestamp and classifies each
// entry as OK, stale, or unknown relative to cfg.RotateAfter.
func CheckRotation(secrets map[string]string, cfg RotateConfig) map[string]RotationStatus {
	result := make(map[string]RotationStatus, len(secrets))
	for k, v := range secrets {
		if k == cfg.RotatedAtKey {
			continue
		}
		ts, ok := secrets[cfg.RotatedAtKey+"."+k]
		if !ok {
			result[k] = RotationUnknown
			continue
		}
		_ = v
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			result[k] = RotationUnknown
			continue
		}
		if time.Since(t) > cfg.RotateAfter {
			result[k] = RotationStale
		} else {
			result[k] = RotationOK
		}
	}
	return result
}

// RotateSummary returns a human-readable summary of rotation statuses.
func RotateSummary(statuses map[string]RotationStatus) string {
	ok, stale, unknown := 0, 0, 0
	for _, s := range statuses {
		switch s {
		case RotationOK:
			ok++
		case RotationStale:
			stale++
		default:
			unknown++
		}
	}
	return fmt.Sprintf("rotation: ok=%d stale=%d unknown=%d", ok, stale, unknown)
}
