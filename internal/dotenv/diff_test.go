package dotenv

import (
	"strings"
	"testing"
)

func TestDiff_AddedKeys(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := Diff(existing, incoming)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added key, got %d", len(result.Added))
	}
	if result.Added["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %s", result.Added["NEW_KEY"])
	}
}

func TestDiff_RemovedKeys(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "OLD_KEY": "old"}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(existing, incoming)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed key, got %d", len(result.Removed))
	}
	if result.Removed["OLD_KEY"] != "old" {
		t.Errorf("expected OLD_KEY=old, got %s", result.Removed["OLD_KEY"])
	}
}

func TestDiff_ChangedKeys(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	result := Diff(existing, incoming)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(result.Changed))
	}
	if result.Changed["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", result.Changed["FOO"])
	}
}

func TestDiff_UnchangedKeys(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(existing, incoming)

	if len(result.Unchanged) != 1 {
		t.Fatalf("expected 1 unchanged key, got %d", len(result.Unchanged))
	}
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestDiff_Summary(t *testing.T) {
	existing := map[string]string{"OLD": "v", "SAME": "v"}
	incoming := map[string]string{"NEW": "v", "SAME": "v"}

	result := Diff(existing, incoming)
	summary := result.Summary()

	if !strings.Contains(summary, "+ NEW") {
		t.Errorf("expected '+ NEW' in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "- OLD") {
		t.Errorf("expected '- OLD' in summary, got: %s", summary)
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	result := Diff(map[string]string{}, map[string]string{})
	if result.HasChanges() {
		t.Error("expected no changes for empty maps")
	}
}
