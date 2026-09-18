package common

import (
	"fmt"
	"log/slog"

	"github.com/google/wire"
	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/app/config"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	entitlementvalidator "github.com/Pototoooo/meterforge/meterforge/entitlement/validators/customer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/registry"
	registrybuilder "github.com/Pototoooo/meterforge/meterforge/registry/builder"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
	"github.com/Pototoooo/meterforge/pkg/framework/lockr"
)

var Entitlement = wire.NewSet(
	NewEntitlementRegistry,
)

func NewEntitlementRegistry(
	logger *slog.Logger,
	db *entdb.Client,
	tracer trace.Tracer,
	entitlementConfig config.EntitlementsConfiguration,
	streamingConnector streaming.Connector,
	meterService meter.Service,
	eventPublisher eventbus.Publisher,
	locker *lockr.Locker,
	customerService customer.Service,
) (*registry.Entitlement, error) {
	entRegistry := registrybuilder.GetEntitlementRegistry(registrybuilder.EntitlementOptions{
		DatabaseClient:            db,
		StreamingConnector:        streamingConnector,
		MeterService:              meterService,
		CustomerService:           customerService,
		Logger:                    logger,
		Publisher:                 eventPublisher,
		EntitlementsConfiguration: entitlementConfig,
		Tracer:                    tracer,
		Locker:                    locker,
	})

	// Create and register the entitlement validator
	validator, err := entitlementvalidator.NewValidator(entRegistry.EntitlementRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to create entitlement validator: %w", err)
	}

	customerService.RegisterRequestValidator(validator)

	return entRegistry, nil
}
