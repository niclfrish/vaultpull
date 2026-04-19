package sync

import (
	"bytes"
	"errors"
	"testing"
)

type mockWriter struct {
	called bool
	data   map[string]string
	err    error
}

func (m *mockWriter) Write(path string, data map[string]string) error {
	m.called = true
	m.data = data
	return m.err
}

func TestApply_NoChanges(t *testing.T) {
	mw := &mockWriter{}
	a := NewApplier(mw)
	var buf bytes.Buffer
	a.out = &buf

	plan := &Plan{}
	err := a.Apply(".env", plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mw.called {
		t.Error("writer should not be called when no changes")
	}
}

func TestApply_WritesCorrectKeys(t *testing.T) {
	mw := &mockWriter{}
	a := NewApplier(mw)
	var buf bytes.Buffer
	a.out = &buf

	plan := &Plan{
		Added:     []DiffEntry{{Key: "NEW_KEY", NewValue: "new"}},
		Changed:   []DiffEntry{{Key: "CHANGED", OldValue: "old", NewValue: "updated"}},
		Removed:   []DiffEntry{{Key: "GONE", OldValue: "bye"}},
		Unchanged: []DiffEntry{{Key: "STABLE", OldValue: "same", NewValue: "same"}},
	}

	if err := a.Apply(".env", plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mw.called {
		t.Fatal("expected writer to be called")
	}
	if mw.data["NEW_KEY"] != "new" {
		t.Errorf("expected NEW_KEY=new, got %s", mw.data["NEW_KEY"])
	}
	if mw.data["CHANGED"] != "updated" {
		t.Errorf("expected CHANGED=updated, got %s", mw.data["CHANGED"])
	}
	if _, ok := mw.data["GONE"]; ok {
		t.Error("removed key should not be present")
	}
	if mw.data["STABLE"] != "same" {
		t.Errorf("expected STABLE=same, got %s", mw.data["STABLE"])
	}
}

func TestApply_WriterError(t *testing.T) {
	mw := &mockWriter{err: errors.New("disk full")}
	a := NewApplier(mw)
	var buf bytes.Buffer
	a.out = &buf

	plan := &Plan{
		Added: []DiffEntry{{Key: "K", NewValue: "v"}},
	}

	if err := a.Apply(".env", plan); err == nil {
		t.Fatal("expected error from writer")
	}
}
