package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSnapshotPath(t *testing.T) {
	got := SnapshotPath("/tmp/myapp/.env")
	want := "/tmp/myapp/.env.snapshot.json"
	if got != want {
		t.Errorf("SnapshotPath() = %q, want %q", got, want)
	}
}

func TestSaveAndLoadSnapshot_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "abc123",
	}

	if err := SaveSnapshot(envPath, secrets); err != nil {
		t.Fatalf("SaveSnapshot() error: %v", err)
	}

	snap, err := LoadSnapshot(envPath)
	if err != nil {
		t.Fatalf("LoadSnapshot() error: %v", err)
	}
	if snap == nil {
		t.Fatal("LoadSnapshot() returned nil, want snapshot")
	}

	if snap.Path != envPath {
		t.Errorf("Path = %q, want %q", snap.Path, envPath)
	}
	if snap.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want %q", snap.Secrets["DB_HOST"], "localhost")
	}
	if time.Since(snap.Timestamp) > 5*time.Second {
		t.Errorf("Timestamp too old: %v", snap.Timestamp)
	}
}

func TestLoadSnapshot_MissingFile_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	snap, err := LoadSnapshot(envPath)
	if err != nil {
		t.Fatalf("LoadSnapshot() unexpected error: %v", err)
	}
	if snap != nil {
		t.Errorf("LoadSnapshot() = %v, want nil", snap)
	}
}

func TestDeleteSnapshot_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	if err := SaveSnapshot(envPath, map[string]string{"X": "1"}); err != nil {
		t.Fatalf("SaveSnapshot() error: %v", err)
	}

	if err := DeleteSnapshot(envPath); err != nil {
		t.Fatalf("DeleteSnapshot() error: %v", err)
	}

	if _, err := os.Stat(SnapshotPath(envPath)); !os.IsNotExist(err) {
		t.Error("snapshot file still exists after DeleteSnapshot()")
	}
}

func TestDeleteSnapshot_NoFile_IsNoop(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	if err := DeleteSnapshot(envPath); err != nil {
		t.Errorf("DeleteSnapshot() on missing file returned error: %v", err)
	}
}
