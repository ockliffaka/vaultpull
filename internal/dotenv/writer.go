// Package dotenv provides utilities for writing secrets to .env files.
package dotenv

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer writes key-value secrets to a .env file.
type Writer struct {
	filePath  string
	overwrite bool
}

// NewWriter creates a new Writer targeting the given file path.
// If overwrite is true, existing values will be replaced.
func NewWriter(filePath string, overwrite bool) *Writer {
	return &Writer{filePath: filePath, overwrite: overwrite}
}

// Write serializes the provided secrets map into the .env file.
// Existing keys are preserved unless overwrite is enabled.
// If overwrite is false, existing keys in the file take precedence over new secrets.
func (w *Writer) Write(secrets map[string]string) error {
	existing := map[string]string{}

	if data, err := os.ReadFile(w.filePath); err == nil {
		existing = parse(string(data))
	}

	merged := make(map[string]string, len(secrets))
	for k, v := range secrets {
		merged[k] = v
	}

	if !w.overwrite {
		for k, v := range existing {
			merged[k] = v
		}
	}

	var sb strings.Builder
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, merged[k])
	}

	return os.WriteFile(w.filePath, []byte(sb.String()), 0600)
}

// parse reads a .env formatted string into a map.
func parse(content string) map[string]string {
	result := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}
