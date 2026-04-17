package dotenv

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintPreview_MasksValues(t *testing.T) {
	diff := DiffResult{
		Added:    map[string]string{"SECRET": "supersecret"},
		Removed:  map[string]string{},
		Changed:  map[string]string{},
		Unchanged: map[string]string{},
	}

	var buf bytes.Buffer
	PrintPreview(&buf, diff, false)
	out := buf.String()

	if strings.Contains(out, "supersecret") {
		t.Error("expected value to be masked")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected key SECRET to appear")
	}
	if !strings.Contains(out, "****") {
		t.Error("expected masked value ****")
	}
}

func TestPrintPreview_ShowValues(t *testing.T) {
	diff := DiffResult{
		Added:    map[string]string{"KEY": "plaintext"},
		Removed:  map[string]string{},
		Changed:  map[string]string{},
		Unchanged: map[string]string{},
	}

	var buf bytes.Buffer
	PrintPreview(&buf, diff, true)
	out := buf.String()

	if !strings.Contains(out, "plaintext") {
		t.Error("expected plaintext value to appear")
	}
}

func TestPrintPreview_NoChanges(t *testing.T) {
	diff := DiffResult{
		Added:     map[string]string{},
		Removed:   map[string]string{},
		Changed:   map[string]string{},
		Unchanged: map[string]string{"FOO": "bar"},
	}

	var buf bytes.Buffer
	PrintPreview(&buf, diff, false)
	out := buf.String()

	if !strings.Contains(out, "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestPrintPreview_AllCategories(t *testing.T) {
	diff := DiffResult{
		Added:     map[string]string{"NEW": "v"},
		Changed:   map[string]string{"MOD": "v"},
		Removed:   map[string]string{"OLD": "v"},
		Unchanged: map[string]string{},
	}

	var buf bytes.Buffer
	PrintPreview(&buf, diff, false)
	out := buf.String()

	for _, prefix := range []string{"+", "~", "-"} {
		if !strings.Contains(out, prefix) {
			t.Errorf("expected prefix %q in output", prefix)
		}
	}
}
