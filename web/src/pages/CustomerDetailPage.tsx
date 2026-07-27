import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  ArrowLeft,
  CalendarClock,
  CheckCircle2,
  CircleDollarSign,
  FileText,
  Gauge,
  KeyRound,
  Mail,
  Plus,
  Radio,
  Rocket,
  ShieldCheck,
  UserRound,
} from 'lucide-react'
import { Link, useParams } from 'react-router-dom'
import { toast } from 'sonner'
import { meterforge, errorMessage } from '../lib/meterforge'
import { formatDate, formatMoney, formatNumber, shortId } from '../lib/format'
import { Badge, Button, EmptyState, ErrorState, Field, LoadingState, Modal, Tabs } from '../components/ui'

export function CustomerDetailPage() {
  const { customerId = '' } = useParams()
  const [tab, setTab] = useState('overview')
  const [subscribeOpen, setSubscribeOpen] = useState(false)
  const customer = useQuery({ queryKey: ['customer', customerId], queryFn: () => meterforge.customers.get(customerId), enabled: Boolean(customerId) })
  const subscriptions = useQuery({ queryKey: ['customer-subscriptions', customerId], queryFn: () => meterforge.customers.listSubscriptions(customerId, { pageSize: 100 }), enabled: Boolean(customerId) })
  const entitlementValues = useQuery({
    queryKey: ['customer-entitlement-values', customerId],
    queryFn: async () => {
      const result = await meterforge.customers.entitlements.list(customerId, { query: { pageSize: 100 } })
      const values = await Promise.all(result.items.map(async (entitlement) => ({
        entitlement,
        value: await meterforge.customers.entitlements.value(customerId, entitlement.featureKey),
      })))
      return values
    },
    enabled: Boolean(customerId),
  })
  const invoices = useQuery({
    queryKey: ['customer-invoices', customerId],
    queryFn: async () => {
      const result = await meterforge.billing.invoices.list({ pageSize: 100 })
      return result.items.filter((invoice) => invoice.customer.id === customerId)
    },
    enabled: Boolean(customerId),
  })

  if (customer.isLoading) return <div className="page"><LoadingState /></div>
  if (customer.error || !customer.data) return <div className="page"><ErrorState message={errorMessage(customer.error ?? new Error('Customer not found'))} /></div>
  const data = customer.data
  return (
    <div className="page customer-detail-page">
      <Link className="back-link" to="/billing"><ArrowLeft size={15} /> Back to customers</Link>
      <header className="customer-detail-header">
        <div className="customer-avatar large">{data.name.slice(0, 2).toUpperCase()}</div>
        <div><div className="eyebrow">Customer</div><h1>{data.name}</h1><div className="customer-header-meta"><code>{data.key}</code><span>{shortId(data.id, 18)}</span>{data.currentSubscriptionId ? <Badge tone="green">Subscribed</Badge> : <Badge tone="amber">No subscription</Badge>}</div></div>
        <div className="page-actions"><Link className="button button-secondary" to={`/metering?tab=events&subject=${data.usageAttribution?.subjectKeys?.[0] ?? data.key}`}><Radio size={16} /> Send usage</Link><Button onClick={() => setSubscribeOpen(true)}><Plus size={16} /> Create Subscription</Button></div>
      </header>
      <Tabs value={tab} onChange={setTab} items={[
        { value: 'overview', label: 'Overview' },
        { value: 'subscriptions', label: 'Subscriptions', count: subscriptions.data?.items.length },
        { value: 'entitlements', label: 'Entitlements', count: entitlementValues.data?.length },
        { value: 'invoicing', label: 'Invoicing', count: invoices.data?.length },
      ]} />
      {tab === 'overview' && <CustomerOverview customer={data} values={entitlementValues.data ?? []} invoiceCount={invoices.data?.length ?? 0} />}
      {tab === 'subscriptions' && <SubscriptionsPanel query={subscriptions} onCreate={() => setSubscribeOpen(true)} />}
      {tab === 'entitlements' && <EntitlementsPanel query={entitlementValues} />}
      {tab === 'invoicing' && <CustomerInvoicesPanel query={invoices} />}
      <SubscribeModal open={subscribeOpen} onClose={() => setSubscribeOpen(false)} customerId={customerId} />
    </div>
  )
}

