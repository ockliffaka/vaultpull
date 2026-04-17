package dotenv

import (
	"strings"
	"testing"
	"time"
)

func TestPrintReport_ContainsSections(t *testing.T) {
	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Path:      ".env",
		Added:     []string{"NEW_KEY"},
		Updated:   []string{"CHANGED_KEY"},
		Skipped:   []string{"SAME_KEY"},
	}
	var sb strings.Builder
	PrintReport(&sb, entry, false)
	out := sb.String()

	for _, want := range []string{"Added", "Updated", "Skipped", "NEW_KEY", "CHANGED_KEY", "SAME_KEY"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestPrintReport_TotalChanges(t *testing.T) {
	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Path:      ".env",
		Added:     []string{"A", "B"},
		Updated:   []string{"C"},
	}
	var sb strings.Builder
	PrintReport(&sb, entry, false)
	out := sb.String()

	if !strings.Contains(out, "Total changes: 3") {
		t.Errorf("expected total changes 3, got:\n%s", out)
	}
}

func TestPrintReport_EmptySections(t *testing.T) {
	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Path:      ".env",
	}
	var sb strings.Builder
	PrintReport(&sb, entry, false)
	out := sb.String()

	if strings.Contains(out, "Added (") {
		t.Errorf("should not print empty Added section\n%s", out)
	}
	if !strings.Contains(out, "Total changes: 0") {
		t.Errorf("expected 0 total changes\n%s", out)
	}
}
