package llmcost

import (
	"context"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/llmcost"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type (
	GetPriceRequest  = llmcost.GetPriceInput
	GetPriceResponse = api.LLMCostPrice
	GetPriceHandler  = httptransport.HandlerWithArgs[GetPriceRequest, GetPriceResponse, api.ULID]
)

func (h *handler) GetPrice() GetPriceHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, priceID api.ULID) (GetPriceRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetPriceRequest{}, err
			}

			return GetPriceRequest{
				ID:        priceID,
				Namespace: ns,
			}, nil
		},
		func(ctx context.Context, request GetPriceRequest) (GetPriceResponse, error) {
			price, err := h.service.GetPrice(ctx, request)
			if err != nil {
				return GetPriceResponse{}, err
			}

			return domainPriceToAPI(price), nil
		},
		commonhttp.JSONResponseEncoderWithStatus[GetPriceResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-llm-cost-price"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
