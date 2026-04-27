package sync

import (
	"strings"
	"testing"
)

func TestGroupSecrets_NilSecrets(t *testing.T) {
	_, err := GroupSecrets(nil, DefaultGroupConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestGroupSecrets_EmptySeparator(t *testing.T) {
	cfg := DefaultGroupConfig()
	cfg.Separator = ""
	_, err := GroupSecrets(map[string]string{"A": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestGroupSecrets_InvalidMaxDepth(t *testing.T) {
	cfg := DefaultGroupConfig()
	cfg.MaxDepth = 0
	_, err := GroupSecrets(map[string]string{"A": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for MaxDepth < 1")
	}
}

func TestGroupSecrets_NoSeparatorInKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "1", "BAR": "2"}
	groups, err := GroupSecrets(secrets, DefaultGroupConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Prefix != "" {
		t.Errorf("expected empty prefix, got %q", groups[0].Prefix)
	}
	if len(groups[0].Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(groups[0].Secrets))
	}
}

func TestGroupSecrets_GroupsByPrefix(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"APP_SECRET": "abc",
		"APP_KEY":    "xyz",
		"STANDALONE": "val",
	}
	groups, err := GroupSecrets(secrets, DefaultGroupConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// expect groups: "", "APP", "DB" sorted
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	if groups[0].Prefix != "" {
		t.Errorf("expected first group to be ungrouped, got %q", groups[0].Prefix)
	}
	if groups[1].Prefix != "APP" {
		t.Errorf("expected APP group, got %q", groups[1].Prefix)
	}
	if groups[2].Prefix != "DB" {
		t.Errorf("expected DB group, got %q", groups[2].Prefix)
	}
	if len(groups[2].Secrets) != 2 {
		t.Errorf("expected 2 DB secrets, got %d", len(groups[2].Secrets))
	}
}

func TestGroupSecrets_MaxDepthTwo(t *testing.T) {
	secrets := map[string]string{
		"AWS_S3_BUCKET": "my-bucket",
		"AWS_S3_REGION": "us-east-1",
		"AWS_EC2_ID":    "i-123",
	}
	cfg := GroupConfig{Separator: "_", MaxDepth: 2}
	groups, err := GroupSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With MaxDepth=2, prefix uses first 2 segments: AWS_S3 and AWS_EC2
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupSummary_NoGroups(t *testing.T) {
	result := GroupSummary(nil)
	if result != "no groups" {
		t.Errorf("unexpected summary: %q", result)
	}
}

func TestGroupSummary_WithGroups(t *testing.T) {
	groups := []SecretGroup{
		{Prefix: "DB", Secrets: map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}},
		{Prefix: "", Secrets: map[string]string{"STANDALONE": "val"}},
	}
	summary := GroupSummary(groups)
	if !strings.Contains(summary, "group=DB keys=2") {
		t.Errorf("expected DB group in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "(ungrouped)") {
		t.Errorf("expected ungrouped label in summary, got: %s", summary)
	}
}
