package dotenv

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Fingerprint holds a content hash and metadata for a secrets map.
type Fingerprint struct {
	Hash      string    `json:"hash"`
	Keys      []string  `json:"keys"`
	CreatedAt time.Time `json:"created_at"`
}

// FingerprintPath returns the path to the fingerprint file for a given .env file.
func FingerprintPath(envPath string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, "."+base+".fingerprint.json")
}

// ComputeFingerprint hashes the sorted key=value pairs of secrets.
func ComputeFingerprint(secrets map[string]string) Fingerprint {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}

	return Fingerprint{
		Hash:      hex.EncodeToString(h.Sum(nil)),
		Keys:      keys,
		CreatedAt: time.Now().UTC(),
	}
}

// SaveFingerprint writes a Fingerprint to disk alongside the .env file.
func SaveFingerprint(envPath string, fp Fingerprint) error {
	data, err := json.MarshalIndent(fp, "", "  ")
	if err != nil {
		return fmt.Errorf("fingerprint: marshal: %w", err)
	}
	return os.WriteFile(FingerprintPath(envPath), data, 0o600)
}

// LoadFingerprint reads a previously saved Fingerprint from disk.
// Returns nil, nil if no fingerprint file exists.
func LoadFingerprint(envPath string) (*Fingerprint, error) {
	data, err := os.ReadFile(FingerprintPath(envPath))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fingerprint: read: %w", err)
	}
	var fp Fingerprint
	if err := json.Unmarshal(data, &fp); err != nil {
		return nil, fmt.Errorf("fingerprint: unmarshal: %w", err)
	}
	return &fp, nil
}

// Changed returns true when secrets differ from the stored fingerprint.
// If no fingerprint exists on disk it is treated as changed.
func Changed(envPath string, secrets map[string]string) (bool, error) {
	stored, err := LoadFingerprint(envPath)
	if err != nil {
		return false, err
	}
	if stored == nil {
		return true, nil
	}
	current := ComputeFingerprint(secrets)
	return current.Hash != stored.Hash, nil
}
