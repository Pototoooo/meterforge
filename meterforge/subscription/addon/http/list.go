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
	"github.com/Pototoooo/meterforge/pkg/slicesx"
)

type (
	ListSubscriptionAddonsParams = struct {
		SubscriptionID string
	}
	ListSubscriptionAddonsRequest = struct {
		SubscriptionID models.NamespacedID
	}
	ListSubscriptionAddonsResponse = []api.SubscriptionAddon
	ListSubscriptionAddonsHandler  = httptransport.HandlerWithArgs[ListSubscriptionAddonsRequest, ListSubscriptionAddonsResponse, ListSubscriptionAddonsParams]
)

func (h *handler) ListSubscriptionAddons() ListSubscriptionAddonsHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, params ListSubscriptionAddonsParams) (ListSubscriptionAddonsRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return ListSubscriptionAddonsRequest{}, fmt.Errorf("failed to resolve namespace: %w", err)
			}

			return ListSubscriptionAddonsRequest{
				SubscriptionID: models.NamespacedID{
					Namespace: ns,
					ID:        params.SubscriptionID,
				},
			}, nil
		},
		func(ctx context.Context, req ListSubscriptionAddonsRequest) (ListSubscriptionAddonsResponse, error) {
			res, err := h.SubscriptionAddonService.List(ctx, req.SubscriptionID.Namespace, subscriptionaddon.ListSubscriptionAddonsInput{
				SubscriptionID: req.SubscriptionID.ID,
			})
			if err != nil {
				return nil, err
			}

			view, err := h.SubscriptionService.GetView(ctx, req.SubscriptionID)
			if err != nil {
				return nil, err
			}

			// Deliberately fails the whole list, not just the unrepresentable item: this is
			// an actual-state surface (what the customer is billed for), which rejects,
			// unlike catalog-list endpoints such as ListAddons, which omit instead.
			return slicesx.MapWithErr(res.Items, func(item subscriptionaddon.SubscriptionAddon) (api.SubscriptionAddon, error) {
				if item.Addon.AsProductCatalogAddon().HasUnitConfig() {
					return api.SubscriptionAddon{}, productcatalog.ErrUnitConfigNotRepresentable
				}

				return MapSubscriptionAddonToResponse(view, item)
			})
		},
		commonhttp.JSONResponseEncoderWithStatus[ListSubscriptionAddonsResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.Options,
			httptransport.WithOperationName("listSubscriptionAddons"),
		)...,
	)
}
