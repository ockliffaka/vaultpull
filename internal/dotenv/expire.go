package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ExpiryRecord holds metadata about when a .env file was last synced.
type ExpiryRecord struct {
	Path      string
	SyncedAt  time.Time
	MaxAge    time.Duration
}

// IsExpired returns true if the sync timestamp is older than MaxAge.
func (e ExpiryRecord) IsExpired() bool {
	return time.Since(e.SyncedAt) > e.MaxAge
}

// StampPath returns the path to the stamp file for a given .env file.
func StampPath(envPath string) string {
	return envPath + ".synced"
}

// WriteStamp writes the current UTC time to a stamp file next to the .env file.
func WriteStamp(envPath string) error {
	stamp := StampPath(envPath)
	if err := os.MkdirAll(filepath.Dir(stamp), 0o755); err != nil {
		return fmt.Errorf("expire: mkdir: %w", err)
	}
	f, err := os.Create(stamp)
	if err != nil {
		return fmt.Errorf("expire: create stamp: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d", time.Now().UTC().Unix())
	return err
}

// ReadStamp reads the stamp file and returns an ExpiryRecord.
func ReadStamp(envPath string, maxAge time.Duration) (ExpiryRecord, error) {
	stamp := StampPath(envPath)
	data, err := os.ReadFile(stamp)
	if err != nil {
		if os.IsNotExist(err) {
			return ExpiryRecord{Path: envPath, MaxAge: maxAge}, nil
		}
		return ExpiryRecord{}, fmt.Errorf("expire: read stamp: %w", err)
	}
	var unix int64
	if _, err := fmt.Sscanf(string(data), "%d", &unix); err != nil {
		return ExpiryRecord{}, fmt.Errorf("expire: parse stamp: %w", err)
	}
	return ExpiryRecord{
		Path:     envPath,
		SyncedAt: time.Unix(unix, 0).UTC(),
		MaxAge:   maxAge,
	}, nil
}
