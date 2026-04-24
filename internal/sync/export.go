package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// ExportFormat defines the output format for secret export.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatDotenv ExportFormat = "dotenv"
	FormatExport ExportFormat = "export"
)

// Exporter writes secrets to a writer in a specified format.
type Exporter struct {
	format ExportFormat
	out    io.Writer
}

// NewExporter creates an Exporter that writes to out in the given format.
func NewExporter(format ExportFormat, out io.Writer) (*Exporter, error) {
	if out == nil {
		out = os.Stdout
	}
	switch format {
	case FormatJSON, FormatDotenv, FormatExport:
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
	return &Exporter{format: format, out: out}, nil
}

// Export writes secrets to the configured writer.
func (e *Exporter) Export(secrets map[string]string) error {
	if secrets == nil {
		secrets = map[string]string{}
	}
	switch e.format {
	case FormatJSON:
		return e.writeJSON(secrets)
	case FormatDotenv:
		return e.writeDotenv(secrets, false)
	case FormatExport:
		return e.writeDotenv(secrets, true)
	}
	return nil
}

func (e *Exporter) writeJSON(secrets map[string]string) error {
	enc := json.NewEncoder(e.out)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}

func (e *Exporter) writeDotenv(secrets map[string]string, withExport bool) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	prefix := ""
	if withExport {
		prefix = "export "
	}
	for _, k := range keys {
		_, err := fmt.Fprintf(e.out, "%s%s=%q\n", prefix, k, secrets[k])
		if err != nil {
			return fmt.Errorf("export write: %w", err)
		}
	}
	return nil
}
