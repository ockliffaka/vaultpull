package dotenv

import (
	"fmt"
	"sort"
)

// PromoteOptions controls how secrets are promoted between environments.
type PromoteOptions struct {
	// Overwrite replaces existing keys in the target environment.
	Overwrite bool
	// DryRun reports what would change without writing anything.
	DryRun bool
	// Keys limits promotion to a specific subset of keys. Empty means all.
	Keys []string
}

// DefaultPromoteOptions returns safe defaults for promotion.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		Overwrite: false,
		DryRun:    false,
		Keys:      nil,
	}
}

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	Source   string
	Target   string
}

// Summary returns a human-readable summary of the promotion result.
func (r PromoteResult) Summary() string {
	return fmt.Sprintf(
		"promote %s -> %s: %d promoted, %d skipped",
		r.Source, r.Target, len(r.Promoted), len(r.Skipped),
	)
}

// Promote copies secrets from src into dst according to opts.
// It returns a PromoteResult describing what was changed.
func Promote(src, dst map[string]string, srcLabel, dstLabel string, opts PromoteOptions) (map[string]string, PromoteResult) {
	result := PromoteResult{Source: srcLabel, Target: dstLabel}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		if _, exists := dst[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		result.Promoted = append(result.Promoted, k)
		if !opts.DryRun {
			out[k] = v
		}
	}

	return out, result
}
