package dotenv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LineageEntry records a single mutation event for a secret key.
type LineageEntry struct {
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value"`
	Source    string    `json:"source"`
	Operation string    `json:"operation"` // "add", "update", "delete"
	Timestamp time.Time `json:"timestamp"`
}

// Lineage holds the full history of mutations for a .env file.
type Lineage struct {
	Entries []LineageEntry `json:"entries"`
}

// LineagePath returns the path to the lineage file for a given .env path.
func LineagePath(envPath string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, "."+base+".lineage.json")
}

// RecordLineage appends a set of diff results to the lineage file.
func RecordLineage(envPath, source string, results []DiffResult) error {
	path := LineagePath(envPath)

	lin := &Lineage{}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, lin)
	}

	now := time.Now().UTC()
	for _, r := range results {
		if r.Status == DiffUnchanged {
			continue
		}
		op := "update"
		if r.Status == DiffAdded {
			op = "add"
		} else if r.Status == DiffRemoved {
			op = "delete"
		}
		lin.Entries = append(lin.Entries, LineageEntry{
			Key:       r.Key,
			OldValue:  r.OldValue,
			NewValue:  r.NewValue,
			Source:    source,
			Operation: op,
			Timestamp: now,
		})
	}

	data, err := json.MarshalIndent(lin, "", "  ")
	if err != nil {
		return fmt.Errorf("lineage: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadLineage reads the lineage file for the given .env path.
// Returns an empty Lineage if the file does not exist.
func LoadLineage(envPath string) (*Lineage, error) {
	path := LineagePath(envPath)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Lineage{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("lineage: read: %w", err)
	}
	var lin Lineage
	if err := json.Unmarshal(data, &lin); err != nil {
		return nil, fmt.Errorf("lineage: parse: %w", err)
	}
	return &lin, nil
}

// KeyHistory returns all lineage entries for a specific key, oldest first.
func (l *Lineage) KeyHistory(key string) []LineageEntry {
	var out []LineageEntry
	for _, e := range l.Entries {
		if e.Key == key {
			out = append(out, e)
		}
	}
	return out
}
