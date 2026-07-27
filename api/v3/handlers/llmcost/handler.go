package llmcost

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/llmcost"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListPrices() ListPricesHandler
	GetPrice() GetPriceHandler
	ListOverrides() ListOverridesHandler
	CreateOverride() CreateOverrideHandler
	DeleteOverride() DeleteOverrideHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	service          llmcost.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	service llmcost.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		service:          service,
		options:          options,
	}
}
