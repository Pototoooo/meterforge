import type { PlanCreate } from './meterforge'

export interface PlanFormValues {
  name: string
  key: string
  featureKey: string
  featureName: string
  unitPrice: number
  includedUsage: number
  baseFee: number
}

export function buildPlan(values: PlanFormValues): PlanCreate {
  const rateCards: PlanCreate['phases'][number]['rateCards'] = []

  if (values.baseFee > 0) {
    rateCards.push({
      type: 'flat_fee',
      key: `${values.key}_base_fee`,
      name: `${values.name} base fee`,
      billingCadence: 'P1M',
      price: {
        type: 'flat',
        amount: String(values.baseFee),
        paymentTerm: 'in_advance',
      },
    })
  }

  rateCards.push({
    type: 'usage_based',
    key: values.featureKey,
    name: values.featureName,
    featureKey: values.featureKey,
    billingCadence: 'P1M',
    price: {
      type: 'unit',
      amount: String(values.unitPrice),
    },
    entitlementTemplate: {
      type: 'metered',
      isSoftLimit: false,
      issueAfterReset: values.includedUsage,
      usagePeriod: 'P1M',
    },
  })

  return {
    name: values.name,
    key: values.key,
    currency: 'USD' as const,
    billingCadence: 'P1M',
    settlementMode: 'credit_then_invoice' as const,
    phases: [
      {
        key: 'default',
        name: 'Default',
        duration: null,
        rateCards,
      },
    ],
  }
}
