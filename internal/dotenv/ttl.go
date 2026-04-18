package dotenv

import (
	"fmt"
	"time"
)

// TTLPolicy defines how long secrets are considered valid.
type TTLPolicy struct {
	MaxAge time.Duration
	WarnAge time.Duration
}

// DefaultTTLPolicy returns a sensible default policy.
func DefaultTTLPolicy() TTLPolicy {
	return TTLPolicy{
		MaxAge:  24 * time.Hour,
		WarnAge: 20 * time.Hour,
	}
}

// TTLStatus represents the freshness state of a secret file.
type TTLStatus int

const (
	TTLFresh TTLStatus = iota
	TTLWarning
	TTLExpired
	TTLUnknown
)

func (s TTLStatus) String() string {
	switch s {
	case TTLFresh:
		return "fresh"
	case TTLWarning:
		return "warning"
	case TTLExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// EvaluateTTL returns the TTLStatus for a given stamp time under a policy.
func EvaluateTTL(stamped time.Time, policy TTLPolicy) TTLStatus {
	if stamped.IsZero() {
		return TTLUnknown
	}
	age := time.Since(stamped)
	if age >= policy.MaxAge {
		return TTLExpired
	}
	if age >= policy.WarnAge {
		return TTLWarning
	}
	return TTLFresh
}

// TTLSummary returns a human-readable summary for a given path.
func TTLSummary(path string, policy TTLPolicy) (string, error) {
	stamp, err := ReadStamp(path)
	if err != nil {
		return "", fmt.Errorf("ttl: read stamp: %w", err)
	}
	status := EvaluateTTL(stamp.SyncedAt, policy)
	age := time.Since(stamp.SyncedAt).Round(time.Second)
	return fmt.Sprintf("%s: age=%s status=%s", path, age, status), nil
}
