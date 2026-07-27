package taxcodes

import (
	"context"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type (
	GetOrganizationDefaultTaxCodesRequest  = taxcode.GetOrganizationDefaultTaxCodesInput
	GetOrganizationDefaultTaxCodesResponse = api.OrganizationDefaultTaxCodes
	GetOrganizationDefaultTaxCodesHandler  = httptransport.Handler[GetOrganizationDefaultTaxCodesRequest, GetOrganizationDefaultTaxCodesResponse]
)

func (h *handler) GetOrganizationDefaultTaxCodes() GetOrganizationDefaultTaxCodesHandler {
	return httptransport.NewHandler(
		func(ctx context.Context, r *http.Request) (GetOrganizationDefaultTaxCodesRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetOrganizationDefaultTaxCodesRequest{}, err
			}
			return GetOrganizationDefaultTaxCodesRequest{Namespace: ns}, nil
		},
		func(ctx context.Context, request GetOrganizationDefaultTaxCodesRequest) (GetOrganizationDefaultTaxCodesResponse, error) {
			cfg, err := h.service.GetOrganizationDefaultTaxCodes(ctx, request)
			if err != nil {
				return GetOrganizationDefaultTaxCodesResponse{}, err
			}
			return ToAPIOrganizationDefaultTaxCodes(cfg)
		},
		commonhttp.JSONResponseEncoderWithStatus[GetOrganizationDefaultTaxCodesResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-organization-default-tax-codes"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
