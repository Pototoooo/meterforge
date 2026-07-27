package featurecost

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/cost"
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	QueryFeatureCost() QueryFeatureCostHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	costService      cost.Service
	featureConnector feature.FeatureConnector
	meterService     meter.Service
	customerService  customer.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	costService cost.Service,
	featureConnector feature.FeatureConnector,
	meterService meter.Service,
	customerService customer.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		costService:      costService,
		featureConnector: featureConnector,
		meterService:     meterService,
		customerService:  customerService,
		options:          options,
	}
}
