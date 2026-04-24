package sync

import (
	"fmt"
	"sort"
)

// EnvTarget represents a named output target with its own path and optional namespace.
type EnvTarget struct {
	Name      string
	Path      string
	Namespace string
}

// MultiEnvWriter writes secrets to multiple .env targets simultaneously.
type MultiEnvWriter struct {
	Targets []EnvTarget
	writer  EnvWriterFunc
}

// EnvWriterFunc is a function that writes a map of secrets to a file path.
type EnvWriterFunc func(path string, secrets map[string]string) error

// NewMultiEnvWriter creates a MultiEnvWriter with the given targets and writer function.
func NewMultiEnvWriter(targets []EnvTarget, writer EnvWriterFunc) (*MultiEnvWriter, error) {
	if len(targets) == 0 {
		return nil, fmt.Errorf("multienv: at least one target is required")
	}
	if writer == nil {
		return nil, fmt.Errorf("multienv: writer function must not be nil")
	}
	return &MultiEnvWriter{Targets: targets, writer: writer}, nil
}

// WriteAll writes the given secrets to all registered targets, applying namespace
// prefixing per target. Returns a map of target name to error (nil on success).
func (m *MultiEnvWriter) WriteAll(secrets map[string]string) map[string]error {
	results := make(map[string]error, len(m.Targets))
	for _, t := range m.Targets {
		keyed := PrefixKeys(secrets, t.Namespace)
		results[t.Name] = m.writer(t.Path, keyed)
	}
	return results
}

// TargetNames returns sorted target names registered in this writer.
func (m *MultiEnvWriter) TargetNames() []string {
	names := make([]string, 0, len(m.Targets))
	for _, t := range m.Targets {
		names = append(names, t.Name)
	}
	sort.Strings(names)
	return names
}

// AnyError returns the first non-nil error from a WriteAll result map, or nil.
func AnyError(results map[string]error) error {
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if results[k] != nil {
			return fmt.Errorf("target %q: %w", k, results[k])
		}
	}
	return nil
}
