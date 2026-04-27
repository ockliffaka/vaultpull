package dotenv

import (
	"fmt"
	"sort"
)

// ChainOptions controls how multiple secret maps are merged in a chain.
type ChainOptions struct {
	// Later sources override earlier ones when Overwrite is true.
	Overwrite bool
	// Sources is an ordered list of named secret maps (first = lowest priority).
	Sources []ChainSource
}

// ChainSource represents a named map of secrets in the chain.
type ChainSource struct {
	Name    string
	Secrets map[string]string
}

// ChainResult holds the merged output and provenance metadata.
type ChainResult struct {
	Merged    map[string]string
	// Provenance maps each key to the name of the source that won.
	Provenance map[string]string
}

// DefaultChainOptions returns sensible defaults: later sources win.
func DefaultChainOptions() ChainOptions {
	return ChainOptions{Overwrite: true}
}

// Chain merges multiple secret sources according to opts.
// Sources are processed in order; with Overwrite=true the last source
// that defines a key wins. With Overwrite=false the first definition wins.
func Chain(opts ChainOptions) (ChainResult, error) {
	if len(opts.Sources) == 0 {
		return ChainResult{}, fmt.Errorf("chain: at least one source is required")
	}

	merged := make(map[string]string)
	provenance := make(map[string]string)

	for _, src := range opts.Sources {
		if src.Name == "" {
			return ChainResult{}, fmt.Errorf("chain: all sources must have a non-empty name")
		}
		for k, v := range src.Secrets {
			_, exists := merged[k]
			if !exists || opts.Overwrite {
				merged[k] = v
				provenance[k] = src.Name
			}
		}
	}

	return ChainResult{Merged: merged, Provenance: provenance}, nil
}

// ProvenanceSummary returns a human-readable summary of which source
// contributed each key, sorted alphabetically.
func (r ChainResult) ProvenanceSummary() []string {
	keys := make([]string, 0, len(r.Provenance))
	for k := range r.Provenance {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s <- %s", k, r.Provenance[k]))
	}
	return lines
}
