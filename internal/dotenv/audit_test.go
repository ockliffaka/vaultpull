package dotenv

import (
	"strings"
	"testing"
)

func TestAuditLog_RecordAndSummary(t *testing.T) {
	log := &AuditLog{}
	d := DiffResult{
		Added:     map[string]string{"NEW_KEY": "val"},
		Changed:   map[string][2]string{"OLD_KEY": {"a", "b"}},
		Unchanged: map[string]string{"SAME": "x"},
	}
	log.Record(".env", d)

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if len(e.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(e.Added))
	}
	if len(e.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d", len(e.Updated))
	}
	if len(e.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(e.Skipped))
	}

	summary := log.Summary()
	if !strings.Contains(summary, ".env") {
		t.Errorf("summary missing path: %s", summary)
	}
	if !strings.Contains(summary, "added:   1") {
		t.Errorf("summary missing added count: %s", summary)
	}
}

func TestAuditLog_EmptySummary(t *testing.T) {
	log := &AuditLog{}
	got := log.Summary()
	if got != "no audit entries recorded" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestAuditLog_MultipleEntries(t *testing.T) {
	log := &AuditLog{}
	d1 := DiffResult{Added: map[string]string{"A": "1"}}
	d2 := DiffResult{Added: map[string]string{"B": "2", "C": "3"}}
	log.Record(".env", d1)
	log.Record(".env.prod", d2)

	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
	summary := log.Summary()
	if !strings.Contains(summary, ".env.prod") {
		t.Errorf("summary should reflect last entry: %s", summary)
	}
}
