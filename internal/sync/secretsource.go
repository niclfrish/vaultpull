package sync

import (
	"fmt"
	"strings"
	"time"
)

// SourceType identifies the origin of a secret.
type SourceType string

const (
	SourceTypeVault SourceType = "vault"
	SourceTypeEnv   SourceType = "env"
	SourceTypeFile  SourceType = "file"
	SourceTypeUnknown SourceType = "unknown"
)

// SecretSource holds provenance metadata for a secret map.
type SecretSource struct {
	Type      SourceType
	Location  string
	FetchedAt time.Time
	Namespace string
}

// DefaultSecretSourceConfig returns a SecretSource with sensible defaults.
func DefaultSecretSourceConfig(sourceType SourceType, location string) SecretSource {
	return SecretSource{
		Type:      sourceType,
		Location:  location,
		FetchedAt: time.Now().UTC(),
	}
}

// AnnotateWithSource injects source metadata keys into the secrets map.
// Keys are prefixed with "__source_" to avoid collisions.
func AnnotateWithSource(secrets map[string]string, src SecretSource) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets)+4)
	for k, v := range secrets {
		out[k] = v
	}
	out["__source_type"] = string(src.Type)
	out["__source_location"] = src.Location
	out["__source_fetched_at"] = src.FetchedAt.Format(time.RFC3339)
	if src.Namespace != "" {
		out["__source_namespace"] = src.Namespace
	}
	return out
}

// StripSourceAnnotations removes all "__source_" prefixed keys from secrets.
func StripSourceAnnotations(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !strings.HasPrefix(k, "__source_") {
			out[k] = v
		}
	}
	return out
}

// SourceSummary returns a human-readable summary of the source.
func SourceSummary(src SecretSource) string {
	ns := src.Namespace
	if ns == "" {
		ns = "(none)"
	}
	return fmt.Sprintf("type=%s location=%s namespace=%s fetched_at=%s",
		src.Type, src.Location, ns, src.FetchedAt.Format(time.RFC3339))
}
