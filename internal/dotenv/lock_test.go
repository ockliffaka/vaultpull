package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAcquireLock_CreatesLockFile(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	lock, err := AcquireLock(env)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer lock.Release()

	if _, err := os.Stat(LockPath(env)); err != nil {
		t.Errorf("lock file not created: %v", err)
	}
}

func TestAcquireLock_FailsIfAlreadyLocked(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	lock, err := AcquireLock(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer lock.Release()

	_, err = AcquireLock(env)
	if err == nil {
		t.Error("expected error acquiring second lock, got nil")
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	lock, _ := AcquireLock(env)
	if err := lock.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}

	if _, err := os.Stat(LockPath(env)); !os.IsNotExist(err) {
		t.Error("lock file should be removed after release")
	}
}

func TestClearStaleLock_RemovesOldLock(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	lp := LockPath(env)

	_ = os.WriteFile(lp, []byte("stale"), 0600)
	past := time.Now().Add(-(StaleLockAge + time.Minute))
	_ = os.Chtimes(lp, past, past)

	cleared, err := ClearStaleLock(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cleared {
		t.Error("expected stale lock to be cleared")
	}
}

func TestClearStaleLock_KeepsFreshLock(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	lp := LockPath(env)

	_ = os.WriteFile(lp, []byte("fresh"), 0600)

	cleared, err := ClearStaleLock(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cleared {
		t.Error("expected fresh lock to be kept")
	}
	os.Remove(lp)
}

func TestClearStaleLock_NoLockFile(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")

	cleared, err := ClearStaleLock(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cleared {
		t.Error("expected false when no lock file exists")
	}
}
