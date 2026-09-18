package customersentitlement

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListCustomerEntitlementAccess() ListCustomerEntitlementAccessHandler
}

type handler struct {
	resolveNamespace   func(ctx context.Context) (string, error)
	customerService    customer.Service
	entitlementService entitlement.Service
	options            []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	customerService customer.Service,
	entitlementService entitlement.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace:   resolveNamespace,
		customerService:    customerService,
		entitlementService: entitlementService,
		options:            options,
	}
}
