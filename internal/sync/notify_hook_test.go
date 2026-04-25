package sync

import (
	"bytes"
	"errors"
	"testing"
)

func TestNotifyOnSync_NilNotifier_NoOp(t *testing.T) {
	fn := NotifyOnSync(nil, "prod", nil)
	if err := fn(map[string]string{"K": "V"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotifyOnSync_EmitsInfoEvent(t *testing.T) {
	cap := &captureSink{}
	n := NewNotifier(cap)
	fn := NotifyOnSync(n, "prod", nil)
	secrets := map[string]string{"A": "1", "B": "2"}
	if err := fn(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cap.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(cap.events))
	}
	ev := cap.events[0]
	if ev.Level != NotifyInfo {
		t.Errorf("expected INFO, got %s", ev.Level)
	}
	if ev.Namespace != "prod" {
		t.Errorf("expected namespace prod, got %s", ev.Namespace)
	}
	if ev.Meta["keys"] != "2" {
		t.Errorf("expected keys=2, got %s", ev.Meta["keys"])
	}
}

func TestNotifyOnSync_SinkError_WritesToWriter(t *testing.T) {
	n := NewNotifier(&errSink{})
	var buf bytes.Buffer
	fn := NotifyOnSync(n, "", &buf)
	if err := fn(map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected warning written to buffer")
	}
}

func TestNotifyOnError_NilNotifier_ReturnsOriginalError(t *testing.T) {
	original := errors.New("vault down")
	fn := NotifyOnError(nil, "dev", nil)
	if err := fn(original); err != original {
		t.Fatalf("expected original error returned")
	}
}

func TestNotifyOnError_NilError_NoOp(t *testing.T) {
	cap := &captureSink{}
	n := NewNotifier(cap)
	fn := NotifyOnError(n, "dev", nil)
	if err := fn(nil); err != nil {
		t.Fatalf("unexpected error")
	}
	if len(cap.events) != 0 {
		t.Error("expected no events for nil error")
	}
}

func TestNotifyOnError_EmitsErrorEvent(t *testing.T) {
	cap := &captureSink{}
	n := NewNotifier(cap)
	original := errors.New("connection refused")
	fn := NotifyOnError(n, "staging", nil)
	returned := fn(original)
	if returned != original {
		t.Fatal("original error must be returned unchanged")
	}
	if len(cap.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(cap.events))
	}
	ev := cap.events[0]
	if ev.Level != NotifyError {
		t.Errorf("expected ERROR level, got %s", ev.Level)
	}
	if ev.Meta["error"] != original.Error() {
		t.Errorf("error meta not set correctly")
	}
}
