package sync

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// TokenRotateHook wraps a secret-fetching function and injects a rotated token
// into the secrets map under a configurable key before writing.

// RotateTokenAndInject fetches the current token from the rotator, injects it
// into secrets under tokenKey, and returns the updated map.
// If tokenKey is empty it defaults to "VAULT_TOKEN".
func RotateTokenAndInject(rotator *TokenRotator, tokenKey string, secrets map[string]string) (map[string]string, error) {
	if rotator == nil {
		return nil, errors.New("tokenrotate_hook: rotator must not be nil")
	}
	if secrets == nil {
		return nil, errors.New("tokenrotate_hook: secrets map must not be nil")
	}
	if tokenKey == "" {
		tokenKey = "VAULT_TOKEN"
	}

	tok, err := rotator.Token()
	if err != nil {
		return nil, fmt.Errorf("tokenrotate_hook: %w", err)
	}

	out := make(map[string]string, len(secrets)+1)
	for k, v := range secrets {
		out[k] = v
	}
	out[tokenKey] = tok
	return out, nil
}

// LogTokenAge writes the current token age to w (defaults to os.Stdout).
func LogTokenAge(rotator *TokenRotator, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	if rotator == nil {
		fmt.Fprintln(w, "tokenrotate_hook: rotator is nil, cannot log age")
		return
	}
	fmt.Fprintf(w, "tokenrotate_hook: current token age: %s\n", rotator.Age().Round(1000000))
}
