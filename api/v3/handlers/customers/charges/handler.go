package charges

import (
	"context"

	billingcharges "github.com/Pototoooo/meterforge/meterforge/billing/charges"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListCustomerCharges() ListCustomerChargesHandler
	CreateCustomerCharge() CreateCustomerChargesHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	service          billingcharges.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	service billingcharges.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		service:          service,
		options:          options,
	}
}
