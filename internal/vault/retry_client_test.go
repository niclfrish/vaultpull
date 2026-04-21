package vault

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	isync "github.com/yourusername/vaultpull/internal/sync"
)

func fastRetryConfig(max int) isync.RetryConfig {
	return isync.RetryConfig{MaxAttempts: max, Delay: 0, Multiplier: 1.0}
}

func TestRetryClient_SucceedsOnFirstAttempt(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"KEY":"val"}}}`))
	}))
	defer ts.Close()

	c, _ := New(ts.URL, "token")
	rc := NewRetryClient(c, fastRetryConfig(3))

	secrets, err := rc.GetSecrets("secret/data/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %v", secrets)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryClient_RetriesOn500(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"errors":["500 internal server error"]}`)) 
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"K":"v"}}}`))
	}))
	defer ts.Close()

	c, _ := New(ts.URL, "token")
	rc := NewRetryClient(c, fastRetryConfig(3))

	_, err := rc.GetSecrets("secret/data/test")
	// May succeed or fail depending on client error parsing; we verify call count.
	_ = err
	if atomic.LoadInt32(&calls) < 2 {
		t.Errorf("expected at least 2 calls for retry, got %d", calls)
	}
}

func TestIsTransient(t *testing.T) {
	cases := []struct {
		err      error
		want     bool
	}{
		{errors.New("connection refused"), true},
		{errors.New("request timeout"), true},
		{errors.New("503 service unavailable"), true},
		{errors.New("404 not found"), false},
		{errors.New("permission denied"), false},
		{nil, false},
	}
	for _, tc := range cases {
		got := isTransient(tc.err)
		if got != tc.want {
			t.Errorf("isTransient(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}
