package sync

import (
	"testing"
)

func TestFlattenSecrets_NilSecrets(t *testing.T) {
	_, err := FlattenSecrets(nil, DefaultFlattenConfig())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestFlattenSecrets_EmptySeparator(t *testing.T) {
	cfg := DefaultFlattenConfig()
	cfg.Separator = ""
	_, err := FlattenSecrets(map[string]string{"a": "1"}, cfg)
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestFlattenSecrets_NoDots_ReturnsUpperCase(t *testing.T) {
	secrets := map[string]string{"dbhost": "localhost"}
	out, err := FlattenSecrets(secrets, DefaultFlattenConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DBHOST"] != "localhost" {
		t.Errorf("expected DBHOST=localhost, got %v", out)
	}
}

func TestFlattenSecrets_DotSeparatedKey(t *testing.T) {
	secrets := map[string]string{"db.host": "localhost", "db.port": "5432"}
	out, err := FlattenSecrets(secrets, DefaultFlattenConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out)
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %v", out)
	}
}

func TestFlattenSecrets_MaxDepth(t *testing.T) {
	cfg := DefaultFlattenConfig()
	cfg.MaxDepth = 2
	secrets := map[string]string{"a.b.c.d": "val"}
	out, err := FlattenSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["A_B"]; !ok {
		t.Errorf("expected key A_B, got %v", out)
	}
}

func TestFlattenSecrets_KeyCollision(t *testing.T) {
	// Both keys flatten to the same result.
	secrets := map[string]string{"db.host": "a", "db_host": "b"}
	_, err := FlattenSecrets(secrets, DefaultFlattenConfig())
	if err == nil {
		t.Fatal("expected collision error")
	}
}

func TestFlattenSecrets_LowerCaseOption(t *testing.T) {
	cfg := DefaultFlattenConfig()
	cfg.UpperCase = false
	secrets := map[string]string{"App.Name": "vaultpull"}
	out, err := FlattenSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["App_Name"] != "vaultpull" {
		t.Errorf("expected App_Name=vaultpull, got %v", out)
	}
}

func TestFlattenSummary(t *testing.T) {
	before := map[string]string{"a.b": "1", "c.d": "2"}
	after := map[string]string{"A_B": "1", "C_D": "2"}
	s := FlattenSummary(before, after)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
