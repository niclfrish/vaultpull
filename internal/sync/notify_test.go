package sync

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// errSink always returns an error from Send.
type errSink struct{}

func (e *errSink) Send(_ NotifyEvent) error {
	return errors.New("sink failure")
}

// captureSink records all received events.
type captureSink struct {
	events []NotifyEvent
}

func (c *captureSink) Send(ev NotifyEvent) error {
	c.events = append(c.events, ev)
	return nil
}

func TestNewNotifier_NoSinks(t *testing.T) {
	n := NewNotifier()
	errs := n.Emit(NotifyInfo, "", "hello", nil)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestNotifier_Emit_ReachesAllSinks(t *testing.T) {
	s1 := &captureSink{}
	s2 := &captureSink{}
	n := NewNotifier(s1, s2)
	n.Emit(NotifyWarn, "prod", "disk full", map[string]string{"host": "srv1"})
	if len(s1.events) != 1 || len(s2.events) != 1 {
		t.Fatal("expected both sinks to receive the event")
	}
	if s1.events[0].Level != NotifyWarn {
		t.Errorf("unexpected level: %s", s1.events[0].Level)
	}
	if s1.events[0].Meta["host"] != "srv1" {
		t.Errorf("meta not propagated")
	}
}

func TestNotifier_Emit_CollectsErrors(t *testing.T) {
	n := NewNotifier(&errSink{}, &errSink{})
	errs := n.Emit(NotifyError, "", "boom", nil)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestWriterSink_Send_FormatsOutput(t *testing.T) {
	var buf bytes.Buffer
	ws := NewWriterSink(&buf)
	sink := &captureSink{}
	n := NewNotifier(ws, sink)
	n.Emit(NotifyInfo, "staging", "sync complete", map[string]string{"keys": "5"})
	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "ns=staging") {
		t.Errorf("expected namespace in output, got: %s", out)
	}
	if !strings.Contains(out, "sync complete") {
		t.Errorf("expected message in output, got: %s", out)
	}
	if !strings.Contains(out, "keys=5") {
		t.Errorf("expected meta in output, got: %s", out)
	}
}

func TestWriterSink_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic.
	ws := NewWriterSink(nil)
	if ws.w == nil {
		t.Fatal("expected fallback to stdout")
	}
}

func TestWriterSink_EmptyNamespace_ShowsDefault(t *testing.T) {
	var buf bytes.Buffer
	ws := NewWriterSink(&buf)
	ws.Send(NotifyEvent{Level: NotifyInfo, Message: "hi", Namespace: ""})
	if !strings.Contains(buf.String(), "ns=default") {
		t.Errorf("expected default namespace, got: %s", buf.String())
	}
}
