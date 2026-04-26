package dotenv

import (
	"strings"
	"testing"
)

func TestPromote_AddsNewKeys(t *testing.T) {
	src := map[string]string{"DB_HOST": "prod-db", "API_KEY": "secret"}
	dst := map[string]string{}

	out, result := Promote(src, dst, "prod", "staging", DefaultPromoteOptions())

	if len(result.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(result.Promoted))
	}
	if out["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %s", out["DB_HOST"])
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"DB_HOST": "prod-db"}
	dst := map[string]string{"DB_HOST": "staging-db"}

	out, result := Promote(src, dst, "prod", "staging", DefaultPromoteOptions())

	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if out["DB_HOST"] != "staging-db" {
		t.Errorf("expected original value preserved, got %s", out["DB_HOST"])
	}
}

func TestPromote_OverwriteReplacesKeys(t *testing.T) {
	src := map[string]string{"DB_HOST": "prod-db"}
	dst := map[string]string{"DB_HOST": "staging-db"}
	opts := PromoteOptions{Overwrite: true}

	out, result := Promote(src, dst, "prod", "staging", opts)

	if len(result.Promoted) != 1 {
		t.Errorf("expected 1 promoted, got %d", len(result.Promoted))
	}
	if out["DB_HOST"] != "prod-db" {
		t.Errorf("expected overwritten value, got %s", out["DB_HOST"])
	}
}

func TestPromote_DryRunDoesNotMutate(t *testing.T) {
	src := map[string]string{"NEW_KEY": "value"}
	dst := map[string]string{}
	opts := PromoteOptions{DryRun: true}

	out, result := Promote(src, dst, "prod", "staging", opts)

	if len(result.Promoted) != 1 {
		t.Errorf("expected 1 in promoted list, got %d", len(result.Promoted))
	}
	if _, ok := out["NEW_KEY"]; ok {
		t.Error("expected dry run to not write key")
	}
}

func TestPromote_FilteredKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	opts := PromoteOptions{Keys: []string{"A", "C"}}

	out, result := Promote(src, dst, "prod", "staging", opts)

	if len(result.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(result.Promoted))
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be excluded")
	}
}

func TestPromoteResult_Summary(t *testing.T) {
	r := PromoteResult{
		Source:   "prod",
		Target:   "staging",
		Promoted: []string{"A", "B"},
		Skipped:  []string{"C"},
	}
	s := r.Summary()
	if !strings.Contains(s, "prod") || !strings.Contains(s, "staging") {
		t.Errorf("summary missing env labels: %s", s)
	}
	if !strings.Contains(s, "2 promoted") {
		t.Errorf("summary missing promoted count: %s", s)
	}
}
