package dotenv_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/dotenv"
)

func TestRenderFile_ThenWrite_Integration(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "service.env.tmpl")
	outPath := filepath.Join(dir, "service.env")

	tmplContent := "API_KEY=${API_KEY}\nBASE_URL=${BASE_URL}\nDEBUG=false\n"
	_ = os.WriteFile(tmplPath, []byte(tmplContent), 0o644)

	secrets := map[string]string{
		"API_KEY":  "supersecret",
		"BASE_URL": "https://api.example.com",
	}

	err := dotenv.RenderFile(tmplPath, secrets, dotenv.DefaultTemplateOptions(), outPath)
	if err != nil {
		t.Fatalf("RenderFile: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "API_KEY=supersecret") {
		t.Errorf("missing API_KEY in output: %q", content)
	}
	if !strings.Contains(content, "BASE_URL=https://api.example.com") {
		t.Errorf("missing BASE_URL in output: %q", content)
	}
	if !strings.Contains(content, "DEBUG=false") {
		t.Errorf("static value lost: %q", content)
	}
}

func TestRenderMap_ThenFormat_Integration(t *testing.T) {
	secrets := map[string]string{
		"SCHEME":   "https",
		"HOST":     "example.com",
		"FULL_URL": "${SCHEME}://${HOST}",
	}

	resolved, err := dotenv.RenderMap(secrets, dotenv.DefaultTemplateOptions())
	if err != nil {
		t.Fatalf("RenderMap: %v", err)
	}

	if resolved["FULL_URL"] != "https://example.com" {
		t.Errorf("unexpected FULL_URL: %q", resolved["FULL_URL"])
	}

	formatted := dotenv.Format(resolved, "")
	if !strings.Contains(formatted, "FULL_URL=https://example.com") {
		t.Errorf("formatted output missing resolved value: %q", formatted)
	}
}
