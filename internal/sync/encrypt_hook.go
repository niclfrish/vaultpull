package sync

import (
	"fmt"
	"io"
	"os"
)

// EncryptHookResult holds the result of an encrypt/decrypt hook operation.
type EncryptHookResult struct {
	Processed int
	Passphrase string
}

// EncryptAndWrite encrypts secrets using the given passphrase and writes them
// via the provided writer function. Reports a summary to w (defaults to os.Stdout).
func EncryptAndWrite(
	secrets map[string]string,
	passphrase string,
	writeFn func(map[string]string) error,
	w io.Writer,
) error {
	if w == nil {
		w = os.Stdout
	}
	enc, err := NewEncryptor(passphrase)
	if err != nil {
		return fmt.Errorf("encrypt hook: %w", err)
	}
	encrypted, err := enc.EncryptSecrets(secrets)
	if err != nil {
		return fmt.Errorf("encrypt hook: %w", err)
	}
	if err := writeFn(encrypted); err != nil {
		return fmt.Errorf("encrypt hook: write failed: %w", err)
	}
	fmt.Fprintf(w, "[encrypt] %d secret(s) encrypted and written\n", len(encrypted))
	return nil
}

// DecryptAndReturn decrypts secrets using the given passphrase and returns the
// plaintext map. Reports a summary to w (defaults to os.Stdout).
func DecryptAndReturn(
	secrets map[string]string,
	passphrase string,
	w io.Writer,
) (map[string]string, error) {
	if w == nil {
		w = os.Stdout
	}
	enc, err := NewEncryptor(passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt hook: %w", err)
	}
	decrypted, err := enc.DecryptSecrets(secrets)
	if err != nil {
		return nil, fmt.Errorf("decrypt hook: %w", err)
	}
	fmt.Fprintf(w, "[decrypt] %d secret(s) decrypted\n", len(decrypted))
	return decrypted, nil
}
