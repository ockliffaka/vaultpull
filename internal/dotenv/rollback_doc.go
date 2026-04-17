// Package dotenv provides utilities for managing .env files used by vaultpull.
//
// # Rollback
//
// The rollback feature allows users to restore a .env file to a previous state
// using backup files created during sync operations.
//
// Backup files follow the naming convention:
//
//	<envfile>.backup.<timestamp>
//
// Example:
//
//	.env.backup.20240101120000
//
// Use [ListBackups] to enumerate available backups sorted newest-first,
// and [Rollback] to restore the most recent backup in place.
package dotenv
