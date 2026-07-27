package e2e

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	api "github.com/Pototoooo/meterforge/api/client/go"
)

func initClient(t testing.TB) *api.ClientWithResponses {
	t.Helper()

	address := os.Getenv("METERFORGE_ADDRESS")
	if address == "" {
		t.Skip("METERFORGE_ADDRESS not set")
	}

	client, err := api.NewClientWithResponses(address)
	require.NoError(t, err)

	return client
}

func shouldRunSlowTests(t *testing.T) bool {
	t.Helper()

	return os.Getenv("RUN_SLOW_TESTS") != ""
}
