package sync

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

// Syncer orchestrates pulling secrets from Vault and writing them to .env files.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
	writer *env.Writer
}

// New creates a new Syncer from the given config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("syncer: init vault client: %w", err)
	}

	writer, err := env.New(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("syncer: init env writer: %w", err)
	}

	return &Syncer{cfg: cfg, client: client, writer: writer}, nil
}

// Run fetches secrets from Vault and writes them to the configured output file.
func (s *Syncer) Run() error {
	secrets, err := s.client.GetSecrets(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("syncer: get secrets: %w", err)
	}

	if len(secrets) == 0 {
		return fmt.Errorf("syncer: no secrets found at path %q", s.cfg.SecretPath)
	}

	if err := s.writer.Write(secrets); err != nil {
		return fmt.Errorf("syncer: write env file: %w", err)
	}

	return nil
}
