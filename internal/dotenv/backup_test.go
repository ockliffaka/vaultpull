package dotenv

import (
	"os"
	"strings"
	"testing"
)

func TestBackup_NonExistentFile(t *testing.T) {
	path, err := Backup("/tmp/vaultpull_nonexistent_file.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}

func TestBackup_CreatesBackupFile(t *testing.T) {
	src := tmpFile(t, "KEY=value\nFOO=bar\n")

	backupPath, err := Backup(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath == "" {
		t.Fatal("expected a backup path, got empty string")
	}
	t.Cleanup(func() { os.Remove(backupPath) })

	if !strings.HasSuffix(backupPath, ".bak") {
		t.Errorf("expected backup path to end with .bak, got %q", backupPath)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("could not read backup file: %v", err)
	}
	if string(data) != "KEY=value\nFOO=bar\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
}

func TestBackup_OriginalUnchanged(t *testing.T) {
	src := tmpFile(t, "SECRET=abc\n")

	backupPath, err := Backup(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { os.Remove(backupPath) })

	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("could not read original file: %v", err)
	}
	if string(data) != "SECRET=abc\n" {
		t.Errorf("original file was modified: %q", string(data))
	}
}
