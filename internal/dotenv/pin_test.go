package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPinPath(t *testing.T) {
	path := PinPath("/tmp/app/.env", "v1")
	expected := "/tmp/app/.env.pin.v1.json"
	// PinPath uses base of envPath
	if filepath.Base(path) != ".env.pin.v1.json" {
		t.Errorf("unexpected pin path base: got %q, want %q", filepath.Base(path), expected)
	}
}

func TestPin_AndLoadPin_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "secret123",
	}

	if err := Pin(envPath, "release-1", secrets); err != nil {
		t.Fatalf("Pin failed: %v", err)
	}

	record, err := LoadPin(envPath, "release-1")
	if err != nil {
		t.Fatalf("LoadPin failed: %v", err)
	}
	if record == nil {
		t.Fatal("expected record, got nil")
	}
	if record.Label != "release-1" {
		t.Errorf("label mismatch: got %q", record.Label)
	}
	if record.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("secret mismatch: got %q", record.Secrets["DB_HOST"])
	}
	if record.PinnedAt.IsZero() {
		t.Error("PinnedAt should not be zero")
	}
}

func TestPin_EmptyLabel_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	err := Pin(envPath, "", map[string]string{"KEY": "val"})
	if err == nil {
		t.Error("expected error for empty label")
	}
}

func TestLoadPin_MissingFile_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	record, err := LoadPin(envPath, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record != nil {
		t.Error("expected nil for missing pin")
	}
}

func TestDeletePin_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	if err := Pin(envPath, "v2", map[string]string{"X": "y"}); err != nil {
		t.Fatalf("Pin failed: %v", err)
	}
	if err := DeletePin(envPath, "v2"); err != nil {
		t.Fatalf("DeletePin failed: %v", err)
	}
	pinFile := PinPath(envPath, "v2")
	if _, err := os.Stat(pinFile); !os.IsNotExist(err) {
		t.Error("expected pin file to be removed")
	}
}

func TestDeletePin_NoFile_IsNoop(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := DeletePin(envPath, "ghost"); err != nil {
		t.Errorf("expected no error deleting nonexistent pin, got: %v", err)
	}
}
