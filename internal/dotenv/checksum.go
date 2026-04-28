package dotenv

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ChecksumPath returns the path to the checksum file for a given .env file.
func ChecksumPath(envPath string) string {
	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	return filepath.Join(dir, "."+base+".sha256")
}

// ComputeChecksum computes a deterministic SHA-256 checksum over the
// key=value pairs in secrets, sorted by key.
func ComputeChecksum(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// SaveChecksum writes the checksum of secrets to the checksum file
// associated with envPath.
func SaveChecksum(envPath string, secrets map[string]string) error {
	sum := ComputeChecksum(secrets)
	return os.WriteFile(ChecksumPath(envPath), []byte(sum+"\n"), 0600)
}

// LoadChecksum reads the previously saved checksum for envPath.
// Returns an empty string and no error if the file does not exist.
func LoadChecksum(envPath string) (string, error) {
	data, err := os.ReadFile(ChecksumPath(envPath))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// ChecksumChanged returns true if the computed checksum of secrets differs
// from the stored checksum for envPath. Returns false (not changed) when no
// stored checksum exists yet.
func ChecksumChanged(envPath string, secrets map[string]string) (bool, error) {
	stored, err := LoadChecksum(envPath)
	if err != nil {
		return false, err
	}
	if stored == "" {
		return false, nil
	}
	current := ComputeChecksum(secrets)
	return current != stored, nil
}
