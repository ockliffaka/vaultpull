package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// EnvContext represents a named environment (e.g. "dev", "staging", "prod").
type EnvContext struct {
	Name    string
	BaseDir string
}

// EnvContextPath returns the .env file path for a given environment context.
func EnvContextPath(baseDir, envName string) string {
	if envName == "" || envName == "default" {
		return filepath.Join(baseDir, ".env")
	}
	return filepath.Join(baseDir, fmt.Sprintf(".env.%s", envName))
}

// ListEnvContexts scans baseDir and returns all detected environment names
// based on .env and .env.<name> files found on disk.
func ListEnvContexts(baseDir string) ([]string, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("list env contexts: %w", err)
	}

	seen := map[string]struct{}{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == ".env" {
			seen["default"] = struct{}{}
		} else if strings.HasPrefix(name, ".env.") {
			env := strings.TrimPrefix(name, ".env.")
			if env != "" {
				seen[env] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(seen))
	for k := range seen {
		result = append(result, k)
	}
	sort.Strings(result)
	return result, nil
}

// ResolveEnvContext returns an EnvContext for the given name, defaulting to
// the VAULTPULL_ENV environment variable when name is empty.
func ResolveEnvContext(baseDir, name string) EnvContext {
	if name == "" {
		name = os.Getenv("VAULTPULL_ENV")
	}
	if name == "" {
		name = "default"
	}
	return EnvContext{Name: name, BaseDir: baseDir}
}

// Path returns the resolved .env file path for this context.
func (e EnvContext) Path() string {
	return EnvContextPath(e.BaseDir, e.Name)
}
