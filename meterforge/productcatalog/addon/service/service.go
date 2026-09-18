package service

import (
	"errors"
	"log/slog"

	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/meterforge/taxcode"
	"github.com/Pototoooo/meterforge/meterforge/watermill/eventbus"
)

type Config struct {
	Adapter   addon.Repository
	TaxCode   taxcode.Service
	Logger    *slog.Logger
	Publisher eventbus.Publisher

	FeatureResolver productcatalog.FeatureResolver
}

func New(config Config) (addon.Service, error) {
	if config.Adapter == nil {
		return nil, errors.New("add-on adapter is required")
	}

	if config.FeatureResolver == nil {
		return nil, errors.New("feature resolver is required")
	}

	if config.TaxCode == nil {
		return nil, errors.New("tax code service is required")
	}

	if config.Logger == nil {
		return nil, errors.New("logger is required")
	}

	if config.Publisher == nil {
		return nil, errors.New("publisher is required")
	}

	return &service{
		adapter:   config.Adapter,
		taxCode:   config.TaxCode,
		logger:    config.Logger,
		publisher: config.Publisher,

		featureResolver: config.FeatureResolver,
	}, nil
}

var _ addon.Service = (*service)(nil)

type service struct {
	adapter   addon.Repository
	taxCode   taxcode.Service
	logger    *slog.Logger
	publisher eventbus.Publisher

	featureResolver productcatalog.FeatureResolver
}
