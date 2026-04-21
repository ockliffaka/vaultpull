package dotenv

import (
	"regexp"
	"testing"
)

func TestRedact_SensitiveKeysAreReplaced(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_NAME":    "vaultpull",
		"API_KEY":     "abc123",
	}
	opts := DefaultRedactOptions()
	got := Redact(secrets, opts)

	if got["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", got["DB_PASSWORD"])
	}
	if got["API_KEY"] != "***" {
		t.Errorf("expected API_KEY to be redacted, got %q", got["API_KEY"])
	}
	if got["APP_NAME"] != "vaultpull" {
		t.Errorf("expected APP_NAME to be preserved, got %q", got["APP_NAME"])
	}
}

func TestRedact_CaseInsensitiveMatch(t *testing.T) {
	secrets := map[string]string{
		"db_Password": "hunter2",
		"MYTOKEN":     "tok",
		"HOST":        "localhost",
	}
	got := Redact(secrets, DefaultRedactOptions())

	if got["db_Password"] != "***" {
		t.Errorf("expected db_Password redacted")
	}
	if got["MYTOKEN"] != "***" {
		t.Errorf("expected MYTOKEN redacted")
	}
	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST preserved")
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	secrets := map[string]string{"SECRET_KEY": "value"}
	opts := RedactOptions{Patterns: DefaultRedactPatterns, Placeholder: "<hidden>"}
	got := Redact(secrets, opts)
	if got["SECRET_KEY"] != "<hidden>" {
		t.Errorf("expected <hidden>, got %q", got["SECRET_KEY"])
	}
}

func TestRedact_EmptySecretsReturnsEmpty(t *testing.T) {
	got := Redact(map[string]string{}, DefaultRedactOptions())
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRedact_OriginalMapUnmodified(t *testing.T) {
	original := map[string]string{"DB_PASSWORD": "secret123", "HOST": "localhost"}
	Redact(original, DefaultRedactOptions())
	if original["DB_PASSWORD"] != "secret123" {
		t.Error("original map should not be modified")
	}
}

func TestRedactPattern_MatchingValues(t *testing.T) {
	secrets := map[string]string{
		"CREDIT_CARD": "4111-1111-1111-1111",
		"APP_ENV":     "production",
	}
	re := regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`)
	got := RedactPattern(secrets, re, "[redacted]")

	if got["CREDIT_CARD"] != "[redacted]" {
		t.Errorf("expected CREDIT_CARD redacted, got %q", got["CREDIT_CARD"])
	}
	if got["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV preserved, got %q", got["APP_ENV"])
	}
}
