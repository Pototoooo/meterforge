package registry

import (
	"github.com/Pototoooo/meterforge/meterforge/credit"
	"github.com/Pototoooo/meterforge/meterforge/credit/grant"
	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	meteredentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/metered"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
)

type Entitlement struct {
	Feature            feature.FeatureConnector
	FeatureRepo        feature.FeatureRepo
	EntitlementOwner   grant.OwnerConnector
	CreditBalance      credit.BalanceConnector
	Grant              credit.GrantConnector
	GrantRepo          grant.Repo
	MeteredEntitlement meteredentitlement.Connector
	Entitlement        entitlement.Service
	EntitlementRepo    entitlement.EntitlementRepo
}
