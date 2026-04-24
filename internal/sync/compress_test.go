package sync

import (
	"strings"
	"testing"
)

func TestCompressValue_ShortValue_Unchanged(t *testing.T) {
	cfg := DefaultCompressConfig()
	input := "short"
	out, err := CompressValue(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected unchanged value %q, got %q", input, out)
	}
}

func TestCompressValue_LongValue_Compressed(t *testing.T) {
	cfg := DefaultCompressConfig()
	input := strings.Repeat("secret-data-", 20) // well over 64 chars
	out, err := CompressValue(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == input {
		t.Error("expected value to be compressed, but got original")
	}
}

func TestCompressDecompress_RoundTrip(t *testing.T) {
	cfg := DefaultCompressConfig()
	original := strings.Repeat("vault-secret-value=", 10)

	compressed, err := CompressValue(original, cfg)
	if err != nil {
		t.Fatalf("compress error: %v", err)
	}

	decompressed, err := DecompressValue(compressed)
	if err != nil {
		t.Fatalf("decompress error: %v", err)
	}

	if decompressed != original {
		t.Errorf("round-trip mismatch: got %q, want %q", decompressed, original)
	}
}

func TestDecompressValue_PlainText_ReturnedAsIs(t *testing.T) {
	input := "plain-text-value"
	out, err := DecompressValue(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected %q, got %q", input, out)
	}
}

func TestCompressSecrets_AllValuesCompressed(t *testing.T) {
	cfg := CompressConfig{MinLength: 5}
	secrets := map[string]string{
		"KEY_A": "short",
		"KEY_B": strings.Repeat("x", 50),
	}

	result, err := CompressSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result["KEY_A"]; !ok {
		t.Error("KEY_A missing from result")
	}
	if result["KEY_B"] == secrets["KEY_B"] {
		t.Error("KEY_B should have been compressed")
	}

	// Verify round-trip for KEY_B
	decompressed, err := DecompressValue(result["KEY_B"])
	if err != nil {
		t.Fatalf("decompress error: %v", err)
	}
	if decompressed != secrets["KEY_B"] {
		t.Errorf("round-trip failed for KEY_B: got %q", decompressed)
	}
}

func TestDefaultCompressConfig(t *testing.T) {
	cfg := DefaultCompressConfig()
	if cfg.MinLength <= 0 {
		t.Errorf("expected positive MinLength, got %d", cfg.MinLength)
	}
}
