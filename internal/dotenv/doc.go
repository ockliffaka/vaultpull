// Package dotenv provides functionality for reading and writing
// .env files used by vaultpull to persist secrets fetched from
// HashiCorp Vault.
//
// The Writer type handles safe serialisation of secret key-value
// pairs, respecting existing file contents and configurable
// overwrite behaviour. Files are written with mode 0600 to
// prevent unintended access by other OS users.
package dotenv
