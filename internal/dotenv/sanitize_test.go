package dotenv

import (
	"testing"
)

func TestSanitize_TrimsSurroundingWhitespace(t *testing.T) {
	in := map[string]string{"KEY": "  hello  "}
	out := Sanitize(in)
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
}

func TestSaPrintable(t *testing.T) {
	in := map[string]string{"KEY": "val\x00ue"}
	out := Sanitize(in)
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestSanitize_PreservesNormalValues(t *testing.T) {
	in := map[string]string{
		"DB_URL": "postgres://user:pass@localhost/db",
		"FLAG":   "true",
	}
	out := Sanitize(in)
	for k, v := range in {
		if out[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestSanitize_PreservesTab(t *testing.T) {
	in := map[string]string{"KEY": "val\tue"}
	out := Sanitize(in)
	if out["KEY"] != "val\tue" {
		t.Errorf("expected tab preserved, got %q", out["KEY"])
	}
}

func TestSanitize_EmptyValue(t *testing.T) {
	in := map[string]string{"KEY": ""}
	out := Sanitize(in)
	if out["KEY"] != "" {
		t.Errorf("expected empty string, got %q", out["KEY"])
	}
}
