package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"vaultpull/internal/dotenv"
)

func TestRollback_AfterBackupAndWrite(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Write original content
	originalContent := "DB_HOST=localhost\nDB_PORT=5432\n"
	if err := os.WriteFile(envPath, []byte(originalContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	// Create backup
	if _, err := dotenv.Backup(envPath); err != nil {
		t.Fatalf("backup: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	// Overwrite with new content
	newContent := "DB_HOST=prod-server\nDB_PORT=5432\nNEW_KEY=value\n"
	if err := os.WriteFile(envPath, []byte(newContent), 0600); err != nil {
		t.Fatalf("overwrite env: %v", err)
	}

	// Rollback
	_, err := dotenv.Rollback(envPath)
	if err != nil {
		t.Fatalf("rollback: %v", err)
	}

	data, _ := os.ReadFile(envPath)
	if string(data) != originalContent {
		t.Errorf("expected original content after rollback, got: %s", data)
	}
}

func TestListBackups_IgnoresUnrelatedFiles(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	os.WriteFile(filepath.Join(dir, ".env.backup.20240101"), []byte("x=1"), 0600)
	os.WriteFile(filepath.Join(dir, ".env.other"), []byte("x=1"), 0600)
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("notes"), 0600)

	backups, err := dotenv.ListBackups(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 1 {
		t.Errorf("expected 1 backup, got %d", len(backups))
	}
}
