package dotenv

import (
	"strings"
	"testing"
)

func TestValidateSchema_AllPresent(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	rules := []SchemaRule{
		{Key: "DB_HOST", Required: true},
		{Key: "DB_PORT", Required: true},
	}
	res := ValidateSchema(secrets, rules)
	if !res.OK() {
		t.Errorf("expected OK, got: %s", res.Summary())
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	rules := []SchemaRule{
		{Key: "DB_HOST", Required: true},
		{Key: "DB_PASS", Required: true},
	}
	res := ValidateSchema(secrets, rules)
	if res.OK() {
		t.Fatal("expected not OK")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_PASS" {
		t.Errorf("unexpected missing: %v", res.Missing)
	}
}

func TestValidateSchema_OptionalWarning(t *testing.T) {
	secrets := map[string]string{}
	rules := []SchemaRule{
		{Key: "LOG_LEVEL", Required: false},
	}
	res := ValidateSchema(secrets, rules)
	if !res.OK() {
		t.Fatal("expected OK for optional key")
	}
	if len(res.Warnings) == 0 {
		t.Error("expected a warning for missing optional key")
	}
}

func TestValidateSchema_PatternMatch(t *testing.T) {
	secrets := map[string]string{"PORT": "8080"}
	rules := []SchemaRule{
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	}
	res := ValidateSchema(secrets, rules)
	if !res.OK() {
		t.Errorf("expected OK, got: %s", res.Summary())
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	secrets := map[string]string{"PORT": "not-a-port"}
	rules := []SchemaRule{
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	}
	res := ValidateSchema(secrets, rules)
	if res.OK() {
		t.Fatal("expected not OK due to pattern mismatch")
	}
	if len(res.Invalid) == 0 {
		t.Error("expected invalid entry for PORT")
	}
}

func TestSchemaResult_Summary_OK(t *testing.T) {
	res := SchemaResult{}
	if res.Summary() != "schema OK" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestSchemaResult_Summary_ContainsDetails(t *testing.T) {
	res := SchemaResult{
		Missing: []string{"SECRET_KEY"},
		Invalid: []string{"PORT does not match pattern"},
	}
	s := res.Summary()
	if !strings.Contains(s, "SECRET_KEY") {
		t.Error("summary should mention missing key")
	}
	if !strings.Contains(s, "PORT") {
		t.Error("summary should mention invalid key")
	}
}
