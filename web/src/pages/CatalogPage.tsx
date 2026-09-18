import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  Boxes,
  Box,
  CheckCircle2,
  CircleDollarSign,
  DatabaseZap,
  Layers3,
  Plus,
  Rocket,
  Search,
  Sparkles,
} from 'lucide-react'
import { useSearchParams } from 'react-router-dom'
import { toast } from 'sonner'
import { meterforge, errorMessage, type Feature } from '../lib/meterforge'
import { buildPlan } from '../lib/builders'
import { formatDate, formatMoney, formatNumber, slugify } from '../lib/format'
import {
  Badge,
  Button,
  EmptyState,
  ErrorState,
  Field,
  LoadingState,
  Modal,
  PageHeader,
  Tabs,
} from '../components/ui'

export function CatalogPage() {
  const [params, setParams] = useSearchParams()
  const tab = params.get('tab') ?? 'features'
  const [featureOpen, setFeatureOpen] = useState(false)
  const [planOpen, setPlanOpen] = useState(false)
  const features = useQuery({ queryKey: ['features'], queryFn: () => meterforge.features.list() })
  const plans = useQuery({ queryKey: ['plans'], queryFn: () => meterforge.plans.list({ pageSize: 100 }) })

  return (
    <div className="page">
      <PageHeader
        eyebrow="Metering & Billing /"
        title="Product Catalog"
        description="把 Meter 封装成可售卖 Feature，并通过 Plan 定义价格、额度与账单周期。"
        actions={
          <Button onClick={() => tab === 'features' ? setFeatureOpen(true) : setPlanOpen(true)}>
            <Plus size={16} /> Create {tab === 'features' ? 'Feature' : 'Plan'}
          </Button>
        }
      />
      <Tabs
        value={tab}
        onChange={(value) => setParams(value === 'features' ? {} : { tab: value })}
        items={[
          { value: 'features', label: 'Features', count: features.data?.length },
          { value: 'plans', label: 'Plans', count: plans.data?.totalCount },
        ]}
      />
      {tab === 'features' ? <FeaturesView query={features} onCreate={() => setFeatureOpen(true)} /> : <PlansView query={plans} onCreate={() => setPlanOpen(true)} />}
      <CreateFeatureModal open={featureOpen} onClose={() => setFeatureOpen(false)} />
      <CreatePlanModal open={planOpen} onClose={() => setPlanOpen(false)} features={features.data ?? []} />
    </div>
  )
}

function FeaturesView({ query, onCreate }: { query: ReturnType<typeof useQuery>; onCreate: () => void }) {
  const [search, setSearch] = useState('')
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const data = (query.data as Feature[] | undefined) ?? []
  const items = data.filter((item) => `${item.name} ${item.key} ${item.meterSlug}`.toLowerCase().includes(search.toLowerCase()))
  return (
    <section className="panel table-panel">
      <div className="table-toolbar"><div className="search-input"><Search size={16} /><input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search features" /></div><span>{items.length} product building blocks</span></div>
      {!items.length ? <EmptyState icon={<Box />} title="Create your first feature" description="Features connect metered usage to product access and pricing." action={<Button onClick={onCreate}><Plus size={16} /> Create Feature</Button>} /> : (
        <div className="data-table-wrap">
          <table className="data-table">
            <thead><tr><th>Feature</th><th>Key</th><th>Meter</th><th>Cost configuration</th><th>Created</th></tr></thead>
            <tbody>{items.map((feature) => (
              <tr key={feature.id}>
                <td><div className="table-resource"><span className="resource-icon"><Sparkles size={16} /></span><strong>{feature.name}</strong></div></td>
                <td><code>{feature.key}</code></td>
                <td>{feature.meterSlug ? <span className="meter-link"><DatabaseZap size={14} /> {feature.meterSlug}</span> : <span className="muted">Unmetered</span>}</td>
                <td>{feature.unitCost ? <Badge tone="blue">{feature.unitCost.type}</Badge> : <span className="muted">Not configured</span>}</td>
                <td>{formatDate(feature.createdAt)}</td>
              </tr>
            ))}</tbody>
          </table>
        </div>
      )}
    </section>
  )
}

