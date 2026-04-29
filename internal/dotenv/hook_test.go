package dotenv

import (
	"runtime"
	"strings"
	"testing"
)

func skipOnWindows(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell-command test on Windows")
	}
}

func TestRunHooks_MatchingEvent(t *testing.T) {
	skipOnWindows(t)
	hooks := []Hook{
		{Event: HookPostSync, Command: "echo", Args: []string{"hello"}},
	}
	results, err := RunHooks(hooks, HookPostSync, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Output != "hello" {
		t.Errorf("expected output 'hello', got %q", results[0].Output)
	}
}

func TestRunHooks_SkipsNonMatchingEvent(t *testing.T) {
	skipOnWindows(t)
	hooks := []Hook{
		{Event: HookPreSync, Command: "echo", Args: []string{"pre"}},
	}
	results, err := RunHooks(hooks, HookPostSync, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRunHooks_StopOnError(t *testing.T) {
	skipOnWindows(t)
	hooks := []Hook{
		{Event: HookPostWrite, Command: "false"},
		{Event: HookPostWrite, Command: "echo", Args: []string{"should not run"}},
	}
	results, err := RunHooks(hooks, HookPostWrite, true)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result before stop, got %d", len(results))
	}
}

func TestRunHooks_ContinueOnError(t *testing.T) {
	skipOnWindows(t)
	hooks := []Hook{
		{Event: HookPostWrite, Command: "false"},
		{Event: HookPostWrite, Command: "echo", Args: []string{"continued"}},
	}
	results, err := RunHooks(hooks, HookPostWrite, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[1].Output != "continued" {
		t.Errorf("expected second hook to run, got %q", results[1].Output)
	}
}

func TestHookSummary_NoResults(t *testing.T) {
	summary := HookSummary(nil)
	if summary != "no hooks ran" {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestHookSummary_WithResults(t *testing.T) {
	results := []HookResult{
		{Event: HookPostSync, Command: "echo hello", Output: "hello", Err: nil},
	}
	summary := HookSummary(results)
	if !strings.Contains(summary, "post-sync") {
		t.Errorf("expected event in summary, got: %q", summary)
	}
	if !strings.Contains(summary, "ok") {
		t.Errorf("expected 'ok' status in summary, got: %q", summary)
	}
}
