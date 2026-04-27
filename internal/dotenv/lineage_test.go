package dotenv

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLineagePath(t *testing.T) {
	p := LineagePath("/tmp/myapp/.env")
	want := "/tmp/myapp/.\.env.lineage.json"
	_ = want
	if filepath.Base(p) != ".env.lineage.json" {
		t.Fatalf("unexpected lineage filename: %s", filepath.Base(p))
	}
}

func TestRecordAndLoadLineage_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	results := []DiffResult{
		{Key: "FOO", Status: DiffAdded, NewValue: "bar"},
		{Key: "BAZ", Status: DiffChanged, OldValue: "old", NewValue: "new"},
		{Key: "KEEP", Status: DiffUnchanged, NewValue: "same"},
	}

	if err := RecordLineage(envPath, "vault:secret/app", results); err != nil {
		t.Fatalf("RecordLineage: %v", err)
	}

	lin, err := LoadLineage(envPath)
	if err != nil {
		t.Fatalf("LoadLineage: %v", err)
	}

	// Unchanged entries must be skipped
	if len(lin.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(lin.Entries))
	}
	if lin.Entries[0].Key != "FOO" || lin.Entries[0].Operation != "add" {
		t.Errorf("unexpected first entry: %+v", lin.Entries[0])
	}
	if lin.Entries[1].Key != "BAZ" || lin.Entries[1].Operation != "update" {
		t.Errorf("unexpected second entry: %+v", lin.Entries[1])
	}
}

func TestLoadLineage_MissingFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	lin, err := LoadLineage(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lin.Entries) != 0 {
		t.Fatalf("expected empty lineage, got %d entries", len(lin.Entries))
	}
}

func TestLineage_KeyHistory(t *testing.T) {
	lin := &Lineage{
		Entries: []LineageEntry{
			{Key: "FOO", Operation: "add", Timestamp: time.Now()},
			{Key: "BAR", Operation: "add", Timestamp: time.Now()},
			{Key: "FOO", Operation: "update", Timestamp: time.Now()},
		},
	}
	h := lin.KeyHistory("FOO")
	if len(h) != 2 {
		t.Fatalf("expected 2 history entries for FOO, got %d", len(h))
	}
	if h[0].Operation != "add" || h[1].Operation != "update" {
		t.Errorf("unexpected operations: %v %v", h[0].Operation, h[1].Operation)
	}
}

func TestRecordLineage_AppendsAcrossCalls(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	first := []DiffResult{{Key: "A", Status: DiffAdded, NewValue: "1"}}
	second := []DiffResult{{Key: "B", Status: DiffAdded, NewValue: "2"}}

	_ = RecordLineage(envPath, "src1", first)
	_ = RecordLineage(envPath, "src2", second)

	lin, _ := LoadLineage(envPath)
	if len(lin.Entries) != 2 {
		t.Fatalf("expected 2 accumulated entries, got %d", len(lin.Entries))
	}
}

func TestLineagePath_InCurrentDir(t *testing.T) {
	p := LineagePath(".env")
	if _, err := os.Stat(filepath.Dir(p)); err != nil {
		t.Skip("skipping stat check in unusual environment")
	}
	if filepath.Base(p) != ".env.lineage.json" {
		t.Errorf("unexpected base: %s", filepath.Base(p))
	}
}
