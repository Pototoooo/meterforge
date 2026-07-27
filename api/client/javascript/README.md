# MeterForge JavaScript SDK

## Install

```sh
npm install --save @meterforge/sdk
```

## Configuration for accessing the MeterForge API

To use the MeterForge SDK on your backend, you need to configure `baseUrl` and `apiKey` for MeterForge Cloud:

```ts
import { MeterForge } from '@meterforge/sdk'

const meterforge = new MeterForge({
  baseUrl: 'http://localhost:8888',
  apiKey: 'mf_...',
})
```

## Configuration for accessing the MeterForge Portal API

To use the MeterForge Portal SDK on your frontend, you need to configure it use a portal token in your configuration:

```ts
import { MeterForge } from '@meterforge/sdk/portal'

const meterforge = new MeterForge({
  baseUrl: 'http://localhost:8888',
  portalToken: 'mf_portal_...',
})
```

## Configuration for accessing the MeterForge React SDK

To use the MeterForge React SDK for the portal API, you need to configure a Portal Client and a React Context:

```ts
import {
  MeterForge,
  MeterForgeProvider,
  useMeterForge,
} from '@meterforge/sdk/react'

function App() {
  // get portal token from your backend
  const meterforge = new MeterForge({
    baseUrl: 'http://localhost:8888',
    portalToken,
  })

  return (
    <MeterForgeProvider value={meterforge}>
      <UsageComponent />
      {/* ... */}
    </MeterForgeProvider>
  )
}

function UsageComponent() {
  // get meterforge client from context
  const meterforge = useMeterForge()

  // ...
}
```

## Ingest usage events

```ts
// Ingest a single AI token usage event
await meterforge.events.ingest({
  type: 'ai-tokens',
  subject: 'customer-acme-corp',
  id: 'evt_01234567', // optional: auto-generated if not provided
  source: 'llm-api-gateway', // optional: defaults to '@meterforge/sdk'
  time: new Date(), // optional: defaults to current time
  data: {
    model: 'gpt-4',
    type: 'input',
    tokens: 1250,
  },
})

// Ingest multiple events in a batch for better performance
await meterforge.events.ingest([
  {
    type: 'ai-tokens',
    subject: 'customer-acme-corp',
    data: { model: 'gpt-4', type: 'input', tokens: 850 },
  },
  {
    type: 'ai-tokens',
    subject: 'customer-acme-corp',
    data: { model: 'gpt-4', type: 'output', tokens: 850 },
  },
])
```

## Client API Reference

The MeterForge SDK provides a comprehensive client interface organized into logical groups. Below is a complete reference of all available methods.

### Overview

