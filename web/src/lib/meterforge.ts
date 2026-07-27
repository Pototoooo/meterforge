import createClient, { createQuerySerializer } from 'openapi-fetch'
import type { components, paths } from '../../../api/client/javascript/src/client/schemas'

export type Meter = components['schemas']['Meter']
export type Feature = components['schemas']['Feature']
export type Plan = components['schemas']['Plan']
export type PlanCreate = components['schemas']['PlanCreate']
export type Customer = components['schemas']['Customer']
export type CustomerCreate = components['schemas']['CustomerCreate']
export type Subscription = components['schemas']['Subscription']
export type EntitlementV2 = components['schemas']['EntitlementV2']
export type EntitlementValueV2 = components['schemas']['EntitlementValueV2']
export type Event = components['schemas']['Event']
export type IngestedEvent = components['schemas']['IngestedEvent']
export type Invoice = components['schemas']['Invoice']
export type MeterQueryRow = components['schemas']['MeterQueryRow']

type Page<T> = { items: T[]; totalCount: number; page: number; pageSize: number }

const client = createClient<paths>({
  baseUrl: '',
  querySerializer: createQuerySerializer({
    array: { explode: true, style: 'form' },
    object: { explode: true, style: 'deepObject' },
  }),
})

function decodeDates<T>(value: T): T {
  if (typeof value === 'string' && /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/.test(value)) {
    return new Date(value) as T
  }
  if (Array.isArray(value)) return value.map(decodeDates) as T
  if (value && typeof value === 'object') {
    for (const [key, item] of Object.entries(value)) {
      ;(value as Record<string, unknown>)[key] = decodeDates(item)
    }
  }
  return value
}

async function unwrap<T>(request: Promise<any>): Promise<T> {
  const result = await request
  if (result.error || !result.response.ok) {
    const problem = result.error as { detail?: string; title?: string } | undefined
    throw new Error(problem?.detail ?? problem?.title ?? `API request failed (${result.response.status})`)
  }
  return decodeDates(result.data) as T
}

export const meterforge = {
  meters: {
    list: () => unwrap<Meter[]>(client.GET('/api/v1/meters')),
    create: (body: components['schemas']['MeterCreate']) => unwrap<Meter>(client.POST('/api/v1/meters', { body })),
    query: (meterIdOrSlug: string, query?: { windowSize?: 'HOUR' | 'DAY' | 'MONTH'; subject?: string[] }) =>
      unwrap<components['schemas']['MeterQueryResult']>(client.GET('/api/v1/meters/{meterIdOrSlug}/query', { params: { path: { meterIdOrSlug }, query }, headers: { Accept: 'application/json' } })),
  },
  features: {
    list: () => unwrap<Feature[]>(client.GET('/api/v1/features')),
    create: (body: components['schemas']['FeatureCreateInputs']) => unwrap<Feature>(client.POST('/api/v1/features', { body })),
  },
  plans: {
    list: (query?: { pageSize?: number; status?: components['schemas']['PlanStatus'][] }) => unwrap<Page<Plan>>(client.GET('/api/v1/plans', { params: { query } })),
    create: (body: PlanCreate) => unwrap<Plan>(client.POST('/api/v1/plans', { body })),
    publish: (planId: string) => unwrap<Plan>(client.POST('/api/v1/plans/{planId}/publish', { params: { path: { planId } } })),
  },
  customers: {
    list: (query?: { pageSize?: number; expand?: components['schemas']['CustomerExpand'][] }) => unwrap<Page<Customer>>(client.GET('/api/v1/customers', { params: { query } })),
    get: (customerIdOrKey: string) => unwrap<Customer>(client.GET('/api/v1/customers/{customerIdOrKey}', { params: { path: { customerIdOrKey } } })),
    create: (body: CustomerCreate) => unwrap<Customer>(client.POST('/api/v1/customers', { body })),
    listSubscriptions: (customerIdOrKey: string, query?: { pageSize?: number }) => unwrap<Page<Subscription>>(client.GET('/api/v1/customers/{customerIdOrKey}/subscriptions', { params: { path: { customerIdOrKey }, query } })),
    entitlements: {
      list: (customerIdOrKey: string, options?: { query?: { pageSize?: number } }) => unwrap<Page<EntitlementV2>>(client.GET('/api/v2/customers/{customerIdOrKey}/entitlements', { params: { path: { customerIdOrKey }, query: options?.query } })),
      value: (customerIdOrKey: string, entitlementIdOrFeatureKey: string) => unwrap<EntitlementValueV2>(client.GET('/api/v2/customers/{customerIdOrKey}/entitlements/{entitlementIdOrFeatureKey}/value', { params: { path: { customerIdOrKey, entitlementIdOrFeatureKey } } })),
    },
  },
  subscriptions: {
    create: (body: components['schemas']['SubscriptionCreate']) => unwrap<Subscription>(client.POST('/api/v1/subscriptions', { body })),
  },
  entitlements: {
    list: (options?: { query?: { pageSize?: number } }) => unwrap<Page<EntitlementV2>>(client.GET('/api/v2/entitlements', { params: { query: options?.query } })),
  },
  events: {
    listV2: (query?: { limit?: number; filter?: string }) => unwrap<components['schemas']['IngestedEventCursorPaginatedResponse']>(client.GET('/api/v2/events', { params: { query } })),
    ingest: (event: Event) => {
      const body: Event[] = [{ ...event, id: event.id ?? crypto.randomUUID(), source: event.source ?? 'meterforge-console', specversion: event.specversion ?? '1.0', time: event.time ?? new Date() }]
      return unwrap<void>(client.POST('/api/v1/events', { body, headers: { 'Content-Type': 'application/cloudevents-batch+json' } }))
    },
  },
  billing: {
    invoices: {
      list: (query?: { pageSize?: number }) => unwrap<Page<Invoice>>(client.GET('/api/v1/billing/invoices', { params: { query } })),
    },
    profiles: {
      list: (query?: { pageSize?: number }) => unwrap<Page<components['schemas']['BillingProfile']>>(client.GET('/api/v1/billing/profiles', { params: { query } })),
    },
  },
  apps: {
    list: () => unwrap<components['schemas']['AppPaginatedResponse']>(client.GET('/api/v1/apps')),
    marketplace: {
      list: (query?: { pageSize?: number }) => unwrap<components['schemas']['MarketplaceListingPaginatedResponse']>(client.GET('/api/v1/marketplace/listings', { params: { query } })),
    },
  },
}

export function errorMessage(error: unknown) {
  if (error instanceof Error) return error.message.replace(/^Failed to fetch$/, '无法连接 MeterForge API')
  return '发生未知错误'
}
