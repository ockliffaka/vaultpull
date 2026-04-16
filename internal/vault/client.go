package vault

import (
	"errors"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	logical *vaultapi.Logical
}

// NewClient creates an authenticated Vault client using the given address and token.
func NewClient(address, token string) (*Client, error) {
	if address == "" {
		return nil, errors.New("vault address must not be empty")
	}
	if token == "" {
		return nil, errors.New("vault token must not be empty")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(token)

	return &Client{logical: client.Logical()}, nil
}

// ReadSecrets reads key/value secrets from the given path.
// Supports both KV v1 and KV v2 (detects "data" wrapper automatically).
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	data := secret.Data

	// KV v2 wraps values under a "data" key.
	if nested, ok := data["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			data = nestedMap
		}
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result, nil
}
