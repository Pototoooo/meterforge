package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ledgeraccountdb "github.com/Pototoooo/meterforge/meterforge/ent/db/ledgeraccount"
	"github.com/Pototoooo/meterforge/meterforge/ledger"
	ledgertestutils "github.com/Pototoooo/meterforge/meterforge/ledger/testutils"
	"github.com/Pototoooo/meterforge/meterforge/testutils"
)

func TestLedgerNamespaceHandlerIsSafeForRepeatedCreateNamespaceCalls(t *testing.T) {
	testDB := testutils.InitPostgresDB(t, testutils.PostgresDBStateAtlasMigrated)
	t.Cleanup(func() {
		require.NoError(t, testDB.EntDriver.Close())
		require.NoError(t, testDB.PGDriver.Close())
	})

	deps, err := ledgertestutils.InitDeps(testDB.EntDriver.Client(), testutils.NewDiscardLogger(t))
	require.NoError(t, err)

	handler := NewLedgerNamespaceHandler(deps.ResolversService)
	namespace := fmt.Sprintf("startup-test-%d", time.Now().UnixNano())

	require.NoError(t, handler.CreateNamespace(t.Context(), namespace))
	require.NoError(t, handler.CreateNamespace(t.Context(), namespace))

	count, err := testDB.EntDriver.Client().LedgerAccount.Query().
		Where(
			ledgeraccountdb.Namespace(namespace),
			ledgeraccountdb.AccountTypeIn(
				ledger.AccountTypeWash,
				ledger.AccountTypeEarnings,
				ledger.AccountTypeBrokerage,
			),
		).
		Count(t.Context())
	require.NoError(t, err)
	require.Equal(t, 3, count)
}
