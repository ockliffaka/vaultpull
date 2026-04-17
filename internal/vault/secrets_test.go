package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestReadSecrets_KVv2(t *testing.T) {
	server := newMockVaultServer(t, map[string]interface{}{
		"data": map[string]interface{}{
			"API_KEY": "abc123",
			"DB_PASS": "secret",
		},
	})
	defer server.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = server.URL
	raw, _ := vaultapi.NewClient(cfg)
	c := &Client{logical: raw.Logical()}

	secrets, err := c.ReadSecrets(context.Background(), "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
	if secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q", secrets["DB_PASS"])
	}
}

func TestReadSecrets_EmptyPath(t *testing.T) {
	server := newMockVaultServer(t, nil)
	defer server.Close()

	cfg := vaultapi.DefaultConfig()
	cfg.Address = server.URL
	raw, _ := vaultapi.NewClient(cfg)
	c := &Client{logical: raw.Logical()}

	_, err := c.ReadSecrets(context.Background(), "secret/data/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestToStringMap_NonStringValues(t *testing.T) {
	input := map[string]interface{}{
		"COUNT": 42,
		"FLAG":  true,
		"NAME":  "vault",
	}
	out, err := toStringMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["COUNT"] != "42" {
		t.Errorf("expected COUNT=42, got %q", out["COUNT"])
	}
	if out["FLAG"] != "true" {
		t.Errorf("expected FLAG=true, got %q", out["FLAG"])
	}
	if out["NAME"] != "vault" {
		t.Errorf("expected NAME=vault, got %q", out["NAME"])
	}
}
