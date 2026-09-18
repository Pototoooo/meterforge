package adapter

import (
	dbchargecreditpurchase "github.com/Pototoooo/meterforge/meterforge/ent/db/chargecreditpurchase"
	dbcustomcurrency "github.com/Pototoooo/meterforge/meterforge/ent/db/customcurrency"
	"github.com/Pototoooo/meterforge/meterforge/ent/db/predicate"
	"github.com/Pototoooo/meterforge/pkg/currencyx"
)

func hasCustomCurrencyCode(namespace string, codes ...currencyx.Code) predicate.ChargeCreditPurchase {
	return dbchargecreditpurchase.HasCustomCurrencyWith(
		dbcustomcurrency.CodeIn(codes...),
		dbcustomcurrency.Namespace(namespace),
		dbcustomcurrency.DeletedAtIsNil(),
	)
}
