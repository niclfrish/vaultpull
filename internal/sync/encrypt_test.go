package sync

import (
	"testing"
)

func TestNewEncryptor_EmptyPassphrase(t *testing.T) {
	_, err := NewEncryptor("")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestNewEncryptor_Success(t *testing.T) {
	enc, err := NewEncryptor("supersecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc == nil {
		t.Fatal("expected non-nil encryptor")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	enc, _ := NewEncryptor("passphrase123")
	plaintext := "my-secret-value"

	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if ciphertext == plaintext {
		t.Fatal("ciphertext should differ from plaintext")
	}

	got, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if got != plaintext {
		t.Fatalf("expected %q, got %q", plaintext, got)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	enc, _ := NewEncryptor("pass")
	_, err := enc.Decrypt("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	enc, _ := NewEncryptor("pass")
	ciphertext, _ := enc.Encrypt("value")
	// Corrupt last character
	tampered := ciphertext[:len(ciphertext)-4] + "AAAA"
	_, err := enc.Decrypt(tampered)
	if err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestEncryptSecrets_RoundTrip(t *testing.T) {
	enc, _ := NewEncryptor("mykey")
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}

	encrypted, err := enc.EncryptSecrets(secrets)
	if err != nil {
		t.Fatalf("EncryptSecrets error: %v", err)
	}

	decrypted, err := enc.DecryptSecrets(encrypted)
	if err != nil {
		t.Fatalf("DecryptSecrets error: %v", err)
	}

	for k, want := range secrets {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: expected %q, got %q", k, want, got)
		}
	}
}

func TestEncrypt_UniqueNonce(t *testing.T) {
	enc, _ := NewEncryptor("nonce-test")
	a, _ := enc.Encrypt("same")
	b, _ := enc.Encrypt("same")
	if a == b {
		t.Fatal("two encryptions of same value should differ due to random nonce")
	}
}
