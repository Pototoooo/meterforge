package addons

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	DeleteAddonRequest  = addon.DeleteAddonInput
	DeleteAddonResponse = any
	DeleteAddonHandler  httptransport.HandlerWithArgs[DeleteAddonRequest, DeleteAddonResponse, string]
)

func (h *handler) DeleteAddon() DeleteAddonHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, addonID string) (DeleteAddonRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return DeleteAddonRequest{}, err
			}

			return DeleteAddonRequest{
				NamespacedID: models.NamespacedID{
					Namespace: ns,
					ID:        addonID,
				},
			}, nil
		},
		func(ctx context.Context, request DeleteAddonRequest) (DeleteAddonResponse, error) {
			err := h.service.DeleteAddon(ctx, request)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
		commonhttp.EmptyResponseEncoder[DeleteAddonResponse](http.StatusNoContent),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("delete-addon"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
