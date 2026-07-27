package taxcodes

import (
	"context"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	models "github.com/Pototoooo/meterforge/pkg/models"
)

type (
	GetTaxCodeRequest  = taxcode.GetTaxCodeInput
	GetTaxCodeResponse = api.BillingTaxCode
	GetTaxCodeParams   = string
	GetTaxCodeHandler  httptransport.HandlerWithArgs[GetTaxCodeRequest, GetTaxCodeResponse, GetTaxCodeParams]
)

func (h *handler) GetTaxCode() GetTaxCodeHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, taxCodeID GetTaxCodeParams) (GetTaxCodeRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetTaxCodeRequest{}, err
			}

			return GetTaxCodeRequest{
				NamespacedID: models.NamespacedID{
					Namespace: ns,
					ID:        taxCodeID,
				},
			}, nil
		},
		func(ctx context.Context, request GetTaxCodeRequest) (GetTaxCodeResponse, error) {
			// Call the service to get the tax code
			taxCode, err := h.service.GetTaxCode(ctx, request)
			if err != nil {
				return GetTaxCodeResponse{}, err
			}
			// Convert to API response type
			return ToAPIBillingTaxCode(taxCode)
		},
		commonhttp.JSONResponseEncoderWithStatus[GetTaxCodeResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-tax-code"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
