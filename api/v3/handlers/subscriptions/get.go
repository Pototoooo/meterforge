package subscriptions

import (
	"context"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	models "github.com/Pototoooo/meterforge/pkg/models"
)

type (
	GetSubscriptionRequest  = models.NamespacedID
	GetSubscriptionResponse = api.BillingSubscription
	GetSubscriptionParams   = string
	GetSubscriptionHandler  httptransport.HandlerWithArgs[GetSubscriptionRequest, GetSubscriptionResponse, GetSubscriptionParams]
)

// GetSubscription returns a handler for getting a subscription.
func (h *handler) GetSubscription() GetSubscriptionHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, subscriptionID GetSubscriptionParams) (GetSubscriptionRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetSubscriptionRequest{}, err
			}

			return GetSubscriptionRequest{
				Namespace: ns,
				ID:        subscriptionID,
			}, nil
		},
		func(ctx context.Context, request GetSubscriptionRequest) (GetSubscriptionResponse, error) {
			// Get the subscription
			m, err := h.subscriptionService.Get(ctx, request)
			if err != nil {
				return GetSubscriptionResponse{}, err
			}

			return ToAPIBillingSubscription(m), nil
		},
		commonhttp.JSONResponseEncoderWithStatus[GetSubscriptionResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-subscription"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
