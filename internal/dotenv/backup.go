// Package dotenv provides utilities for reading and writing .env files.
package dotenv

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Backup creates a timestamped backup of the given file.
// Returns the backup path or an empty string if the source does not exist.
func Backup(path string) (string, error) {
	src, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("backup: open source: %w", err)
	}
	defer src.Close()

	backupPath := fmt.Sprintf("%s.%s.bak", path, time.Now().UTC().Format("20060102T150405Z"))

	dst, err := os.OpenFile(backupPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return "", fmt.Errorf("backup: create destination: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(backupPath)
		return "", fmt.Errorf("backup: copy: %w", err)
	}

	return backupPath, nil
}
