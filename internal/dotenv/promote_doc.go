// Package dotenv provides utilities for managing .env files, including
// reading, writing, merging, diffing, encrypting, and promoting secrets
// between environments.
//
// # Promote
//
// The Promote function copies secrets from a source environment map into a
// destination environment map. It supports:
//
//   - Selective key promotion via PromoteOptions.Keys
//   - Safe defaults that skip existing keys (Overwrite: false)
//   - Dry-run mode that reports changes without applying them
//
// Example:
//
//	out, result := dotenv.Promote(prodSecrets, stagingSecrets, "prod", "staging",
//		dotenv.DefaultPromoteOptions())
//	fmt.Println(result.Summary())
package dotenv
