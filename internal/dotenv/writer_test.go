package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if content != "" {
		if _, err := f.WriteString(content); err != nil {
			t.Fatal(err)
		}
	}
	f.Close()
	return f.Name()
}

func TestWrite_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	w := NewWriter(path, false)
	if err := w.Write(map[string]string{"KEY": "value"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "KEY=value\n" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestWrite_PreservesExistingKeys(t *testing.T) {
	path := tmpFile(t, "EXISTING=old\n")
	w := NewWriter(path, false)
	if err := w.Write(map[string]string{"NEW": "val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := parse(readFile(t, path))
	if got["EXISTING"] != "old" {
		t.Errorf("expected EXISTING=old, got %q", got["EXISTING"])
	}
	if got["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %q", got["NEW"])
	}
}

func TestWrite_OverwriteReplacesKeys(t *testing.T) {
	path := tmpFile(t, "KEY=old\n")
	w := NewWriter(path, true)
	if err := w.Write(map[string]string{"KEY": "new"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := parse(readFile(t, path))
	if got["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %q", got["KEY"])
	}
}

func TestParse_SkipsComments(t *testing.T) {
	result := parse("# comment\nKEY=val\n")
	if _, ok := result["# comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if result["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", result["KEY"])
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
