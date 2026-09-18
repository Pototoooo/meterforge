import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/openapi3";
import "./shared/index.tsp";
import "./events/index.tsp";
import "./meters/index.tsp";
import "./customers/index.tsp";
import "./subscriptions/index.tsp";
import "./apps/index.tsp";
import "./billing/index.tsp";
import "./currencies/index.tsp";
import "./features/index.tsp";
import "./llmcost/index.tsp";
import "./productcatalog/index.tsp";
import "./defaults/index.tsp";
import "./governance/index.tsp";
import "./invoices/index.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

/**
 * MeterForge is a cloud-native usage metering and billing service. The API allows
 * you to ingest events, query meter usage, and manage resources.
 */
@service(#{ title: "MeterForge API" })
@info(#{ version: "0.0.1" })
@server(
  "http://localhost:{port}/api/v3",
  "Local development",
  {
    port?: string = "8888",
  }
)
@tagMetadata(Shared.MetersTag, #{ description: Shared.MetersDescription })
@tagMetadata(Shared.EventsTag, #{ description: Shared.EventsDescription })
@tagMetadata(Shared.CustomersTag, #{ description: Shared.CustomersDescription })
@tagMetadata(
  Shared.SubscriptionsTag,
  #{ description: Shared.SubscriptionsDescription }
)
@tagMetadata(
  Shared.EntitlementsTag,
  #{ description: Shared.EntitlementsDescription }
)
@tagMetadata(Shared.AppsTag, #{ description: Shared.AppsDescription })
@tagMetadata(Shared.BillingTag, #{ description: Shared.BillingDescription })
@tagMetadata(
  Shared.CurrenciesTag,
  #{ description: Shared.CurrenciesDescription }
)
@tagMetadata(Shared.TaxTag, #{ description: Shared.TaxDescription })
@tagMetadata(Shared.FeaturesTag, #{ description: Shared.FeaturesDescription })
@tagMetadata(Shared.LLMCostTag, #{ description: Shared.LLMCostDescription })
@tagMetadata(
  Shared.ProductCatalogTag,
  #{ description: Shared.ProductCatalogDescription }
)
@tagMetadata(Shared.DefaultsTag, #{ description: Shared.DefaultsDescription })
@tagMetadata(
  Shared.GovernanceTag,
  #{ description: Shared.GovernanceDescription }
)
namespace MeterForge;

@route("/meterforge/events")
@tag(Shared.EventsTag)
interface EventsEndpoints extends Events.EventsOperations {}

@route("/meterforge/meters")
@tag(Shared.MetersTag)
interface MetersEndpoints extends Meters.MetersOperations {}

@route("/meterforge/meters/{meterId}/query")
@tag(Shared.MetersTag)
interface MetersQueryEndpoints extends Meters.MetersQueryOperations {}

@route("/meterforge/customers")
@tag(Shared.CustomersTag)
interface CustomersEndpoints extends Customers.CustomersOperations {}

@route("/meterforge/customers/{customerId}/billing")
@tag(Shared.CustomersTag)
interface CustomerBillingEndpoints
  extends Customers.Billing.CustomerBillingOperations {}

@route("/meterforge/customers/{customerId}/entitlement-access")
@tag(Shared.EntitlementsTag)
interface CustomerEntitlementsEndpoints
  extends Entitlements.CustomerEntitlementsOperations {}

@route("/meterforge/customers/{customerId}/credits/grants")
@tag(Shared.CustomersTag)
interface CustomerCreditGrantsEndpoints
  extends Customers.Credits.Grants.CustomerCreditGrantsOperations {}

@route("/meterforge/customers/{customerId}/credits/balance")
@tag(Shared.CustomersTag)
interface CustomerCreditBalanceEndpoints
  extends Customers.Credits.Balance.CustomerCreditBalancesOperations {}

@route("/meterforge/customers/{customerId}/credits/adjustments")
@tag(Shared.CustomersTag)
interface CustomerCreditAdjustmentsEndpoints
  extends Customers.Credits.Adjustments.CustomerCreditAdjustmentsOperations {}

@route("/meterforge/customers/{customerId}/credits/grants/{creditGrantId}")
@tag(Shared.CustomersTag)
interface CustomerCreditGrantEndpoints {
  @route("/void")
  voidGrant is Customers.Credits.Grants.CustomerCreditGrantVoidOperations.voidGrant;

  @route("/settlement/external")
  updateExternalSettlement is Customers.Credits.Grants.CustomerCreditGrantExternalSettlementOperations.updateExternalSettlement;
}

@route("/meterforge/customers/{customerId}/credits/transactions")
@tag(Shared.CustomersTag)
interface CustomerCreditTransactionEndpoints
  extends Customers.Credits.Transactions.CustomerCreditTransactionOperations {}

@route("/meterforge/customers/{customerId}/charges")
@tag(Shared.CustomersTag)
interface CustomerChargesEndpoints
  extends Customers.Charges.CustomerChargesOperations {}

@route("/meterforge/subscriptions")
@tag(Shared.SubscriptionsTag)
interface SubscriptionsEndpoints
  extends Subscriptions.SubscriptionsOperations {}

@route("/meterforge/subscriptions/{subscriptionId}/addons")
@tag(Shared.SubscriptionsTag)
interface SubscriptionAddonEndpoints
  extends Subscriptions.SubscriptionAddonOperations {}

@route("/meterforge/apps")
@tag(Shared.AppsTag)
interface AppsEndpoints extends Apps.AppsOperations {}

@route("/meterforge/app-catalog")
@tag(Shared.AppsTag)
interface AppCatalogEndpoints extends Apps.AppCatalogOperations {}

@route("/meterforge/profiles")
@tag(Shared.BillingTag)
interface BillingProfilesEndpoints extends Billing.BillingProfilesOperations {}

@route("/meterforge/billing/invoices")
@tag(Shared.BillingTag)
interface BillingInvoicesEndpoints extends Invoices.InvoicesOperations {}

@route("/meterforge/tax-codes")
@tag(Shared.TaxTag)
interface TaxCodesEndpoints extends Tax.TaxCodesOperations {}

@route("/meterforge/currencies")
@tag(Shared.CurrenciesTag)
interface CurrenciesEndpoints extends Currencies.CurrenciesOperations {}

@route("/meterforge/currencies/custom")
@tag(Shared.CurrenciesTag)
interface CurrenciesCustomEndpoints
  extends Currencies.CurrenciesCustomOperations {}

@route("/meterforge/currencies/custom/{currencyId}/cost-bases")
@tag(Shared.CurrenciesTag)
interface CurrenciesCustomCostBasesEndpoints
  extends Currencies.CurrenciesCustomCostBasesOperations {}

@route("/meterforge/features")
@tag(Shared.FeaturesTag)
interface FeaturesEndpoints extends Features.FeatureOperations {}

@route("/meterforge/features/{featureId}/cost")
@tag(Shared.FeaturesTag)
interface FeatureCostEndpoints extends Features.FeatureCostOperations {}

@route("/meterforge/llm-cost/prices")
@tag(Shared.LLMCostTag)
interface LLMCostPricesEndpoints extends LLMCost.LLMCostPricesOperations {}

@route("/meterforge/llm-cost/overrides")
@tag(Shared.LLMCostTag)
interface LLMCostOverridesEndpoints
  extends LLMCost.LLMCostOverridesOperations {}

@route("/meterforge/plans")
@tag(Shared.ProductCatalogTag)
interface PlansEndpoints extends ProductCatalog.PlanOperations {}

@route("/meterforge/addons")
@tag(Shared.ProductCatalogTag)
interface AddonsEndpoints extends ProductCatalog.AddonOperations {}

@route("/meterforge/plans/{planId}/addons")
@tag(Shared.ProductCatalogTag)
interface PlanAddonEndpoints extends ProductCatalog.PlanAddonOperations {}

@route("/meterforge/defaults/tax-codes")
@tag(Shared.DefaultsTag)
interface OrganizationDefaultTaxCodesEndpoints
  extends Defaults.OrganizationDefaultTaxCodesOperations {}

@route("/meterforge/governance")
@tag(Shared.GovernanceTag)
interface GovernanceEndpoints extends Governance.GovernanceOperations {}
