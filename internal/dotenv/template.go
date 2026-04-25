package dotenv

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateOptions controls how template rendering behaves.
type TemplateOptions struct {
	// Strict causes Render to return an error if any placeholder is unresolved.
	Strict bool
	// Placeholder is the default value for unresolved keys when Strict is false.
	Placeholder string
}

// DefaultTemplateOptions returns sensible defaults.
func DefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{
		Strict:      false,
		Placeholder: "",
	}
}

// varPattern matches ${KEY} or $KEY style references.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Render replaces variable references in the template string with values from
// secrets. References take the form ${KEY} or $KEY.
func Render(tmpl string, secrets map[string]string, opts TemplateOptions) (string, error) {
	var missing []string
	result := varPattern.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := strings.TrimPrefix(strings.TrimPrefix(match, "${"), "$")
		key = strings.TrimSuffix(key, "}")
		if val, ok := secrets[key]; ok {
			return val
		}
		missing = append(missing, key)
		return opts.Placeholder
	})
	if opts.Strict && len(missing) > 0 {
		return "", fmt.Errorf("template: unresolved variables: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// RenderFile reads a template file, renders it with secrets, and writes the
// result to outPath. If outPath is empty the result is written to stdout.
func RenderFile(templatePath string, secrets map[string]string, opts TemplateOptions, outPath string) error {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template: read %s: %w", templatePath, err)
	}
	output, err := Render(string(data), secrets, opts)
	if err != nil {
		return err
	}
	if outPath == "" {
		fmt.Print(output)
		return nil
	}
	return os.WriteFile(outPath, []byte(output), 0o600)
}

// RenderMap applies Render to every value in secrets itself, allowing
// cross-referencing between keys (single pass, no cycles).
func RenderMap(secrets map[string]string, opts TemplateOptions) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	var buf bytes.Buffer
	for k, v := range secrets {
		buf.Reset()
		resolved, err := Render(v, secrets, opts)
		if err != nil {
			return nil, fmt.Errorf("template: key %s: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}
