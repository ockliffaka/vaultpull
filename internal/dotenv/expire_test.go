package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteAndReadStamp_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	if err := WriteStamp(envPath); err != nil {
		t.Fatalf("WriteStamp: %v", err)
	}

	rec, err := ReadStamp(envPath, 24*time.Hour)
	if err != nil {
		t.Fatalf("ReadStamp: %v", err)
	}

	if rec.SyncedAt.IsZero() {
		t.Error("expected non-zero SyncedAt")
	}
	if rec.IsExpired() {
		t.Error("freshly written stamp should not be expired")
	}
}

func TestReadStamp_MissingFile_NotExpired(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	rec, err := ReadStamp(envPath, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rec.SyncedAt.IsZero() {
		t.Error("expected zero SyncedAt for missing stamp")
	}
	// Zero time means synced "never" — duration since epoch is huge, so expired.
	if !rec.IsExpired() {
		t.Error("missing stamp should be considered expired")
	}
}

func TestIsExpired_OldStamp(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Write a stamp with a very old unix time.
	stamp := StampPath(envPath)
	if err := os.WriteFile(stamp, []byte("1"), 0o600); err != nil {
		t.Fatal(err)
	}

	rec, err := ReadStamp(envPath, time.Hour)
	if err != nil {
		t.Fatalf("ReadStamp: %v", err)
	}
	if !rec.IsExpired() {
		t.Error("old stamp should be expired")
	}
}

func TestStampPath(t *testing.T) {
	got := StampPath("/project/.env")
	want := "/project/.env.synced"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
