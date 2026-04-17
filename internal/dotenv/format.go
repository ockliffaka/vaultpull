package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// FormatOptions controls how the .env output is rendered.
type FormatOptions struct {
	// SortKeys sorts keys alphabetically in the output.
	SortKeys bool
	// Header is an optional comment block prepended to the file.
	Header string
}

// Format renders a map of key/value pairs into a .env-compatible string.
func Format(secrets map[string]string, opts FormatOptions) string {
	var sb strings.Builder

	if opts.Header != "" {
		for _, line := range strings.Split(opts.Header, "\n") {
			fmt.Fprintf(&sb, "# %s\n", line)
		}
		sb.WriteString("\n")
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	for _, k := range keys {
		v := secrets[k]
		if strings.ContainsAny(v, " \t\n") {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}

	return sb.String()
}
