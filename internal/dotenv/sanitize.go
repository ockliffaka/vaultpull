package dotenv

import (
	"strings"
	"unicode"
)

// Sanitize cleans up secret values to make them safe for .env files.
// It trims surrounding whitespace and removes non-printable characters.
func Sanitize(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = sanitizeValue(v)
	}
	return out
}

func sanitizeValue(s string) string {
	s = strings.TrimSpace(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsPrint(r) || r == '\t' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
