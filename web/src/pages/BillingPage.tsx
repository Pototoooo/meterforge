import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  CircleDollarSign,
  FileText,
  KeyRound,
  Plus,
  Search,
  ShieldCheck,
  UserRound,
  Users,
} from 'lucide-react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { toast } from 'sonner'
import { meterforge, errorMessage } from '../lib/meterforge'
import { formatDate, formatMoney, formatNumber, shortId, slugify } from '../lib/format'
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

export function BillingPage() {
  const [params, setParams] = useSearchParams()
  const tab = params.get('tab') ?? 'customers'
  const [createOpen, setCreateOpen] = useState(false)
  const customers = useQuery({ queryKey: ['customers'], queryFn: () => meterforge.customers.list({ pageSize: 100, expand: ['subscriptions'] }) })
  const entitlements = useQuery({ queryKey: ['entitlements'], queryFn: () => meterforge.entitlements.list({ query: { pageSize: 100 } }) })
  const invoices = useQuery({ queryKey: ['invoices'], queryFn: () => meterforge.billing.invoices.list({ pageSize: 100 }) })
  return (
    <div className="page">
      <div className="sandbox-banner"><span>Sandbox</span><p>You're using the Sandbox App for billing. No payment will be collected.</p><Link to="/settings">View settings</Link></div>
      <PageHeader
        eyebrow="Metering & Billing /"
        title="Billing"
        description="管理客户订阅、使用额度和从聚合用量生成的账单。"
        actions={tab === 'customers' ? <Button onClick={() => setCreateOpen(true)}><Plus size={16} /> Create Customer</Button> : undefined}
      />
      <Tabs
        value={tab}
        onChange={(value) => setParams(value === 'customers' ? {} : { tab: value })}
        items={[
          { value: 'customers', label: 'Customers', count: customers.data?.totalCount },
          { value: 'entitlements', label: 'Entitlements', count: entitlements.data?.totalCount },
          { value: 'invoices', label: 'Invoices', count: invoices.data?.totalCount },
        ]}
      />
      {tab === 'customers' && <CustomersView query={customers} onCreate={() => setCreateOpen(true)} />}
      {tab === 'entitlements' && <EntitlementsView query={entitlements} />}
      {tab === 'invoices' && <InvoicesView query={invoices} />}
      <CreateCustomerModal open={createOpen} onClose={() => setCreateOpen(false)} />
    </div>
  )
}

function CustomersView({ query, onCreate }: { query: ReturnType<typeof useQuery>; onCreate: () => void }) {
  const [search, setSearch] = useState('')
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const data = query.data as Awaited<ReturnType<typeof meterforge.customers.list>>
  const items = data.items.filter((item) => `${item.name} ${item.key} ${item.primaryEmail}`.toLowerCase().includes(search.toLowerCase()))
  return (
    <section className="panel table-panel">
      <div className="table-toolbar"><div className="search-input"><Search size={16} /><input value={search} onChange={(event) => setSearch(event.target.value)} placeholder="Search customers" /></div><span>{data.totalCount} billing accounts</span></div>
      {!items.length ? <EmptyState icon={<Users />} title="Create your first customer" description="Customers connect usage subjects to plans, entitlements and invoices." action={<Button onClick={onCreate}><Plus size={16} /> Create Customer</Button>} /> : (
        <div className="data-table-wrap">
          <table className="data-table customer-table">
            <thead><tr><th>Name</th><th>Key</th><th>Subjects</th><th>Subscription</th><th>Primary Email</th><th>Created</th></tr></thead>
            <tbody>{items.map((customer) => (
              <tr key={customer.id}>
                <td><Link className="customer-name" to={`/billing/customers/${customer.id}`}><span className="customer-avatar">{customer.name.slice(0, 2).toUpperCase()}</span><div><strong>{customer.name}</strong><small>{shortId(customer.id, 15)}</small></div></Link></td>
                <td><code>{customer.key ?? '—'}</code></td>
                <td>{customer.usageAttribution?.subjectKeys?.length ? customer.usageAttribution.subjectKeys.map((subject) => <Badge tone="blue" key={subject}>{subject}</Badge>) : <span className="muted">No subjects</span>}</td>
                <td>{customer.currentSubscriptionId ? <Badge tone="green">active</Badge> : <Badge tone="amber">none</Badge>}</td>
                <td>{customer.primaryEmail ?? '—'}</td>
                <td>{formatDate(customer.createdAt)}</td>
              </tr>
            ))}</tbody>
          </table>
        </div>
      )}
    </section>
  )
}

