package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/dotenv"
)

func TestEnvContextPath_Default(t *testing.T) {
	got := dotenv.EnvContextPath("/app", "")
	want := "/app/.env"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEnvContextPath_Named(t *testing.T) {
	got := dotenv.EnvContextPath("/app", "staging")
	want := "/app/.env.staging"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestListEnvContexts_ReturnsAll(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{".env", ".env.dev", ".env.prod", "unrelated.txt"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("K=V"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	contexts, err := dotenv.ListEnvContexts(dir)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"default", "dev", "prod"}
	if len(contexts) != len(want) {
		t.Fatalf("got %v, want %v", contexts, want)
	}
	for i, w := range want {
		if contexts[i] != w {
			t.Errorf("index %d: got %q, want %q", i, contexts[i], w)
		}
	}
}

func TestListEnvContexts_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	contexts, err := dotenv.ListEnvContexts(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(contexts) != 0 {
		t.Errorf("expected empty, got %v", contexts)
	}
}

func TestResolveEnvContext_FallsBackToEnvVar(t *testing.T) {
	t.Setenv("VAULTPULL_ENV", "qa")
	ctx := dotenv.ResolveEnvContext("/app", "")
	if ctx.Name != "qa" {
		t.Errorf("got %q, want %q", ctx.Name, "qa")
	}
}

func TestResolveEnvContext_DefaultWhenEmpty(t *testing.T) {
	t.Setenv("VAULTPULL_ENV", "")
	ctx := dotenv.ResolveEnvContext("/app", "")
	if ctx.Name != "default" {
		t.Errorf("got %q, want %q", ctx.Name, "default")
	}
}

func TestEnvContext_Path(t *testing.T) {
	ctx := dotenv.EnvContext{Name: "prod", BaseDir: "/secrets"}
	got := ctx.Path()
	want := "/secrets/.env.prod"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
