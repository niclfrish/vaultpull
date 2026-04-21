package sync

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func fastWatchConfig(ticks int) WatchConfig {
	return WatchConfig{
		Interval: 10 * time.Millisecond,
		MaxTicks: ticks,
	}
}

func TestWatcher_SyncsSecretsOnTick(t *testing.T) {
	var written map[string]string
	fetch := func(_ context.Context) (map[string]string, error) {
		return map[string]string{"KEY": "val"}, nil
	}
	write := func(s map[string]string) error {
		written = s
		return nil
	}
	var buf bytes.Buffer
	w := NewWatcher(fastWatchConfig(1), fetch, write, &buf)
	if err := w.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if written["KEY"] != "val" {
		t.Errorf("expected written[KEY]=val, got %q", written["KEY"])
	}
	if !strings.Contains(buf.String(), "secrets synced") {
		t.Errorf("expected synced message, got: %s", buf.String())
	}
}

func TestWatcher_LogsFetchError(t *testing.T) {
	fetch := func(_ context.Context) (map[string]string, error) {
		return nil, errors.New("vault unavailable")
	}
	write := func(_ map[string]string) error { return nil }
	var buf bytes.Buffer
	w := NewWatcher(fastWatchConfig(1), fetch, write, &buf)
	_ = w.Run(context.Background())
	if !strings.Contains(buf.String(), "fetch error") {
		t.Errorf("expected fetch error in output, got: %s", buf.String())
	}
}

func TestWatcher_LogsWriteError(t *testing.T) {
	fetch := func(_ context.Context) (map[string]string, error) {
		return map[string]string{"A": "1"}, nil
	}
	write := func(_ map[string]string) error { return errors.New("disk full") }
	var buf bytes.Buffer
	w := NewWatcher(fastWatchConfig(1), fetch, write, &buf)
	_ = w.Run(context.Background())
	if !strings.Contains(buf.String(), "write error") {
		t.Errorf("expected write error in output, got: %s", buf.String())
	}
}

func TestWatcher_CancelStopsLoop(t *testing.T) {
	fetch := func(_ context.Context) (map[string]string, error) {
		return map[string]string{}, nil
	}
	write := func(_ map[string]string) error { return nil }
	ctx, cancel := context.WithCancel(context.Background())
	w := NewWatcher(WatchConfig{Interval: 5 * time.Millisecond, MaxTicks: 0}, fetch, write, nil)
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	err := w.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestDefaultWatchConfig(t *testing.T) {
	cfg := DefaultWatchConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.MaxTicks != 0 {
		t.Errorf("expected MaxTicks=0, got %d", cfg.MaxTicks)
	}
}
