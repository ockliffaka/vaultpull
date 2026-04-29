package dotenv

import (
	"fmt"
	"os/exec"
	"strings"
)

// HookEvent represents the lifecycle event that triggers a hook.
type HookEvent string

const (
	HookPreSync  HookEvent = "pre-sync"
	HookPostSync HookEvent = "post-sync"
	HookPreWrite HookEvent = "pre-write"
	HookPostWrite HookEvent = "post-write"
)

// Hook defines a shell command to run on a given lifecycle event.
type Hook struct {
	Event   HookEvent
	Command string
	Args    []string
}

// HookResult captures the outcome of running a hook.
type HookResult struct {
	Event   HookEvent
	Command string
	Output  string
	Err     error
}

// RunHooks executes all hooks matching the given event in order.
// It returns a slice of results, one per hook. Execution continues
// even if a hook fails unless stopOnError is true.
func RunHooks(hooks []Hook, event HookEvent, stopOnError bool) ([]HookResult, error) {
	var results []HookResult

	for _, h := range hooks {
		if h.Event != event {
			continue
		}

		result := runHook(h)
		results = append(results, result)

		if result.Err != nil && stopOnError {
			return results, fmt.Errorf("hook %q failed: %w", h.Command, result.Err)
		}
	}

	return results, nil
}

func runHook(h Hook) HookResult {
	cmd := exec.Command(h.Command, h.Args...) //nolint:gosec
	out, err := cmd.CombinedOutput()
	return HookResult{
		Event:   h.Event,
		Command: strings.Join(append([]string{h.Command}, h.Args...), " "),
		Output:  strings.TrimSpace(string(out)),
		Err:     err,
	}
}

// HookSummary returns a human-readable summary of hook results.
func HookSummary(results []HookResult) string {
	if len(results) == 0 {
		return "no hooks ran"
	}

	var sb strings.Builder
	for _, r := range results {
		status := "ok"
		if r.Err != nil {
			status = fmt.Sprintf("error: %v", r.Err)
		}
		fmt.Fprintf(&sb, "[%s] %s -> %s\n", r.Event, r.Command, status)
	}
	return strings.TrimRight(sb.String(), "\n")
}
