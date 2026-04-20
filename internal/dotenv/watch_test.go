package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchOnce_NotExpired(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	stampPath := StampPath(envPath)

	// Write a fresh stamp
	if err := WriteStamp(stampPath); err != nil {
		t.Fatal(err)
	}

	policy := DefaultTTLPolicy()
	refreshCalled := false
	result := WatchOnce(envPath, policy, func(p string) error {
		refreshCalled = true
		return nil
	})

	if result.Expired {
		t.Error("expected not expired for fresh stamp")
	}
	if refreshCalled {
		t.Error("refresh should not be called when not expired")
	}
}

func TestWatchOnce_Expired_CallsRefresh(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	stampPath := StampPath(envPath)

	// Write an old stamp
	old := time.Now().Add(-48 * time.Hour)
	data := old.Format(time.RFC3339)
	if err := os.WriteFile(stampPath, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}

	policy := DefaultTTLPolicy()
	refreshCalled := false
	result := WatchOnce(envPath, policy, func(p string) error {
		refreshCalled = true
		return nil
	})

	if !result.Expired {
		t.Error("expected expired for old stamp")
	}
	if !refreshCalled {
		t.Error("expected refresh to be called")
	}
	if !result.Refreshed {
		t.Error("expected Refreshed=true after successful refresh")
	}
}

func TestWatch_MaxCycles(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	// Write fresh stamp so no refresh triggered
	if err := WriteStamp(StampPath(envPath)); err != nil {
		t.Fatal(err)
	}

	opts := WatchOptions{
		Interval:  10 * time.Millisecond,
		MaxCycles: 3,
		OnRefresh: nil,
	}
	stop := make(chan struct{})
	results := Watch(envPath, DefaultTTLPolicy(), opts, stop)

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}
