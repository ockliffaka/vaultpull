package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCheckExpiry_NoStamp_ReturnsFresh(t *testing.T) {
	result, err := CheckExpiry("/nonexistent/stamp.json", time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != StatusFresh {
		t.Errorf("expected Fresh, got %v", result.Status)
	}
}

func TestCheckExpiry_FreshStamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stamp.json")

	if err := WriteStamp(path, time.Now()); err != nil {
		t.Fatal(err)
	}

	result, err := CheckExpiry(path, 10*time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != StatusFresh {
		t.Errorf("expected Fresh, got %v", result.Status)
	}
}

func TestCheckExpiry_StaleStamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stamp.json")

	past := time.Now().Add(-30 * time.Minute)
	if err := WriteStamp(path, past); err != nil {
		t.Fatal(err)
	}

	result, err := CheckExpiry(path, 20*time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != StatusStale {
		t.Errorf("expected Stale, got %v", result.Status)
	}
}

func TestCheckExpiry_ExpiredStamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stamp.json")

	past := time.Now().Add(-2 * time.Hour)
	if err := WriteStamp(path, past); err != nil {
		t.Fatal(err)
	}

	result, err := CheckExpiry(path, 10*time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != StatusExpired {
		t.Errorf("expected Expired, got %v", result.Status)
	}
}

func TestExpiryResult_String(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stamp.json")
	os.Remove(path)

	r := ExpiryResult{Status: StatusFresh, Age: 5 * time.Second}
	if r.String() == "" {
		t.Error("expected non-empty string")
	}

	r2 := ExpiryResult{Status: StatusExpired, Age: 2 * time.Hour, TTL: time.Hour}
	if r2.String() == "" {
		t.Error("expected non-empty string for expired")
	}
}
