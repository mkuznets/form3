package testutils

import "testing"

func SetupIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}
