import { useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  BarChart3,
  Braces,
  CircleDot,
  Clock3,
  DatabaseZap,
  Gauge,
  Plus,
  Radio,
  Search,
  Send,
} from 'lucide-react'
import { useSearchParams } from 'react-router-dom'
import { toast } from 'sonner'
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { meterforge, errorMessage, type Meter } from '../lib/meterforge'
import { formatDate, formatNumber, shortId, slugify } from '../lib/format'
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

const tabItems = [
  { value: 'meters', label: 'Meters' },
  { value: 'query', label: 'Query' },
  { value: 'events', label: 'Events' },
]

export function MeteringPage() {
  const [params, setParams] = useSearchParams()
  const tab = params.get('tab') ?? 'meters'
  const [createOpen, setCreateOpen] = useState(false)
  const meters = useQuery({ queryKey: ['meters'], queryFn: () => meterforge.meters.list() })

  const changeTab = (value: string) => setParams(value === 'meters' ? {} : { tab: value })

  return (
    <div className="page">
      <PageHeader
        eyebrow="Metering & Billing /"
        title="Metering"
        description="定义事件聚合规则、发送 CloudEvents，并把实时用量可视化。"
        actions={
          tab === 'meters' ? (
            <Button onClick={() => setCreateOpen(true)}><Plus size={16} /> Create Meter</Button>
          ) : undefined
        }
      />
      <Tabs items={tabItems.map((item) => ({ ...item, count: item.value === 'meters' ? meters.data?.length : undefined }))} value={tab} onChange={changeTab} />

      {tab === 'meters' && <MetersView query={meters} onCreate={() => setCreateOpen(true)} />}
      {tab === 'query' && <QueryView meters={meters.data ?? []} />}
      {tab === 'events' && <EventsView meters={meters.data ?? []} />}

      <CreateMeterModal open={createOpen} onClose={() => setCreateOpen(false)} />
    </div>
  )
}

