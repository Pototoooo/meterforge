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
	DeleteOverrideRequest  = llmcost.DeleteOverrideInput
	DeleteOverrideResponse = interface{}
	DeleteOverrideHandler  = httptransport.HandlerWithArgs[DeleteOverrideRequest, DeleteOverrideResponse, api.ULID]
)

func (h *handler) DeleteOverride() DeleteOverrideHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, overrideID api.ULID) (DeleteOverrideRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return DeleteOverrideRequest{}, err
			}

			return DeleteOverrideRequest{
				ID:        overrideID,
				Namespace: ns,
			}, nil
		},
		func(ctx context.Context, req DeleteOverrideRequest) (DeleteOverrideResponse, error) {
			err := h.service.DeleteOverride(ctx, req)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
		commonhttp.EmptyResponseEncoder[DeleteOverrideResponse](http.StatusNoContent),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("delete-llm-cost-override"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
		)...,
	)
}
