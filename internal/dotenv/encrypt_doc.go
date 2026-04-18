// Package dotenv provides utilities for reading, writing, and managing
// .env files used by vaultpull to store secrets locally.
//
// # Encryption
//
// The encrypt.go file provides AES-GCM symmetric encryption helpers for
// protecting secret values at rest. These are intended for use when writing
// an encrypted snapshot of synced secrets to disk.
//
// Usage:
//
//	enc, err := dotenv.Encrypt("my-secret", encryptionKey)
//	plain, err := dotenv.Decrypt(enc, encryptionKey)
//
// Keys must be exactly 16, 24, or 32 bytes long (AES-128/192/256).
// Each call to Encrypt produces a unique ciphertext due to a random nonce.
package dotenv
