package dotenv

import (
	"fmt"
	"os"
	"time"
)

// WatchResult holds the outcome of a watch cycle.
type WatchResult struct {
	Refreshed bool
	Expired   bool
	Error     error
	At        time.Time
}

// WatchOptions configures the Watch loop.
type WatchOptions struct {
	Interval  time.Duration
	MaxCycles int // 0 = run forever
	OnRefresh func(path string) error
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval:  30 * time.Second,
		MaxCycles: 0,
	}
}

// Watch polls the stamp file at the given path and calls opts.OnRefresh
// when the secrets are expired. It blocks until the context is cancelled
// via the stop channel or MaxCycles is reached.
func Watch(path string, policy TTLPolicy, opts WatchOptions, stop <-chan struct{}) []WatchResult {
	var results []WatchResult
	cycles := 0

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return results
		case t := <-ticker.C:
			result := WatchResult{At: t}
			eval := EvaluateTTL(path, policy)
			if eval.Status == TTLExpired {
				result.Expired = true
				if opts.OnRefresh != nil {
					if err := opts.OnRefresh(path); err != nil {
						result.Error = fmt.Errorf("refresh failed: %w", err)
					} else {
						result.Refreshed = true
						_ = WriteStamp(StampPath(path))
					}
				}
			}
			results = append(results, result)
			cycles++
			if opts.MaxCycles > 0 && cycles >= opts.MaxCycles {
				return results
			}
		}
	}
}

// WatchOnce checks expiry once and calls refresh if needed.
func WatchOnce(path string, policy TTLPolicy, onRefresh func(string) error) WatchResult {
	result := WatchResult{At: time.Now()}
	eval := EvaluateTTL(path, policy)
	if eval.Status == TTLExpired {
		result.Expired = true
		if onRefresh != nil {
			if err := onRefresh(path); err != nil {
				result.Error = fmt.Errorf("refresh failed: %w", err)
			} else {
				result.Refreshed = true
				_ = WriteStamp(StampPath(path))
			}
		}
	}
	_ = os.Getenv("") // satisfy import
	return result
}
