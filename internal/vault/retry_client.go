package vault

import (
	"fmt"
	"strings"

	isync "github.com/yourusername/vaultpull/internal/sync"
)

// RetryClient wraps a Client and retries transient Vault errors.
type RetryClient struct {
	client  *Client
	retrier *isync.Retrier
}

// NewRetryClient creates a RetryClient using the provided Client and RetryConfig.
func NewRetryClient(c *Client, cfg isync.RetryConfig) *RetryClient {
	return &RetryClient{
		client:  c,
		retrier: isync.NewRetrier(cfg),
	}
}

// GetSecrets fetches secrets from Vault, retrying on transient errors.
func (rc *RetryClient) GetSecrets(path string) (map[string]string, error) {
	var result map[string]string

	err := rc.retrier.Run(func() error {
		secrets, err := rc.client.GetSecrets(path)
		if err != nil {
			if isTransient(err) {
				return err
			}
			return isync.NonRetryable(err)
		}
		result = secrets
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("vault get secrets: %w", err)
	}
	return result, nil
}

// isTransient returns true for errors that are worth retrying (5xx, timeouts).
func isTransient(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	transientPhrases := []string{
		"connection refused",
		"timeout",
		"503",
		"502",
		"500",
		"i/o timeout",
		"EOF",
	}
	for _, phrase := range transientPhrases {
		if strings.Contains(msg, phrase) {
			return true
		}
	}
	return false
}
