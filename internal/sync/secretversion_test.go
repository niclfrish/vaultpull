package sync

import (
	"testing"
	"time"
)

func TestParseVersionHeader_Valid(t *testing.T) {
	v, err := ParseVersionHeader("3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 3 {
		t.Errorf("expected 3, got %d", v)
	}
}

func TestParseVersionHeader_WithSpaces(t *testing.T) {
	v, err := ParseVersionHeader("  7  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 7 {
		t.Errorf("expected 7, got %d", v)
	}
}

func TestParseVersionHeader_Empty(t *testing.T) {
	_, err := ParseVersionHeader("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

func TestParseVersionHeader_NonNumeric(t *testing.T) {
	_, err := ParseVersionHeader("abc")
	if err == nil {
		t.Fatal("expected error for non-numeric string")
	}
}

func TestParseVersionHeader_Zero(t *testing.T) {
	_, err := ParseVersionHeader("0")
	if err == nil {
		t.Fatal("expected error for version 0")
	}
}

func TestVersionSummary_Active(t *testing.T) {
	sv := SecretVersion{Version: 2, CreatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)}
	s := VersionSummary(sv)
	if s != "v2 created=2024-01-15T10:00:00Z status=active" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestVersionSummary_Destroyed(t *testing.T) {
	sv := SecretVersion{Version: 1, CreatedAt: time.Now(), Destroyed: true}
	s := VersionSummary(sv)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
	if !containsStr(s, "destroyed") {
		t.Errorf("expected 'destroyed' in summary, got: %s", s)
	}
}

func TestVersionSummary_Deleted(t *testing.T) {
	now := time.Now()
	sv := SecretVersion{Version: 4, CreatedAt: now, DeletedAt: &now}
	s := VersionSummary(sv)
	if !containsStr(s, "deleted") {
		t.Errorf("expected 'deleted' in summary, got: %s", s)
	}
}

func TestFilterActiveVersions(t *testing.T) {
	now := time.Now()
	versions := SecretVersionMap{
		"KEY_A": {Version: 1, CreatedAt: now},
		"KEY_B": {Version: 2, CreatedAt: now, Destroyed: true},
		"KEY_C": {Version: 3, CreatedAt: now, DeletedAt: &now},
		"KEY_D": {Version: 4, CreatedAt: now},
	}
	active := FilterActiveVersions(versions)
	if len(active) != 2 {
		t.Errorf("expected 2 active versions, got %d", len(active))
	}
	if _, ok := active["KEY_A"]; !ok {
		t.Error("expected KEY_A in active versions")
	}
	if _, ok := active["KEY_D"]; !ok {
		t.Error("expected KEY_D in active versions")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
