// Package dotenv provides utilities for managing .env files.
//
// # File Locking
//
// The lock module prevents concurrent writes to the same .env file by
// creating a lightweight lock file alongside the target file.
//
// Usage:
//
//	lock, err := dotenv.AcquireLock(".env")
//	if err != nil {
//	    log.Fatal("could not acquire lock:", err)
//	}
//	defer lock.Release()
//
// Stale locks (older than StaleLockAge) can be cleared automatically:
//
//	cleared, err := dotenv.ClearStaleLock(".env")
//
// Lock files are named <envfile>.lock and contain the PID and timestamp
// of the process that acquired them.
package dotenv
