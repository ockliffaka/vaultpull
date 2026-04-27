package dotenv

import (
	"os"
	"testing"
)

func TestMerge_AddsNewKeys(t *testing.T) {
	f := tmpFile(t, "")
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	res, err := Merge(f, incoming, MergeKeepExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added != 2 || res.Updated != 0 || res.Skipped != 0 {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestMerge_KeepExistingSkipsConflicts(t *testing.T) {
	f := tmpFile(t, "FOO=original\n")
	incoming := map[string]string{"FOO": "new", "BAR": "added"}

	res, err := Merge(f, incoming, MergeKeepExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 || res.Added != 1 {
		t.Errorf("unexpected result: %+v", res)
	}

	data, _ := os.ReadFile(f)
	parsed := parse(string(data))
	if parsed["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %s", parsed["FOO"])
	}
}

func TestMerge_OverwriteUpdatesConflicts(t *testing.T) {
	f := tmpFile(t, "FOO=original\n")
	incoming := map[string]string{"FOO": "new"}

	res, err := Merge(f, incoming, MergeOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Updated != 1 {
		t.Errorf("expected 1 updated, got %+v", res)
	}

	data, _ := os.ReadFile(f)
	parsed := parse(string(data))
	if parsed["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", parsed["FOO"])
	}
}

func TestMerge_NonExistentFileCreatesIt(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.env"

	_, err := Merge(path, map[string]string{"KEY": "val"}, MergeKeepExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}
