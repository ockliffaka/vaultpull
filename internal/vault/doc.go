// Package vault provides a thin wrapper around the HashiCorp Vault API client
// for use by vaultpull. It handles authentication via a static token and
// supports reading secrets from both KV v1 and KV v2 secret engines.
//
// Usage:
//
//	client, err := vault.NewClient("https://vault.example.com", "s.mytoken")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	secrets, err := client.ReadSecrets("secret/data/myapp")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for k, v := range secrets {
//		fmt.Printf("%s=%s\n", k, v)
//	}
package vault
