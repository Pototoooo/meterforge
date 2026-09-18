import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/openapi3";

import "..";

import "./auth.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

/**
 * MeterForge is a cloud native usage metering service.
 * The MeterForge API allows you to ingest events, query meter usage, and manage resources.
 */
@service(#{ title: "MeterForge Legacy Cloud-Compatible API" })
@info(#{
  version: "1.0.0",
})
@server("http://localhost:8888", "Self-hosted")
@useAuth(CloudTokenAuth | CloudCookieAuth)
// Tags will appear in alphabetical order in the API documentation
@tagMetadata(
  "Apps",
  #{
    description: "Manage integrations for extending MeterForge's functionality. ",
  }
)
@tagMetadata("App: Stripe", #{ description: "Support for Stripe billing." })
@tagMetadata(
  "App: Custom Invoicing",
  #{ description: "Interface third party invoicing and payment systems." }
)
@tagMetadata(
  "Billing",
  #{ description: "Manage your billing profiles and invoices. " }
)
@tagMetadata(
  "Customers",
  #{
    description: "Manage customer subscription lifecycles and plan assignments. ",
  }
)
@tagMetadata("Debug", #{ description: "Debugging and testing endpoints." })
@tagMetadata(
  "Entitlements",
  #{
    description: "With Entitlements, you can define and enforce usage limits, implement quota-based pricing, and manage access to features in your application.",
  }
)
@tagMetadata(
  "Events",
  #{
    description: "Events are used to track usage of your product or service. Events are processed asynchronously by the meters, so they may not be immediately available for querying.",
  }
)
@tagMetadata(
  "Lookup Information",
  #{ description: "Lookup information for static data like currencies" }
)
@tagMetadata(
  "Meters",
  #{
    description: "Meters specify how to aggregate events for billing and analytics purposes. Meters can be configured with multiple aggregation methods and groupings. Multiple meters can be created for the same event type, enabling flexible metering scenarios.",
  }
)
@tagMetadata(
  "Notifications",
  #{
    description: "Notifications provide automated triggers when specific entitlement balances and usage thresholds are reached, ensuring that your customers and sales teams are always informed. Notify customers and internal teams when specific conditions are met, like reaching 75%, 100%, and 150% of their monthly usage allowance. [Read more](https://github.com/Pototoooo/meterforge/tree/main/docs/guides/notifications/overview).",
  }
)
@tagMetadata(
  "Portal",
  #{
    description: "With the Consumer Portal, you can build in-app user-facing dashboards where your users can track their usage in real-time. Subject scoped portal tokens can be generated on your behalf to allow restricted access to the MeterForge API.",
  }
)
@tagMetadata(
  "Product Catalog",
  #{
    description: "Configure and manage your product plans, pricing tiers, and subscription offerings. ",
  }
)
@tagMetadata(
  "Subjects",
  #{
    description: "⚠️ __Deprecated__: Subjects as managable entities are being depracated, use customers with subject key usage attribution instead. \n\n Subjects are entities that consume resources you wish to meter. These can range from users, servers, and services to devices. The design of subjects is intentionally generic, enabling flexible application across various metering scenarios. Typically, a subject acts as a unique identifier within your system for a user or customer. Meters are aggregating events for each subject.",
  }
)
@tagMetadata(
  "Subscriptions",
  #{
    description: "With Subscriptions, you can easily start, cancel, and manage customer subscriptions. For example, provisioning them on a specific plan or assigning custom rate cards.",
  }
)
namespace MeterForgeCloud;

@route("/api/v1/events")
@tag("Events")
@friendlyName("Events")
interface EventsEndpoints extends MeterForge.EventsEndpoints {}

@route("/api/v2/events")
@tag("Events")
@friendlyName("EventsV2")
interface EventsV2Endpoints extends MeterForge.EventsV2Endpoints {}

@route("/api/v1/meters")
@tag("Meters")
@friendlyName("Meters")
interface MetersEndpoints extends MeterForge.MetersEndpoints {}

