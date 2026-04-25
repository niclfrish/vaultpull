package sync

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SecretVersion holds metadata about a specific version of a secret.
type SecretVersion struct {
	Version   int
	CreatedAt time.Time
	DeletedAt *time.Time
	Destroyed bool
}

// SecretVersionMap maps secret keys to their version metadata.
type SecretVersionMap map[string]SecretVersion

// ParseVersionHeader parses a Vault KV v2 version string like "version=3" into an int.
func ParseVersionHeader(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("secretversion: empty version string")
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("secretversion: invalid version %q: %w", s, err)
	}
	if v < 1 {
		return 0, fmt.Errorf("secretversion: version must be >= 1, got %d", v)
	}
	return v, nil
}

// VersionSummary returns a human-readable summary of a SecretVersion.
func VersionSummary(sv SecretVersion) string {
	status := "active"
	if sv.Destroyed {
		status = "destroyed"
	} else if sv.DeletedAt != nil {
		status = "deleted"
	}
	return fmt.Sprintf("v%d created=%s status=%s",
		sv.Version,
		sv.CreatedAt.UTC().Format(time.RFC3339),
		status,
	)
}

// FilterActiveVersions returns only the versions that are neither deleted nor destroyed.
func FilterActiveVersions(versions SecretVersionMap) SecretVersionMap {
	out := make(SecretVersionMap, len(versions))
	for k, v := range versions {
		if !v.Destroyed && v.DeletedAt == nil {
			out[k] = v
		}
	}
	return out
}
