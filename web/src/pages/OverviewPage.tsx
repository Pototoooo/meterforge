import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  ArrowRight,
  Check,
  CircleDollarSign,
  Gauge,
  Radio,
  Sparkles,
  Users,
} from 'lucide-react'
import { Link } from 'react-router-dom'
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { meterforge, errorMessage } from '../lib/meterforge'
import { formatDate, formatNumber } from '../lib/format'
import { Badge, ErrorState, LoadingState, PageHeader } from '../components/ui'

export function OverviewPage() {
  const meters = useQuery({ queryKey: ['meters'], queryFn: () => meterforge.meters.list() })
  const features = useQuery({ queryKey: ['features'], queryFn: () => meterforge.features.list() })
  const plans = useQuery({ queryKey: ['plans'], queryFn: () => meterforge.plans.list({ pageSize: 100 }) })
  const customers = useQuery({
    queryKey: ['customers'],
    queryFn: () => meterforge.customers.list({ pageSize: 100, expand: ['subscriptions'] }),
  })
  const invoices = useQuery({
    queryKey: ['invoices'],
    queryFn: () => meterforge.billing.invoices.list({ pageSize: 100 }),
  })
  const events = useQuery({
    queryKey: ['events', 'overview'],
    queryFn: () => meterforge.events.listV2({
      limit: 8,
      filter: JSON.stringify({ source: { $eq: 'meterforge-console' } }),
    }),
  })

  const preferredMeter = meters.data?.find((meter) => meter.slug === 'api_requests_total') ?? meters.data?.[0]
  const usage = useQuery({
    queryKey: ['meter-query', preferredMeter?.slug, 'overview'],
    queryFn: () => meterforge.meters.query(preferredMeter!.slug, { windowSize: 'DAY' }),
    enabled: Boolean(preferredMeter),
  })

  const chartData = useMemo(() => {
    const grouped = new Map<string, number>()
    for (const row of usage.data?.data ?? []) {
      const key = formatDate(row.windowStart)
      grouped.set(key, (grouped.get(key) ?? 0) + row.value)
    }
    return [...grouped.entries()].map(([date, value]) => ({ date, value })).slice(-14)
  }, [usage.data])

  const activePlans = plans.data?.items.filter((plan) => plan.status === 'active').length ?? 0
  const allInvoices = invoices.data?.items ?? []
  const validEvents = events.data?.items.filter((item) => !item.validationError) ?? []
  const invoiceTotal = allInvoices.reduce((sum, invoice) => sum + Number(invoice.totals.total), 0)
  const completed = [
    Boolean(meters.data?.length),
    Boolean(features.data?.length),
    Boolean(activePlans),
    Boolean(customers.data?.items.length),
    Boolean(validEvents.length),
  ].filter(Boolean).length

  const firstError = [meters, features, plans, customers, invoices].find((query) => query.error)?.error

  return (
    <div className="page overview-page">
      <PageHeader
        eyebrow="Metering & Billing"
        title="Overview"
        description="从用量事件到订阅和发票，实时查看本地计费闭环。"
        actions={
          <Link className="button button-primary" to="/metering?tab=events">
            <Radio size={16} /> Send event
          </Link>
        }
      />

      {firstError && <ErrorState message={errorMessage(firstError)} />}

      <section className="metric-grid">
        <MetricCard
          icon={<Gauge size={19} />}
          label="Active meters"
          value={meters.data?.length}
          caption="实时聚合规则"
          tone="violet"
        />
        <MetricCard
          icon={<Users size={19} />}
          label="Customers"
          value={customers.data?.totalCount}
          caption={`${customers.data?.items.filter((item) => item.currentSubscriptionId).length ?? 0} 个已有订阅`}
          tone="blue"
        />
        <MetricCard
          icon={<Sparkles size={19} />}
          label="Published plans"
          value={activePlans}
          caption={`${features.data?.length ?? 0} 个产品功能`}
          tone="cyan"
        />
        <MetricCard
          icon={<CircleDollarSign size={19} />}
          label="Invoice total"
          value={`$${formatNumber(invoiceTotal)}`}
          caption={`${allInvoices.length} 张 Sandbox 发票`}
          tone="green"
        />
      </section>

      <section className="overview-grid">
        <article className="panel usage-panel">
          <div className="panel-heading">
            <div>
              <span className="panel-kicker">Usage trend</span>
              <h2>{preferredMeter?.name ?? 'Meter usage'}</h2>
            </div>
            <Link to="/metering?tab=query">Query details <ArrowRight size={14} /></Link>
          </div>
          {usage.isLoading ? (
            <LoadingState />
          ) : chartData.length ? (
            <div className="chart-wrap">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={chartData} margin={{ top: 12, right: 8, left: -18, bottom: 0 }}>
                  <defs>
                    <linearGradient id="usageGradient" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stopColor="#4f46e5" stopOpacity={0.32} />
                      <stop offset="100%" stopColor="#4f46e5" stopOpacity={0.02} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e9ebf3" />
                  <XAxis dataKey="date" tickLine={false} axisLine={false} tick={{ fill: '#7a8092', fontSize: 11 }} />
                  <YAxis tickLine={false} axisLine={false} tick={{ fill: '#7a8092', fontSize: 11 }} />
                  <Tooltip contentStyle={{ border: '1px solid #e4e6ef', borderRadius: 10, boxShadow: '0 8px 24px rgba(20,24,40,.08)' }} />
                  <Area type="monotone" dataKey="value" stroke="#4f46e5" strokeWidth={2.5} fill="url(#usageGradient)" />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <div className="chart-empty">发送事件后，这里将展示每日用量趋势。</div>
          )}
        </article>

        <article className="panel setup-panel">
          <div className="panel-heading">
            <div>
              <span className="panel-kicker">Demo readiness</span>
              <h2>核心闭环</h2>
            </div>
            <strong>{completed} / 5</strong>
          </div>
          <div className="progress-track"><span style={{ width: `${completed * 20}%` }} /></div>
          <div className="checklist">
            <SetupItem done={Boolean(meters.data?.length)} label="创建 Meter" to="/metering" />
            <SetupItem done={Boolean(features.data?.length)} label="创建 Feature" to="/catalog" />
            <SetupItem done={Boolean(activePlans)} label="发布 Plan" to="/catalog?tab=plans" />
            <SetupItem done={Boolean(customers.data?.items.length)} label="创建 Customer 与订阅" to="/billing" />
            <SetupItem done={Boolean(validEvents.length)} label="发送并查询事件" to="/metering?tab=events" />
          </div>
        </article>
      </section>

      <section className="panel recent-panel">
        <div className="panel-heading">
          <div>
            <span className="panel-kicker">Console ingestion</span>
            <h2>Recent events</h2>
          </div>
          <Badge tone={validEvents.length ? 'green' : 'red'}>{validEvents.length ? 'pipeline healthy' : 'validation required'}</Badge>
        </div>
        <div className="event-list compact-events">
          {(events.data?.items ?? []).slice(0, 6).map((item) => (
            <div className="event-row" key={item.event.id}>
              <span className={item.validationError ? 'event-dot error' : 'event-dot'} />
              <div>
                <strong>{item.event.type}</strong>
                <small>{item.event.subject}</small>
              </div>
              <code>{item.event.id}</code>
              <span>{formatDate(item.ingestedAt)}</span>
              {item.validationError ? <Badge tone="red">invalid</Badge> : <Badge tone="green">stored</Badge>}
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}

function MetricCard({ icon, label, value, caption, tone }: { icon: React.ReactNode; label: string; value: React.ReactNode; caption: string; tone: string }) {
  return (
    <article className="metric-card">
      <div className={`metric-icon ${tone}`}>{icon}</div>
      <div className="metric-label">{label}</div>
      <strong>{value ?? '—'}</strong>
      <small>{caption}</small>
    </article>
  )
}

function SetupItem({ done, label, to }: { done: boolean; label: string; to: string }) {
  return (
    <Link className={done ? 'setup-item done' : 'setup-item'} to={to}>
      <span>{done ? <Check size={14} /> : null}</span>
      <div>{label}</div>
      <ArrowRight size={14} />
    </Link>
  )
}