function CustomerOverview({ customer, values, invoiceCount }: { customer: NonNullable<Awaited<ReturnType<typeof meterforge.customers.get>>>; values: Array<{ entitlement: any; value: any }>; invoiceCount: number }) {
  const used = values.reduce((sum, item) => sum + Number(item.value.usage ?? 0), 0)
  const allowance = values.reduce((sum, item) => sum + Number(item.value.totalAvailableGrantAmount ?? item.value.balance ?? 0), 0)
  return (
    <div className="customer-overview-grid">
      <section className="panel detail-card">
        <div className="panel-heading"><div><span className="panel-kicker">Account</span><h2>Customer details</h2></div><UserRound size={20} /></div>
        <dl className="detail-list"><div><dt><Mail size={15} /> Primary email</dt><dd>{customer.primaryEmail ?? 'Not set'}</dd></div><div><dt><CircleDollarSign size={15} /> Currency</dt><dd>{customer.currency ?? 'USD'}</dd></div><div><dt><CalendarClock size={15} /> Created</dt><dd>{formatDate(customer.createdAt)}</dd></div><div><dt><KeyRound size={15} /> Customer ID</dt><dd><code>{shortId(customer.id, 20)}</code></dd></div></dl>
      </section>
      <section className="panel detail-card">
        <div className="panel-heading"><div><span className="panel-kicker">Usage attribution</span><h2>Event subjects</h2></div><Radio size={20} /></div>
        <div className="subject-list">{customer.usageAttribution?.subjectKeys?.map((subject) => <div key={subject}><span className="event-dot" /><code>{subject}</code><Badge tone="green">active</Badge></div>) ?? <span className="muted">No subjects</span>}</div>
      </section>
      <section className="panel usage-summary-card">
        <div><span className="panel-kicker">Current period</span><h2>Entitlement usage</h2></div>
        <div className="usage-ring" style={{ '--progress': `${allowance ? Math.min(100, used / allowance * 100) : 0}%` } as React.CSSProperties}><div><strong>{formatNumber(used)}</strong><span>used</span></div></div>
        <div className="usage-summary-stats"><div><small>Granted</small><strong>{formatNumber(allowance)}</strong></div><div><small>Features</small><strong>{values.length}</strong></div><div><small>Invoices</small><strong>{invoiceCount}</strong></div></div>
      </section>
    </div>
  )
}

function SubscriptionsPanel({ query, onCreate }: { query: ReturnType<typeof useQuery>; onCreate: () => void }) {
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const data = query.data as Awaited<ReturnType<typeof meterforge.customers.listSubscriptions>>
  if (!data.items.length) return <section className="panel"><EmptyState icon={<Rocket />} title="No active subscription" description="Choose a published plan to provision pricing and entitlements." action={<Button onClick={onCreate}><Plus size={16} /> Create Subscription</Button>} /></section>
  return <section className="subscription-grid">{data.items.map((subscription) => <article className="panel subscription-card" key={subscription.id}><div className="subscription-card-head"><span className="plan-icon"><Rocket size={19} /></span><Badge>{subscription.status}</Badge></div><h2>{subscription.name}</h2><p>{subscription.description || 'Created from a published product plan.'}</p><div className="subscription-facts"><div><span>Plan</span><strong>{subscription.plan?.key ?? 'Custom'}</strong></div><div><span>Billing cadence</span><strong>{subscription.billingCadence}</strong></div><div><span>Currency</span><strong>{subscription.currency}</strong></div><div><span>Started</span><strong>{formatDate(subscription.activeFrom)}</strong></div></div><code>{subscription.id}</code></article>)}</section>
}

function EntitlementsPanel({ query }: { query: ReturnType<typeof useQuery> }) {
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const items = query.data as Array<{ entitlement: any; value: any }>
  if (!items.length) return <section className="panel"><EmptyState icon={<ShieldCheck />} title="No entitlements" description="Entitlements are provisioned from a plan's rate cards." /></section>
  return <section className="entitlement-card-grid">{items.map(({ entitlement, value }) => {
    const grant = Number(value.totalAvailableGrantAmount ?? 0)
    const usage = Number(value.usage ?? 0)
    const ratio = grant ? Math.min(100, usage / grant * 100) : 0
    return <article className="panel entitlement-card" key={entitlement.id}><div className="entitlement-head"><span className="resource-icon"><Gauge size={17} /></span><Badge tone={value.hasAccess ? 'green' : 'red'}>{value.hasAccess ? 'access granted' : 'access denied'}</Badge></div><h3>{entitlement.featureKey}</h3><code>{entitlement.type}</code><div className="balance-row"><div><small>Balance</small><strong>{formatNumber(value.balance)}</strong></div><div><small>Usage</small><strong>{formatNumber(value.usage)}</strong></div><div><small>Overage</small><strong>{formatNumber(value.overage)}</strong></div></div><div className="entitlement-progress"><span style={{ width: `${ratio}%` }} /></div><div className="entitlement-period"><CalendarClock size={14} /> {formatDate(entitlement.currentUsagePeriod?.from)} – {formatDate(entitlement.currentUsagePeriod?.to)}</div></article>
  })}</section>
}

