package entitlementdriverv2

import (
	"github.com/Pototoooo/meterforge/meterforge/customer"
	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	meteredentitlement "github.com/Pototoooo/meterforge/meterforge/entitlement/metered"
	"github.com/Pototoooo/meterforge/meterforge/namespace/namespacedriver"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
)

// EntitlementHandler exposes V2 customer entitlement endpoints
type EntitlementHandler interface {
	CreateCustomerEntitlement() CreateCustomerEntitlementHandler
	ListCustomerEntitlements() ListCustomerEntitlementsHandler
	GetCustomerEntitlement() GetCustomerEntitlementHandler
	DeleteCustomerEntitlement() DeleteCustomerEntitlementHandler
	OverrideCustomerEntitlement() OverrideCustomerEntitlementHandler
	ListCustomerEntitlementGrants() ListCustomerEntitlementGrantsHandler
	CreateCustomerEntitlementGrant() CreateCustomerEntitlementGrantHandler
	GetCustomerEntitlementHistory() GetCustomerEntitlementHistoryHandler
	ResetCustomerEntitlementUsage() ResetCustomerEntitlementUsageHandler
	ListEntitlements() ListEntitlementsHandler
	GetEntitlement() GetEntitlementHandler
}

type entitlementHandler struct {
	namespaceDecoder namespacedriver.NamespaceDecoder
	options          []httptransport.HandlerOption
	connector        entitlement.Service
	balanceConnector meteredentitlement.Connector
	customerService  customer.Service
}

func NewEntitlementHandler(
	connector entitlement.Service,
	balanceConnector meteredentitlement.Connector,
	customerService customer.Service,
	namespaceDecoder namespacedriver.NamespaceDecoder,
	options ...httptransport.HandlerOption,
) EntitlementHandler {
	return &entitlementHandler{
		namespaceDecoder: namespaceDecoder,
		options:          options,
		connector:        connector,
		balanceConnector: balanceConnector,
		customerService:  customerService,
	}
}
