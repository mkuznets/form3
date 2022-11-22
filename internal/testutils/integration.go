package testutils

import (
	"os"
	"testing"
)

const BaseUrlEnvName = "FORM3_API_BASE_URL"

func EnsureIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("-short is enabled, skipping integration test")
	}

	if _, ok := os.LookupEnv(BaseUrlEnvName); !ok {
		t.Skipf("%s environment variable is required", BaseUrlEnvName)
	}
}

func BaseUrl() string {
	return os.Getenv(BaseUrlEnvName)
}
