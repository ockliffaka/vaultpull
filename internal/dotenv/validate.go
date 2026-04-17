package dotenv

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError holds all issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d issue(s):\n  - %s", len(e.Issues), strings.Join(e.Issues, "\n  - "))
}

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks a map of key/value pairs for common .env issues.
// It returns a *ValidationError if any issues are found, or nil.
func Validate(secrets map[string]string) error {
	var issues []string

	for k, v := range secrets {
		if k == "" {
			issues = append(issues, "empty key found")
			continue
		}
		if !validKeyRe.MatchString(k) {
			issues = append(issues, fmt.Sprintf("invalid key %q: must match [A-Za-z_][A-Za-z0-9_]*", k))
		}
		if strings.ContainsAny(v, "\n\r") {
			issues = append(issues, fmt.Sprintf("key %q: value contains newline characters", k))
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
