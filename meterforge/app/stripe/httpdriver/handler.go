package httpdriver

import (
	"context"
	"errors"
	"net/http"

	appstripe "github.com/Pototoooo/meterforge/meterforge/app/stripe"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	AppStripeHandler
}

type AppStripeHandler interface {
	AppStripeWebhook() AppStripeWebhookHandler
	UpdateStripeAPIKey() UpdateStripeAPIKeyHandler
	CreateAppStripeCheckoutSession() CreateAppStripeCheckoutSessionHandler

	// Customer Stripe Data handlers
	GetCustomerStripeAppData() GetCustomerStripeAppDataHandler
	UpsertCustomerStripeAppData() UpsertCustomerStripeAppDataHandler

	// Customer Stripe Portal handlers
	CreateStripeCustomerPortalSession() CreateStripeCustomerPortalSessionHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	service          appstripe.Service
	billingService   billing.Service
	customerService  customer.Service
	namespaceDecoder namespacedriver.NamespaceDecoder
	options          []httptransport.HandlerOption
}

func (h *handler) resolveNamespace(ctx context.Context) (string, error) {
	ns, ok := h.namespaceDecoder.GetNamespace(ctx)
	if !ok {
		return "", commonhttp.NewHTTPError(http.StatusInternalServerError, errors.New("internal server error"))
	}

	return ns, nil
}

func New(
	namespaceDecoder namespacedriver.NamespaceDecoder,
	service appstripe.Service,
	billingService billing.Service,
	customerService customer.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		service:          service,
		billingService:   billingService,
		customerService:  customerService,
		namespaceDecoder: namespaceDecoder,
		options:          options,
	}
}
