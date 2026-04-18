package dotenv

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockFile represents a file-based lock to prevent concurrent writes.
type LockFile struct {
	path string
}

// LockPath returns the lock file path for a given env file.
func LockPath(envPath string) string {
	return envPath + ".lock"
}

// AcquireLock creates a lock file for the given env file path.
// Returns an error if a lock already exists.
func AcquireLock(envPath string) (*LockFile, error) {
	lp := LockPath(envPath)

	if _, err := os.Stat(lp); err == nil {
		data, _ := os.ReadFile(lp)
		return nil, fmt.Errorf("lock already held: %s", string(data))
	}

	info := fmt.Sprintf("pid=%d time=%s", os.Getpid(), time.Now().Format(time.RFC3339))
	if err := os.WriteFile(lp, []byte(info), 0600); err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return &LockFile{path: lp}, nil
}

// Release removes the lock file.
func (l *LockFile) Release() error {
	if err := os.Remove(l.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return nil
}

// StaleLockAge is the duration after which a lock is considered stale.
const StaleLockAge = 5 * time.Minute

// ClearStaleLock removes a lock file if it is older than StaleLockAge.
func ClearStaleLock(envPath string) (bool, error) {
	lp := LockPath(envPath)
	info, err := os.Stat(lp)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if time.Since(info.ModTime()) > StaleLockAge {
		if err := os.Remove(lp); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// ensure LockPath uses filepath for cross-platform safety
var _ = filepath.Join
