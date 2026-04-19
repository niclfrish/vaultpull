package config

import (
	"os"
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	_, err := Load(func(c *Config) {
		c.SecretPath = "secret/data/myapp"
	})
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN")
	}
}

func TestLoad_MissingSecretPath(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test-token")
	defer os.Unsetenv("VAULT_TOKEN")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing secret path")
	}
}

func TestLoad_Success(t *testing.T) {
	os.Setenv("VAULT_TOKEN", "test-token")
	os.Setenv("VAULT_NAMESPACE", "dev")
	defer os.Unsetenv("VAULT_TOKEN")
	defer os.Unsetenv("VAULT_NAMESPACE")

	cfg, err := Load(func(c *Config) {
		c.SecretPath = "secret/data/myapp"
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Namespace != "dev" {
		t.Errorf("expected namespace 'dev', got '%s'", cfg.Namespace)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output file '.env', got '%s'", cfg.OutputFile)
	}
}

func TestGetEnv_Fallback(t *testing.T) {
	os.Unsetenv("SOME_MISSING_VAR")
	val := getEnv("SOME_MISSING_VAR", "default")
	if val != "default" {
		t.Errorf("expected 'default', got '%s'", val)
	}
}
