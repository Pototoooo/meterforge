package common

import (
	"github.com/google/wire"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/governance"
	governanceservice "github.com/Pototoooo/meterforge/meterforge/governance/service"
	"github.com/Pototoooo/meterforge/meterforge/registry"
)

var Governance = wire.NewSet(
	NewGovernanceService,
)

func NewGovernanceService(
	customer customer.Service,
	entitlementRegistry *registry.Entitlement,
	tracer trace.Tracer,
	meter metric.Meter,
) (governance.Service, error) {
	return governanceservice.New(governanceservice.Config{
		Customer:    customer,
		Entitlement: entitlementRegistry.Entitlement,
		Feature:     entitlementRegistry.Feature,
		Tracer:      tracer,
		Meter:       meter,
	})
}
