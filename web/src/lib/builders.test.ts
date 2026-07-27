import { describe, expect, it } from 'vitest'
import { buildPlan } from './builders'

describe('buildPlan', () => {
  it('builds an interview-demo plan with entitlement and pricing', () => {
    const plan = buildPlan({
      name: 'Starter',
      key: 'starter',
      featureKey: 'api_calls',
      featureName: 'API Calls',
      unitPrice: 0.01,
      includedUsage: 1000,
      baseFee: 19,
    })

    expect(plan.phases[0].rateCards).toHaveLength(2)
    expect(plan.phases[0].rateCards[1]).toMatchObject({
      type: 'usage_based',
      featureKey: 'api_calls',
      entitlementTemplate: { issueAfterReset: 1000 },
    })
  })
})
