package sync

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

// CompressConfig holds configuration for secret value compression.
type CompressConfig struct {
	// MinLength is the minimum value length before compression is applied.
	MinLength int
}

// DefaultCompressConfig returns a sensible default compression config.
func DefaultCompressConfig() CompressConfig {
	return CompressConfig{
		MinLength: 64,
	}
}

// CompressValue gzip-compresses a string value and returns a base64-encoded result.
// Values shorter than cfg.MinLength are returned unchanged.
func CompressValue(value string, cfg CompressConfig) (string, error) {
	if len(value) < cfg.MinLength {
		return value, nil
	}

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := io.WriteString(w, value); err != nil {
		return "", fmt.Errorf("compress: write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("compress: close failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// DecompressValue decodes a base64+gzip-compressed value.
// If decoding or decompression fails, the original value is returned unchanged.
func DecompressValue(value string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		// Not a compressed value — return as-is.
		return value, nil
	}

	r, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return value, nil
	}
	defer r.Close()

	out, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("decompress: read failed: %w", err)
	}

	return string(out), nil
}

// CompressSecrets applies CompressValue to every value in the secrets map.
func CompressSecrets(secrets map[string]string, cfg CompressConfig) (map[string]string, error) {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		compressed, err := CompressValue(v, cfg)
		if err != nil {
			return nil, fmt.Errorf("compress: key %q: %w", k, err)
		}
		result[k] = compressed
	}
	return result, nil
}
