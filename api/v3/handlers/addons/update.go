package addons

import (
	"context"
	"fmt"
	"net/http"

	apiv3 "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/api/v3/request"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	UpdateAddonRequest  = addon.UpdateAddonInput
	UpdateAddonResponse = apiv3.Addon
	UpdateAddonParams   = string
	UpdateAddonHandler  httptransport.HandlerWithArgs[UpdateAddonRequest, UpdateAddonResponse, string]
)

func (h *handler) UpdateAddon() UpdateAddonHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, addonID UpdateAddonParams) (UpdateAddonRequest, error) {
			body := apiv3.UpsertAddonRequest{}
			if err := request.ParseBody(r, &body); err != nil {
				return UpdateAddonRequest{}, err
			}

			// NOTE: We gate the addon authoring behind this config flag. It is applied for both create and update and will be removed when unit config is feature complete.
			if !h.unitConfigEnabled {
				for _, rc := range body.RateCards {
					if rc.UnitConfig != nil {
						return UpdateAddonRequest{}, models.NewGenericValidationError(fmt.Errorf("unit_config is not enabled on this deployment of MeterForge"))
					}
				}
			}

			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return UpdateAddonRequest{}, err
			}

			req, err := FromAPIUpsertAddonRequest(ns, addonID, body)
			if err != nil {
				return UpdateAddonRequest{}, err
			}

			req.IgnoreNonCriticalIssues = true

			return req, nil
		},
		func(ctx context.Context, request UpdateAddonRequest) (UpdateAddonResponse, error) {
			a, err := h.service.UpdateAddon(ctx, request)
			if err != nil {
				return UpdateAddonResponse{}, err
			}

			if a == nil {
				return UpdateAddonResponse{}, fmt.Errorf("failed to update add-on")
			}

			return ToAPIAddon(*a)
		},
		commonhttp.JSONResponseEncoderWithStatus[UpdateAddonResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("update-addon"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
