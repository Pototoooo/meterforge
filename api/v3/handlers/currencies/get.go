package currencies

import (
	"context"
	"fmt"
	"net/http"

	v3 "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/meterforge/currencies"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	GetCurrencyRequest  = currencies.GetCurrencyInput
	GetCurrencyResponse = v3.BillingCurrency
	GetCurrencyHandler  = httptransport.HandlerWithArgs[GetCurrencyRequest, GetCurrencyResponse, GetCurrencyParams]
	GetCurrencyParams   = string
)

func (h *handler) GetCurrency() GetCurrencyHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, params GetCurrencyParams) (GetCurrencyRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetCurrencyRequest{}, fmt.Errorf("failed to resolve namespace: %w", err)
			}

			return GetCurrencyRequest{
				NamespacedID: models.NamespacedID{
					Namespace: ns,
					ID:        params,
				},
				CurrencyExpandOptions: currencies.CurrencyExpandOptions{
					CostBasis: true,
				},
			}, nil
		},
		func(ctx context.Context, request GetCurrencyRequest) (GetCurrencyResponse, error) {
			resp, err := h.service.GetCurrency(ctx, request)
			if err != nil {
				return GetCurrencyResponse{}, err
			}

			return ToAPIBillingCurrency(resp)
		},
		commonhttp.JSONResponseEncoderWithStatus[GetCurrencyResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-custom-currency"),
		)...,
	)
}
