package httpdriver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Pototoooo/meterforge/api"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	subscriptionaddon "github.com/Pototoooo/meterforge/meterforge/subscription/addon"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	GetSubscriptionAddonParams = struct {
		SubscriptionID      string
		SubscriptionAddonID string
	}
	GetSubscriptionAddonRequest = struct {
		SubscriptionID      models.NamespacedID
		SubscriptionAddonID models.NamespacedID
	}
	GetSubscriptionAddonResponse = api.SubscriptionAddon
	GetSubscriptionAddonHandler  = httptransport.HandlerWithArgs[GetSubscriptionAddonRequest, GetSubscriptionAddonResponse, GetSubscriptionAddonParams]
)

func (h *handler) GetSubscriptionAddon() GetSubscriptionAddonHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, params GetSubscriptionAddonParams) (GetSubscriptionAddonRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetSubscriptionAddonRequest{}, fmt.Errorf("failed to resolve namespace: %w", err)
			}

			return GetSubscriptionAddonRequest{
				SubscriptionID:      models.NamespacedID{Namespace: ns, ID: params.SubscriptionID},
				SubscriptionAddonID: models.NamespacedID{Namespace: ns, ID: params.SubscriptionAddonID},
			}, nil
		},
		func(ctx context.Context, req GetSubscriptionAddonRequest) (GetSubscriptionAddonResponse, error) {
			res, err := h.SubscriptionAddonService.Get(ctx, subscriptionaddon.GetSubscriptionAddonInput{
				NamespacedID: req.SubscriptionAddonID,
			})
			if err != nil {
				return GetSubscriptionAddonResponse{}, err
			}

			view, err := h.SubscriptionService.GetView(ctx, req.SubscriptionID)
			if err != nil {
				return GetSubscriptionAddonResponse{}, err
			}

			if res.Addon.AsProductCatalogAddon().HasUnitConfig() {
				return GetSubscriptionAddonResponse{}, productcatalog.ErrUnitConfigNotRepresentable
			}

			return MapSubscriptionAddonToResponse(view, *res)
		},
		commonhttp.JSONResponseEncoderWithStatus[GetSubscriptionAddonResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.Options,
			httptransport.WithOperationName("getSubscriptionAddon"),
		)...,
	)
}
