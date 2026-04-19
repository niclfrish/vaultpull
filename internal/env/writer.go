package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
}

// New creates a new Writer for the given file path.
func New(filePath string) *Writer {
	return &Writer{filePath: filePath}
}

// Write writes the provided secrets map to the .env file.
// Existing file will be overwritten.
func (w *Writer) Write(secrets map[string]string) error {
	f, err := os.Create(w.filePath)
	if err != nil {
		return fmt.Errorf("env: failed to create file %q: %w", w.filePath, err)
	}
	defer f.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		line := fmt.Sprintf("%s=%s\n", sanitizeKey(k), escapeValue(secrets[k]))
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("env: failed to write key %q: %w", k, err)
		}
	}
	return nil
}

// sanitizeKey uppercases and replaces hyphens with underscores.
func sanitizeKey(k string) string {
	return strings.ToUpper(strings.ReplaceAll(k, "-", "_"))
}

// escapeValue wraps values containing spaces or special chars in double quotes.
func escapeValue(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