function EntitlementsView({ query }: { query: ReturnType<typeof useQuery> }) {
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const data = query.data as Awaited<ReturnType<typeof meterforge.entitlements.list>>
  return (
    <section className="panel table-panel">
      {!data.items.length ? <EmptyState icon={<ShieldCheck />} title="No entitlements yet" description="Entitlements are provisioned automatically when a customer subscribes to a plan." /> : (
        <div className="data-table-wrap"><table className="data-table"><thead><tr><th>Customer</th><th>Feature</th><th>Type</th><th>Allowance</th><th>Usage period</th><th>Access policy</th></tr></thead><tbody>
          {data.items.map((item) => <tr key={item.id}>
            <td><Link to={`/billing/customers/${item.customerId}`}><strong>{item.customerKey}</strong></Link></td>
            <td><code>{item.featureKey}</code></td>
            <td><Badge tone="blue">{item.type}</Badge></td>
            <td>{item.type === 'metered' ? formatNumber(item.issue?.amount ?? item.issueAfterReset) : 'Unlimited'}</td>
            <td>{item.type === 'metered' ? item.usagePeriod.intervalISO : '—'}</td>
            <td>{item.type === 'metered' && item.isSoftLimit ? 'Soft limit' : 'Hard limit'}</td>
          </tr>)}
        </tbody></table></div>
      )}
    </section>
  )
}

function InvoicesView({ query }: { query: ReturnType<typeof useQuery> }) {
  if (query.isLoading) return <LoadingState />
  if (query.error) return <ErrorState message={errorMessage(query.error)} />
  const data = query.data as Awaited<ReturnType<typeof meterforge.billing.invoices.list>>
  return (
    <section className="panel table-panel">
      {!data.items.length ? <EmptyState icon={<FileText />} title="View your invoices" description="Invoices are created automatically when subscriptions start and billing periods advance." /> : (
        <div className="data-table-wrap"><table className="data-table"><thead><tr><th>Invoice</th><th>Customer</th><th>Status</th><th>Period</th><th>Total</th><th>Collection</th></tr></thead><tbody>
          {data.items.map((invoice) => <tr key={invoice.id}>
            <td><div className="table-resource"><span className="resource-icon"><FileText size={16} /></span><div><strong>{invoice.number ?? shortId(invoice.id)}</strong><small>{invoice.type}</small></div></div></td>
            <td><Link to={`/billing/customers/${invoice.customer.id}`}><strong>{invoice.customer.name}</strong></Link></td>
            <td><Badge>{invoice.status}</Badge></td>
            <td>{formatDate(invoice.period?.from)} – {formatDate(invoice.period?.to)}</td>
            <td><strong>{formatMoney(invoice.totals.total, invoice.currency)}</strong></td>
            <td>{formatDate(invoice.collectionAt)}</td>
          </tr>)}
        </tbody></table></div>
      )}
    </section>
  )
}

function CreateCustomerModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const navigate = useNavigate()
  const client = useQueryClient()
  const [name, setName] = useState('Acme AI')
  const [key, setKey] = useState('acme_ai')
  const [subject, setSubject] = useState('acme-ai-prod')
  const [email, setEmail] = useState('billing@acme.example')
  const create = useMutation({
    mutationFn: () => meterforge.customers.create({ name, key, primaryEmail: email || undefined, currency: 'USD', usageAttribution: { subjectKeys: subject.split(',').map((item) => item.trim()).filter(Boolean) } }),
    onSuccess: (customer) => {
      client.invalidateQueries({ queryKey: ['customers'] })
      toast.success('Customer 已创建', { description: '下一步为客户选择并创建订阅。' })
      onClose()
      if (customer) navigate(`/billing/customers/${customer.id}`)
    },
    onError: (error) => toast.error(errorMessage(error)),
  })
  return (
    <Modal open={open} onClose={onClose} title="Create Customer" description="Connect external usage subjects to one billable customer.">
      <form className="form-stack" onSubmit={(event) => { event.preventDefault(); create.mutate() }}>
        <div className="form-section-title"><span>1</span><div><strong>General Information</strong><small>Basic information about the customer.</small></div></div>
        <div className="form-grid two"><Field label="Name"><input required value={name} onChange={(event) => { setName(event.target.value); if (!key) setKey(slugify(event.target.value)) }} /></Field><Field label="Key"><input required value={key} onChange={(event) => setKey(slugify(event.target.value))} /></Field></div>
        <div className="form-section-title"><span>2</span><div><strong>Usage Attribution</strong><small>Events with these subjects belong to this customer.</small></div></div>
        <Field label="Include usage from" hint="多个 Subject 使用逗号分隔"><div className="input-with-icon"><KeyRound size={16} /><input required value={subject} onChange={(event) => setSubject(event.target.value)} /></div></Field>
        <div className="form-section-title"><span>3</span><div><strong>Details</strong><small>Optional billing contact information.</small></div></div>
        <Field label="Primary Email"><input type="email" value={email} onChange={(event) => setEmail(event.target.value)} /></Field>
        <div className="modal-actions"><Button type="button" variant="secondary" onClick={onClose}>Cancel</Button><Button type="submit" disabled={create.isPending}><UserRound size={16} /> {create.isPending ? 'Creating…' : 'Create Customer'}</Button></div>
      </form>
    </Modal>
  )
}