function PlansView({ query, onCreate }: { query: ReturnType<typeof useQuery>; onCreate: () => void }) {
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const items = (query.data as Awaited<ReturnType<typeof meterforge.plans.list>> | undefined)?.items ?? []
  return (
    <section className="catalog-grid">
      {!items.length ? <div className="panel"><EmptyState icon={<Boxes />} title="Create your first plan" description="Plans combine pricing and entitlements into a subscribable product." action={<Button onClick={onCreate}><Plus size={16} /> Create Plan</Button>} /></div> : items.map((plan) => {
        const rateCards = plan.phases.flatMap((phase) => phase.rateCards)
        const flat = rateCards.find((card) => card.price?.type === 'flat')
        const usage = rateCards.find((card) => card.type === 'usage_based')
        return (
          <article className="plan-card" key={plan.id}>
            <div className="plan-card-header">
              <span className="plan-icon"><Layers3 size={20} /></span>
              <Badge>{plan.status}</Badge>
            </div>
            <div><h3>{plan.name}</h3><code>{plan.key} · v{plan.version}</code></div>
            <div className="plan-price">
              <strong>{flat?.price && 'amount' in flat.price ? formatMoney(flat.price.amount, plan.currency) : 'Usage based'}</strong>
              <span>/ month</span>
            </div>
            <div className="plan-summary">
              <div><CircleDollarSign size={15} /><span>{usage?.price?.type === 'unit' ? `${formatMoney(usage.price.amount, plan.currency)} per unit` : 'Configured pricing'}</span></div>
              <div><CheckCircle2 size={15} /><span>{rateCards.length} rate cards across {plan.phases.length} phase</span></div>
              <div><Rocket size={15} /><span>{plan.settlementMode === 'credit_then_invoice' ? 'Credits, then invoice' : plan.settlementMode}</span></div>
            </div>
            <div className="plan-card-footer"><span>{plan.billingCadence} billing</span><span>Updated {formatDate(plan.updatedAt)}</span></div>
          </article>
        )
      })}
    </section>
  )
}

function CreateFeatureModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const client = useQueryClient()
  const meters = useQuery({ queryKey: ['meters'], queryFn: () => meterforge.meters.list(), enabled: open })
  const [name, setName] = useState('API Requests')
  const [key, setKey] = useState('api_requests')
  const [meterSlug, setMeterSlug] = useState('api_requests_total')
  const create = useMutation({
    mutationFn: () => meterforge.features.create({ name, key, meterSlug: meterSlug || undefined }),
    onSuccess: () => { client.invalidateQueries({ queryKey: ['features'] }); toast.success('Feature 已创建'); onClose() },
    onError: (error) => toast.error(errorMessage(error)),
  })
  return (
    <Modal open={open} onClose={onClose} title="New Feature" description="A product building block backed by metered usage.">
      <form className="form-stack" onSubmit={(event) => { event.preventDefault(); create.mutate() }}>
        <div className="form-section-title"><span>1</span><div><strong>General Information</strong><small>Basic information about the feature.</small></div></div>
        <div className="form-grid two">
          <Field label="Name"><input required value={name} onChange={(event) => { setName(event.target.value); if (!key) setKey(slugify(event.target.value)) }} /></Field>
          <Field label="Key"><input required value={key} onChange={(event) => setKey(slugify(event.target.value))} /></Field>
        </div>
        <div className="form-section-title"><span>2</span><div><strong>Event Processing</strong><small>Associate this feature with a meter.</small></div></div>
        <Field label="Meter"><select value={meterSlug} onChange={(event) => setMeterSlug(event.target.value)}><option value="">Unmetered feature</option>{meters.data?.map((meter) => <option value={meter.slug} key={meter.id}>{meter.name} · {meter.aggregation}</option>)}</select></Field>
        <div className="modal-actions"><Button type="button" variant="secondary" onClick={onClose}>Cancel</Button><Button type="submit" disabled={create.isPending}>{create.isPending ? 'Creating…' : 'Create Feature'}</Button></div>
      </form>
    </Modal>
  )
}

