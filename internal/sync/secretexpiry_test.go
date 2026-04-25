package sync

import (
	"testing"
	"time"
)

func TestParseExpiryHeader_Empty(t *testing.T) {
	_, err := ParseExpiryHeader("")
	if err == nil {
		t.Fatal("expected error for empty header")
	}
}

func TestParseExpiryHeader_UnixTimestamp(t *testing.T) {
	expected := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := ParseExpiryHeader("1893456000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func TestParseExpiryHeader_RFC3339(t *testing.T) {
	raw := "2030-01-01T00:00:00Z"
	got, err := ParseExpiryHeader(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Year() != 2030 {
		t.Errorf("expected year 2030, got %d", got.Year())
	}
}

func TestParseExpiryHeader_Invalid(t *testing.T) {
	_, err := ParseExpiryHeader("not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestClassifyExpiry_NoMetadata(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := ClassifyExpiry(secrets, time.Now())
	if len(result) != 0 {
		t.Errorf("expected 0 results, got %d", len(result))
	}
}

func TestClassifyExpiry_ExpiredAndActive(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)
	secrets := map[string]string{
		"DB_PASS__expires_at__":  strconv.FormatInt(past.Unix(), 10),
		"API_KEY__expires_at__": strconv.FormatInt(future.Unix(), 10),
	}
	infos := ClassifyExpiry(secrets, now)
	if len(infos) != 2 {
		t.Fatalf("expected 2 infos, got %d", len(infos))
	}
	expiredCount := 0
	for _, i := range infos {
		if i.Expired {
			expiredCount++
		}
	}
	if expiredCount != 1 {
		t.Errorf("expected 1 expired, got %d", expiredCount)
	}
}

func TestExpirySummary_NoInfos(t *testing.T) {
	s := ExpirySummary(nil)
	if s != "no expiry metadata found" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestExpirySummary_WithInfos(t *testing.T) {
	infos := []ExpiryInfo{
		{Key: "A", Expired: true},
		{Key: "B", Expired: false},
	}
	s := ExpirySummary(infos)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
