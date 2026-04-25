package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRender_BasicSubstitution(t *testing.T) {
	secrets := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out, err := Render("postgres://${HOST}:${PORT}/db", secrets, DefaultTemplateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out != "postgres://localhost:5432/db" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_DollarSignStyle(t *testing.T) {
	secrets := map[string]string{"NAME": "world"}
	out, err := Render("hello $NAME", secrets, DefaultTemplateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello world" {
		t.Errorf("unexpected: %q", out)
	}
}

func TestRender_UnresolvedNonStrict(t *testing.T) {
	opts := DefaultTemplateOptions()
	opts.Placeholder = "MISSING"
	out, err := Render("value=${UNKNOWN}", map[string]string{}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if out != "value=MISSING" {
		t.Errorf("unexpected: %q", out)
	}
}

func TestRender_StrictMode_ReturnsError(t *testing.T) {
	opts := DefaultTemplateOptions()
	opts.Strict = true
	_, err := Render("${MISSING_KEY}", map[string]string{}, opts)
	if err == nil {
		t.Fatal("expected error for unresolved key in strict mode")
	}
}

func TestRenderFile_WritesOutput(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "app.env.tmpl")
	outPath := filepath.Join(dir, "app.env")

	_ = os.WriteFile(tmplPath, []byte("DB_URL=${DB_HOST}:${DB_PORT}"), 0o644)
	secrets := map[string]string{"DB_HOST": "db.local", "DB_PORT": "3306"}

	if err := RenderFile(tmplPath, secrets, DefaultTemplateOptions(), outPath); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(outPath)
	if string(data) != "DB_URL=db.local:3306" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestRenderFile_MissingTemplate(t *testing.T) {
	err := RenderFile("/nonexistent/file.tmpl", nil, DefaultTemplateOptions(), "")
	if err == nil {
		t.Fatal("expected error for missing template file")
	}
}

func TestRenderMap_CrossReference(t *testing.T) {
	secrets := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	out, err := RenderMap(secrets, DefaultTemplateOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("cross-reference not resolved: %q", out["API_URL"])
	}
}
