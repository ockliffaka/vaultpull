package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a label attached to a secret key for grouping or filtering.
type Tag struct {
	Key   string
	Value string
}

// TagMap maps secret keys to a list of tags.
type TagMap map[string][]Tag

// AddTag attaches a tag to the given secret key in the TagMap.
func (tm TagMap) AddTag(secretKey, tagKey, tagValue string) {
	tm[secretKey] = append(tm[secretKey], Tag{Key: tagKey, Value: tagValue})
}

// GetTags returns all tags for the given secret key.
func (tm TagMap) GetTags(secretKey string) []Tag {
	return tm[secretKey]
}

// FilterByTag returns all secret keys that have a tag matching the given key and value.
func (tm TagMap) FilterByTag(tagKey, tagValue string) []string {
	var matches []string
	for secretKey, tags := range tm {
		for _, t := range tags {
			if t.Key == tagKey && t.Value == tagValue {
				matches = append(matches, secretKey)
				break
			}
		}
	}
	sort.Strings(matches)
	return matches
}

// Summary returns a human-readable summary of all tags in the map.
func (tm TagMap) Summary() string {
	if len(tm) == 0 {
		return "no tags defined"
	}
	keys := make([]string, 0, len(tm))
	for k := range tm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		tagParts := make([]string, 0, len(tm[k]))
		for _, t := range tm[k] {
			tagParts = append(tagParts, fmt.Sprintf("%s=%s", t.Key, t.Value))
		}
		fmt.Fprintf(&sb, "%s: [%s]\n", k, strings.Join(tagParts, ", "))
	}
	return strings.TrimRight(sb.String(), "\n")
}
