package httpdriver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/app"
	appstripe "github.com/Pototoooo/meterforge/meterforge/app/stripe"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	billingcharges "github.com/Pototoooo/meterforge/meterforge/billing/charges"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/pkg/featuregate"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ProfileHandler
	InvoiceLineHandler
	InvoiceHandler
	CustomerOverrideHandler
}

type ProfileHandler interface {
	CreateProfile() CreateProfileHandler
	GetProfile() GetProfileHandler
	DeleteProfile() DeleteProfileHandler
	UpdateProfile() UpdateProfileHandler
	ListProfiles() ListProfilesHandler
}

type InvoiceLineHandler interface {
	CreatePendingLine() CreatePendingLineHandler
}

type InvoiceHandler interface {
	ListInvoices() ListInvoicesHandler
	InvoicePendingLinesAction() InvoicePendingLinesActionHandler
	DeleteInvoice() DeleteInvoiceHandler
	GetInvoice() GetInvoiceHandler
	UpdateInvoice() UpdateInvoiceHandler
	ProgressInvoice(ProgressAction) ProgressInvoiceHandler
	SimulateInvoice() SimulateInvoiceHandler
}

type CustomerOverrideHandler interface {
	ListCustomerOverrides() ListCustomerOverridesHandler
	UpsertCustomerOverride() UpsertCustomerOverrideHandler
	GetCustomerOverride() GetCustomerOverrideHandler
	DeleteCustomerOverride() DeleteCustomerOverrideHandler
}

type handler struct {
	service          billing.Service
	chargeService    billingcharges.ChargeService
	appService       app.Service
	logger           *slog.Logger
	namespaceDecoder namespacedriver.NamespaceDecoder
	featureSwitches  config.BillingFeatureSwitchesConfiguration
	credits          config.CreditsConfiguration
	featureGate      *featuregate.FeatureGateChecker
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
	featureSwitches config.BillingFeatureSwitchesConfiguration,
	service billing.Service,
	appService app.Service,
	stripeAppService appstripe.Service,
	chargeService billingcharges.ChargeService,
	credits config.CreditsConfiguration,
	featureGate *featuregate.FeatureGateChecker,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		service:          service,
		chargeService:    chargeService,
		appService:       appService,
		logger:           logger,
		namespaceDecoder: namespaceDecoder,
		options:          options,
		featureSwitches:  featureSwitches,
		credits:          credits,
		featureGate:      featureGate,
	}
}