namespace MeterForgeCloud.Portal {
  @route("/api/v1/portal/tokens")
  @tag("Portal")
  @friendlyName("PortalTokens")
  interface TokensEndpoints extends MeterForge.Portal.TokensEndpoints {}

  @route("/api/v1/portal/meters")
  @tag("Portal")
  @useAuth(CloudPortalTokenAuth)
  @friendlyName("PortalMeters")
  interface MetersEndpoints extends MeterForge.Portal.MetersEndpoints {}
}

@route("/api/v1/debug")
@tag("Debug")
@friendlyName("Debug")
interface DebugEndpoints extends MeterForge.DebugEndpoints {}

namespace MeterForgeCloud {
  @route("/api/v1/marketplace")
  @tag("Apps")
  @friendlyName("Marketplace")
  interface MarketplaceEndpoints extends MeterForge.MarketplaceEndpoints {}

  @route("/api/v1/apps")
  @tag("Apps")
  @friendlyName("Apps")
  interface AppsEndpoints extends MeterForge.AppsEndpoints {}

  @tag("App: Stripe")
  @friendlyName("AppStripe")
  interface AppStripeEndpoints extends MeterForge.AppStripeEndpoints {}

  @tag("App: Custom Invoicing")
  @friendlyName("AppCustomInvoicing")
  interface AppCustomInvoicingEndpoints
    extends MeterForge.AppCustomInvoicingEndpoints {}
}

namespace MeterForgeCloud {
  @route("/api/v1/notification/channels")
  @tag("Notifications")
  @friendlyName("NotificationChannels")
  interface ChannelsEndpoints
    extends MeterForge.NotificationChannelsEndpoints {}

  @route("/api/v1/notification/rules")
  @tag("Notifications")
  @friendlyName("NotificationRules")
  interface RulesEndpoints extends MeterForge.NotificationRulesEndpoints {}

  @route("/api/v1/notification/events")
  @tag("Notifications")
  @friendlyName("NotificationEvents")
  interface EventsEndpoints extends MeterForge.NotificationEventsEndpoints {}
}

namespace MeterForgeCloud {
  @route("/api/v1/entitlements")
  @tag("Entitlements")
  @friendlyName("Entitlements")
  interface EntitlementsEndpoints extends MeterForge.EntitlementsEndpoints {}

  @route("/api/v1/grants")
  @tag("Entitlements")
  @friendlyName("Grants")
  interface GrantsEndpoints extends MeterForge.GrantsEndpoints {}

  @route("/api/v1/subjects/{subjectIdOrKey}/entitlements")
  @tag("Entitlements")
  @friendlyName("Subjects")
  interface SubjectEntitlementsEndpoints
    extends MeterForge.SubjectEntitlementsEndpoints {}

  @route("/api/v1/customers/{customerIdOrKey}")
  @tag("Entitlements")
  @tag("Customers")
  @friendlyName("Customer")
  interface CustomerEndpoints extends MeterForge.CustomerEndpoints {}

  @route("/api/v1/customers/{customerIdOrKey}/apps")
  @tag("Customers")
  @tag("Apps")
  @friendlyName("CustomerApps")
  interface CustomerAppsEndpoints extends MeterForge.CustomerAppsEndpoints {}

  @route("/api/v1/customers/{customerIdOrKey}/stripe")
  @tag("Customers")
  @friendlyName("CustomerStripe")
  interface CustomerStripeEndpoints
    extends MeterForge.CustomerStripeEndpoints {}

  @route("/api/v1/customers/{customerIdOrKey}/entitlements/{featureKey}")
  @tag("Entitlements")
  @tag("Customers")
  @friendlyName("CustomerEntitlement")
  interface CustomerEntitlementEndpoints
    extends MeterForge.CustomerEntitlementEndpoints {}

  namespace V2 {
    @route("/api/v2/customers/{customerIdOrKey}/entitlements")
    @tag("Entitlements")
    @tag("Customers")
    @friendlyName("CustomerEntitlementsV2")
    interface CustomerEntitlementsV2Endpoints
      extends MeterForge.CustomerEntitlementsV2Endpoints {}

