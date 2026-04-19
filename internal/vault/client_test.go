package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, path string, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
}

func TestNew_MissingAddress(t *testing.T) {
	_, err := New("", "token", "")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNew_MissingToken(t *testing.T) {
	_, err := New("http://127.0.0.1:8200", "", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNew_Success(t *testing.T) {
	c, err := New("http://127.0.0.1:8200", "s.test", "mynamespace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.namespace != "mynamespace" {
		t.Errorf("expected namespace %q, got %q", "mynamespace", c.namespace)
	}
}

func TestGetSecrets_EmptyPath(t *testing.T) {
	c, _ := New("http://127.0.0.1:8200", "s.test", "")
	_, err := c.GetSecrets("")
	if err == nil {
		t.Fatal("expected error for empty secret path")
	}
}

func TestGetSecrets_KVv2(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"API_KEY": "abc123",
				"DB_PASS": "secret",
			},
		},
	}
	ts := newTestServer(t, "/v1/secret/data/myapp", payload)
	defer ts.Close()

	c, err := New(ts.URL, "s.test", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, err := c.GetSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secrets: %v", err)
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
	if secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q", secrets["DB_PASS"])
	}
}
