package config

import (
	"errors"
	"os"
)

// Config holds the runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	Namespace  string
	OutputFile string
	SecretPath string
}

// Load reads configuration from environment variables and applies overrides.
func Load(overrides ...func(*Config)) (*Config, error) {
	cfg := &Config{
		VaultAddr:  getEnv("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		Namespace:  os.Getenv("VAULT_NAMESPACE"),
		OutputFile: ".env",
	}

	for _, override := range overrides {
		override(cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.VaultToken == "" {
		return errors.New("VAULT_TOKEN is required but not set")
	}
	if c.SecretPath == "" {
		return errors.New("secret path is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
