package subscriptiontestutils

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	plansubscription "github.com/Pototoooo/meterforge/meterforge/productcatalog/subscription"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	"github.com/Pototoooo/meterforge/meterforge/testutils"
	"github.com/Pototoooo/meterforge/pkg/clock"
	"github.com/Pototoooo/meterforge/pkg/datetime"
)

func GetExamplePlanInput(t *testing.T) plan.CreatePlanInput {
	b := BuildTestPlanInput(t)

	b.AddPhase(lo.ToPtr(datetime.MustParseDuration(t, "P1M")), ExampleRateCard1.Clone())
	b.AddPhase(lo.ToPtr(datetime.MustParseDuration(t, "P2M")), ExampleRateCard1.Clone(), ExampleRateCard2.Clone(), ExampleRateCard3ForAddons.Clone())
	b.AddPhase(nil, ExampleRateCard1.Clone(), ExampleRateCard3ForAddons.Clone())

	return b.Build()
}

// PlanHelper simply creates and returns a plan
type planHelper struct {
	planService plan.Service
}

func NewPlanHelper(planService plan.Service) *planHelper {
	return &planHelper{
		planService: planService,
	}
}

func (h *planHelper) CreatePlan(t *testing.T, input plan.CreatePlanInput) subscription.Plan {
	t.Helper()
	ctx := context.Background()

	p, err := h.planService.CreatePlan(ctx, input)
	require.Nil(t, err)
	require.NotNil(t, p)

	p, err = h.planService.PublishPlan(ctx, plan.PublishPlanInput{
		NamespacedID: p.NamespacedID,
		EffectivePeriod: productcatalog.EffectivePeriod{
			EffectiveFrom: lo.ToPtr(clock.Now()),
			EffectiveTo:   lo.ToPtr(testutils.GetRFC3339Time(t, "2030-01-01T00:00:00Z")),
		},
	})

	require.Nil(t, err)
	require.NotNil(t, p)

	require.Nilf(t, p.DeletedAt, "plan %s should not be deleted", p.NamespacedID.ID)

	return &plansubscription.Plan{
		Plan: p.AsProductCatalogPlan(),
		Ref:  &p.NamespacedID,
	}
}
