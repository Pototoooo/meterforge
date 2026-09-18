package currencies

import (
	"github.com/Pototoooo/meterforge/pkg/currencyx"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type CostBasis struct {
	models.ManagedModel
	models.NamespacedID
	currencyx.CostBasis

	CurrencyID string `json:"currency_id"`

	// CustomCurrency is included only if the CostBasis is expanded
	CustomCurrency *Currency `json:"-"`
}
