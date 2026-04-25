// Package dotenv provides utilities for reading, writing, and managing
// .env files used by the vaultpull CLI.
//
// # Environment Contexts
//
// vaultpull supports multiple named environments (e.g. dev, staging, prod)
// by mapping each to a distinct .env file on disk:
//
//	.env          → "default" context
//	.env.staging  → "staging" context
//	.env.prod     → "prod" context
//
// The active context is resolved in the following order:
//  1. Explicit --env flag passed to the CLI command
//  2. VAULTPULL_ENV environment variable
//  3. Falls back to "default"
//
// Use EnvContextPath to compute the file path for a given context, and
// ListEnvContexts to discover all contexts present in a directory.
package dotenv
