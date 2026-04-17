package dotenv

import (
	"fmt"
	"strings"
	"time"
)

// AuditEntry records a single sync operation result.
type AuditEntry struct {
	Timestamp time.Time
	Path      string
	Added     []string
	Updated   []string
	Skipped   []string
	Removed   []string
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// Record appends a new entry derived from a DiffResult to the log.
func (a *AuditLog) Record(path string, d DiffResult) {
	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Path:      path,
	}
	for k := range d.Added {
		entry.Added = append(entry.Added, k)
	}
	for k := range d.Changed {
		entry.Updated = append(entry.Updated, k)
	}
	for k := range d.Unchanged {
		entry.Skipped = append(entry.Skipped, k)
	}
	a.Entries = append(a.Entries, entry)
}

// Summary returns a human-readable summary of the last audit entry.
func (a *AuditLog) Summary() string {
	if len(a.Entries) == 0 {
		return "no audit entries recorded"
	}
	e := a.Entries[len(a.Entries)-1]
	var sb strings.Builder
	fmt.Fprintf(&sb, "[%s] sync to %s\n", e.Timestamp.Format(time.RFC3339), e.Path)
	fmt.Fprintf(&sb, "  added:   %d\n", len(e.Added))
	fmt.Fprintf(&sb, "  updated: %d\n", len(e.Updated))
	fmt.Fprintf(&sb, "  skipped: %d\n", len(e.Skipped))
	return sb.String()
}
