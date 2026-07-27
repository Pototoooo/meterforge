package customerscredits

import (
	"context"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/api/v3/request"
	"github.com/Pototoooo/meterforge/meterforge/billing/creditgrant"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type (
	UpdateCreditGrantExternalSettlementRequest  = creditgrant.UpdateExternalSettlementInput
	UpdateCreditGrantExternalSettlementResponse = api.BillingCreditGrant
	UpdateCreditGrantExternalSettlementParams   struct {
		CustomerID    api.ULID
		CreditGrantID api.ULID
	}
	UpdateCreditGrantExternalSettlementHandler = httptransport.HandlerWithArgs[UpdateCreditGrantExternalSettlementRequest, UpdateCreditGrantExternalSettlementResponse, UpdateCreditGrantExternalSettlementParams]
)

func (h *handler) UpdateCreditGrantExternalSettlement() UpdateCreditGrantExternalSettlementHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, args UpdateCreditGrantExternalSettlementParams) (UpdateCreditGrantExternalSettlementRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return UpdateCreditGrantExternalSettlementRequest{}, err
			}

			var body api.UpdateCreditGrantExternalSettlementRequest
			if err := request.ParseBody(r, &body); err != nil {
				return UpdateCreditGrantExternalSettlementRequest{}, err
			}

			req, err := fromAPIUpdateCreditGrantExternalSettlementRequest(ns, args.CustomerID, args.CreditGrantID, body)
			if err != nil {
				return UpdateCreditGrantExternalSettlementRequest{}, err
			}

			return req, nil
		},
		func(ctx context.Context, request UpdateCreditGrantExternalSettlementRequest) (UpdateCreditGrantExternalSettlementResponse, error) {
			charge, err := h.creditGrantService.UpdateExternalSettlement(ctx, request)
			if err != nil {
				return UpdateCreditGrantExternalSettlementResponse{}, err
			}

			return toAPIBillingCreditGrant(charge)
		},
		commonhttp.JSONResponseEncoderWithStatus[UpdateCreditGrantExternalSettlementResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("update-credit-grant-external-settlement"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
