// Package dotenv provides utilities for reading and writing .env files.
package dotenv

import "sort"

// DiffResult holds the result of comparing two secret maps.
type DiffResult struct {
	Added    map[string]string
	Removed  map[string]string
	Changed  map[string]string
	Unchanged map[string]string
}

// Diff compares existing local keys against incoming vault secrets.
// existing is the current .env map, incoming is the vault secrets map.
func Diff(existing, incoming map[string]string) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, newVal := range incoming {
		oldVal, exists := existing[k]
		if !exists {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = newVal
		} else {
			result.Unchanged[k] = newVal
		}
	}

	for k, oldVal := range existing {
		if _, exists := incoming[k]; !exists {
			result.Removed[k] = oldVal
		}
	}

	return result
}

// Summary returns a human-readable summary of the diff.
func (d DiffResult) Summary() string {
	keys := func(m map[string]string) []string {
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return ks
	}

	var out string
	for _, k := range keys(d.Added) {
		out += "+ " + k + "\n"
	}
	for _, k := range keys(d.Changed) {
		out += "~ " + k + "\n"
	}
	for _, k := range keys(d.Removed) {
		out += "- " + k + "\n"
	}
	return out
}

// HasChanges returns true if there are any added, changed, or removed keys.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Changed) > 0 || len(d.Removed) > 0
}
