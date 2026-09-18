package apps

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/app"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

// GetAppHandler is a handler to get an app by id
type (
	GetAppRequest  = app.GetAppInput
	GetAppResponse = api.BillingApp
	GetAppHandler  httptransport.HandlerWithArgs[GetAppRequest, GetAppResponse, string]
)

// GetApp returns an app handler
func (h *handler) GetApp() GetAppHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, appId string) (GetAppRequest, error) {
			// Resolve namespace
			namespace, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetAppRequest{}, fmt.Errorf("failed to resolve namespace: %w", err)
			}

			return GetAppRequest{
				Namespace: namespace,
				ID:        appId,
			}, nil
		},
		func(ctx context.Context, request GetAppRequest) (GetAppResponse, error) {
			app, err := h.appService.GetApp(ctx, request)
			if err != nil {
				return GetAppResponse{}, fmt.Errorf("failed to get app: %w", err)
			}

			return ToAPIBillingApp(app)
		},
		commonhttp.JSONResponseEncoder[GetAppResponse],
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-app"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
