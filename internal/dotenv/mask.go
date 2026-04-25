package dotenv

import (
	"fmt"
	"strings"
)

// MaskOptions controls how secret values are masked in output.
type MaskOptions struct {
	// ShowChars is the number of characters to reveal at the start of a value.
	ShowChars int
	// Placeholder replaces the hidden portion of the value.
	Placeholder string
}

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		ShowChars:   4,
		Placeholder: "****",
	}
}

// MaskValue masks a single secret value according to the given options.
// If the value is shorter than or equal to ShowChars, the entire value is replaced.
func MaskValue(value string, opts MaskOptions) string {
	if len(value) == 0 {
		return ""
	}
	if opts.ShowChars <= 0 || len(value) <= opts.ShowChars {
		return opts.Placeholder
	}
	return value[:opts.ShowChars] + opts.Placeholder
}

// MaskMap returns a copy of secrets with all values masked.
// Keys listed in revealKeys are left unmasked.
func MaskMap(secrets map[string]string, opts MaskOptions, revealKeys ...string) map[string]string {
	reveal := make(map[string]bool, len(revealKeys))
	for _, k := range revealKeys {
		reveal[strings.ToUpper(k)] = true
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if reveal[strings.ToUpper(k)] {
			result[k] = v
		} else {
			result[k] = MaskValue(v, opts)
		}
	}
	return result
}

// MaskSummary returns a human-readable summary of how many keys were masked.
func MaskSummary(original, masked map[string]string) string {
	count := 0
	for k, v := range original {
		if mv, ok := masked[k]; ok && mv != v {
			count++
		}
	}
	return fmt.Sprintf("%d of %d value(s) masked", count, len(original))
}

// IsMasked reports whether a value appears to have been masked using the given options.
// It checks whether the value ends with the placeholder and the visible prefix is
// shorter than or equal to ShowChars, or whether the value equals the placeholder entirely.
func IsMasked(value string, opts MaskOptions) bool {
	if value == opts.Placeholder {
		return true
	}
	if opts.ShowChars > 0 && strings.HasSuffix(value, opts.Placeholder) {
		prefix := strings.TrimSuffix(value, opts.Placeholder)
		return len(prefix) <= opts.ShowChars
	}
	return false
}