function CreatePlanModal({ open, onClose, features }: { open: boolean; onClose: () => void; features: Feature[] }) {
  const client = useQueryClient()
  const [name, setName] = useState('Starter')
  const [key, setKey] = useState('starter')
  const [featureKey, setFeatureKey] = useState(features[0]?.key ?? '')
  const [baseFee, setBaseFee] = useState(19)
  const [unitPrice, setUnitPrice] = useState(0.01)
  const [includedUsage, setIncludedUsage] = useState(1000)
  const feature = features.find((item) => item.key === featureKey) ?? features[0]
  const create = useMutation({
    mutationFn: async () => {
      if (!feature) throw new Error('请先创建一个 Feature')
      const plan = await meterforge.plans.create(buildPlan({ name, key, featureKey: feature.key, featureName: feature.name, baseFee, unitPrice, includedUsage }))
      if (!plan) throw new Error('Plan 创建后未返回资源')
      return meterforge.plans.publish(plan.id)
    },
    onSuccess: () => { client.invalidateQueries({ queryKey: ['plans'] }); toast.success('Plan 已创建并发布', { description: '现在可以为 Customer 创建订阅。' }); onClose() },
    onError: (error) => toast.error(errorMessage(error)),
  })
  return (
    <Modal open={open} onClose={onClose} title="New Plan" description="Create a monthly plan with pricing and a metered entitlement.">
      <form className="form-stack" onSubmit={(event) => { event.preventDefault(); create.mutate() }}>
        <div className="form-section-title"><span>1</span><div><strong>General Information</strong><small>Basic details about the plan.</small></div></div>
        <div className="form-grid two">
          <Field label="Name"><input required value={name} onChange={(event) => { setName(event.target.value); if (!key) setKey(slugify(event.target.value)) }} /></Field>
          <Field label="Key"><input required value={key} onChange={(event) => setKey(slugify(event.target.value))} /></Field>
        </div>
        <div className="form-section-title"><span>2</span><div><strong>Pricing & Entitlement</strong><small>Monthly USD billing, credit then invoice.</small></div></div>
        <Field label="Feature"><select required value={featureKey} onChange={(event) => setFeatureKey(event.target.value)}><option value="">Select a feature</option>{features.filter((item) => item.meterSlug).map((item) => <option key={item.id} value={item.key}>{item.name} · {item.key}</option>)}</select></Field>
        <div className="form-grid three">
          <Field label="Base fee ($)"><input type="number" min="0" step="0.01" value={baseFee} onChange={(event) => setBaseFee(Number(event.target.value))} /></Field>
          <Field label="Price / unit ($)"><input type="number" min="0" step="0.0001" value={unitPrice} onChange={(event) => setUnitPrice(Number(event.target.value))} /></Field>
          <Field label="Included units"><input type="number" min="0" value={includedUsage} onChange={(event) => setIncludedUsage(Number(event.target.value))} /></Field>
        </div>
        <div className="pricing-preview"><div><small>Subscription price</small><strong>{formatMoney(baseFee)}</strong></div><Plus size={16} /><div><small>Usage overage</small><strong>{formatMoney(unitPrice)} / unit</strong></div><div className="pricing-credit"><small>Included</small><strong>{formatNumber(includedUsage)} units</strong></div></div>
        <div className="modal-actions"><Button type="button" variant="secondary" onClick={onClose}>Cancel</Button><Button type="submit" disabled={!feature || create.isPending}><Rocket size={16} /> {create.isPending ? 'Publishing…' : 'Create & Publish'}</Button></div>
      </form>
    </Modal>
  )
}
