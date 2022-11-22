package form3_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mkuznets.com/go/form3"
)

func TestDefaultBackOffProvider_IncresingDelayAndReset(t *testing.T) {
	b := form3.DefaultBackOffProvider()
	var prev time.Duration
	for i := 0; i < 3; i++ {
		assert.Less(t, prev, b.NextBackOff(), "Retry delay should increase")
		prev = b.NextBackOff()
	}

	b.Reset()
	assert.Less(t, b.NextBackOff(), prev, "Retry delay should be reset")
}
