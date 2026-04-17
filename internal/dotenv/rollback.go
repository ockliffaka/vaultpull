package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListBackups returns all backup files for the given env file, sorted newest first.
func ListBackups(envPath string) ([]string, error) {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var backups []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), base+".backup.") {
			backups = append(backups, filepath.Join(dir, e.Name()))
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(backups)))
	return backups, nil
}

// Rollback replaces envPath with the most recent backup file.
// Returns an error if no backup exists.
func Rollback(envPath string) (string, error) {
	backups, err := ListBackups(envPath)
	if err != nil {
		return "", err
	}
	if len(backups) == 0 {
		return "", fmt.Errorf("no backup found for %s", envPath)
	}

	latest := backups[0]
	data, err := os.ReadFile(latest)
	if err != nil {
		return "", fmt.Errorf("read backup: %w", err)
	}

	if err := os.WriteFile(envPath, data, 0600); err != nil {
		return "", fmt.Errorf("restore file: %w", err)
	}

	return latest, nil
}
