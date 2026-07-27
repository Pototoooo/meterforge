package flatfee

import "github.com/Pototoooo/meterforge/meterforge/productcatalog"

type ProRatingModeAdapterEnum string

const (
	ProratePricesProratingAdapterMode ProRatingModeAdapterEnum = ProRatingModeAdapterEnum(productcatalog.ProRatingModeProratePrices)
	NoProratingAdapterMode            ProRatingModeAdapterEnum = "no_prorate"
)

func (e ProRatingModeAdapterEnum) Values() []string {
	return []string{
		string(ProratePricesProratingAdapterMode),
		string(NoProratingAdapterMode),
	}
}
