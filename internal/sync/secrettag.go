package sync

import (
	"fmt"
	"strings"
	"time"
)

// SecretTagConfig holds configuration for secret tagging.
type SecretTagConfig struct {
	Prefix    string
	Timestamp bool
	Source    string
}

// DefaultSecretTagConfig returns a default configuration.
func DefaultSecretTagConfig() SecretTagConfig {
	return SecretTagConfig{
		Prefix:    "__meta",
		Timestamp: true,
		Source:    "vault",
	}
}

// TagSecrets injects metadata tag keys into the secrets map.
// Tags are added as special keys using the configured prefix.
func TagSecrets(secrets map[string]string, cfg SecretTagConfig) (map[string]string, error) {
	if secrets == nil {
		return nil, fmt.Errorf("tagSecrets: secrets map is nil")
	}
	if cfg.Prefix == "" {
		return nil, fmt.Errorf("tagSecrets: prefix must not be empty")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	if cfg.Source != "" {
		key := fmt.Sprintf("%s_source", cfg.Prefix)
		out[key] = cfg.Source
	}

	if cfg.Timestamp {
		key := fmt.Sprintf("%s_synced_at", cfg.Prefix)
		out[key] = time.Now().UTC().Format(time.RFC3339)
	}

	out[fmt.Sprintf("%s_count", cfg.Prefix)] = fmt.Sprintf("%d", len(secrets))

	return out, nil
}

// StripTagKeys removes all injected tag keys from the secrets map.
func StripTagKeys(secrets map[string]string, prefix string) map[string]string {
	if secrets == nil {
		return nil
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if !strings.HasPrefix(k, prefix+"_") {
			out[k] = v
		}
	}
	return out
}
