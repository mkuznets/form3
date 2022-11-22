package form3_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mkuznets.com/go/form3"
	"mkuznets.com/go/form3/internal/testutils"
)

func TestIntegration_Client_Api(t *testing.T) {
	testutils.EnsureIntegration(t)

	t.Run("invalid endpoint", func(t *testing.T) {
		api := newClient("").Api()
		err := api.Do(context.Background(), &form3.Call{
			Method: "GET",
			Path:   "/v1/random",
		})
		require.ErrorAs(t, err, &form3.Error{})
		assert.ErrorContains(t, err, "404")
	})

}
