package testutils

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

type TestBackOff struct {
	i int
	c int
}

func NewTestBackOff(maxRetries int) *TestBackOff {
	return &TestBackOff{c: maxRetries}
}

func (b *TestBackOff) Reset() {
	b.i = 0
}

func (b *TestBackOff) NextBackOff() time.Duration {
	if b.i < b.c {
		b.i++
		return 0
	}
	return backoff.Stop
}
