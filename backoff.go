package form3

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Parameters of DefaultBackOffProvider.
const (
	BackOffInitialInterval     = 500 * time.Millisecond
	BackOffRandomizationFactor = 0.5
	BackOffMultiplier          = 1.5
	BackOffMaxInterval         = 60 * time.Second
	BackOffMaxElapsedTime      = 3 * time.Minute
)

// BackOff is an interface for the retry policy. It is compatible with types provided by github.com/cenkalti/backoff/v4.
type BackOff interface {
	// NextBackOff returns the duration to wait before retrying the operation.
	// Negative duration indicates that no more retries should be made.
	NextBackOff() time.Duration

	// Reset to initial state.
	Reset()
}

// DefaultBackOffProvider returns the retry policy used in form3.Client by default. Implements exponential backoff with randomized intervals.
func DefaultBackOffProvider() BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = BackOffInitialInterval
	b.RandomizationFactor = BackOffRandomizationFactor
	b.Multiplier = BackOffMultiplier
	b.MaxInterval = BackOffMaxInterval
	b.MaxElapsedTime = BackOffMaxElapsedTime
	return b
}
