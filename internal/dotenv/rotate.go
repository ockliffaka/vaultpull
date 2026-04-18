package dotenv

import (
	"fmt"
	"os"
	"time"
)

// RotateResult holds the outcome of a secret rotation.
type RotateResult struct {
	Path      string
	BackupPath string
	KeysUpdated int
	RotatedAt  time.Time
}

// Rotate backs up the existing .env file, merges new secrets (overwrite mode),
// writes a fresh expiry stamp, and returns a RotateResult.
func Rotate(path string, incoming map[string]string, ttl time.Duration) (*RotateResult, error) {
	backupPath, err := Backup(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("rotate: backup failed: %w", err)
	}

	existing, err := parse(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("rotate: parse failed: %w", err)
	}

	merged := Merge(existing, incoming, true)

	w, err := NewWriter(path)
	if err != nil {
		return nil, fmt.Errorf("rotate: writer failed: %w", err)
	}
	if err := w.Write(merged); err != nil {
		return nil, fmt.Errorf("rotate: write failed: %w", err)
	}

	stampPath := StampPath(path)
	if err := WriteStamp(stampPath, ttl); err != nil {
		return nil, fmt.Errorf("rotate: stamp failed: %w", err)
	}

	return &RotateResult{
		Path:        path,
		BackupPath:  backupPath,
		KeysUpdated: len(incoming),
		RotatedAt:   time.Now(),
	}, nil
}
