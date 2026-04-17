package dotenv

import (
	"strings"
	"testing"
)

func TestValidate_ValidSecrets(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"_PRIVATE_KEY": "abc123",
	}
	if err := Validate(secrets); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_InvalidKey(t *testing.T) {
	secrets := map[string]string{
		"1INVALID": "value",
	}
	err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
	if !strings.Contains(err.Error(), "invalid key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	secrets := map[string]string{
		"": "value",
	}
	err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	if !strings.Contains(err.Error(), "empty key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_NewlineInValue(t *testing.T) {
	secrets := map[string]string{
		"SECRET": "line1\nline2",
	}
	err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error for newline in value")
	}
	if !strings.Contains(err.Error(), "newline") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	secrets := map[string]string{
		"bad-key": "ok",
		"ANOTHER_BAD": "val\r\n",
	}
	err := Validate(secrets)
	if err == nil {
		t.Fatal("expected error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(ve.Issues))
	}
}
