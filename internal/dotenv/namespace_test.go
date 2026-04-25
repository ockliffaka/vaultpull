package dotenv

import (
	"testing"
)

func TestGroupByNamespace_SingleDepth(t *testing.T) {
	secrets := map[string]string{
		"APP_DB_HOST": "localhost",
		"APP_DB_PORT": "5432",
		"APP_API_KEY": "abc123",
		"DEBUG":       "true",
	}

	nm := GroupByNamespace(secrets, "_", 1)

	if _, ok := nm["APP"]; !ok {
		t.Fatal("expected namespace APP")
	}
	if len(nm["APP"]) != 3 {
		t.Errorf("expected 3 keys under APP, got %d", len(nm["APP"]))
	}
	if _, ok := nm["_default"]; !ok {
		t.Fatal("expected _default namespace for keys without delimiter")
	}
	if nm["_default"]["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true in _default namespace")
	}
}

func TestGroupByNamespace_DoubleDepth(t *testing.T) {
	secrets := map[string]string{
		"APP_DB_HOST": "localhost",
		"APP_DB_PORT": "5432",
		"APP_API_KEY": "abc123",
	}

	nm := GroupByNamespace(secrets, "_", 2)

	if _, ok := nm["APP_DB"]; !ok {
		t.Fatal("expected namespace APP_DB")
	}
	if len(nm["APP_DB"]) != 2 {
		t.Errorf("expected 2 keys under APP_DB, got %d", len(nm["APP_DB"]))
	}
	if _, ok := nm["APP_API"]; !ok {
		t.Fatal("expected namespace APP_API")
	}
}

func TestGroupByNamespace_EmptyDelimiterDefaults(t *testing.T) {
	secrets := map[string]string{"FOO_BAR": "baz"}
	nm := GroupByNamespace(secrets, "", 1)
	if _, ok := nm["FOO"]; !ok {
		t.Error("expected FOO namespace with default delimiter")
	}
}

func TestListNamespaces_Sorted(t *testing.T) {
	nm := NamespaceMap{
		"Z_NS": {"Z_KEY": "v"},
		"A_NS": {"A_KEY": "v"},
		"M_NS": {"M_KEY": "v"},
	}
	ns := ListNamespaces(nm)
	if ns[0] != "A_NS" || ns[1] != "M_NS" || ns[2] != "Z_NS" {
		t.Errorf("expected sorted namespaces, got %v", ns)
	}
}

func TestFilterNamespace_ReturnsCorrectKeys(t *testing.T) {
	nm := NamespaceMap{
		"APP": {"APP_HOST": "localhost", "APP_PORT": "8080"},
	}
	filtered := FilterNamespace(nm, "APP")
	if filtered["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost")
	}
}

func TestFilterNamespace_MissingReturnsNil(t *testing.T) {
	nm := NamespaceMap{}
	if FilterNamespace(nm, "MISSING") != nil {
		t.Error("expected nil for missing namespace")
	}
}

func TestNamespaceSummary_Format(t *testing.T) {
	nm := NamespaceMap{
		"APP": {"APP_HOST": "localhost", "APP_PORT": "8080"},
		"DB":  {"DB_URL": "postgres://"},
	}
	summary := NamespaceSummary(nm)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	for _, ns := range []string{"APP", "DB"} {
		if !contains(summary, ns) {
			t.Errorf("expected summary to contain namespace %s", ns)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
