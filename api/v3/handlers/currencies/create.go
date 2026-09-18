package currencies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/samber/lo"

	v3 "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/request"
	"github.com/Pototoooo/meterforge/meterforge/currencies"
	"github.com/Pototoooo/meterforge/pkg/currencyx"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type (
	CreateCurrencyRequest  = currencies.CreateCurrencyInput
	CreateCurrencyResponse = v3.BillingCurrency
	CreateCurrencyHandler  = httptransport.Handler[CreateCurrencyRequest, CreateCurrencyResponse]
)

func (h *handler) CreateCurrency() CreateCurrencyHandler {
	return httptransport.NewHandler(
		func(ctx context.Context, r *http.Request) (CreateCurrencyRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return CreateCurrencyRequest{}, fmt.Errorf("failed to resolve namespace: %w", err)
			}

			body := v3.CreateCurrencyCustomRequest{}
			if err = request.ParseBody(r, &body); err != nil {
				return CreateCurrencyRequest{}, fmt.Errorf("failed to parse create custom currency request: %w", err)
			}

			return CreateCurrencyRequest{
				CurrencyDetails: currencyx.CurrencyDetails{
					Code:               currencyx.Code(body.Code),
					Name:               body.Name,
					Symbol:             lo.FromPtr(body.Symbol),
					Precision:          body.Precision,
					DecimalMark:        body.DecimalMark,
					ThousandsSeparator: body.ThousandSeparator,
				},
				Namespace: ns,
			}, nil
		},
		func(ctx context.Context, request CreateCurrencyRequest) (CreateCurrencyResponse, error) {
			resp, err := h.service.CreateCurrency(ctx, request)
			if err != nil {
				return CreateCurrencyResponse{}, err
			}

			return ToAPIBillingCurrency(resp)
		},
		commonhttp.JSONResponseEncoderWithStatus[CreateCurrencyResponse](http.StatusCreated),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("create-custom-currency"),
		)...,
	)
}
