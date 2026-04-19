package sync

import (
	"fmt"
	"strings"
)

// NamespacedPath builds a Vault KV v2 secret path with an optional namespace prefix.
// If namespace is empty, path is returned as-is.
func NamespacedPath(namespace, path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("namespace: path must not be empty")
	}

	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		return path, nil
	}

	namespace = strings.Trim(namespace, "/")
	path = strings.Trim(path, "/")

	return namespace + "/" + path, nil
}

// PrefixKeys returns a new map with all keys prefixed by the given namespace
// in UPPER_SNAKE_CASE style, e.g. namespace "app" + key "db_pass" -> "APP_DB_PASS".
func PrefixKeys(namespace string, secrets map[string]string) map[string]string {
	if namespace == "" {
		return secrets
	}

	prefix := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(namespace), "-", "_")) + "_"
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[prefix+k] = v
	}
	return result
}
