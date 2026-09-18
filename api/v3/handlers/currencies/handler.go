package currencies

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/currencies"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListCurrencies() ListCurrenciesHandler
	CreateCurrency() CreateCurrencyHandler
	CreateCostBasis() CreateCostBasisHandler
	ListCostBases() ListCostBasesHandler
	GetCurrency() GetCurrencyHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	options          []httptransport.HandlerOption
	service          currencies.Service
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	currencyService currencies.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		options:          options,
		service:          currencyService,
	}
}
