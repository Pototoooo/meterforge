namespace Shared;

const MeteringAndBillingTitle = "Metering & Billing";
const MeterForgeTitle = "MeterForge";

const MetersTag = "Meters";
const MetersDescription = "Meters specify how to aggregate events for billing and analytics purposes. Meters can be configured with multiple aggregation methods and groupings. Multiple meters can be created for the same event type, enabling flexible metering scenarios.";

const EventsTag = "Metering Events";
const EventsDescription = "Metering events are used to track usage of your product or service. Events are processed asynchronously by the meters, so they may not be immediately available for querying.";

const CustomersTag = "MeterForge Customers";
const CustomersDescription = "Customers are used to track usage of your product or service. Customers can be individuals or organizations that can subscribe to plans and have access to features.";

const EntitlementsTag = "MeterForge Entitlements";
const EntitlementsDescription = "Entitlements are used to control access to features for customers.";

const SubscriptionsTag = "MeterForge Subscriptions";
const SubscriptionsDescription = "Subscriptions are used to track usage of your product or service. Subscriptions can be individuals or organizations that can subscribe to plans and have access to features.";

const AppsTag = "MeterForge Apps";
const AppsDescription = "Apps enable you to extend and customize billing and usage workflows by integrating with external systems and services. Apps can automate and enhance your billing ecosystem by supporting capabilities such as synchronizing usage data with third-party platforms, calculating taxes, generating and delivering invoices, handling payment collection, and other billing-related tasks.";

const BillingTag = "MeterForge Billing Settings";
const BillingDescription = "Billing settings manages the billing profiles and invoices for customers.";

const CurrenciesTag = "MeterForge Currencies";
const CurrenciesDescription = "Currencies and cost bases for customers.";

const TaxTag = "MeterForge Tax";
const TaxDescription = "Tax codes are used to calculate taxes for customers.";

const LLMCostTag = "MeterForge LLM Cost";
const LLMCostDescription = "LLM cost database providing normalized pricing data for language models across providers. Prices are synced from multiple external sources and auto-confirmed when they agree.";

const FeaturesTag = "MeterForge Features";
const FeaturesDescription = "Features represent product capabilities backed by meters, with optional per-unit cost configuration.";

const ProductCatalogTag = "MeterForge Product Catalog";
const ProductCatalogDescription = "Product catalog manages plans, add-ons, and their associations for subscription-based billing.";

const DefaultsTag = "MeterForge Defaults";
const DefaultsDescription = "Organization-level default configuration.";

const GovernanceTag = "MeterForge Governance";
const GovernanceDescription = "Governance evaluation of customers to check their feature access.";

// The functionality should be considered unstable.
const UnstableExtension = "x-unstable";
// The functionality is internal and should not be used by customers.
const InternalExtension = "x-internal";
// The functionality is private and should not be exposed to customers.
const PrivateExtension = "x-private";
