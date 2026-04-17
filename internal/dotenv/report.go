package dotenv

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// PrintReport writes a detailed sync report to w based on the AuditEntry.
func PrintReport(w io.Writer, e AuditEntry, showValues bool) {
	fmt.Fprintf(w, "Sync Report — %s\n", e.Path)
	fmt.Fprintf(w, "Timestamp: %s\n", e.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintln(w, strings.Repeat("-", 40))

	printSection(w, "Added", e.Added)
	printSection(w, "Updated", e.Updated)
	printSection(w, "Skipped", e.Skipped)

	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintf(w, "Total changes: %d\n", len(e.Added)+len(e.Updated))
}

func printSection(w io.Writer, label string, keys []string) {
	if len(keys) == 0 {
		return
	}
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	fmt.Fprintf(w, "\n%s (%d):\n", label, len(sorted))
	for _, k := range sorted {
		fmt.Fprintf(w, "  • %s\n", k)
	}
}
