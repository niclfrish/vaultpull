package vault

import (
	"errors"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with project-specific helpers.
type Client struct {
	api       *vaultapi.Client
	namespace string
}

// New creates a new Vault client using the provided address, token, and optional namespace.
func New(address, token, namespace string) (*Client, error) {
	if address == "" {
		return nil, errors.New("vault address must not be empty")
	}
	if token == "" {
		return nil, errors.New("vault token must not be empty")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault api client: %w", err)
	}

	client.SetToken(token)

	if namespace != "" {
		client.SetNamespace(namespace)
	}

	return &Client{api: client, namespace: namespace}, nil
}

// GetSecrets reads a KV v2 secret at the given path and returns the key/value data.
func (c *Client) GetSecrets(secretPath string) (map[string]string, error) {
	if secretPath == "" {
		return nil, errors.New("secret path must not be empty")
	}

	// Normalise path: strip leading slash
	secretPath = strings.TrimPrefix(secretPath, "/")

	secret, err := c.api.Logical().Read(secretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", secretPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", secretPath)
	}

	// KV v2 nests data under secret.Data["data"]
	raw, ok := secret.Data["data"]
	if !ok {
		// KV v1 — data is at the top level
		raw = secret.Data
	}

	dataMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at path %q", secretPath)
	}

	result := make(map[string]string, len(dataMap))
	for k, v := range dataMap {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}

// ListSecrets returns the keys available under the given path using the
// Vault LIST operation. It is compatible with both KV v1 and KV v2 mounts.
func (c *Client) ListSecrets(secretPath string) ([]string, error) {
	if secretPath == "" {
		return nil, errors.New("secret path must not be empty")
	}

	secretPath = strings.TrimPrefix(secretPath, "/")

	secret, err := c.api.Logical().List(secretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets at %q: %w", secretPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secrets found at path %q", secretPath)
	}

	rawKeys, ok := secret.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("no keys returned at path %q", secretPath)
	}

	ifaces, ok := rawKeys.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected keys format at path %q", secretPath)
	}

	keys := make([]string, len(ifaces))
	for i, v := range ifaces {
		keys[i] = fmt.Sprintf("%v", v)
	}
	return keys, nil
}
