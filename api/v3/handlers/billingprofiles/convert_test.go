package billingprofiles

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/meterforge/billing"
)

func TestFromAPIBillingWorkflowSubscriptionEndProrationModeDefaults(t *testing.T) {
	createConfig, err := FromAPIBillingWorkflowCreate(api.BillingWorkflow{})
	require.NoError(t, err)
	require.Equal(t, billing.SubscriptionEndProrationModeBillActualPeriod, createConfig.Invoicing.SubscriptionEndProrationMode)

	updateConfig, err := FromAPIBillingWorkflow(api.BillingWorkflow{})
	require.NoError(t, err)
	require.Empty(t, updateConfig.Invoicing.SubscriptionEndProrationMode)

	explicitMode := api.BillingWorkflowInvoicingSubscriptionEndProrationModeBillFullPeriod
	explicitConfig, err := FromAPIBillingWorkflowCreate(api.BillingWorkflow{
		Invoicing: &api.BillingWorkflowInvoicingSettings{
			SubscriptionEndProrationMode: lo.ToPtr(explicitMode),
		},
	})
	require.NoError(t, err)
	require.Equal(t, billing.SubscriptionEndProrationModeBillFullPeriod, explicitConfig.Invoicing.SubscriptionEndProrationMode)
}
