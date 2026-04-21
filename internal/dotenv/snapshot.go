package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot captures the state of a secrets map at a point in time.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
}

// SnapshotPath returns the path to the snapshot file for a given env file.
func SnapshotPath(envPath string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, "."+base+".snapshot.json")
}

// SaveSnapshot writes a snapshot of the given secrets map to disk.
func SaveSnapshot(envPath string, secrets map[string]string) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Path:      envPath,
		Secrets:   secrets,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	dest := SnapshotPath(envPath)
	if err := os.WriteFile(dest, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// LoadSnapshot reads the most recent snapshot for the given env file.
// Returns nil, nil if no snapshot exists.
func LoadSnapshot(envPath string) (*Snapshot, error) {
	src := SnapshotPath(envPath)

	data, err := os.ReadFile(src)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	return &snap, nil
}

// DeleteSnapshot removes the snapshot file for the given env path, if present.
func DeleteSnapshot(envPath string) error {
	src := SnapshotPath(envPath)
	if err := os.Remove(src); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot: delete failed: %w", err)
	}
	return nil
}