| Namespace | Resource | Operation | Method | Description |
|-----------|----------|-----------|--------|-------------|
| **[Events](#events)** | | | | Track usage by ingesting events |
| | Events | Create | [`meterforge.events.ingest(events)`](./src/client/events.ts#L19) | Ingest a single event or batch of events |
| | Events | Read | [`meterforge.events.list(params?)`](./src/client/events.ts#L41) | List ingested events with advanced filtering |
| | Events | Read | [`meterforge.events.listV2(params?)`](./src/client/events.ts#L60) | List ingested events with advanced filtering (V2) |
| **[Meters](#meters)** | | | | Track and aggregate usage data from events |
| | Meters | Create | [`meterforge.meters.create(meter)`](./src/client/meters.ts#L19) | Create a new meter |
| | Meters | Read | [`meterforge.meters.get(idOrSlug)`](./src/client/meters.ts#L34) | Get a meter by ID or slug |
| | Meters | Read | [`meterforge.meters.list()`](./src/client/meters.ts#L55) | List all meters |
| | Meters | Read | [`meterforge.meters.query(idOrSlug, query?)`](./src/client/meters.ts#L70) | Query usage data |
| | Meters | Update | [`meterforge.meters.update(idOrSlug, meter)`](./src/client/meters.ts#L100) | Update a meter by ID or slug |
| | Meters | Delete | [`meterforge.meters.delete(idOrSlug)`](./src/client/meters.ts#L124) | Delete a meter by ID or slug |
| **[Subjects](#subjects)** | | | | Manage entities that consume resources |
| | Subjects | Create | [`meterforge.subjects.upsert(subjects)`](./src/client/subjects.ts#L21) | Create or update one or multiple subjects |
| | Subjects | Read | [`meterforge.subjects.get(idOrKey)`](./src/client/subjects.ts#L39) | Get a subject by ID or key |
| | Subjects | Read | [`meterforge.subjects.list()`](./src/client/subjects.ts#L60) | List all subjects |
| | Subjects | Delete | [`meterforge.subjects.delete(idOrKey)`](./src/client/subjects.ts#L74) | Delete a subject by ID or key |
| **[Customers](#customers)** | | | | Manage customer subscription lifecycles and plan assignments |
| | Customers | Create | [`meterforge.customers.create(customer)`](./src/client/customers.ts#L37) | Create a new customer |
| | Customers | Read | [`meterforge.customers.get(customerIdOrKey)`](./src/client/customers.ts#L52) | Get a customer by ID or key |
| | Customers | Read | [`meterforge.customers.list(query?)`](./src/client/customers.ts#L123) | List all customers |
| | Customers | Read | [`meterforge.customers.getAccess(customerIdOrKey)`](./src/client/customers.ts#L143) | Get customer access information |
| | Customers | Read | [`meterforge.customers.listSubscriptions(customerIdOrKey, query?)`](./src/client/customers.ts#L169) | List customer subscriptions |
| | Customers | Update | [`meterforge.customers.update(customerIdOrKey, customer)`](./src/client/customers.ts#L75) | Update a customer |
| | Customers | Delete | [`meterforge.customers.delete(customerIdOrKey)`](./src/client/customers.ts#L99) | Delete a customer |
| | Apps | Update | [`meterforge.customers.apps.upsert(customerIdOrKey, appData)`](./src/client/customers.ts#L200) | Upsert app data |
| | Apps | Read | [`meterforge.customers.apps.list(customerIdOrKey)`](./src/client/customers.ts#L228) | List app data |
| | Apps | Delete | [`meterforge.customers.apps.delete(customerIdOrKey, appId)`](./src/client/customers.ts#L254) | Delete app data |
| | Stripe | Update | [`meterforge.customers.stripe.upsert(customerIdOrKey, appDataBase)`](./src/client/customers.ts#L285) | Upsert Stripe app data |
| | Stripe | Read | [`meterforge.customers.stripe.get(customerIdOrKey)`](./src/client/customers.ts#L313) | Get Stripe app data |
| | Stripe | Create | [`meterforge.customers.stripe.createPortalSession(customerIdOrKey, params)`](./src/client/customers.ts#L337) | Create a Stripe customer portal session |
| | Entitlements V1 | Read | [`meterforge.customers.entitlementsV1.value(customerIdOrKey, featureKey)`](./src/client/customers.ts#L372) | Get entitlement value (V1 API) |
| | Entitlements | Read | [`meterforge.customers.entitlements.list(customerIdOrKey)`](./src/client/customers.ts#L401) | List entitlements |
| | Entitlements | Create | [`meterforge.customers.entitlements.create(customerIdOrKey, entitlement)`](./src/client/customers.ts#L428) | Create an entitlement |
| | Entitlements | Read | [`meterforge.customers.entitlements.get(customerIdOrKey, featureKeyOrId)`](./src/client/customers.ts#L454) | Get an entitlement |
| | Entitlements | Delete | [`meterforge.customers.entitlements.delete(customerIdOrKey, entitlementId)`](./src/client/customers.ts#L479) | Delete an entitlement |
| | Entitlements | Update | [`meterforge.customers.entitlements.override(customerIdOrKey, featureKeyOrId, entitlement)`](./src/client/customers.ts#L505) | Override an entitlement |
| | Entitlements | Read | [`meterforge.customers.entitlements.value(customerIdOrKey, featureKeyOrId, query?)`](./src/client/customers.ts#L588) | Get entitlement value |
| | Entitlements | Read | [`meterforge.customers.entitlements.history(customerIdOrKey, featureKeyOrId, query?)`](./src/client/customers.ts#L617) | Get entitlement history |
| | Entitlements | Update | [`meterforge.customers.entitlements.resetUsage(customerIdOrKey, entitlementId, body?)`](./src/client/customers.ts#L653) | Reset usage |
| | Entitlements | Read | [`meterforge.customers.entitlements.listGrants(customerIdOrKey, featureKeyOrId, query?)`](./src/client/customers.ts#L532) | List grants |
| | Entitlements | Create | [`meterforge.customers.entitlements.createGrant(customerIdOrKey, featureKeyOrId, grant)`](./src/client/customers.ts#L561) | Create a grant |
| **[Features](#features)** | | | | Define application capabilities and services |
| | Features | Create | [`meterforge.features.create(feature)`](./src/client/features.ts#L24) | Create a new feature |
| | Features | Read | [`meterforge.features.get(featureIdOrKey)`](./src/client/features.ts#L39) | Get a feature by ID |
| | Features | Read | [`meterforge.features.list(params?)`](./src/client/features.ts#L61) | List all features |
| | Features | Delete | [`meterforge.features.delete(featureIdOrKey)`](./src/client/features.ts#L84) | Delete a feature by ID |
| **[Entitlements (V1)](#entitlements-v1)** | | | | Subject-based usage limits and access controls |
| | Entitlements | Create | [`meterforge.entitlementsV1.create(subjectIdOrKey, entitlement)`](./src/client/entitlements.ts#L40) | Create an entitlement for a subject |
| | Entitlements | Read | [`meterforge.entitlementsV1.get(entitlementId)`](./src/client/entitlements.ts#L68) | Get an entitlement by ID |
| | Entitlements | Read | [`meterforge.entitlementsV1.list(query?)`](./src/client/entitlements.ts#L91) | List all entitlements |
| | Entitlements | Read | [`meterforge.entitlementsV1.value(subjectIdOrKey, featureIdOrKey, query?)`](./src/client/entitlements.ts#L147) | Get the value of an entitlement |
| | Entitlements | Read | [`meterforge.entitlementsV1.history(subjectIdOrKey, entitlementIdOrFeatureKey, query?)`](./src/client/entitlements.ts#L180) | Get the history of an entitlement |
| | Entitlements | Update | [`meterforge.entitlementsV1.override(subjectIdOrKey, entitlementIdOrFeatureKey, override)`](./src/client/entitlements.ts#L213) | Override an entitlement |
| | Entitlements | Update | [`meterforge.entitlementsV1.reset(subjectIdOrKey, entitlementIdOrFeatureKey, reset?)`](./src/client/entitlements.ts#L247) | Reset entitlement usage |
| | Entitlements | Delete | [`meterforge.entitlementsV1.delete(subjectIdOrKey, entitlementId)`](./src/client/entitlements.ts#L116) | Delete an entitlement |
| | Grants | Create | [`meterforge.entitlementsV1.grants.create(subjectIdOrKey, entitlementIdOrFeatureKey, grant)`](./src/client/entitlements.ts#L283) | Create a grant for an entitlement |
| | Grants | Read | [`meterforge.entitlementsV1.grants.list(subjectIdOrKey, entitlementIdOrFeatureKey, query?)`](./src/client/entitlements.ts#L314) | List grants for an entitlement |
| | Grants | Read | [`meterforge.entitlementsV1.grants.listAll(query?)`](./src/client/entitlements.ts#L345) | List all grants |
| | Grants | Delete | [`meterforge.entitlementsV1.grants.void(entitlementId, grantId)`](./src/client/entitlements.ts#L369) | Void a grant |
| **[Entitlements](#entitlements)** | | | | Customer-based entitlements and access controls |
| | Entitlements | Read | [`meterforge.entitlements.list(query?)`](./src/client/entitlements.ts#L404) | List all entitlements (admin purposes) |
| | Entitlements | Read | [`meterforge.entitlements.get(entitlementId)`](./src/client/entitlements.ts#L425) | Get an entitlement by ID |
| | Grants | Read | [`meterforge.entitlements.grants.list(query?)`](./src/client/entitlements.ts#L453) | List all grants (admin purposes) |
| | Grants | Delete | [`meterforge.entitlements.grants.void(grantId)`](./src/client/entitlements.ts#L478) | Void a grant |
| **[Plans](#plans)** | | | | Manage subscription plans and pricing|
| | Plans | Create | [`meterforge.plans.create(plan)`](./src/client/plans.ts#L28) | Create a new plan|
| | Plans | Read | [`meterforge.plans.get(planId)`](./src/client/plans.ts#L44) | Get a plan by ID|
| | Plans | Read | [`meterforge.plans.list(query?)`](./src/client/plans.ts#L66) | List all plans|
| | Plans | Update | [`meterforge.plans.update(planId, plan)`](./src/client/plans.ts#L85) | Update a plan|
| | Plans | Delete | [`meterforge.plans.delete(planId)`](./src/client/plans.ts#L105) | Delete a plan by ID|
| | Plans | Other | [`meterforge.plans.archive(planId)`](./src/client/plans.ts#L123) | Archive a plan|
| | Plans | Other | [`meterforge.plans.publish(planId)`](./src/client/plans.ts#L141) | Publish a plan|
| | Addons | Read | [`meterforge.plans.addons.list(planId)`](./src/client/plans.ts#L168) | List addons|
| | Addons | Create | [`meterforge.plans.addons.create(planId, addon)`](./src/client/plans.ts#L191) | Create an addon|
| | Addons | Read | [`meterforge.plans.addons.get(planId, planAddonId)`](./src/client/plans.ts#L212) | Get an addon by ID|
| | Addons | Update | [`meterforge.plans.addons.update(planId, planAddonId, addon)`](./src/client/plans.ts#L238) | Update an addon|
| | Addons | Delete | [`meterforge.plans.addons.delete(planId, planAddonId)`](./src/client/plans.ts#L263) | Delete an addon by ID|
| **[Addons](#addons)** | | | | Manage standalone addons available across plans|
| | Addons | Create | [`meterforge.addons.create(addon)`](./src/client/addons.ts#L15) | Create a new addon|
| | Addons | Read | [`meterforge.addons.get(addonId)`](./src/client/addons.ts#L48) | Get an addon by ID|
| | Addons | Read | [`meterforge.addons.list(query?)`](./src/client/addons.ts#L30) | List all addons|
| | Addons | Update | [`meterforge.addons.update(addonId, addon)`](./src/client/addons.ts#L64) | Update an addon|
| | Addons | Delete | [`meterforge.addons.delete(addonId)`](./src/client/addons.ts#L84) | Delete an addon by ID|
| | Addons | Other | [`meterforge.addons.publish(addonId)`](./src/client/addons.ts#L99) | Publish an addon|
| | Addons | Other | [`meterforge.addons.archive(addonId)`](./src/client/addons.ts#L114) | Archive an addon|
| **[Subscriptions](#subscriptions)** | | | | Manage customer subscriptions|
| | Subscriptions | Create | [`meterforge.subscriptions.create(body)`](./src/client/subscriptions.ts#L24) | Create a new subscription|
| | Subscriptions | Read | [`meterforge.subscriptions.get(subscriptionId)`](./src/client/subscriptions.ts#L39) | Get a subscription by ID|
| | Subscriptions | Update | [`meterforge.subscriptions.edit(subscriptionId, body)`](./src/client/subscriptions.ts#L61) | Edit a subscription|
| | Subscriptions | Delete | [`meterforge.subscriptions.delete(subscriptionId)`](./src/client/subscriptions.ts#L180) | Delete a subscription (only scheduled)|
| | Subscriptions | Other | [`meterforge.subscriptions.cancel(subscriptionId, body?)`](./src/client/subscriptions.ts#L85) | Cancel a subscription|
| | Subscriptions | Other | [`meterforge.subscriptions.change(subscriptionId, body)`](./src/client/subscriptions.ts#L110) | Change a subscription (upgrade/downgrade)|
| | Subscriptions | Other | [`meterforge.subscriptions.migrate(subscriptionId, body)`](./src/client/subscriptions.ts#L135) | Migrate to a new plan version|
| | Subscriptions | Other | [`meterforge.subscriptions.unscheduleCancelation(subscriptionId)`](./src/client/subscriptions.ts#L158) | Unschedule a subscription cancelation|
| **[Subscription Addons](#subscription-addons)** | | | | Manage addons attached to specific subscriptions|
| | Subscription Addons | Create | [`meterforge.subscriptionAddons.create(subscriptionId, body)`](./src/client/subscription-addons.ts#L16) | Create a new subscription addon|
| | Subscription Addons | Read | [`meterforge.subscriptionAddons.get(subscriptionId, subscriptionAddonId)`](./src/client/subscription-addons.ts#L58) | Get a subscription addon by ID|
| | Subscription Addons | Read | [`meterforge.subscriptionAddons.list(subscriptionId)`](./src/client/subscription-addons.ts#L39) | List all addons of a subscription|
| | Subscription Addons | Update | [`meterforge.subscriptionAddons.update(subscriptionId, subscriptionAddonId, body)`](./src/client/subscription-addons.ts#L82) | Update a subscription addon|
| **[Billing](#billing)** | | | | Comprehensive billing management (profiles, invoices, overrides)|
| | Profiles | Create | [`meterforge.billing.profiles.create(profile)`](./src/client/billing.ts#L42) | Create a billing profile|
| | Profiles | Read | [`meterforge.billing.profiles.get(id)`](./src/client/billing.ts#L60) | Get a billing profile by ID|
| | Profiles | Read | [`meterforge.billing.profiles.list(query?)`](./src/client/billing.ts#L80) | List billing profiles|
| | Profiles | Update | [`meterforge.billing.profiles.update(id, profile)`](./src/client/billing.ts#L101) | Update a billing profile|
| | Profiles | Delete | [`meterforge.billing.profiles.delete(id)`](./src/client/billing.ts#L123) | Delete a billing profile|
| | Invoices | Read | [`meterforge.billing.invoices.list(query?)`](./src/client/billing.ts#L150) | List invoices|
| | Invoices | Read | [`meterforge.billing.invoices.get(id, query?)`](./src/client/billing.ts#L170) | Get an invoice by ID|
| | Invoices | Update | [`meterforge.billing.invoices.update(id, invoice)`](./src/client/billing.ts#L192) | Update an invoice (draft or earlier)|
| | Invoices | Delete | [`meterforge.billing.invoices.delete(id)`](./src/client/billing.ts#L213) | Delete an invoice (draft or earlier)|
| | Invoices | Other | [`meterforge.billing.invoices.advance(id)`](./src/client/billing.ts#L235) | Advance invoice to next status|
| | Invoices | Other | [`meterforge.billing.invoices.approve(id)`](./src/client/billing.ts#L257) | Approve an invoice (sends to customer)|
| | Invoices | Other | [`meterforge.billing.invoices.retry(id, body?)`](./src/client/billing.ts#L278) | Retry advancing after failure|
| | Invoices | Other | [`meterforge.billing.invoices.void(id)`](./src/client/billing.ts#L302) | Void an invoice|
| | Invoices | Other | [`meterforge.billing.invoices.recalculateTax(id)`](./src/client/billing.ts#L325) | Recalculate invoice tax amounts|
| | Invoices | Other | [`meterforge.billing.invoices.simulate(customerId, query?)`](./src/client/billing.ts#L346) | Simulate an invoice for a customer|
| | Invoices | Create | [`meterforge.billing.invoices.createLineItems(customerId, body)`](./src/client/billing.ts#L377) | Create pending line items|
| | Invoices | Create | [`meterforge.billing.invoices.invoicePendingLines(customerId)`](./src/client/billing.ts#L401) | Invoice pending lines|
| | Customers | Create | [`meterforge.billing.customers.createOverride(customerId, body)`](./src/client/billing.ts#L427) | Create or update a customer override|
| | Customers | Read | [`meterforge.billing.customers.getOverride(customerId, id)`](./src/client/billing.ts#L450) | Get a customer override|
| | Customers | Read | [`meterforge.billing.customers.listOverrides(customerId)`](./src/client/billing.ts#L471) | List customer overrides|
| | Customers | Delete | [`meterforge.billing.customers.deleteOverride(customerId, id)`](./src/client/billing.ts#L489) | Delete a customer override|
| **[Apps](#apps)** | | | | Manage integrations and app marketplace|
| | Apps | Read | [`meterforge.apps.list(query?)`](./src/client/apps.ts#L32) | List installed apps|
| | Apps | Read | [`meterforge.apps.get(id)`](./src/client/apps.ts#L50) | Get an app by ID|
| | Apps | Update | [`meterforge.apps.update(id, body)`](./src/client/apps.ts#L69) | Update an app|
| | Apps | Delete | [`meterforge.apps.uninstall(id)`](./src/client/apps.ts#L89) | Uninstall an app|
| | Marketplace | Read | [`meterforge.apps.marketplace.list(query?)`](./src/client/apps.ts#L115) | List available marketplace apps|
| | Marketplace | Read | [`meterforge.apps.marketplace.get(id)`](./src/client/apps.ts#L133) | Get marketplace listing details|
| | Marketplace | Read | [`meterforge.apps.marketplace.getOauth2InstallUrl(id, redirectUrl)`](./src/client/apps.ts#L151) | Get OAuth2 install URL|
| | Marketplace | Other | [`meterforge.apps.marketplace.authorizeOauth2(id, body)`](./src/client/apps.ts#L172) | Authorize OAuth2 code|
| | Marketplace | Create | [`meterforge.apps.marketplace.installWithAPIKey(id, body)`](./src/client/apps.ts#L193) | Install app with API key|
| | Stripe | Create | [`meterforge.apps.stripe.createCheckoutSession(body)`](./src/client/apps.ts#L223) | Create a Stripe checkout session|
| | Stripe | Update | [`meterforge.apps.stripe.updateApiKey(body)`](./src/client/apps.ts#L243) | Update Stripe API key|
| | Custom Invoicing | Other | [`meterforge.apps.customInvoicing.draftSynchronized(body)`](./src/client/apps.ts#L271) | Submit draft synchronization results|
| | Custom Invoicing | Other | [`meterforge.apps.customInvoicing.issuingSynchronized(body)`](./src/client/apps.ts#L295) | Submit issuing synchronization results|
| | Custom Invoicing | Update | [`meterforge.apps.customInvoicing.updatePaymentStatus(invoiceId, body)`](./src/client/apps.ts#L319) | Update payment status|
| **[Notifications](#notifications)** | | | | Set up automated notifications for usage thresholds|
| | Channels | Create | [`meterforge.notifications.channels.create(channel)`](./src/client/notifications.ts#L40) | Create a notification channel|
| | Channels | Read | [`meterforge.notifications.channels.get(channelId)`](./src/client/notifications.ts#L58) | Get a notification channel by ID|
| | Channels | Update | [`meterforge.notifications.channels.update(channelId, channel)`](./src/client/notifications.ts#L84) | Update a notification channel|
| | Channels | Read | [`meterforge.notifications.channels.list(query?)`](./src/client/notifications.ts#L111) | List notification channels|
| | Channels | Delete | [`meterforge.notifications.channels.delete(channelId)`](./src/client/notifications.ts#L131) | Delete a notification channel|
| | Rules | Create | [`meterforge.notifications.rules.create(rule)`](./src/client/notifications.ts#L164) | Create a notification rule|
| | Rules | Read | [`meterforge.notifications.rules.get(ruleId)`](./src/client/notifications.ts#L182) | Get a notification rule by ID|
| | Rules | Update | [`meterforge.notifications.rules.update(ruleId, rule)`](./src/client/notifications.ts#L205) | Update a notification rule|
| | Rules | Read | [`meterforge.notifications.rules.list(query?)`](./src/client/notifications.ts#L229) | List notification rules|
| | Rules | Delete | [`meterforge.notifications.rules.delete(ruleId)`](./src/client/notifications.ts#L249) | Delete a notification rule|
| | Events | Read | [`meterforge.notifications.events.get(eventId)`](./src/client/notifications.ts#L282) | Get a notification event by ID|
| | Events | Read | [`meterforge.notifications.events.list(query?)`](./src/client/notifications.ts#L307) | List notification events|
| **[Portal](#portal)** | | | | Manage consumer portal tokens for customer-facing interfaces|
| | Portal | Create | [`meterforge.portal.create(body)`](./src/client/portal.ts#L19) | Create a consumer portal token|
| | Portal | Read | [`meterforge.portal.list(query?)`](./src/client/portal.ts#L34) | List consumer portal tokens|
| | Portal | Other | [`meterforge.portal.invalidate(query?)`](./src/client/portal.ts#L52) | Invalidate consumer portal tokens|
| **[Info](#info)** | | | | Utility endpoints for system information|
| | Info | Read | [`meterforge.info.listCurrencies()`](./src/client/info.ts#L18) | List all supported currencies|
| | Info | Read | [`meterforge.info.getProgress(id)`](./src/client/info.ts#L32) | Get progress of a long-running operation|
| **[Debug](#debug)** | | | | Debug utilities for monitoring and troubleshooting|
| | Debug | Read | [`meterforge.debug.getMetrics()`](./src/client/debug.ts#L18) | Get event ingestion metrics (OpenMetrics format)|

