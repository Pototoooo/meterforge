# MeterForge Go SDK

Go client for the MeterForge API — usage metering and billing for
AI and DevTool companies. This package is generated from the MeterForge
TypeSpec definitions and ships typed request and response models.

> [!IMPORTANT]
> This SDK is a work in progress.
>
> This SDK targets the [MeterForge API v3](https://github.com/Pototoooo/meterforge/tree/main/docs/api/v3),
> a rewrite of the MeterForge API following AIP (API Improvement Proposal)
> standardization.

## Table of Contents

- [Installation](#installation)
- [Initialization](#initialization)
- [Usage](#usage)
- [Available Resources and Operations](#available-resources-and-operations)
  - [Events](#events)
  - [Meters](#meters)
  - [Customers](#customers)
  - [Entitlements](#entitlements)
  - [Subscriptions](#subscriptions)
  - [Apps](#apps)
  - [Billing](#billing)
  - [Invoices](#invoices)
  - [Tax](#tax)
  - [Currencies](#currencies)
  - [Features](#features)
  - [LLMCost](#llmcost)
  - [Plans](#plans)
  - [Addons](#addons)
  - [PlanAddons](#planaddons)
  - [Defaults](#defaults)
  - [Governance](#governance)
- [Error Handling](#error-handling)
- [Pagination and Streaming](#pagination-and-streaming)

## Installation

```bash
go get github.com/Pototoooo/meterforge/api/v3/client
```

## Initialization

Create a client with a base URL and an API key. The API key is sent as a
`Bearer` token on every request.

```go
package main

import (
	"log"
	"os"

	"github.com/Pototoooo/meterforge/api/v3/client"
)

func main() {
	om, err := meterforge.New(
		"http://localhost:8888/api/v3",
		meterforge.WithToken(os.Getenv("METERFORGE_API_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = mf
}
```

For region-specific deployments, pass the concrete API base URL for that
region to `New`.

## Usage

Every operation is reachable through a namespaced service on the client and
returns a typed response plus an `error`.

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/Pototoooo/meterforge/api/v3/client"
)

func main() {
	om, err := meterforge.New(
		"http://localhost:8888/api/v3",
		meterforge.WithToken(os.Getenv("METERFORGE_API_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	meter, err := mf.Meters.Create(ctx, meterforge.CreateMeterRequest{
		Name:          "Tokens",
		Key:           "tokens",
		Aggregation:   meterforge.MeterAggregationSum,
		EventType:     "request",
		ValueProperty: meterforge.String("$.tokens"),
	})
	if err != nil {
		log.Fatal(err)
	}

	meters, err := mf.Meters.List(ctx, meterforge.MeterListParams{})
	if err != nil {
		log.Fatal(err)
	}

	_, _ = meter, meters
}
```

Operation arguments follow the generated method signature: path parameters
come first, then a typed request body when present, then typed query params
when present.

## Available Resources and Operations

Operations are grouped by resource and exposed as services on the client.
The full call path, HTTP route, and a short description are listed below.

### Events

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Events.List` | `GET /meterforge/events` | List ingested events. |
| `mf.Events.IngestEvent` | `POST /meterforge/events` | Ingests an event or batch of events following the CloudEvents specification. |
| `mf.Events.IngestEvents` | `POST /meterforge/events` |  |
| `mf.Events.IngestEventsJSON` | `POST /meterforge/events` |  |

### Meters

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Meters.Create` | `POST /meterforge/meters` | Create a meter. |
| `mf.Meters.Get` | `GET /meterforge/meters/{meterId}` | Get a meter by ID. |
| `mf.Meters.List` | `GET /meterforge/meters` | List meters. |
| `mf.Meters.Update` | `PUT /meterforge/meters/{meterId}` | Update a meter. |
| `mf.Meters.Delete` | `DELETE /meterforge/meters/{meterId}` | Delete a meter. |
| `mf.Meters.Query` | `POST /meterforge/meters/{meterId}/query` | Query a meter for usage. Set `Accept: application/json` (the default) to get a structured JSON response. Set `Accept: text/csv` to download the same data as a CSV file suitable for spreadsheets. The CSV columns, in order, are: `from, to, [subject,] [customer_id, customer_key, customer_name,] <dimensions...>, value` The `subject` column is emitted only when `subject` is in the query's `group_by_dimensions`. The three `customer_*` columns are emitted together only when `customer_id` is in the query's `group_by_dimensions`. |
| `mf.Meters.QueryCSV` | `POST /meterforge/meters/{meterId}/query` |  |
| `mf.Meters.QueryCSVStream` | `POST /meterforge/meters/{meterId}/query` | Streaming variant of `QueryCSV` returning an `io.ReadCloser`. |

### Customers

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Customers.Create` | `POST /meterforge/customers` |  |
| `mf.Customers.Get` | `GET /meterforge/customers/{customerId}` |  |
| `mf.Customers.List` | `GET /meterforge/customers` |  |
| `mf.Customers.Upsert` | `PUT /meterforge/customers/{customerId}` |  |
| `mf.Customers.Delete` | `DELETE /meterforge/customers/{customerId}` |  |
| `mf.Customers.Billing.Get` | `GET /meterforge/customers/{customerId}/billing` |  |
| `mf.Customers.Billing.Update` | `PUT /meterforge/customers/{customerId}/billing` |  |
| `mf.Customers.Billing.UpdateAppData` | `PUT /meterforge/customers/{customerId}/billing/app-data` |  |
| `mf.Customers.Billing.CreateStripeCheckoutSession` | `POST /meterforge/customers/{customerId}/billing/stripe/checkout-sessions` | Create a [Stripe Checkout Session](https://docs.stripe.com/payments/checkout) for the customer. Creates a Checkout Session for collecting payment method information from customers. The session operates in "setup" mode, which collects payment details without charging the customer immediately. The collected payment method can be used for future subscription billing. For hosted checkout sessions, redirect customers to the returned URL. For embedded sessions, use the client_secret to initialize Stripe.js in your application. |
| `mf.Customers.Billing.CreateStripePortalSession` | `POST /meterforge/customers/{customerId}/billing/stripe/portal-sessions` | Create Stripe Customer Portal Session. Useful to redirect the customer to the Stripe Customer Portal to manage their payment methods, change their billing address and access their invoice history. Only returns URL if the customer billing profile is linked to a stripe app and customer. |
| `mf.Customers.Credits.Grants.Create` | `POST /meterforge/customers/{customerId}/credits/grants` | Create a new credit grant. A credit grant represents an allocation of prepaid credits to a customer. |
| `mf.Customers.Credits.Grants.Get` | `GET /meterforge/customers/{customerId}/credits/grants/{creditGrantId}` | Get a credit grant. |
| `mf.Customers.Credits.Grants.List` | `GET /meterforge/customers/{customerId}/credits/grants` | List credit grants. |
| `mf.Customers.Credits.Grants.Void` | `POST /meterforge/customers/{customerId}/credits/grants/{creditGrantId}/void` | Void a credit grant, forfeiting the remaining unused balance. Voiding is a forward-looking, irreversible operation. Credits already consumed by usage remain unaffected — only the remaining balance is forfeited. The grant reads as `voided` status afterwards. Payment state is not adjusted when `payment_adjustment` is `none`, so invoice-backed or externally collected payments may still collect the original amount. Only `active` grants can be voided; voiding a pending, expired, or fully consumed grant returns a conflict. Retrying a successful void is an idempotent success. |
| `mf.Customers.Credits.Grants.UpdateExternalSettlement` | `POST /meterforge/customers/{customerId}/credits/grants/{creditGrantId}/settlement/external` | Update the payment settlement status of an externally funded credit grant. Use this endpoint to synchronize the payment state of an external payment with the system so that revenue recognition and credit availability work as expected. |
| `mf.Customers.Credits.Balance.Get` | `GET /meterforge/customers/{customerId}/credits/balance` | Get a credit balance. |
| `mf.Customers.Credits.Adjustments.Create` | `POST /meterforge/customers/{customerId}/credits/adjustments` | A credit adjustment can be used to make manual adjustments to a customer's credit balance. Supported use-cases: - Usage correction |
| `mf.Customers.Credits.Transactions.List` | `GET /meterforge/customers/{customerId}/credits/transactions` | List credit transactions for a customer. Returns an immutable, chronological record of credit movements: funded credits and consumed credits. Transactions are returned in reverse chronological order by default. |
| `mf.Customers.Charges.List` | `GET /meterforge/customers/{customerId}/charges` | List customer charges. Returns the customer's charges that are represented as either flat fee or usage-based charges. |
| `mf.Customers.Charges.Create` | `POST /meterforge/customers/{customerId}/charges` | Create customer charge. |

### Entitlements

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Entitlements.ListCustomerAccess` | `GET /meterforge/customers/{customerId}/entitlement-access` |  |

### Subscriptions

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Subscriptions.Create` | `POST /meterforge/subscriptions` |  |
| `mf.Subscriptions.List` | `GET /meterforge/subscriptions` |  |
| `mf.Subscriptions.Get` | `GET /meterforge/subscriptions/{subscriptionId}` |  |
| `mf.Subscriptions.Cancel` | `POST /meterforge/subscriptions/{subscriptionId}/cancel` | Cancels the subscription. Will result in a scheduling conflict if there are other subscriptions scheduled to start after the cancelation time. |
| `mf.Subscriptions.UnscheduleCancelation` | `POST /meterforge/subscriptions/{subscriptionId}/unschedule-cancelation` | Unschedules the subscription cancelation. |
| `mf.Subscriptions.Change` | `POST /meterforge/subscriptions/{subscriptionId}/change` | Closes a running subscription and starts a new one according to the specification. Can be used for upgrades, downgrades, and plan changes. |
| `mf.Subscriptions.CreateAddon` | `POST /meterforge/subscriptions/{subscriptionId}/addons` | Add add-on to a subscription. |
| `mf.Subscriptions.ListAddons` | `GET /meterforge/subscriptions/{subscriptionId}/addons` | List the add-ons of a subscription. |
| `mf.Subscriptions.GetAddon` | `GET /meterforge/subscriptions/{subscriptionId}/addons/{subscriptionAddonId}` | Get an add-on association for a subscription. |

### Apps

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Apps.List` | `GET /meterforge/apps` | List installed apps. |
| `mf.Apps.Get` | `GET /meterforge/apps/{appId}` | Get an installed app. |
| `mf.Apps.Uninstall` | `DELETE /meterforge/apps/{appId}` | Uninstall an app by ID. |
| `mf.Apps.ListCatalog` | `GET /meterforge/app-catalog` | List available apps. |
| `mf.Apps.GetCatalogItem` | `GET /meterforge/app-catalog/{appType}` | Get an app catalog item by type. |
| `mf.Apps.Install` | `POST /meterforge/app-catalog/install` | Install an app from the catalog. |

### Billing

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Billing.ListProfiles` | `GET /meterforge/profiles` | List billing profiles. |
| `mf.Billing.CreateProfile` | `POST /meterforge/profiles` | Create a new billing profile. Billing profiles contain the settings for billing and controls invoice generation. An organization can have multiple billing profiles defined. A billing profile is linked to a specific app. This association is established during the billing profile's creation and remains immutable. |
| `mf.Billing.GetProfile` | `GET /meterforge/profiles/{id}` | Get a billing profile. |
| `mf.Billing.UpdateProfile` | `PUT /meterforge/profiles/{id}` | Update a billing profile. |
| `mf.Billing.DeleteProfile` | `DELETE /meterforge/profiles/{id}` | Delete a billing profile. Only such billing profiles can be deleted that are: - not the default profile - not pinned to any customer using customer overrides - only have finalized invoices |

### Invoices

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Invoices.List` | `GET /meterforge/billing/invoices` | List billing invoices. Returns a page of invoices. Gathering invoices are never included. Use `filter` to narrow by status, customer, dates, or service period start. Use `sort` to control ordering. |
| `mf.Invoices.Get` | `GET /meterforge/billing/invoices/{invoiceId}` | Get a billing invoice by ID. Returns the full invoice resource including line items, status details, totals, and workflow configuration snapshot. |
| `mf.Invoices.Update` | `PUT /meterforge/billing/invoices/{invoiceId}` | Update a billing invoice. Only the mutable fields of the invoice can be edited: description, labels, supplier, customer, workflow settings, and top-level lines. Top-level lines are matched by `id`; lines without an `id` are created, and existing lines omitted from `lines` are deleted. Detailed (child) lines are always computed and cannot be edited directly. Only invoices in draft status can be updated. |
| `mf.Invoices.Delete` | `DELETE /meterforge/billing/invoices/{invoiceId}` | Delete a billing invoice. Only standard invoices in draft status can be deleted. Deleting an invoice will also delete all associated line items and workflow configuration. |
| `mf.Invoices.Advance` | `POST /meterforge/billing/invoices/{invoiceId}/advance` | Advance a billing invoice. Advances the invoice to the next workflow state. The next state is determined by the invoice's current status and workflow configuration. Only invoices in draft or issued status can be advanced. |
| `mf.Invoices.Approve` | `POST /meterforge/billing/invoices/{invoiceId}/approve` | Approve a billing invoice. This call instantly sends the invoice to the customer using the configured billing profile app. This call is valid in two invoice statuses: - draft: the invoice will be sent to the customer, the invoice state becomes issued - manual_approval_needed: the invoice will be sent to the customer, the invoice state becomes issued |
| `mf.Invoices.Retry` | `POST /meterforge/billing/invoices/{invoiceId}/retry` | Retry sending a billing invoice. Retry advancing the invoice after a failed attempt. The action can be called when the invoice's statusDetails' actions field contain the "retry" action. |
| `mf.Invoices.SnapshotQuantities` | `POST /meterforge/billing/invoices/{invoiceId}/snapshot-quantities` | Snapshot quantities for usage-based line items. This call will snapshot the quantities for all usage based line items in the invoice. This call is only valid in draft.waiting_for_collection status, where the collection period can be skipped using this action. |

### Tax

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Tax.CreateCode` | `POST /meterforge/tax-codes` |  |
| `mf.Tax.GetCode` | `GET /meterforge/tax-codes/{taxCodeId}` |  |
| `mf.Tax.ListCodes` | `GET /meterforge/tax-codes` |  |
| `mf.Tax.UpsertCode` | `PUT /meterforge/tax-codes/{taxCodeId}` |  |
| `mf.Tax.DeleteCode` | `DELETE /meterforge/tax-codes/{taxCodeId}` |  |

### Currencies

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Currencies.List` | `GET /meterforge/currencies` | List currencies supported by the billing system. |
| `mf.Currencies.CreateCustomCurrency` | `POST /meterforge/currencies/custom` | Create a custom currency. This operation allows defining your own custom currency for billing purposes. |
| `mf.Currencies.GetCustomCurrency` | `GET /meterforge/currencies/custom/{currencyId}` | Get a custom currency. |
| `mf.Currencies.ListCostBases` | `GET /meterforge/currencies/custom/{currencyId}/cost-bases` | List cost bases for a currency. For custom currencies, there can be multiple cost bases with different `effective_from` dates. |
| `mf.Currencies.CreateCostBasis` | `POST /meterforge/currencies/custom/{currencyId}/cost-bases` | Create a cost basis for a currency. |

### Features

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Features.List` | `GET /meterforge/features` | List all features. |
| `mf.Features.Create` | `POST /meterforge/features` | Create a feature. |
| `mf.Features.Get` | `GET /meterforge/features/{featureId}` | Get a feature by id. |
| `mf.Features.Update` | `PATCH /meterforge/features/{featureId}` | Update a feature by id. Currently only the unit_cost field can be updated. |
| `mf.Features.Delete` | `DELETE /meterforge/features/{featureId}` | Delete a feature by id. |
| `mf.Features.QueryCost` | `POST /meterforge/features/{featureId}/cost/query` | Query the cost of a feature. |

### LLMCost

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.LLMCost.ListPrices` | `GET /meterforge/llm-cost/prices` | List global LLM cost prices. Returns prices with overrides applied if any. |
| `mf.LLMCost.GetPrice` | `GET /meterforge/llm-cost/prices/{priceId}` | Get a specific LLM cost price by ID. Returns the price with overrides applied if any. |
| `mf.LLMCost.ListOverrides` | `GET /meterforge/llm-cost/overrides` | List per-namespace price overrides. |
| `mf.LLMCost.CreateOverride` | `POST /meterforge/llm-cost/overrides` | Create a per-namespace price override. |
| `mf.LLMCost.DeleteOverride` | `DELETE /meterforge/llm-cost/overrides/{priceId}` | Delete a per-namespace price override. |

### Plans

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Plans.List` | `GET /meterforge/plans` | List all plans. |
| `mf.Plans.Create` | `POST /meterforge/plans` | Create a new plan. |
| `mf.Plans.Update` | `PUT /meterforge/plans/{planId}` | Update a plan by id. |
| `mf.Plans.Get` | `GET /meterforge/plans/{planId}` | Get a plan by id. |
| `mf.Plans.Delete` | `DELETE /meterforge/plans/{planId}` | Delete a plan by id. |
| `mf.Plans.Archive` | `POST /meterforge/plans/{planId}/archive` | Archive a plan version. |
| `mf.Plans.Publish` | `POST /meterforge/plans/{planId}/publish` | Publish a plan version. |

### Addons

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Addons.List` | `GET /meterforge/addons` | List all add-ons. |
| `mf.Addons.Create` | `POST /meterforge/addons` | Create a new add-on. |
| `mf.Addons.Update` | `PUT /meterforge/addons/{addonId}` | Update an add-on by id. |
| `mf.Addons.Get` | `GET /meterforge/addons/{addonId}` | Get add-on by id. |
| `mf.Addons.Delete` | `DELETE /meterforge/addons/{addonId}` | Soft delete add-on by id. |
| `mf.Addons.Archive` | `POST /meterforge/addons/{addonId}/archive` | Archive an add-on version. |
| `mf.Addons.Publish` | `POST /meterforge/addons/{addonId}/publish` | Publish an add-on version. |

### PlanAddons

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.PlanAddons.List` | `GET /meterforge/plans/{planId}/addons` | List add-ons associated with a plan. |
| `mf.PlanAddons.Create` | `POST /meterforge/plans/{planId}/addons` | Add an add-on to a plan. |
| `mf.PlanAddons.Get` | `GET /meterforge/plans/{planId}/addons/{planAddonId}` | Get an add-on association for a plan. |
| `mf.PlanAddons.Update` | `PUT /meterforge/plans/{planId}/addons/{planAddonId}` | Update an add-on association for a plan. |
| `mf.PlanAddons.Delete` | `DELETE /meterforge/plans/{planId}/addons/{planAddonId}` | Remove an add-on from a plan. |

### Defaults

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Defaults.GetOrganizationTaxCodes` | `GET /meterforge/defaults/tax-codes` |  |
| `mf.Defaults.UpdateOrganizationTaxCodes` | `PUT /meterforge/defaults/tax-codes` |  |

### Governance

| Method | HTTP | Description |
| --- | --- | --- |
| `mf.Governance.QueryAccess` | `POST /meterforge/governance/query` | Query feature access for a list of customers. The endpoint resolves each provided identifier to a customer and returns the access status for the requested features, plus optional credit balance availability. _Designed to be called on a fixed refresh interval and the query response is intended to be cached._ |

## Error Handling

A non-2xx response returns an `*APIError` carrying the problem-details
fields (`StatusCode`, `Status`, `Type`, `Title`, `Detail`, `Instance`) from
the response where available. Client-side validation errors such as an empty
path ID are returned before any HTTP request is made.

```go
package main

import (
	"context"
	"errors"
	"log"

	"github.com/Pototoooo/meterforge/api/v3/client"
)

func main() {
	om, err := meterforge.New("http://localhost:8888/api/v3", meterforge.WithToken("mf_..."))
	if err != nil {
		log.Fatal(err)
	}

	_, err = mf.Meters.Get(context.Background(), "unknown")
	if err != nil {
		var apiErr *meterforge.APIError
		if errors.As(err, &apiErr) {
			log.Printf("%d %s %s", apiErr.StatusCode, apiErr.Title, apiErr.Type)
			return
		}
		log.Fatal(err)
	}
}
```

## Pagination and Streaming

Paginated list operations also emit `ListAll` helpers that return
`iter.Seq2[T, error]`. Text responses such as meter CSV export emit a byte
returning method and a `Stream` variant for callers that want an
`io.ReadCloser`.

Cursor-paginated responses report their position as `Next` and `Previous`
on `CursorMetaPage`. Both are opaque cursor tokens: pass them back verbatim
as the `page[after]` / `page[before]` query parameters
(`CursorPageParams.After` / `CursorPageParams.Before`); do not parse or
construct them.

Iterating with `Before` set walks pages backward while the items within
each page stay in forward order, so the resulting stream is not globally
sorted.

```go
for meter, err := range mf.Meters.ListAll(ctx, meterforge.MeterListParams{}) {
	if err != nil {
		log.Fatal(err)
	}
	log.Println(meter.Key)
}

stream, err := mf.Meters.QueryCSVStream(ctx, "meter-id", meterforge.MeterQueryRequest{})
if err != nil {
	log.Fatal(err)
}
defer stream.Close()
```
