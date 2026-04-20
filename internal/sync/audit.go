package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AuditEntry records a single sync operation result.
type AuditEntry struct {
	Timestamp  time.Time         `json:"timestamp"`
	EnvFile    string            `json:"env_file"`
	Namespace  string            `json:"namespace,omitempty"`
	SecretPath string            `json:"secret_path"`
	Changes    map[string]string `json:"changes"`
	DryRun     bool              `json:"dry_run"`
}

// AuditLogger writes audit entries to a JSON log file.
type AuditLogger struct {
	path string
}

// NewAuditLogger creates an AuditLogger that appends to the given file path.
func NewAuditLogger(path string) *AuditLogger {
	return &AuditLogger{path: path}
}

// Log appends an AuditEntry to the audit log file.
func (a *AuditLogger) Log(entry AuditEntry) error {
	f, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// BuildAuditEntry constructs an AuditEntry from a Plan.
func BuildAuditEntry(plan Plan, envFile, secretPath, namespace string, dryRun bool) AuditEntry {
	changes := make(map[string]string, len(plan.Changes))
	for _, c := range plan.Changes {
		changes[c.Key] = string(c.Action)
	}
	return AuditEntry{
		Timestamp:  time.Now().UTC(),
		EnvFile:    envFile,
		Namespace:  namespace,
		SecretPath: secretPath,
		Changes:    changes,
		DryRun:     dryRun,
	}
}
