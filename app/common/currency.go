package common

import (
	"fmt"

	"github.com/google/wire"

	"github.com/Pototoooo/meterforge/meterforge/currencies"
	currencyadapter "github.com/Pototoooo/meterforge/meterforge/currencies/adapter"
	"github.com/Pototoooo/meterforge/meterforge/currencies/currencyresolver"
	"github.com/Pototoooo/meterforge/meterforge/currencies/service"
	entdb "github.com/Pototoooo/meterforge/meterforge/ent/db"
)

var Currency = wire.NewSet(
	NewCurrencyAdapter,
	NewCurrencyService,
	NewCurrencyResolver,
)

var NewCurrencyResolver = currencyresolver.New

func NewCurrencyAdapter(db *entdb.Client) (currencies.Repository, error) {
	repo, err := currencyadapter.New(currencyadapter.Config{
		Client: db,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create currency adapter: %w", err)
	}

	return repo, nil
}

func NewCurrencyService(repo currencies.Repository) (currencies.Service, error) {
	s, err := service.New(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create currency service: %w", err)
	}

	return s, nil
}
