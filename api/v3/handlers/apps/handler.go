package apps

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/app"
	appstripe "github.com/Pototoooo/meterforge/meterforge/app/stripe"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	// App handlers
	ListApps() ListAppsHandler
	GetApp() GetAppHandler
	UninstallApp() UninstallAppHandler

	CatalogHandler
}

type CatalogHandler interface {
	ListAppCatalog() ListAppCatalogHandler
	GetAppCatalog() GetAppCatalogHandler
	InstallApp() InstallAppHandler
}

var _ Handler = (*handler)(nil)

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	appService       app.Service
	billingService   billing.Service
	stripeAppService appstripe.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	appService app.Service,
	billingService billing.Service,
	stripeAppService appstripe.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		appService:       appService,
		billingService:   billingService,
		stripeAppService: stripeAppService,
		options:          options,
	}
}
