package dotenv

import (
	"fmt"
	"io"
	"sort"
)

// PrintPreview writes a human-readable diff preview to w.
// It masks secret values by default unless showValues is true.
func PrintPreview(w io.Writer, diff DiffResult, showValues bool) {
	mask := func(v string) string {
		if showValues {
			return v
		}
		if len(v) == 0 {
			return ""
		}
		return "****"
	}

	printSection := func(label string, m map[string]string) {
		if len(m) == 0 {
			return
		}
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "%s %s=%s\n", label, k, mask(m[k]))
		}
	}

	printSection("+", diff.Added)
	printSection("~", diff.Changed)
	printSection("-", diff.Removed)

	if !diff.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
	}
}
