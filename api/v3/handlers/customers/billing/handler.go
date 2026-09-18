package customersbilling

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/app"
	appstripe "github.com/Pototoooo/meterforge/meterforge/app/stripe"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	GetCustomerBilling() GetCustomerBillingHandler
	UpdateCustomerBilling() UpdateCustomerBillingHandler
	UpdateCustomerBillingAppData() UpdateCustomerBillingAppDataHandler
	CreateCustomerStripeCheckoutSession() CreateCustomerStripeCheckoutSessionHandler
	CreateCustomerStripePortalSession() CreateCustomerStripePortalSessionHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	billingService   billing.Service
	customerService  customer.Service
	appService       app.Service
	stripeService    appstripe.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	billingService billing.Service,
	customerService customer.Service,
	stripeService appstripe.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		billingService:   billingService,
		customerService:  customerService,
		stripeService:    stripeService,
		options:          options,
	}
}
