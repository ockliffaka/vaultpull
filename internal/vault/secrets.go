package vault

import (
	"context"
	"fmt"
	"strings"
)

// ReadSecrets reads key-value secrets from the given path.
// It supports both KV v1 and KV v2 secret engines.
func (c *Client) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	secret, err := c.logical.ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	data, ok := secret.Data["data"]
	if ok {
		// KV v2: secret.Data["data"] is the actual map
		nested, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data format at path %q (kv v2)", path)
		}
		return toStringMap(nested)
	}

	// KV v1: secret.Data is the actual map
	return toStringMap(secret.Data)
}

// ListSecretPaths lists all secret keys under the given prefix path.
func (c *Client) ListSecretPaths(ctx context.Context, prefix string) ([]string, error) {
	prefix = strings.TrimSuffix(prefix, "/") + "/"
	secret, err := c.logical.ListWithContext(ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", prefix, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secrets found under prefix %q", prefix)
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected keys format at prefix %q", prefix)
	}

	paths := make([]string, 0, len(keys))
	for _, k := range keys {
		if s, ok := k.(string); ok {
			paths = append(paths, prefix+s)
		}
	}
	return paths, nil
}

func toStringMap(m map[string]interface{}) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case string:
			out[k] = val
		default:
			out[k] = fmt.Sprintf("%v", v)
		}
	}
	return out, nil
}
