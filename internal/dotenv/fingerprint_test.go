package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFingerprintPath(t *testing.T) {
	got := FingerprintPath("/tmp/project/.env")
	want := "/tmp/project/."+".env.fingerprint.json"
	if got != want {
		t.Errorf("FingerprintPath = %q, want %q", got, want)
	}
}

func TestComputeFingerprint_Deterministic(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	fp1 := ComputeFingerprint(secrets)
	fp2 := ComputeFingerprint(secrets)
	if fp1.Hash != fp2.Hash {
		t.Error("expected identical hashes for same secrets")
	}
	if len(fp1.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(fp1.Keys))
	}
}

func TestComputeFingerprint_DifferentSecrets(t *testing.T) {
	a := ComputeFingerprint(map[string]string{"KEY": "value1"})
	b := ComputeFingerprint(map[string]string{"KEY": "value2"})
	if a.Hash == b.Hash {
		t.Error("expected different hashes for different secrets")
	}
}

func TestSaveAndLoadFingerprint_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	secrets := map[string]string{"API_KEY": "abc123", "REGION": "us-east-1"}
	fp := ComputeFingerprint(secrets)

	if err := SaveFingerprint(envPath, fp); err != nil {
		t.Fatalf("SaveFingerprint: %v", err)
	}

	loaded, err := LoadFingerprint(envPath)
	if err != nil {
		t.Fatalf("LoadFingerprint: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil fingerprint")
	}
	if loaded.Hash != fp.Hash {
		t.Errorf("hash mismatch: got %q, want %q", loaded.Hash, fp.Hash)
	}
}

func TestLoadFingerprint_MissingFile_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	fp, err := LoadFingerprint(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp != nil {
		t.Error("expected nil for missing fingerprint file")
	}
}

func TestChanged_NoFingerprint_ReturnsTrue(t *testing.T) {
	dir := t.TempDir()
	changed, err := Changed(filepath.Join(dir, ".env"), map[string]string{"X": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true when no fingerprint exists")
	}
}

func TestChanged_SameSecrets_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	_ = os.WriteFile(envPath, []byte(""), 0o600)

	secrets := map[string]string{"TOKEN": "secret"}
	fp := ComputeFingerprint(secrets)
	if err := SaveFingerprint(envPath, fp); err != nil {
		t.Fatalf("SaveFingerprint: %v", err)
	}

	changed, err := Changed(envPath, secrets)
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if changed {
		t.Error("expected changed=false for identical secrets")
	}
}

func TestChanged_DifferentSecrets_ReturnsTrue(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	old := map[string]string{"TOKEN": "old"}
	if err := SaveFingerprint(envPath, ComputeFingerprint(old)); err != nil {
		t.Fatalf("SaveFingerprint: %v", err)
	}

	changed, err := Changed(envPath, map[string]string{"TOKEN": "new"})
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for different secrets")
	}
}
