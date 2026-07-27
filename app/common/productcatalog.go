package common

import (
	"fmt"
	"log/slog"

	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/meterforge/cost"
	costadapter "github.com/Pototoooo/meterforge/meterforge/cost/adapter"
	costservice "github.com/Pototoooo/meterforge/meterforge/cost/service"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
	"github.com/Pototoooo/meterforge/meterforge/llmcost"
	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	productcatalogpgadapter "github.com/Pototoooo/meterforge/meterforge/productcatalog/adapter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	addonadapter "github.com/Pototoooo/meterforge/meterforge/productcatalog/addon/adapter"
	addonservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/addon/service"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/featureresolver"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	planadapter "github.com/Pototoooo/meterforge/meterforge/productcatalog/plan/adapter"
	planservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/plan/service"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon"
	planaddonadapter "github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon/adapter"
	planaddonservice "github.com/Pototoooo/meterforge/meterforge/productcatalog/planaddon/service"
	"github.com/Pototoooo/meterforge/meterforge/streaming"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
)

var ProductCatalog = wire.NewSet(
	Feature,
	Cost,
	Plan,
	Addon,
	PlanAddon,
)

var Feature = wire.NewSet(
	NewFeatureConnector,
	NewFeatureResolver,
)

var Cost = wire.NewSet(
	NewCostService,
)

var Plan = wire.NewSet(
	NewPlanService,
)

var Addon = wire.NewSet(
	NewAddonService,
)

var PlanAddon = wire.NewSet(
	NewPlanAddonService,
)

func NewFeatureConnector(
	logger *slog.Logger,
	db *entdb.Client,
	meterService meter.Service,
	publisher eventbus.Publisher,
) feature.FeatureConnector {
	featureRepo := productcatalogpgadapter.NewPostgresFeatureRepo(db, logger)
	return feature.NewFeatureConnector(featureRepo, meterService, publisher)
}

var NewFeatureResolver = featureresolver.New

func NewCostService(
	featureConnector feature.FeatureConnector,
	meterService meter.Service,
	streamingConnector streaming.Connector,
	llmcostService llmcost.Service,
) (cost.Service, error) {
	adapter := costadapter.New(featureConnector, meterService, streamingConnector, llmcostService)

	return costservice.New(costservice.Config{
		Adapter: adapter,
	})
}

func NewPlanService(
	logger *slog.Logger,
	db *entdb.Client,
	featureResolver productcatalog.FeatureResolver,
	taxCodeService taxcode.Service,
	publisher eventbus.Publisher,
) (plan.Service, error) {
	adapter, err := planadapter.New(planadapter.Config{
		Client: db,
		Logger: logger.With("subsystem", "productcatalog.plan"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize plan adapter: %w", err)
	}

	return planservice.New(planservice.Config{
		Adapter:         adapter,
		FeatureResolver: featureResolver,
		TaxCode:         taxCodeService,
		Logger:          logger.With("subsystem", "productcatalog.plan"),
		Publisher:       publisher,
	})
}

func NewAddonService(
	logger *slog.Logger,
	db *entdb.Client,
	featureResolver productcatalog.FeatureResolver,
	taxCodeService taxcode.Service,
	publisher eventbus.Publisher,
) (addon.Service, error) {
	adapter, err := addonadapter.New(addonadapter.Config{
		Client: db,
		Logger: logger.With("subsystem", "productcatalog.addon"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize add-on adapter: %w", err)
	}

	return addonservice.New(addonservice.Config{
		Adapter:         adapter,
		FeatureResolver: featureResolver,
		TaxCode:         taxCodeService,
		Logger:          logger.With("subsystem", "productcatalog.addon"),
		Publisher:       publisher,
	})
}

func NewPlanAddonService(
	logger *slog.Logger,
	db *entdb.Client,
	planService plan.Service,
	addonService addon.Service,
	publisher eventbus.Publisher,
) (planaddon.Service, error) {
	adapter, err := planaddonadapter.New(planaddonadapter.Config{
		Client: db,
		Logger: logger.With("subsystem", "productcatalog.planaddon"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize add-on adapter: %w", err)
	}

	return planaddonservice.New(planaddonservice.Config{
		Adapter:   adapter,
		Plan:      planService,
		Addon:     addonService,
		Logger:    logger.With("subsystem", "productcatalog.addon"),
		Publisher: publisher,
	})
}
