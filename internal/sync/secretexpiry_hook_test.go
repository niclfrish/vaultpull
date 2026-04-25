package sync

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCheckExpiry_NoMetadata(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	var buf bytes.Buffer
	res, err := CheckExpiry(secrets, true, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Infos) != 0 {
		t.Errorf("expected 0 infos, got %d", len(res.Infos))
	}
}

func TestCheckExpiry_AllActive_NoError(t *testing.T) {
	future := time.Now().Add(time.Hour).UTC()
	secrets := map[string]string{
		"TOKEN__expires_at__": strconv.FormatInt(future.Unix(), 10),
	}
	var buf bytes.Buffer
	res, err := CheckExpiry(secrets, true, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Expired) != 0 {
		t.Errorf("expected no expired, got %d", len(res.Expired))
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK in output, got: %s", buf.String())
	}
}

func TestCheckExpiry_ExpiredFailOnExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour).UTC()
	secrets := map[string]string{
		"OLD_KEY__expires_at__": strconv.FormatInt(past.Unix(), 10),
	}
	var buf bytes.Buffer
	_, err := CheckExpiry(secrets, true, &buf)
	if err == nil {
		t.Fatal("expected error for expired secret")
	}
	if !strings.Contains(buf.String(), "EXPIRED") {
		t.Errorf("expected EXPIRED in output, got: %s", buf.String())
	}
}

func TestCheckExpiry_ExpiredNoFail(t *testing.T) {
	past := time.Now().Add(-time.Hour).UTC()
	secrets := map[string]string{
		"OLD_KEY__expires_at__": strconv.FormatInt(past.Unix(), 10),
	}
	res, err := CheckExpiry(secrets, false, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Expired) != 1 {
		t.Errorf("expected 1 expired, got %d", len(res.Expired))
	}
}

func TestExpiryFilterStage_RemovesExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour).UTC()
	future := time.Now().Add(time.Hour).UTC()
	secrets := map[string]string{
		"LIVE":                   "yes",
		"DEAD":                   "no",
		"DEAD__expires_at__":     strconv.FormatInt(past.Unix(), 10),
		"LIVE__expires_at__":     strconv.FormatInt(future.Unix(), 10),
	}
	stage := ExpiryFilterStage()
	out, err := stage.Fn(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DEAD"]; ok {
		t.Error("expected DEAD to be removed")
	}
	if _, ok := out["LIVE"]; !ok {
		t.Error("expected LIVE to remain")
	}
}
