package api

import (
	"github.com/cenkalti/backoff/v4"
	"time"
)

const (
	BackOffInitialInterval     = 500 * time.Millisecond
	BackOffRandomizationFactor = 0.5
	BackOffMultiplier          = 1.5
	BackOffMaxInterval         = 60 * time.Second
	BackOffMaxElapsedTime      = 3 * time.Minute
)

type BackOff interface {
	NextBackOff() time.Duration
	Reset()
}

func DefaultBackOffProvider() BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = BackOffInitialInterval
	b.RandomizationFactor = BackOffRandomizationFactor
	b.Multiplier = BackOffMultiplier
	b.MaxInterval = BackOffMaxInterval
	b.MaxElapsedTime = BackOffMaxElapsedTime
	return b
}
