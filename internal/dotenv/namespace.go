package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// Namespace groups secret keys by a prefix delimiter, allowing structured
// access to flat key=value maps as if they were namespaced.
//
// For example, keys like APP_DB_HOST and APP_DB_PORT share the "APP_DB" namespace.

// NamespaceMap maps a namespace prefix to its constituent keys and values.
type NamespaceMap map[string]map[string]string

// GroupByNamespace partitions secrets into namespaces based on the given
// delimiter and depth. depth=1 groups by the first segment, depth=2 by the
// first two segments, etc.
//
// Keys that do not contain the delimiter are placed under the "_default"
// namespace.
func GroupByNamespace(secrets map[string]string, delimiter string, depth int) NamespaceMap {
	if delimiter == "" {
		delimiter = "_"
	}
	if depth < 1 {
		depth = 1
	}

	result := make(NamespaceMap)

	for k, v := range secrets {
		parts := strings.SplitN(k, delimiter, depth+1)
		var ns string
		if len(parts) <= depth {
			ns = "_default"
		} else {
			ns = strings.Join(parts[:depth], delimiter)
		}
		if result[ns] == nil {
			result[ns] = make(map[string]string)
		}
		result[ns][k] = v
	}

	return result
}

// ListNamespaces returns the sorted namespace keys present in a NamespaceMap.
func ListNamespaces(nm NamespaceMap) []string {
	keys := make([]string, 0, len(nm))
	for k := range nm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FilterNamespace returns only the secrets belonging to the given namespace.
// Returns nil if the namespace does not exist.
func FilterNamespace(nm NamespaceMap, ns string) map[string]string {
	return nm[ns]
}

// NamespaceSummary returns a human-readable summary of the namespace map.
func NamespaceSummary(nm NamespaceMap) string {
	var sb strings.Builder
	for _, ns := range ListNamespaces(nm) {
		fmt.Fprintf(&sb, "[%s] %d key(s)\n", ns, len(nm[ns]))
	}
	return strings.TrimRight(sb.String(), "\n")
}
