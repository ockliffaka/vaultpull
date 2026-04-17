// Package dotenv provides utilities for reading, writing, merging,
// formatting, diffing, previewing, validating, and sanitizing .env files.
//
// Validate checks that secret keys conform to standard env var naming rules
// and that values do not contain characters that would break .env parsing.
//
// Sanitize trims whitespace and strips non-printable characters from secret
// values before they are written to disk, ensuring the resulting .env file
// is safe to source in shell environments.
package dotenv
