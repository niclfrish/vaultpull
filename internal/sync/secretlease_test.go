package sync

import (
	"strings"
	"testing"
	"time"
)

func TestParseLeaseHeader_Empty(t *testing.T) {
	_, err := ParseLeaseHeader("")
	if err == nil {
		t.Fatal("expected error for empty header")
	}
}

func TestParseLeaseHeader_Valid(t *testing.T) {
	raw := "lease_id=abc-123,ttl=3600,renewable=true"
	info, err := ParseLeaseHeader(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.LeaseID != "abc-123" {
		t.Errorf("expected lease_id abc-123, got %q", info.LeaseID)
	}
	if info.Duration != 3600*time.Second {
		t.Errorf("expected 3600s, got %s", info.Duration)
	}
	if !info.Renewable {
		t.Error("expected renewable=true")
	}
}

func TestParseLeaseHeader_DefaultTTL(t *testing.T) {
	info, err := ParseLeaseHeader("lease_id=x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Duration != DefaultLeaseTTL {
		t.Errorf("expected default TTL %s, got %s", DefaultLeaseTTL, info.Duration)
	}
}

func TestParseLeaseHeader_InvalidTTL(t *testing.T) {
	_, err := ParseLeaseHeader("ttl=notanumber")
	if err == nil || !strings.Contains(err.Error(), "invalid ttl") {
		t.Fatalf("expected invalid ttl error, got %v", err)
	}
}

func TestParseLeaseHeader_ZeroTTL(t *testing.T) {
	_, err := ParseLeaseHeader("ttl=0")
	if err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("expected positive TTL error, got %v", err)
	}
}

func TestParseLeaseHeader_InvalidRenewable(t *testing.T) {
	_, err := ParseLeaseHeader("renewable=maybe")
	if err == nil || !strings.Contains(err.Error(), "invalid renewable") {
		t.Fatalf("expected invalid renewable error, got %v", err)
	}
}

func TestLeaseInfo_IsExpired_NotYet(t *testing.T) {
	info := &LeaseInfo{ExpiresAt: time.Now().Add(10 * time.Minute)}
	if info.IsExpired() {
		t.Error("expected lease to not be expired")
	}
}

func TestLeaseInfo_IsExpired_Past(t *testing.T) {
	info := &LeaseInfo{ExpiresAt: time.Now().Add(-1 * time.Second)}
	if !info.IsExpired() {
		t.Error("expected lease to be expired")
	}
}

func TestLeaseSummary_Nil(t *testing.T) {
	s := LeaseSummary(nil)
	if s != "no lease" {
		t.Errorf("expected 'no lease', got %q", s)
	}
}

func TestLeaseSummary_WithLease(t *testing.T) {
	info := &LeaseInfo{
		LeaseID:   "tok-42",
		Duration:  5 * time.Minute,
		Renewable: true,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	s := LeaseSummary(info)
	if !strings.Contains(s, "tok-42") {
		t.Errorf("summary missing lease ID: %s", s)
	}
	if !strings.Contains(s, "renewable") {
		t.Errorf("summary missing renewable: %s", s)
	}
}
