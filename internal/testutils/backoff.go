package testutils

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

type ZeroMaxCountBackOff struct {
	i int
	c int
}

func NewMaxRetriesBackOff(maxRetries int) *ZeroMaxCountBackOff {
	return &ZeroMaxCountBackOff{c: maxRetries}
}

func (b *ZeroMaxCountBackOff) Reset() {
	b.i = 0
}

func (b *ZeroMaxCountBackOff) NextBackOff() time.Duration {
	if b.i < b.c {
		b.i++
		return 0
	}
	return backoff.Stop
}
