package mutator

import (
	"github.com/Pototoooo/meterforge/meterforge/billing/rating"
	"github.com/Pototoooo/meterforge/meterforge/billing/rating/service/rate"
)

type PostCalculationMutator interface {
	Mutate(rate.PricerCalculateInput, rating.DetailedLines) (rating.DetailedLines, error)
}

type PreCalculationMutator interface {
	Mutate(rate.PricerCalculateInput) (rate.PricerCalculateInput, error)
}
