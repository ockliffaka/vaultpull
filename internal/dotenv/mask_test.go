package dotenv

import (
	"strings"
	"testing"
)

func TestMaskValue_ShortValue(t *testing.T) {
	opts := DefaultMaskOptions()
	result := MaskValue("abc", opts)
	if result != opts.Placeholder {
		t.Errorf("expected %q, got %q", opts.Placeholder, result)
	}
}

func TestMaskValue_LongValue(t *testing.T) {
	opts := DefaultMaskOptions()
	result := MaskValue("supersecret", opts)
	if !strings.HasPrefix(result, "supe") {
		t.Errorf("expected prefix 'supe', got %q", result)
	}
	if !strings.HasSuffix(result, opts.Placeholder) {
		t.Errorf("expected suffix %q, got %q", opts.Placeholder, result)
	}
}

func TestMaskValue_EmptyValue(t *testing.T) {
	opts := DefaultMaskOptions()
	result := MaskValue("", opts)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestMaskValue_ZeroShowChars(t *testing.T) {
	opts := MaskOptions{ShowChars: 0, Placeholder: "[hidden]"}
	result := MaskValue("anyvalue", opts)
	if result != "[hidden]" {
		t.Errorf("expected '[hidden]', got %q", result)
	}
}

func TestMaskMap_AllMasked(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "s3cr3t",
		"API_KEY": "abcdefgh",
	}
	opts := DefaultMaskOptions()
	masked := MaskMap(secrets, opts)
	for k, v := range masked {
		if v == secrets[k] {
			t.Errorf("key %q was not masked", k)
		}
	}
}

func TestMaskMap_RevealKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "s3cr3t",
		"APP_ENV": "production",
	}
	opts := DefaultMaskOptions()
	masked := MaskMap(secrets, opts, "APP_ENV")
	if masked["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be revealed, got %q", masked["APP_ENV"])
	}
	if masked["DB_PASS"] == "s3cr3t" {
		t.Error("expected DB_PASS to be masked")
	}
}

func TestMaskMap_OriginalUnmodified(t *testing.T) {
	secrets := map[string]string{"TOKEN": "mytoken"}
	opts := DefaultMaskOptions()
	MaskMap(secrets, opts)
	if secrets["TOKEN"] != "mytoken" {
		t.Error("original map was modified")
	}
}

func TestMaskSummary(t *testing.T) {
	original := map[string]string{"A": "val1", "B": "val2", "C": "val3"}
	masked := map[string]string{"A": "va****", "B": "val2", "C": "va****"}
	summary := MaskSummary(original, masked)
	if !strings.Contains(summary, "2 of 3") {
		t.Errorf("unexpected summary: %q", summary)
	}
}