    @route("/api/v2/customers/{customerIdOrKey}/entitlements/{entitlementIdOrFeatureKey}")
    @tag("Entitlements")
    @tag("Customers")
    @friendlyName("CustomerEntitlementV2")
    interface CustomerEntitlementEndpoints
      extends MeterForge.CustomerEntitlementV2Endpoints {}

    @route("/api/v2/grants")
    @tag("Entitlements")
    @friendlyName("GrantsV2")
    interface GrantsV2Endpoints extends MeterForge.GrantsV2Endpoints {}

    @route("/api/v2/entitlements")
    @tag("Entitlements")
    @friendlyName("EntitlementsV2")
    interface EntitlementsV2Endpoints
      extends MeterForge.EntitlementsV2Endpoints {}
  }
}

namespace MeterForgeCloud {
  @route("/api/v1/customers")
  @tag("Customers")
  @friendlyName("Customers")
  interface CustomersEndpoints extends MeterForge.CustomersEndpoints {}

  @route("/api/v1/billing/customers")
  @tag("Billing")
  @friendlyName("CustomerOverrides")
  interface CustomerOverridesEndpoints
    extends MeterForge.CustomerOverridesEndpoints {}

  @route("/api/v1/billing/profiles")
  @tag("Billing")
  @friendlyName("BillingProfiles")
  interface BillingProfilesEndpoints
    extends MeterForge.BillingProfilesEndpoints {}

  @route("/api/v1/billing/invoices")
  @tag("Billing")
  @friendlyName("Invoices")
  interface InvoicesEndpoints extends MeterForge.InvoicesEndpoints {}

  @route("/api/v1/billing/invoices/{invoiceId}")
  @tag("Billing")
  @friendlyName("Invoice")
  interface InvoiceEndpoints extends MeterForge.InvoiceEndpoints {}

  @route("/api/v1/billing/customers/{customerId}/invoices")
  @tag("Billing")
  @friendlyName("CustomerInvoice")
  interface CustomerInvoiceEndpoints
    extends MeterForge.CustomerInvoiceEndpoints {}
}

namespace MeterForgeCloud {
  @route("/api/v1/features")
  @tag("Product Catalog")
  @friendlyName("Features")
  interface FeaturesEndpoints extends MeterForge.FeaturesEndpoints {}

  @route("/api/v1/plans")
  @tag("Product Catalog")
  @friendlyName("Plans")
  interface PlansEndpoints extends MeterForge.PlansEndpoints {}

  @route("/api/v1/addons")
  @tag("Product Catalog")
  @friendlyName("Addons")
  interface AddonsEndpoints extends MeterForge.AddonsEndpoints {}

  @route("/api/v1/plans/{planId}/addons")
  @tag("Product Catalog")
  @friendlyName("PlanAddons")
  interface PlanAddonsEndpoints extends MeterForge.PlanAddonsEndpoints {}

  @route("/api/v1/subscriptions")
  @tag("Subscriptions")
  @friendlyName("Subscriptions")
  interface SubscriptionsEndpoints extends MeterForge.SubscriptionsEndpoints {}

  @route("/api/v1/subscriptions/{subscriptionId}/addons")
  @tag("Subscriptions")
  @friendlyName("SubscriptionAddons")
  interface SubscriptionAddonsEndpoints
    extends MeterForge.SubscriptionAddonsEndpoints {}
}

namespace MeterForgeCloud {
  @route("/api/v1/subjects")
  @tag("Subjects")
  @friendlyName("Subjects")
  interface SubjectsEndpoints extends MeterForge.SubjectsEndpoints {}
}

namespace MeterForgeCloud {
  @route("/api/v1/info/currencies")
  @tag("Lookup Information")
  @friendlyName("Currencies")
  interface CurrenciesEndpoints extends MeterForge.CurrenciesEndpoints {}

  @route("/api/v1/info/progress")
  @tag("Lookup Information")
  @friendlyName("Progress")
  interface ProgressEndpoints extends MeterForge.ProgressEndpoints {}
}
