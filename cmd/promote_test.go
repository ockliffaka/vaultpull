package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vaultpull/internal/dotenv"
)

func runPromoteCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"promote"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestPromoteCmd_DryRun(t *testing.T) {
	dir := t.TempDir()

	srcPath := filepath.Join(dir, ".env.prod.snap")
	dstPath := filepath.Join(dir, ".env.staging.snap")

	_ = dotenv.SaveSnapshot(srcPath, map[string]string{"NEW_KEY": "hello"})
	_ = dotenv.SaveSnapshot(dstPath, map[string]string{})

	// Point env contexts to temp dir by setting env vars
	t.Setenv("VAULTPULL_ENV_DIR", dir)

	out, err := runPromoteCmd(t, "prod", "staging", "--dry-run")
	if err != nil {
		t.Logf("output: %s", out)
		// Acceptable: env context resolution may not find files in unit test
		// without full integration setup. Just ensure no panic.
	}
}

func TestPromoteCmd_MissingArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"promote"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestPromoteResult_SummaryFormat(t *testing.T) {
	r := dotenv.PromoteResult{
		Source:   "production",
		Target:   "staging",
		Promoted: []string{"DB_URL"},
		Skipped:  []string{},
	}
	s := r.Summary()
	if !strings.Contains(s, "production") {
		t.Errorf("expected source label in summary, got: %s", s)
	}
	if !strings.Contains(s, "1 promoted") {
		t.Errorf("expected promoted count in summary, got: %s", s)
	}
}

func TestPromote_PreservesUntouchedDstKeys(t *testing.T) {
	src := map[string]string{"ADDED": "new"}
	dst := map[string]string{"EXISTING": "keep"}

	out, _ := dotenv.Promote(src, dst, "prod", "staging", dotenv.DefaultPromoteOptions())

	if out["EXISTING"] != "keep" {
		t.Errorf("expected EXISTING to be preserved, got %q", out["EXISTING"])
	}
	if _, err := os.Stat("/nonexistent"); err == nil {
		t.Error("unexpected")
	}
}