function MetersView({ query, onCreate }: { query: ReturnType<typeof useQuery<Meter[]>>; onCreate: () => void }) {
  const [search, setSearch] = useState('')
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const items = (query.data ?? []).filter((meter) => `${meter.name} ${meter.slug} ${meter.eventType}`.toLowerCase().includes(search.toLowerCase()))

  return (
    <section className="panel table-panel">
      <div className="table-toolbar">
        <div className="search-input"><Search size={16} /><input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search meters" /></div>
        <span>{items.length} aggregation rules</span>
      </div>
      {!items.length ? (
        <EmptyState icon={<Gauge />} title="Create your first meter" description="Meters aggregate usage events into queryable business metrics." action={<Button onClick={onCreate}><Plus size={16} /> Create Meter</Button>} />
      ) : (
        <div className="meter-card-grid">
          {items.map((meter) => (
            <article className="meter-card" key={meter.id}>
              <div className="meter-card-top">
                <span className="resource-icon"><DatabaseZap size={18} /></span>
                <Badge tone="blue">{meter.aggregation}</Badge>
              </div>
              <h3>{meter.name}</h3>
              <code>{meter.slug}</code>
              <p>{meter.description || `Aggregates ${meter.eventType} events.`}</p>
              <div className="meter-meta">
                <span><Radio size={14} /> {meter.eventType}</span>
                <span><Braces size={14} /> {Object.keys(meter.groupBy ?? {}).length} dimensions</span>
              </div>
              <div className="meter-card-footer">
                <span>Created {formatDate(meter.createdAt)}</span>
                <a href={`/metering?tab=query&meter=${meter.slug}`}>Query usage →</a>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  )
}

function QueryView({ meters }: { meters: Meter[] }) {
  const [params, setParams] = useSearchParams()
  const [meterSlug, setMeterSlug] = useState(params.get('meter') ?? meters[0]?.slug ?? '')
  const [subject, setSubject] = useState('')
  const [windowSize, setWindowSize] = useState<'HOUR' | 'DAY' | 'MONTH'>('DAY')
  const [queryArgs, setQueryArgs] = useState({ meter: meterSlug, subject: '', window: windowSize })

  const query = useQuery({
    queryKey: ['meter-query', queryArgs],
    queryFn: () => meterforge.meters.query(queryArgs.meter, {
      windowSize: queryArgs.window,
      subject: queryArgs.subject ? [queryArgs.subject] : undefined,
    }),
    enabled: Boolean(queryArgs.meter),
  })

  const chartData = useMemo(() => (query.data?.data ?? []).map((row) => ({
    label: formatDate(row.windowStart),
    value: row.value,
    subject: row.subject ?? 'all',
  })).slice(-40), [query.data])
  const total = chartData.reduce((sum, row) => sum + row.value, 0)

  if (!meters.length) return <EmptyState icon={<BarChart3 />} title="Create a meter first" description="Usage queries need at least one meter definition." />

  return (
    <div className="query-layout">
      <aside className="panel query-builder">
        <div className="panel-heading"><div><span className="panel-kicker">Query builder</span><h2>Usage filters</h2></div></div>
        <Field label="Meter">
          <select value={meterSlug} onChange={(event) => { setMeterSlug(event.target.value); setParams({ tab: 'query', meter: event.target.value }) }}>
            {meters.map((meter) => <option key={meter.id} value={meter.slug}>{meter.name}</option>)}
          </select>
        </Field>
        <Field label="Subject" hint="留空查询全部客户">
          <input value={subject} onChange={(event) => setSubject(event.target.value)} placeholder="customer-1" />
        </Field>
        <Field label="Window size">
          <select value={windowSize} onChange={(event) => setWindowSize(event.target.value as typeof windowSize)}>
            <option value="HOUR">Hourly</option>
            <option value="DAY">Daily</option>
            <option value="MONTH">Monthly</option>
          </select>
        </Field>
        <Button onClick={() => setQueryArgs({ meter: meterSlug, subject, window: windowSize })}><BarChart3 size={16} /> Run query</Button>
      </aside>

      <section className="panel query-result">
        <div className="panel-heading">
          <div><span className="panel-kicker">Aggregated usage</span><h2>{meters.find((meter) => meter.slug === queryArgs.meter)?.name}</h2></div>
          <div className="query-total"><small>Total</small><strong>{formatNumber(total)}</strong></div>
        </div>
        {query.isLoading ? <LoadingState label="正在聚合 ClickHouse 用量" /> : query.error ? <ErrorState message={errorMessage(query.error)} /> : chartData.length ? (
          <>
            <div className="query-chart">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData} margin={{ top: 10, right: 10, left: -18, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#eaecf3" />
                  <XAxis dataKey="label" tickLine={false} axisLine={false} tick={{ fill: '#7c8293', fontSize: 10 }} />
                  <YAxis tickLine={false} axisLine={false} tick={{ fill: '#7c8293', fontSize: 11 }} />
                  <Tooltip cursor={{ fill: '#f3f4f9' }} />
                  <Bar dataKey="value" fill="#4f46e5" radius={[5, 5, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
            <div className="query-rows">
              {chartData.slice(-6).reverse().map((row, index) => <div key={`${row.label}-${index}`}><span>{row.label}</span><code>{row.subject}</code><strong>{formatNumber(row.value)}</strong></div>)}
            </div>
          </>
        ) : <EmptyState icon={<BarChart3 />} title="No usage in this range" description="尝试清空 Subject，或先发送一条匹配该 Meter 的事件。" />}
      </section>
    </div>
  )
}

function EventsView({ meters }: { meters: Meter[] }) {
  const [sendOpen, setSendOpen] = useState(false)
  const events = useQuery({ queryKey: ['events'], queryFn: () => meterforge.events.listV2({ limit: 50 }) })
  return (
    <section className="panel table-panel">
      <div className="table-toolbar">
        <div><strong>Event stream</strong><small> CloudEvents ingestion and validation</small></div>
        <Button onClick={() => setSendOpen(true)}><Send size={16} /> Send event</Button>
      </div>
      {events.isLoading ? <LoadingState /> : events.error ? <ErrorState message={errorMessage(events.error)} /> : (
        <div className="data-table-wrap">
          <table className="data-table">
            <thead><tr><th>Status</th><th>Type</th><th>Subject</th><th>Event ID</th><th>Ingested</th><th>Validation</th></tr></thead>
            <tbody>
              {(events.data?.items ?? []).map((item) => (
                <tr key={item.event.id}>
                  <td><Badge tone={item.validationError ? 'red' : 'green'}>{item.validationError ? 'invalid' : 'stored'}</Badge></td>
                  <td><strong>{item.event.type}</strong></td>
                  <td><code>{item.event.subject}</code></td>
                  <td title={item.event.id}><code>{shortId(item.event.id, 13)}</code></td>
                  <td>{formatDate(item.ingestedAt)}</td>
                  <td className="validation-cell">{item.validationError || 'Matched meter and customer'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      <SendEventModal open={sendOpen} onClose={() => setSendOpen(false)} meters={meters} />
    </section>
  )
}

function CreateMeterModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const client = useQueryClient()
  const [name, setName] = useState('LLM Tokens')
  const [slug, setSlug] = useState('llm_tokens_total')
  const [eventType, setEventType] = useState('llm.request')
  const [aggregation, setAggregation] = useState<'COUNT' | 'SUM' | 'UNIQUE_COUNT'>('SUM')
  const [valueProperty, setValueProperty] = useState('$.tokens')
  const [groupBy, setGroupBy] = useState('model=$.model\nprovider=$.provider')
  const create = useMutation({
    mutationFn: () => meterforge.meters.create({
      name,
      slug,
      eventType,
      aggregation,
      valueProperty: aggregation === 'COUNT' ? undefined : valueProperty,
      groupBy: Object.fromEntries(groupBy.split('\n').map((line) => line.trim()).filter(Boolean).map((line) => line.split('=').map((part) => part.trim())).filter((parts) => parts.length === 2)),
      description: `Tracks ${name.toLowerCase()} usage.`,
    }),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ['meters'] })
      toast.success('Meter 已创建', { description: '事件消费者大约需要 15 秒刷新聚合规则。' })
      onClose()
    },
    onError: (error) => toast.error(errorMessage(error)),
  })

  return (
    <Modal open={open} onClose={onClose} title="New Meter" description="Define which events to track and how their values are aggregated.">
      <form className="form-stack" onSubmit={(event) => { event.preventDefault(); create.mutate() }}>
        <div className="form-section-title"><span>1</span><div><strong>General Information</strong><small>Identify and manage this meter.</small></div></div>
        <div className="form-grid two">
          <Field label="Name"><input required value={name} onChange={(event) => { setName(event.target.value); if (!slug) setSlug(slugify(event.target.value)) }} /></Field>
          <Field label="Key"><input required value={slug} onChange={(event) => setSlug(slugify(event.target.value))} /></Field>
        </div>
        <div className="form-section-title"><span>2</span><div><strong>Event Processing</strong><small>Configure matching and aggregation.</small></div></div>
        <Field label="Event Type Filter"><input required value={eventType} onChange={(event) => setEventType(event.target.value)} /></Field>
        <div className="form-grid two">
          <Field label="Aggregation"><select value={aggregation} onChange={(event) => setAggregation(event.target.value as typeof aggregation)}><option>SUM</option><option>COUNT</option><option>UNIQUE_COUNT</option></select></Field>
          <Field label="Value property" hint="COUNT 聚合不需要此字段"><input disabled={aggregation === 'COUNT'} value={valueProperty} onChange={(event) => setValueProperty(event.target.value)} /></Field>
        </div>
        <Field label="Labels to group by" hint="每行格式：name=$.property"><textarea rows={3} value={groupBy} onChange={(event) => setGroupBy(event.target.value)} /></Field>
        <div className="info-banner"><Clock3 size={17} /><span>创建后 Sink Worker 的 Meter 缓存会在约 15 秒内刷新。等待就绪后再发送首条事件，可避免事件无法回溯聚合。</span></div>
        <div className="modal-actions"><Button type="button" variant="secondary" onClick={onClose}>Cancel</Button><Button type="submit" disabled={create.isPending}>{create.isPending ? 'Creating…' : 'Create Meter'}</Button></div>
      </form>
    </Modal>
  )
}

function SendEventModal({ open, onClose, meters }: { open: boolean; onClose: () => void; meters: Meter[] }) {
  const client = useQueryClient()
  const [meterSlug, setMeterSlug] = useState(meters[0]?.slug ?? '')
  const meter = meters.find((item) => item.slug === meterSlug) ?? meters[0]
  const customers = useQuery({
    queryKey: ['customers', 'event-subjects'],
    queryFn: () => meterforge.customers.list({ pageSize: 100 }),
    enabled: open,
  })
  const subjects = useMemo(() => customers.data?.items.flatMap((customer) => customer.usageAttribution?.subjectKeys ?? []) ?? [], [customers.data])
  const [subject, setSubject] = useState('')
  const [data, setData] = useState('{}')

  useEffect(() => {
    if (!meterSlug && meters[0]) setMeterSlug(meters[0].slug)
  }, [meterSlug, meters])

  useEffect(() => {
    if (meter) setData(JSON.stringify(sampleDataForMeter(meter), null, 2))
  }, [meter])

  useEffect(() => {
    if (!subject && subjects[0]) setSubject(subjects[0])
  }, [subject, subjects])
  const send = useMutation({
    mutationFn: async () => {
      const parsed = JSON.parse(data) as Record<string, unknown>
      return meterforge.events.ingest({ type: meter.eventType, subject, source: 'meterforge-console', data: parsed })
    },
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ['events'] })
      client.invalidateQueries({ queryKey: ['events', 'overview'] })
      client.invalidateQueries({ queryKey: ['meter-query'] })
      toast.success('Event accepted', { description: '异步聚合通常在几秒内可查询。' })
      onClose()
    },
    onError: (error) => toast.error(error instanceof SyntaxError ? 'Event data 不是有效 JSON' : errorMessage(error)),
  })

  return (
    <Modal open={open} onClose={onClose} title="Send usage event" description="Send a CloudEvent into the local Kafka ingestion pipeline.">
      <form className="form-stack" onSubmit={(event) => { event.preventDefault(); send.mutate() }}>
        <Field label="Meter"><select value={meter?.slug ?? ''} onChange={(event) => setMeterSlug(event.target.value)}>{meters.map((item) => <option value={item.slug} key={item.id}>{item.name} · {item.eventType}</option>)}</select></Field>
        <Field label="Subject" hint={subjects.length ? '已从 Customer 的 Usage Attribution 自动选择' : '请先为 Customer 配置 Usage Attribution'}>
          <input required list="event-subjects" value={subject} onChange={(event) => setSubject(event.target.value)} />
          <datalist id="event-subjects">{subjects.map((item) => <option value={item} key={item} />)}</datalist>
        </Field>
        <Field label="Event data"><textarea className="code-input" rows={8} value={data} onChange={(event) => setData(event.target.value)} /></Field>
        <div className="event-preview"><CircleDot size={16} /><span>type</span><code>{meter?.eventType}</code><span>source</span><code>meterforge-console</code></div>
        <div className="modal-actions"><Button type="button" variant="secondary" onClick={onClose}>Cancel</Button><Button type="submit" disabled={!meter || send.isPending}><Send size={16} /> {send.isPending ? 'Sending…' : 'Send Event'}</Button></div>
      </form>
    </Modal>
  )
}

function sampleDataForMeter(meter: Meter) {
  const data: Record<string, unknown> = {}
  if (meter.aggregation !== 'COUNT' && meter.valueProperty) setJsonPath(data, meter.valueProperty, 125)
  for (const [name, path] of Object.entries(meter.groupBy ?? {})) {
    const value = name.includes('method') ? 'POST'
      : name.includes('route') ? '/v1/chat/completions'
        : name.includes('model') ? 'gpt-4o-mini'
          : name.includes('provider') ? 'openai'
            : 'demo'
    setJsonPath(data, path, value)
  }
  return data
}

function setJsonPath(target: Record<string, unknown>, path: string, value: unknown) {
  const parts = path.replace(/^\$\.?/, '').split('.').filter(Boolean)
  if (!parts.length) return
  let cursor = target
  for (const part of parts.slice(0, -1)) {
    const child = cursor[part]
    if (!child || typeof child !== 'object' || Array.isArray(child)) cursor[part] = {}
    cursor = cursor[part] as Record<string, unknown>
  }
  cursor[parts.at(-1)!] = value
}
