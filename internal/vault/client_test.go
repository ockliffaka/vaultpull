package vault_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/vaultpull/internal/vault"
)

func TestNewClient_MissingAddress(t *testing.T) {
	_, err := vault.NewClient("", "token")
	if err == nil {
		t.Fatal("expected error for empty address")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	_, err := vault.NewClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestReadSecrets_KVv1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}`))
	}))
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secrets: %v", err)
	}

	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
	if secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q", secrets["DB_PASS"])
	}
}

func TestReadSecrets_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}
