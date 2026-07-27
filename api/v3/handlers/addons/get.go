package addons

import (
	"context"
	"fmt"
	"net/http"

	apiv3 "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	GetAddonRequest       = addon.GetAddonInput
	GetAddonRequestParams = string
	GetAddonResponse      = apiv3.Addon
	GetAddonHandler       httptransport.HandlerWithArgs[GetAddonRequest, GetAddonResponse, GetAddonRequestParams]
)

func (h *handler) GetAddon() GetAddonHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, addonID GetAddonRequestParams) (GetAddonRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetAddonRequest{}, err
			}

			return GetAddonRequest{
				NamespacedID: models.NamespacedID{
					Namespace: ns,
					ID:        addonID,
				},
			}, nil
		},
		func(ctx context.Context, request GetAddonRequest) (GetAddonResponse, error) {
			a, err := h.service.GetAddon(ctx, request)
			if err != nil {
				return GetAddonResponse{}, err
			}

			if a == nil {
				return GetAddonResponse{}, fmt.Errorf("failed to get add-on: %s", request.ID)
			}

			return ToAPIAddon(*a)
		},
		commonhttp.JSONResponseEncoderWithStatus[GetAddonResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-addon"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
