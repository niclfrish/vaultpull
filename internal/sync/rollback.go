package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Backup holds the path to a backup file and the original destination.
type Backup struct {
	OriginalPath string
	BackupPath   string
}

// CreateBackup copies the current .env file to a timestamped backup file.
// If the original file does not exist, no backup is created and nil is returned.
func CreateBackup(envPath string) (*Backup, error) {
	data, err := os.ReadFile(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("rollback: read original file: %w", err)
	}

	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backupPath := filepath.Join(dir, fmt.Sprintf("%s.%s.bak", base, timestamp))

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return nil, fmt.Errorf("rollback: write backup file: %w", err)
	}

	return &Backup{
		OriginalPath: envPath,
		BackupPath:   backupPath,
	}, nil
}

// Restore copies the backup file back to the original path and removes the backup.
func (b *Backup) Restore() error {
	if b == nil {
		return nil
	}

	data, err := os.ReadFile(b.BackupPath)
	if err != nil {
		return fmt.Errorf("rollback: read backup: %w", err)
	}

	if err := os.WriteFile(b.OriginalPath, data, 0600); err != nil {
		return fmt.Errorf("rollback: restore file: %w", err)
	}

	_ = os.Remove(b.BackupPath)
	return nil
}

// Discard removes the backup file without restoring it.
func (b *Backup) Discard() error {
	if b == nil {
		return nil
	}
	return os.Remove(b.BackupPath)
}
