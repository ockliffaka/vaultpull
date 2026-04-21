package dotenv

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaRule defines an expected key and optional constraints.
type SchemaRule struct {
	Key      string
	Required bool
	Pattern  string // optional regex the value must match
}

// SchemaResult holds the outcome of a schema validation check.
type SchemaResult struct {
	Missing  []string
	Invalid  []string
	Warnings []string
}

// ValidateSchema checks that the provided secrets satisfy all schema rules.
func ValidateSchema(secrets map[string]string, rules []SchemaRule) SchemaResult {
	result := SchemaResult{}

	for _, rule := range rules {
		val, exists := secrets[rule.Key]

		if !exists || strings.TrimSpace(val) == "" {
			if rule.Required {
				result.Missing = append(result.Missing, rule.Key)
			} else {
				result.Warnings = append(result.Warnings, fmt.Sprintf("%s is optional but not set", rule.Key))
			}
			continue
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("%s has invalid pattern: %v", rule.Key, err))
				continue
			}
			if !re.MatchString(val) {
				result.Invalid = append(result.Invalid, fmt.Sprintf("%s does not match pattern %q", rule.Key, rule.Pattern))
			}
		}
	}

	return result
}

// OK returns true when there are no missing or invalid keys.
func (r SchemaResult) OK() bool {
	return len(r.Missing) == 0 && len(r.Invalid) == 0
}

// Summary returns a human-readable summary of the schema result.
func (r SchemaResult) Summary() string {
	var sb strings.Builder
	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("missing required keys: %s\n", strings.Join(r.Missing, ", ")))
	}
	if len(r.Invalid) > 0 {
		sb.WriteString(fmt.Sprintf("invalid values: %s\n", strings.Join(r.Invalid, ", ")))
	}
	if len(r.Warnings) > 0 {
		sb.WriteString(fmt.Sprintf("warnings: %s\n", strings.Join(r.Warnings, ", ")))
	}
	if sb.Len() == 0 {
		return "schema OK"
	}
	return strings.TrimRight(sb.String(), "\n")
}
