package subscriptions

import (
	"context"
	"net/http"

	"github.com/samber/lo"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/api/v3/request"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	models "github.com/Pototoooo/meterforge/pkg/models"
)

type (
	CancelSubscriptionRequest struct {
		ID     models.NamespacedID
		Timing subscription.Timing
	}
	CancelSubscriptionResponse = api.BillingSubscription
	CancelSubscriptionParams   = string
	CancelSubscriptionHandler  httptransport.HandlerWithArgs[CancelSubscriptionRequest, CancelSubscriptionResponse, CancelSubscriptionParams]
)

func (h *handler) CancelSubscription() CancelSubscriptionHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, subscriptionID CancelSubscriptionParams) (CancelSubscriptionRequest, error) {
			// Parse body
			body := api.BillingSubscriptionCancel{}
			if err := request.ParseBody(r, &body); err != nil {
				return CancelSubscriptionRequest{}, err
			}

			// Resolve namespace
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return CancelSubscriptionRequest{}, err
			}

			// Timing (defaults to immediate)
			timing := subscription.Timing{}
			if body.Timing == nil {
				timing.Enum = lo.ToPtr(subscription.TimingImmediate)
			} else {
				timing, err = FromAPIBillingSubscriptionEditTiming(*body.Timing)
				if err != nil {
					return CancelSubscriptionRequest{}, err
				}
			}

			return CancelSubscriptionRequest{
				ID: models.NamespacedID{
					Namespace: ns,
					ID:        subscriptionID,
				},
				Timing: timing,
			}, nil
		},
		func(ctx context.Context, req CancelSubscriptionRequest) (CancelSubscriptionResponse, error) {
			sub, err := h.subscriptionService.Cancel(ctx, req.ID, req.Timing)
			if err != nil {
				return CancelSubscriptionResponse{}, err
			}

			return ToAPIBillingSubscription(sub), nil
		},
		commonhttp.JSONResponseEncoderWithStatus[CancelSubscriptionResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("cancel-subscription"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
