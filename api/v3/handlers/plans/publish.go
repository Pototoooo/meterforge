package plans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/samber/lo"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	"github.com/Pototoooo/meterforge/pkg/clock"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	PublishPlanRequest  = plan.PublishPlanInput
	PublishPlanResponse = api.BillingPlan
	PublishPlanParams   = string
	PublishPlanHandler  httptransport.HandlerWithArgs[PublishPlanRequest, PublishPlanResponse, PublishPlanParams]
)

func (h *handler) PublishPlan() PublishPlanHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, planID PublishPlanParams) (PublishPlanRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return PublishPlanRequest{}, err
			}

			return PublishPlanRequest{
				NamespacedID: models.NamespacedID{
					Namespace: ns,
					ID:        planID,
				},
				EffectivePeriod: productcatalog.EffectivePeriod{
					EffectiveFrom: lo.ToPtr(clock.Now()),
				},
			}, nil
		},
		func(ctx context.Context, request PublishPlanRequest) (PublishPlanResponse, error) {
			p, err := h.service.PublishPlan(ctx, request)
			if err != nil {
				return PublishPlanResponse{}, err
			}

			if p == nil {
				return PublishPlanResponse{}, fmt.Errorf("failed to publish plan")
			}

			return ToAPIBillingPlan(*p)
		},
		commonhttp.JSONResponseEncoderWithStatus[PublishPlanResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("publish-plan"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
