package sync

import (
	"context"
	"io"
	"os"
	"time"
)

// WatchConfig holds configuration for the secret watcher.
type WatchConfig struct {
	// Interval is how often to poll Vault for changes.
	Interval time.Duration
	// MaxTicks is the maximum number of poll cycles (0 = unlimited).
	MaxTicks int
}

// DefaultWatchConfig returns a WatchConfig with sensible defaults.
func DefaultWatchConfig() WatchConfig {
	return WatchConfig{
		Interval: 30 * time.Second,
		MaxTicks: 0,
	}
}

// FetchFunc fetches secrets from Vault and returns them as a map.
type FetchFunc func(ctx context.Context) (map[string]string, error)

// WriteFunc writes secrets to the local .env file.
type WriteFunc func(secrets map[string]string) error

// Watcher polls Vault at a fixed interval and writes changes to disk.
type Watcher struct {
	cfg    WatchConfig
	fetch  FetchFunc
	write  WriteFunc
	out    io.Writer
}

// NewWatcher creates a new Watcher.
func NewWatcher(cfg WatchConfig, fetch FetchFunc, write WriteFunc, out io.Writer) *Watcher {
	if out == nil {
		out = os.Stdout
	}
	return &Watcher{cfg: cfg, fetch: fetch, write: write, out: out}
}

// Run starts the watch loop. It blocks until ctx is cancelled or MaxTicks is reached.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	ticks := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			secrets, err := w.fetch(ctx)
			if err != nil {
				_, _ = io.WriteString(w.out, "[watch] fetch error: "+err.Error()+"\n")
			} else {
				if werr := w.write(secrets); werr != nil {
					_, _ = io.WriteString(w.out, "[watch] write error: "+werr.Error()+"\n")
				} else {
					_, _ = io.WriteString(w.out, "[watch] secrets synced\n")
				}
			}
			ticks++
			if w.cfg.MaxTicks > 0 && ticks >= w.cfg.MaxTicks {
				return nil
			}
		}
	}
}
