package features

import (
	"context"

	"github.com/Pototoooo/meterforge/meterforge/llmcost"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

type Handler interface {
	ListFeatures() ListFeaturesHandler
	GetFeature() GetFeatureHandler
	CreateFeature() CreateFeatureHandler
	UpdateFeature() UpdateFeatureHandler
	DeleteFeature() DeleteFeatureHandler
}

type handler struct {
	resolveNamespace func(ctx context.Context) (string, error)
	connector        feature.FeatureConnector
	meterService     meter.Service
	llmcostService   llmcost.Service
	options          []httptransport.HandlerOption
}

func New(
	resolveNamespace func(ctx context.Context) (string, error),
	connector feature.FeatureConnector,
	meterService meter.Service,
	llmcostService llmcost.Service,
	options ...httptransport.HandlerOption,
) Handler {
	return &handler{
		resolveNamespace: resolveNamespace,
		connector:        connector,
		meterService:     meterService,
		llmcostService:   llmcostService,
		options:          options,
	}
}
