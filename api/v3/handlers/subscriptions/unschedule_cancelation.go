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
	UnscheduleCancelationRequest  = models.NamespacedID
	UnscheduleCancelationResponse = api.BillingSubscription
	UnscheduleCancelationParams   = string
	UnscheduleCancelationHandler  httptransport.HandlerWithArgs[UnscheduleCancelationRequest, UnscheduleCancelationResponse, UnscheduleCancelationParams]
)

func (h *handler) UnscheduleCancelation() UnscheduleCancelationHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, subscriptionID UnscheduleCancelationParams) (UnscheduleCancelationRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return UnscheduleCancelationRequest{}, err
			}

			return UnscheduleCancelationRequest{
				Namespace: ns,
				ID:        subscriptionID,
			}, nil
		},
		func(ctx context.Context, req UnscheduleCancelationRequest) (UnscheduleCancelationResponse, error) {
			sub, err := h.subscriptionService.Continue(ctx, req)
			if err != nil {
				return UnscheduleCancelationResponse{}, err
			}

			return ToAPIBillingSubscription(sub), nil
		},
		commonhttp.JSONResponseEncoderWithStatus[UnscheduleCancelationResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("unschedule-cancelation"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
