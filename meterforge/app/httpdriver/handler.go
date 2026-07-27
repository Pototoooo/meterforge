package httpdriver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/app"
	stripeapp "github.com/Pototoooo/meterforge/meterforge/app/stripe"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	AppHandler
}

type AppHandler interface {
	// App handlers
	ListApps() ListAppsHandler
	GetApp() GetAppHandler
	UninstallApp() UninstallAppHandler
	UpdateApp() UpdateAppHandler

	// Customer Data handlers
	ListCustomerData() ListCustomerDataHandler
	UpsertCustomerData() UpsertCustomerDataHandler
	DeleteCustomerData() DeleteCustomerDataHandler

	// Marketplace handlers
	ListMarketplaceListings() ListMarketplaceListingsHandler
	GetMarketplaceListing() GetMarketplaceListingHandler
	MarketplaceAppAPIKeyInstall() MarketplaceAppAPIKeyInstallHandler
	MarketplaceAppInstall() MarketplaceAppInstallHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	service app.Service

	stripeAppService stripeapp.Service
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
	logger *slog.Logger,
	namespaceDecoder namespacedriver.NamespaceDecoder,
	appService app.Service,
	appStripeService stripeapp.Service,
	billingService billing.Service,
	customerService customer.Service,

	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		service:          appService,
		namespaceDecoder: namespaceDecoder,
		stripeAppService: appStripeService,
		billingService:   billingService,
		customerService:  customerService,
		options:          options,
	}
}
