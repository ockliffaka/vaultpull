package dotenv

import (
	"os"
	"testing"
	"time"
)

func TestEvaluateTTL_Fresh(t *testing.T) {
	policy := DefaultTTLPolicy()
	stamped := time.Now().Add(-1 * time.Hour)
	if got := EvaluateTTL(stamped, policy); got != TTLFresh {
		t.Errorf("expected fresh, got %s", got)
	}
}

func TestEvaluateTTL_Warning(t *testing.T) {
	policy := DefaultTTLPolicy()
	stamped := time.Now().Add(-21 * time.Hour)
	if got := EvaluateTTL(stamped, policy); got != TTLWarning {
		t.Errorf("expected warning, got %s", got)
	}
}

func TestEvaluateTTL_Expired(t *testing.T) {
	policy := DefaultTTLPolicy()
	stamped := time.Now().Add(-25 * time.Hour)
	if got := EvaluateTTL(stamped, policy); got != TTLExpired {
		t.Errorf("expected expired, got %s", got)
	}
}

func TestEvaluateTTL_ZeroTime(t *testing.T) {
	policy := DefaultTTLPolicy()
	if got := EvaluateTTL(time.Time{}, policy); got != TTLUnknown {
		t.Errorf("expected unknown, got %s", got)
	}
}

func TestTTLStatus_String(t *testing.T) {
	cases := map[TTLStatus]string{
		TTLFresh:   "fresh",
		TTLWarning: "warning",
		TTLExpired: "expired",
		TTLUnknown: "unknown",
	}
	for status, want := range cases {
		if got := status.String(); got != want {
			t.Errorf("status %d: got %q want %q", status, got, want)
		}
	}
}

func TestTTLSummary_ReturnsString(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.env"
	if err := os.WriteFile(path, []byte("KEY=val\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := WriteStamp(path); err != nil {
		t.Fatal(err)
	}
	policy := DefaultTTLPolicy()
	summary, err := TTLSummary(path, policy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
