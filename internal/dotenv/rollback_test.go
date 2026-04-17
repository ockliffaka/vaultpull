package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListBackups_NoBackups(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	backups, err := ListBackups(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("expected 0 backups, got %d", len(backups))
	}
}

func TestListBackups_ReturnsNewestFirst(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	names := []string{".env.backup.20240101", ".env.backup.20240303", ".env.backup.20240202"}
	for _, n := range names {
		os.WriteFile(filepath.Join(dir, n), []byte("x=1"), 0600)
	}

	backups, err := ListBackups(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 3 {
		t.Fatalf("expected 3 backups, got %d", len(backups))
	}
	if filepath.Base(backups[0]) != ".env.backup.20240303" {
		t.Errorf("expected newest first, got %s", backups[0])
	}
}

func TestRollback_NoBackup(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	_, err := Rollback(envPath)
	if err == nil {
		t.Fatal("expected error for missing backup")
	}
}

func TestRollback_RestoresContent(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	original := []byte("KEY=original\n")
	backupPath := envPath + ".backup.20240101120000"
	os.WriteFile(backupPath, original, 0600)
	os.WriteFile(envPath, []byte("KEY=changed\n"), 0600)

	restored, err := Rollback(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if restored != backupPath {
		t.Errorf("expected %s, got %s", backupPath, restored)
	}

	data, _ := os.ReadFile(envPath)
	if string(data) != string(original) {
		t.Errorf("content not restored: %s", data)
	}
}
