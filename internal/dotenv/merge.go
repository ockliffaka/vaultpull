// Package dotenv provides utilities for reading and writing .env files.
package dotenv

import (
	"fmt"
	"os"
	"strings"
)

// MergeStrategy defines how conflicts between existing and new keys are handled.
type MergeStrategy int

const (
	// MergeKeepExisting preserves existing values when a key already exists.
	MergeKeepExisting MergeStrategy = iota
	// MergeOverwrite replaces existing values with new ones.
	MergeOverwrite
)

// MergeResult holds statistics about a merge operation.
type MergeResult struct {
	Added    int
	Updated  int
	Skipped  int
}

// Merge combines secrets into an existing .env file using the given strategy.
// It returns a MergeResult summarising what changed.
func Merge(path string, incoming map[string]string, strategy MergeStrategy) (MergeResult, error) {
	existing := map[string]string{}

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return MergeResult{}, fmt.Errorf("merge: read %s: %w", path, err)
	}
	if err == nil {
		existing = parse(string(data))
	}

	result := MergeResult{}
	for k, v := range incoming {
		if _, exists := existing[k]; exists {
			if strategy == MergeOverwrite {
				existing[k] = v
				result.Updated++
			} else {
				result.Skipped++
			}
		} else {
			existing[k] = v
			result.Added++
		}
	}

	var sb strings.Builder
	for k, v := range existing {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0600); err != nil {
		return MergeResult{}, fmt.Errorf("merge: write %s: %w", path, err)
	}

	return result, nil
}
