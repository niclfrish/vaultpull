package sync

import (
	"context"
	"testing"
	"time"
)

func TestNewThrottler_InvalidMinInterval(t *testing.T) {
	_, err := NewThrottler(ThrottleConfig{MinInterval: -1 * time.Millisecond})
	if err == nil {
		t.Fatal("expected error for negative MinInterval")
	}
}

func TestNewThrottler_InvalidMaxBatchSize(t *testing.T) {
	_, err := NewThrottler(ThrottleConfig{MaxBatchSize: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxBatchSize")
	}
}

func TestNewThrottler_Success(t *testing.T) {
	th, err := NewThrottler(DefaultThrottleConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil Throttler")
	}
}

func TestThrottler_Wait_ZeroInterval(t *testing.T) {
	th, _ := NewThrottler(ThrottleConfig{MinInterval: 0})
	ctx := context.Background()
	start := time.Now()
	if err := th.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time.Since(start) > 50*time.Millisecond {
		t.Error("Wait with zero interval should return immediately")
	}
}

func TestThrottler_Wait_ContextCancelled(t *testing.T) {
	th, _ := NewThrottler(ThrottleConfig{MinInterval: 5 * time.Second})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := th.Wait(ctx)
	if err == nil {
		t.Fatal("expected error when context is cancelled")
	}
}

func TestThrottler_Wait_RespectsInterval(t *testing.T) {
	th, _ := NewThrottler(ThrottleConfig{MinInterval: 50 * time.Millisecond})
	ctx := context.Background()
	_ = th.Wait(ctx)
	start := time.Now()
	_ = th.Wait(ctx)
	if time.Since(start) < 40*time.Millisecond {
		t.Error("second Wait returned too quickly; interval not respected")
	}
}

func TestThrottler_Batch_NoLimit(t *testing.T) {
	th, _ := NewThrottler(ThrottleConfig{MaxBatchSize: 0})
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	batches := th.Batch(secrets)
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}
	if len(batches[0]) != 3 {
		t.Errorf("expected 3 keys in batch, got %d", len(batches[0]))
	}
}

func TestThrottler_Batch_SplitsCorrectly(t *testing.T) {
	th, _ := NewThrottler(ThrottleConfig{MaxBatchSize: 2})
	secrets := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5"}
	batches := th.Batch(secrets)
	total := 0
	for _, b := range batches {
		if len(b) > 2 {
			t.Errorf("batch size %d exceeds MaxBatchSize 2", len(b))
		}
		total += len(b)
	}
	if total != 5 {
		t.Errorf("expected 5 total keys across batches, got %d", total)
	}
}

func TestThrottler_Batch_ExactBatchSize(t *testing.T) {
	th, _ := NewThrottler(ThrottleConfig{MaxBatchSize: 3})
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	batches := th.Batch(secrets)
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch for exact size, got %d", len(batches))
	}
}
