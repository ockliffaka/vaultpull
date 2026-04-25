package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpull/cmd"
)

func runEnvCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := cmd.NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestEnvList_NoContexts(t *testing.T) {
	dir := t.TempDir()
	out, err := runEnvCmd(t, "env", "list", "--dir", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No environment contexts found") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestEnvList_WithContexts(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{".env", ".env.dev", ".env.prod"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("K=V"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	out, err := runEnvCmd(t, "env", "list", "--dir", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"default", "dev", "prod"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}

func TestEnvShow_ExplicitName(t *testing.T) {
	dir := t.TempDir()
	out, err := runEnvCmd(t, "env", "show", "staging", "--dir", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected staging in output: %q", out)
	}
	if !strings.Contains(out, ".env.staging") {
		t.Errorf("expected .env.staging path in output: %q", out)
	}
}

func TestEnvShow_DefaultContext(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("VAULTPULL_ENV", "")
	out, err := runEnvCmd(t, "env", "show", "--dir", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "default") {
		t.Errorf("expected default context in output: %q", out)
	}
}

// Ensure cobra wiring compiles.
var _ *cobra.Command
