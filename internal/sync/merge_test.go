package sync

import "testing"

func TestMerge_VaultOverridesLocal(t *testing.T) {
	local := map[string]string{
		"KEY_A": "local_a",
		"KEY_B": "local_b",
	}
	vault := map[string]string{
		"KEY_A": "vault_a",
		"KEY_C": "vault_c",
	}

	result := Merge(local, vault)

	if result["KEY_A"] != "vault_a" {
		t.Errorf("expected vault_a, got %s", result["KEY_A"])
	}
	if result["KEY_B"] != "local_b" {
		t.Errorf("expected local_b, got %s", result["KEY_B"])
	}
	if result["KEY_C"] != "vault_c" {
		t.Errorf("expected vault_c, got %s", result["KEY_C"])
	}
}

func TestMerge_EmptyLocal(t *testing.T) {
	result := Merge(nil, map[string]string{"X": "1"})
	if result["X"] != "1" {
		t.Errorf("expected 1, got %s", result["X"])
	}
}

func TestMerge_EmptyVault(t *testing.T) {
	result := Merge(map[string]string{"Y": "2"}, nil)
	if result["Y"] != "2" {
		t.Errorf("expected 2, got %s", result["Y"])
	}
}

func TestMerge_BothEmpty(t *testing.T) {
	result := Merge(nil, nil)
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}
