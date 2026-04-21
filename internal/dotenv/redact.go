package dotenv

import (
	"regexp"
	"strings"
)

// RedactOptions controls which keys are redacted and how.
type RedactOptions struct {
	// Patterns is a list of key substring patterns (case-insensitive) that
	// trigger redaction. Defaults to DefaultRedactPatterns.
	Patterns []string
	// Placeholder replaces the secret value. Defaults to "***".
	Placeholder string
}

// DefaultRedactPatterns lists common key substrings considered sensitive.
var DefaultRedactPatterns = []string{
	"password", "passwd", "secret", "token", "api_key", "apikey",
	"private", "credential", "auth", "cert", "key",
}

// DefaultRedactOptions returns a RedactOptions with sensible defaults.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Patterns:    DefaultRedactPatterns,
		Placeholder: "***",
	}
}

// Redact returns a copy of secrets where values whose keys match any of the
// configured patterns are replaced with the placeholder string.
func Redact(secrets map[string]string, opts RedactOptions) map[string]string {
	if opts.Placeholder == "" {
		opts.Placeholder = "***"
	}
	if len(opts.Patterns) == 0 {
		opts.Patterns = DefaultRedactPatterns
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if isSensitiveKey(k, opts.Patterns) {
			result[k] = opts.Placeholder
		} else {
			result[k] = v
		}
	}
	return result
}

// isSensitiveKey reports whether key matches any of the given patterns
// (case-insensitive substring match).
func isSensitiveKey(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// RedactPattern returns a copy of secrets where values matching the given
// regular expression are replaced with the placeholder.
func RedactPattern(secrets map[string]string, pattern *regexp.Regexp, placeholder string) map[string]string {
	if placeholder == "" {
		placeholder = "***"
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if pattern.MatchString(v) {
			result[k] = placeholder
		} else {
			result[k] = v
		}
	}
	return result
}
