package sync

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncryptAndWrite_Success(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	var written map[string]string
	var buf bytes.Buffer

	err := EncryptAndWrite(secrets, "pass", func(m map[string]string) error {
		written = m
		return nil
	}, &buf)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if written["KEY"] == "value" {
		t.Error("expected encrypted value, got plaintext")
	}
	if !bytes.Contains(buf.Bytes(), []byte("1 secret(s) encrypted")) {
		t.Errorf("expected summary in output, got: %s", buf.String())
	}
}

func TestEncryptAndWrite_EmptyPassphrase(t *testing.T) {
	err := EncryptAndWrite(map[string]string{"K": "v"}, "", func(m map[string]string) error {
		return nil
	}, nil)
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestEncryptAndWrite_WriterError(t *testing.T) {
	wantErr := errors.New("disk full")
	err := EncryptAndWrite(map[string]string{"K": "v"}, "pass", func(m map[string]string) error {
		return wantErr
	}, nil)
	if err == nil {
		t.Fatal("expected error from writeFn")
	}
}

func TestDecryptAndReturn_Success(t *testing.T) {
	enc, _ := NewEncryptor("pass")
	encrypted, _ := enc.EncryptSecrets(map[string]string{"FOO": "bar"})

	var buf bytes.Buffer
	decrypted, err := DecryptAndReturn(encrypted, "pass", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decrypted["FOO"] != "bar" {
		t.Errorf("expected bar, got %q", decrypted["FOO"])
	}
	if !bytes.Contains(buf.Bytes(), []byte("1 secret(s) decrypted")) {
		t.Errorf("expected summary in output, got: %s", buf.String())
	}
}

func TestDecryptAndReturn_WrongPassphrase(t *testing.T) {
	enc, _ := NewEncryptor("correct")
	encrypted, _ := enc.EncryptSecrets(map[string]string{"K": "v"})

	_, err := DecryptAndReturn(encrypted, "wrong", nil)
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestDecryptAndReturn_EmptyPassphrase(t *testing.T) {
	_, err := DecryptAndReturn(map[string]string{"K": "v"}, "", nil)
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}
