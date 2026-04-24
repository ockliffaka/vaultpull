package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PinRecord stores a pinned snapshot of secrets at a specific version.
type PinRecord struct {
	Label     string            `json:"label"`
	PinnedAt  time.Time         `json:"pinned_at"`
	Secrets   map[string]string `json:"secrets"`
	SourcePath string           `json:"source_path"`
}

// PinPath returns the path to the pin file for a given env file and label.
func PinPath(envPath, label string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, fmt.Sprintf(".%s.pin.%s.json", base, label))
}

// Pin saves the current secrets map as a named pin record.
func Pin(envPath, label string, secrets map[string]string) error {
	if label == "" {
		return fmt.Errorf("pin label must not be empty")
	}
	record := PinRecord{
		Label:      label,
		PinnedAt:   time.Now().UTC(),
		Secrets:    secrets,
		SourcePath: envPath,
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal failed: %w", err)
	}
	path := PinPath(envPath, label)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("pin: write failed: %w", err)
	}
	return nil
}

// LoadPin reads a named pin record for the given env file.
func LoadPin(envPath, label string) (*PinRecord, error) {
	path := PinPath(envPath, label)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("pin: read failed: %w", err)
	}
	var record PinRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("pin: unmarshal failed: %w", err)
	}
	return &record, nil
}

// DeletePin removes a named pin record.
func DeletePin(envPath, label string) error {
	path := PinPath(envPath, label)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("pin: delete failed: %w", err)
	}
	return nil
}
