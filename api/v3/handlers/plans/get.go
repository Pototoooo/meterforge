package plans

import (
	"context"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	GetPlanRequest  = plan.GetPlanInput
	GetPlanResponse = api.BillingPlan
	GetPlanParams   = string
	GetPlanHandler  httptransport.HandlerWithArgs[GetPlanRequest, GetPlanResponse, GetPlanParams]
)

func (h *handler) GetPlan() GetPlanHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, planID GetPlanParams) (GetPlanRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetPlanRequest{}, err
			}

			return GetPlanRequest{
				NamespacedID: models.NamespacedID{
					Namespace: ns,
					ID:        planID,
				},
			}, nil
		},
		func(ctx context.Context, request GetPlanRequest) (GetPlanResponse, error) {
			p, err := h.service.GetPlan(ctx, request)
			if err != nil {
				return GetPlanResponse{}, err
			}

			return ToAPIBillingPlan(*p)
		},
		commonhttp.JSONResponseEncoderWithStatus[GetPlanResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-plan"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
