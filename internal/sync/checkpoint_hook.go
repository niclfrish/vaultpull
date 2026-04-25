package sync

import (
	"crypto/sha256"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// CheckpointAfterSync saves a checkpoint after a successful sync operation.
// It records key count, checksum of the secret values, and sync timestamp.
func CheckpointAfterSync(
	store *CheckpointStore,
	vaultPath, namespace string,
	secrets map[string]string,
	w io.Writer,
) error {
	if store == nil {
		return nil
	}

	cp := Checkpoint{
		Path:      vaultPath,
		Namespace: namespace,
		SyncedAt:  time.Now().UTC(),
		KeyCount:  len(secrets),
		Checksum:  checksumSecrets(secrets),
	}

	if err := store.Save(cp); err != nil {
		return fmt.Errorf("checkpoint_hook: save: %w", err)
	}

	if w != nil {
		fmt.Fprintf(w, "[checkpoint] saved: path=%s namespace=%s keys=%d checksum=%s\n",
			cp.Path, cp.Namespace, cp.KeyCount, cp.Checksum[:8])
	}
	return nil
}

// LoadCheckpoint retrieves a previous checkpoint and writes a summary to w.
// Returns nil without error when no checkpoint exists.
func LoadCheckpoint(store *CheckpointStore, vaultPath, namespace string, w io.Writer) (*Checkpoint, error) {
	if store == nil {
		return nil, nil
	}
	cp, err := store.Load(vaultPath, namespace)
	if err != nil {
		return nil, fmt.Errorf("checkpoint_hook: load: %w", err)
	}
	if cp == nil {
		if w != nil {
			fmt.Fprintf(w, "[checkpoint] no previous checkpoint found for path=%s\n", vaultPath)
		}
		return nil, nil
	}
	if w != nil {
		fmt.Fprintf(w, "[checkpoint] last sync: %s (%d keys)\n",
			cp.SyncedAt.Format(time.RFC3339), cp.KeyCount)
	}
	return cp, nil
}

// checksumSecrets produces a stable SHA-256 hex digest over sorted key=value pairs.
func checksumSecrets(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SecretsChangedSinceCheckpoint returns true when the current secrets differ
// from what was recorded in the checkpoint (by checksum comparison).
func SecretsChangedSinceCheckpoint(cp *Checkpoint, secrets map[string]string) bool {
	if cp == nil {
		return true
	}
	current := checksumSecrets(secrets)
	return !strings.EqualFold(current, cp.Checksum)
}