function CustomerInvoicesPanel({ query }: { query: ReturnType<typeof useQuery> }) {
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const items = query.data as Awaited<ReturnType<typeof meterforge.billing.invoices.list>>['items']
  if (!items.length) return <section className="panel"><EmptyState icon={<FileText />} title="No invoices yet" description="A gathering invoice is created automatically when an invoiceable subscription starts." /></section>
  return <section className="panel table-panel"><div className="data-table-wrap"><table className="data-table"><thead><tr><th>Invoice</th><th>Status</th><th>Period</th><th>Charges</th><th>Total</th></tr></thead><tbody>{items.map((invoice) => <tr key={invoice.id}><td><strong>{invoice.number ?? shortId(invoice.id)}</strong></td><td><Badge>{invoice.status}</Badge></td><td>{formatDate(invoice.period?.from)} – {formatDate(invoice.period?.to)}</td><td>{formatMoney(invoice.totals.chargesTotal, invoice.currency)}</td><td><strong>{formatMoney(invoice.totals.total, invoice.currency)}</strong></td></tr>)}</tbody></table></div></section>
}

function SubscribeModal({ open, onClose, customerId }: { open: boolean; onClose: () => void; customerId: string }) {
  const client = useQueryClient()
  const plans = useQuery({ queryKey: ['plans'], queryFn: () => meterforge.plans.list({ pageSize: 100, status: ['active'] }), enabled: open })
  const activePlans = plans.data?.items.filter((plan) => plan.status === 'active') ?? []
  const [planKey, setPlanKey] = useState('')
  const selected = activePlans.find((plan) => plan.key === planKey) ?? activePlans[0]
  const subscribe = useMutation({
    mutationFn: () => meterforge.subscriptions.create({ customerId, plan: { key: selected.key } }),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ['customer', customerId] })
      client.invalidateQueries({ queryKey: ['customer-subscriptions', customerId] })
      client.invalidateQueries({ queryKey: ['customer-entitlement-values', customerId] })
      client.invalidateQueries({ queryKey: ['customers'] })
      client.invalidateQueries({ queryKey: ['invoices'] })
      toast.success('Subscription 已创建', { description: '价格、额度和发票工作流已开始生效。' })
      onClose()
    },
    onError: (error) => toast.error(errorMessage(error)),
  })
  return <Modal open={open} onClose={onClose} title="Create Subscription" description="Assign a published plan and provision its billing artifacts."><form className="form-stack" onSubmit={(event) => { event.preventDefault(); subscribe.mutate() }}><Field label="Published Plan"><select value={planKey || selected?.key || ''} onChange={(event) => setPlanKey(event.target.value)}><option value="">Select a plan</option>{activePlans.map((plan) => <option key={plan.id} value={plan.key}>{plan.name} · {plan.currency} · {plan.billingCadence}</option>)}</select></Field>{selected && <div className="subscription-preview"><div><span className="plan-icon"><Rocket size={18} /></span><div><strong>{selected.name}</strong><code>{selected.key} · v{selected.version}</code></div></div><div><CheckCircle2 size={15} /> {selected.phases.flatMap((phase) => phase.rateCards).length} rate cards</div><div><CircleDollarSign size={15} /> {selected.settlementMode}</div></div>}<div className="info-banner"><ShieldCheck size={17} /><span>创建订阅会自动生成 Entitlement，并为可计费项目启动 gathering invoice。</span></div><div className="modal-actions"><Button type="button" variant="secondary" onClick={onClose}>Cancel</Button><Button type="submit" disabled={!selected || subscribe.isPending}><Rocket size={16} /> {subscribe.isPending ? 'Creating…' : 'Start Subscription'}</Button></div></form></Modal>
}
