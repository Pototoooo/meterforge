package httpdriver

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/api"
	plansubscription "github.com/Pototoooo/meterforge/meterforge/productcatalog/subscription"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	MigrateSubscriptionRequest  = plansubscription.MigrateSubscriptionRequest
	MigrateSubscriptionResponse = api.SubscriptionChangeResponseBody
	MigrateSubscriptionParams   = struct {
		ID string
	}
	MigrateSubscriptionHandler = httptransport.HandlerWithArgs[MigrateSubscriptionRequest, MigrateSubscriptionResponse, MigrateSubscriptionParams]
)

func (h *handler) MigrateSubscription() MigrateSubscriptionHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, params MigrateSubscriptionParams) (MigrateSubscriptionRequest, error) {
			var body api.MigrateSubscriptionJSONRequestBody

			if err := commonhttp.JSONRequestBodyDecoder(r, &body); err != nil {
				return MigrateSubscriptionRequest{}, err
			}

			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return MigrateSubscriptionRequest{}, err
			}

			var timing *subscription.Timing

			if body.Timing != nil {
				t, err := MapAPITimingToTiming(*body.Timing)
				if err != nil {
					return MigrateSubscriptionRequest{}, err
				}

				timing = &t
			}

			return MigrateSubscriptionRequest{
				ID: models.NamespacedID{
					Namespace: ns,
					ID:        params.ID,
				},
				TargetVersion:    body.TargetVersion,
				StartingPhase:    body.StartingPhase,
				Timing:           timing,
				BillingAnchor:    body.BillingAnchor,
				RejectUnitConfig: true,
			}, nil
		},
		func(ctx context.Context, request MigrateSubscriptionRequest) (MigrateSubscriptionResponse, error) {
			res, err := h.PlanSubscriptionService.Migrate(ctx, request)
			if err != nil {
				return MigrateSubscriptionResponse{}, err
			}

			v, err := MapSubscriptionViewToAPI(res.Next)

			return MigrateSubscriptionResponse{
				Current: MapSubscriptionToAPI(res.Current),
				Next:    v,
			}, err
		},
		commonhttp.JSONResponseEncoderWithStatus[MigrateSubscriptionResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.Options,
			httptransport.WithOperationName("MigrateSubscription"),
			httptransport.WithErrorEncoder(errorEncoder()),
		)...,
	)
}
