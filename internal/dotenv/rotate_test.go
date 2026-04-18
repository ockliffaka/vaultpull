package dotenv

import (
	"os"
	"testing"
	"time"
)

func TestRotate_CreatesFileAndStamp(t *testing.T) {
	f := tmpFile(t, "")
	incoming := map[string]string{"API_KEY": "secret", "DB_URL": "postgres://localhost"}

	res, err := Rotate(f, incoming, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.KeysUpdated != 2 {
		t.Errorf("expected 2 keys updated, got %d", res.KeysUpdated)
	}

	parsed, err := parse(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if parsed["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %q", parsed["API_KEY"])
	}

	stamp := StampPath(f)
	if _, err := os.Stat(stamp); os.IsNotExist(err) {
		t.Error("expected stamp file to exist")
	}
	os.Remove(stamp)
}

func TestRotate_OverwritesExistingKeys(t *testing.T) {
	f := tmpFile(t, "API_KEY=old\nKEEP=yes\n")
	incoming := map[string]string{"API_KEY": "new"}

	_, err := Rotate(f, incoming, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsed, err := parse(f)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if parsed["API_KEY"] != "new" {
		t.Errorf("expected API_KEY=new, got %q", parsed["API_KEY"])
	}
	if parsed["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes to be preserved, got %q", parsed["KEEP"])
	}
	os.Remove(StampPath(f))
}

func TestRotate_CreatesBackup(t *testing.T) {
	f := tmpFile(t, "OLD=value\n")

	res, err := Rotate(f, map[string]string{"NEW": "val"}, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BackupPath == "" {
		t.Error("expected a backup path")
	}
	if _, err := os.Stat(res.BackupPath); os.IsNotExist(err) {
		t.Error("backup file does not exist")
	}
	os.Remove(StampPath(f))
}
