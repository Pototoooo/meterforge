package common

import (
	"fmt"
	"log/slog"

	"github.com/google/wire"
	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	customeradapter "github.com/Pototoooo/meterforge/meterforge/customer/adapter"
	customerservice "github.com/Pototoooo/meterforge/meterforge/customer/service"
	customerservicehooks "github.com/Pototoooo/meterforge/meterforge/customer/service/hooks"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	ledgerresolvers "github.com/Pototoooo/meterforge/meterforge/ledger/resolvers"
	"github.com/Pototoooo/meterforge/meterforge/registry"
	"github.com/Pototoooo/meterforge/meterforge/subject"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
)

var Customer = wire.NewSet(
	NewCustomerService,
)

func NewCustomerService(
	logger *slog.Logger,
	db *entdb.Client,
	eventPublisher eventbus.Publisher,
) (customer.Service, error) {
	customerAdapter, err := customeradapter.New(customeradapter.Config{
		Client: db,
		Logger: logger.WithGroup("customer.postgres"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer adapter: %w", err)
	}

	service, err := customerservice.New(customerservice.Config{
		Adapter:   customerAdapter,
		Publisher: eventPublisher,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer service: %w", err)
	}

	return service, nil
}

type CustomerSubjectHook customerservicehooks.SubjectCustomerHook

type CustomerLedgerHook ledgerresolvers.CustomerLedgerHook

func NewCustomerLedgerServiceHook(
	creditsConfig config.CreditsConfiguration,
	tracer trace.Tracer,
	accountResolver customerLedgerProvisioner,
	customerService customer.Service,
) (CustomerLedgerHook, error) {
	if !creditsConfig.Enabled {
		return ledgerresolvers.NoopCustomerLedgerHook{}, nil
	}

	h, err := ledgerresolvers.NewCustomerLedgerHook(ledgerresolvers.CustomerLedgerHookConfig{
		Service: accountResolver,
		Tracer:  tracer,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer ledger hook: %w", err)
	}

	customerService.RegisterHooks(h)

	return h, nil
}

func NewCustomerSubjectServiceHook(
	config config.CustomerConfiguration,
	logger *slog.Logger,
	tracer trace.Tracer,
	subjectService subject.Service,
	customerService customer.Service,
	customerOverrideService billing.CustomerOverrideService,
) (CustomerSubjectHook, error) {
	if !config.EnableSubjectHook {
		return new(customerservicehooks.NoopSubjectCustomerHook), nil
	}

	// Initialize the subject customer hook and register it for Subject service
	h, err := customerservicehooks.NewSubjectCustomerHook(customerservicehooks.SubjectCustomerHookConfig{
		Customer:         customerService,
		CustomerOverride: customerOverrideService,
		Logger:           logger,
		Tracer:           tracer,
		IgnoreErrors:     config.IgnoreErrors,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer subject hook: %w", err)
	}

	subjectService.RegisterHooks(h)

	return h, nil
}

type CustomerEntitlementValidatorHook customerservicehooks.EntitlementValidatorHook

func NewCustomerEntitlementValidatorServiceHook(
	logger *slog.Logger,
	entitlementRegistry *registry.Entitlement,
	customerService customer.Service,
) (CustomerEntitlementValidatorHook, error) {
	h, err := customerservicehooks.NewEntitlementValidatorHook(customerservicehooks.EntitlementValidatorHookConfig{
		EntitlementService: entitlementRegistry.Entitlement,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create customer entitlement validator hook: %w", err)
	}

	customerService.RegisterHooks(h)

	return h, nil
}
