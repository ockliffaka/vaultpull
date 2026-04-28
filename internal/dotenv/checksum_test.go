package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChecksumPath(t *testing.T) {
	got := ChecksumPath("/tmp/project/.env")
	want := "/tmp/project/."+".env.sha256"
	if got != want {
		t.Errorf("ChecksumPath = %q, want %q", got, want)
	}
}

func TestComputeChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	a := ComputeChecksum(secrets)
	b := ComputeChecksum(secrets)
	if a != b {
		t.Errorf("checksum not deterministic: %q vs %q", a, b)
	}
}

func TestComputeChecksum_OrderIndependent(t *testing.T) {
	a := ComputeChecksum(map[string]string{"A": "1", "B": "2"})
	b := ComputeChecksum(map[string]string{"B": "2", "A": "1"})
	if a != b {
		t.Errorf("checksum should be order-independent: %q vs %q", a, b)
	}
}

func TestComputeChecksum_DifferentSecrets(t *testing.T) {
	a := ComputeChecksum(map[string]string{"KEY": "value1"})
	b := ComputeChecksum(map[string]string{"KEY": "value2"})
	if a == b {
		t.Error("expected different checksums for different secrets")
	}
}

func TestSaveAndLoadChecksum_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	secrets := map[string]string{"TOKEN": "abc123", "HOST": "localhost"}
	if err := SaveChecksum(envPath, secrets); err != nil {
		t.Fatalf("SaveChecksum: %v", err)
	}

	loaded, err := LoadChecksum(envPath)
	if err != nil {
		t.Fatalf("LoadChecksum: %v", err)
	}

	expected := ComputeChecksum(secrets)
	if loaded != expected {
		t.Errorf("loaded checksum %q, want %q", loaded, expected)
	}
}

func TestLoadChecksum_MissingFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	sum, err := LoadChecksum(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != "" {
		t.Errorf("expected empty string, got %q", sum)
	}
}

func TestChecksumChanged_NoStoredChecksum(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	changed, err := ChecksumChanged(envPath, map[string]string{"A": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected false when no stored checksum exists")
	}
}

func TestChecksumChanged_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	original := map[string]string{"KEY": "old"}
	if err := SaveChecksum(envPath, original); err != nil {
		t.Fatalf("SaveChecksum: %v", err)
	}

	updated := map[string]string{"KEY": "new"}
	changed, err := ChecksumChanged(envPath, updated)
	if err != nil {
		t.Fatalf("ChecksumChanged: %v", err)
	}
	if !changed {
		t.Error("expected changed=true after updating secrets")
	}
}

func TestChecksumChanged_NoChangeWhenSame(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	secrets := map[string]string{"KEY": "value"}
	if err := SaveChecksum(envPath, secrets); err != nil {
		t.Fatalf("SaveChecksum: %v", err)
	}

	changed, err := ChecksumChanged(envPath, secrets)
	if err != nil {
		t.Fatalf("ChecksumChanged: %v", err)
	}
	if changed {
		t.Error("expected changed=false when secrets are identical")
	}
}

func TestChecksumPath_HiddenFile(t *testing.T) {
	p := ChecksumPath(".env")
	_ = os.Remove(p) // cleanup if exists
	if filepath.Base(p) != "."+".env.sha256" {
		t.Errorf("unexpected checksum filename: %s", filepath.Base(p))
	}
}
