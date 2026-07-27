package planaddons

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/api/v3/labels"
	"github.com/Pototoooo/meterforge/api/v3/request"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	CreatePlanAddonRequest  = planaddon.CreatePlanAddonInput
	CreatePlanAddonResponse = api.PlanAddon
	CreatePlanAddonParams   = string
	CreatePlanAddonHandler  httptransport.HandlerWithArgs[CreatePlanAddonRequest, CreatePlanAddonResponse, CreatePlanAddonParams]
)

func (h *handler) CreatePlanAddon() CreatePlanAddonHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, planID CreatePlanAddonParams) (CreatePlanAddonRequest, error) {
			body := api.CreatePlanAddonRequest{}
			if err := request.ParseBody(r, &body); err != nil {
				return CreatePlanAddonRequest{}, err
			}

			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return CreatePlanAddonRequest{}, err
			}

			meta, err := labels.ToMetadata(body.Labels)
			if err != nil {
				return CreatePlanAddonRequest{}, err
			}

			return CreatePlanAddonRequest{
				NamespacedModel: models.NamespacedModel{
					Namespace: ns,
				},
				PlanID:        planID,
				AddonID:       body.Addon.Id,
				FromPlanPhase: body.FromPlanPhase,
				MaxQuantity:   body.MaxQuantity,
				Metadata:      meta,
			}, nil
		},
		func(ctx context.Context, request CreatePlanAddonRequest) (CreatePlanAddonResponse, error) {
			a, err := h.addonService.CreatePlanAddon(ctx, request)
			if err != nil {
				return CreatePlanAddonResponse{}, err
			}

			if a == nil {
				return CreatePlanAddonResponse{}, fmt.Errorf("failed to create plan addon")
			}

			return ToAPIPlanAddon(*a)
		},
		commonhttp.JSONResponseEncoderWithStatus[CreatePlanAddonResponse](http.StatusCreated),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("create-plan-addon"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
