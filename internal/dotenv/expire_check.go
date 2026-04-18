package dotenv

import (
	"fmt"
	"time"
)

// ExpiryStatus describes whether secrets are fresh, stale, or expired.
type ExpiryStatus int

const (
	StatusFresh   ExpiryStatus = iota
	StatusStale                // within warning window
	StatusExpired              // past TTL
)

// ExpiryResult holds the result of an expiry check.
type ExpiryResult struct {
	Status    ExpiryStatus
	Age       time.Duration
	TTL       time.Duration
	WarnAfter time.Duration
	StampFile string
}

// String returns a human-readable summary of the expiry result.
func (r ExpiryResult) String() string {
	switch r.Status {
	case StatusExpired:
		return fmt.Sprintf("secrets expired (age: %s, ttl: %s)", r.Age.Round(time.Second), r.TTL)
	case StatusStale:
		return fmt.Sprintf("secrets stale — synced %s ago (warn after: %s)", r.Age.Round(time.Second), r.WarnAfter)
	default:
		return fmt.Sprintf("secrets fresh — synced %s ago", r.Age.Round(time.Second))
	}
}

// CheckExpiry reads the stamp at stampFile and evaluates freshness.
// warnAfter: duration after which status becomes Stale.
// ttl: duration after which status becomes Expired.
// If no stamp exists, StatusFresh is returned (no sync has occurred yet).
func CheckExpiry(stampFile string, warnAfter, ttl time.Duration) (ExpiryResult, error) {
	stamp, err := ReadStamp(stampFile)
	if err != nil {
		// Missing stamp — treat as fresh (never synced warning handled by caller)
		return ExpiryResult{Status: StatusFresh, StampFile: stampFile}, nil
	}

	age := time.Since(stamp)
	result := ExpiryResult{
		Age:       age,
		TTL:       ttl,
		WarnAfter: warnAfter,
		StampFile: stampFile,
	}

	switch {
	case age >= ttl:
		result.Status = StatusExpired
	case age >= warnAfter:
		result.Status = StatusStale
	default:
		result.Status = StatusFresh
	}

	return result, nil
}
