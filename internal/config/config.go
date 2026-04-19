package config

import (
	"errors"
	"os"
)

// Config holds all runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	SecretPath string
	OutputFile string
	Namespace  string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	token := getEnv("VAULT_TOKEN", "")
	if token == "" {
		return nil, errors.New("config: VAULT_TOKEN is required")
	}

	secretPath := getEnv("VAULT_SECRET_PATH", "")
	if secretPath == "" {
		return nil, errors.New("config: VAULT_SECRET_PATH is required")
	}

	return &Config{
		VaultAddr:  getEnv("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken: token,
		SecretPath: secretPath,
		OutputFile: getEnv("VAULTPULL_OUTPUT", ".env"),
		Namespace:  getEnv("VAULT_NAMESPACE", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
