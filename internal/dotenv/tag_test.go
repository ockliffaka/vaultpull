package dotenv

import (
	"strings"
	"testing"
)

func TestAddTag_AndGetTags(t *testing.T) {
	tm := make(TagMap)
	tm.AddTag("DB_PASSWORD", "sensitivity", "high")
	tm.AddTag("DB_PASSWORD", "owner", "infra")

	tags := tm.GetTags("DB_PASSWORD")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Key != "sensitivity" || tags[0].Value != "high" {
		t.Errorf("unexpected first tag: %+v", tags[0])
	}
}

func TestGetTags_UnknownKey(t *testing.T) {
	tm := make(TagMap)
	tags := tm.GetTags("NONEXISTENT")
	if tags != nil {
		t.Errorf("expected nil tags for unknown key, got %v", tags)
	}
}

func TestFilterByTag_ReturnsMatchingKeys(t *testing.T) {
	tm := make(TagMap)
	tm.AddTag("DB_PASSWORD", "sensitivity", "high")
	tm.AddTag("API_KEY", "sensitivity", "high")
	tm.AddTag("LOG_LEVEL", "sensitivity", "low")

	result := tm.FilterByTag("sensitivity", "high")
	if len(result) != 2 {
		t.Fatalf("expected 2 matches, got %d: %v", len(result), result)
	}
	if result[0] != "API_KEY" || result[1] != "DB_PASSWORD" {
		t.Errorf("unexpected order or values: %v", result)
	}
}

func TestFilterByTag_NoMatches(t *testing.T) {
	tm := make(TagMap)
	tm.AddTag("DB_PASSWORD", "sensitivity", "high")

	result := tm.FilterByTag("owner", "frontend")
	if len(result) != 0 {
		t.Errorf("expected no matches, got %v", result)
	}
}

func TestSummary_EmptyTagMap(t *testing.T) {
	tm := make(TagMap)
	out := tm.Summary()
	if out != "no tags defined" {
		t.Errorf("unexpected summary for empty map: %q", out)
	}
}

func TestSummary_ContainsKeyAndTags(t *testing.T) {
	tm := make(TagMap)
	tm.AddTag("API_KEY", "sensitivity", "high")
	tm.AddTag("API_KEY", "owner", "backend")

	out := tm.Summary()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("summary missing key: %q", out)
	}
	if !strings.Contains(out, "sensitivity=high") {
		t.Errorf("summary missing tag: %q", out)
	}
	if !strings.Contains(out, "owner=backend") {
		t.Errorf("summary missing tag: %q", out)
	}
}
