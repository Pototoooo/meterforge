import { useQuery } from '@tanstack/react-query'
import { Blocks, CheckCircle2, CreditCard, ReceiptText, Settings, Webhook } from 'lucide-react'
import { meterforge, errorMessage } from '../lib/meterforge'
import { formatDate } from '../lib/format'
import { Badge, EmptyState, ErrorState, LoadingState, PageHeader, Tabs } from '../components/ui'

export function SettingsPage() {
  const apps = useQuery({ queryKey: ['apps'], queryFn: () => meterforge.apps.list() })
  const marketplace = useQuery({ queryKey: ['marketplace'], queryFn: () => meterforge.apps.marketplace.list({ pageSize: 100 }) })
  const profiles = useQuery({ queryKey: ['billing-profiles'], queryFn: () => meterforge.billing.profiles.list({ pageSize: 100 }) })
  const error = apps.error ?? marketplace.error ?? profiles.error
  return (
    <div className="page">
      <PageHeader eyebrow="Metering & Billing /" title="Settings" description="查看本地账单应用、Billing Profile 和集成能力。" />
      <Tabs value="apps" onChange={() => undefined} items={[{ value: 'apps', label: 'Apps', count: apps.data?.items.length }, { value: 'profiles', label: 'Billing Profiles', count: profiles.data?.totalCount }, { value: 'notifications', label: 'Notifications' }]} />
      {error && <ErrorState message={errorMessage(error)} />}
      {apps.isLoading ? <LoadingState /> : (
        <>
          <section className="settings-section"><div className="section-heading"><div><h2>Installed</h2><p>Applications currently participating in the billing workflow.</p></div></div><div className="app-grid">{apps.data?.items.map((app) => <article className="app-card installed" key={app.id}><div className="app-logo sandbox"><Blocks size={24} /></div><div className="app-card-title"><div><h3>{app.name}</h3><span>{app.type}</span></div><Badge tone="green">{app.status}</Badge></div><p>{app.listing.description}</p><div className="capability-list">{app.listing.capabilities.map((capability) => <span key={capability.key}><CheckCircle2 size={14} /> {capability.name}</span>)}</div><footer>Installed {formatDate(app.createdAt)}</footer></article>)}</div></section>
          <section className="settings-section"><div className="section-heading"><div><h2>Available</h2><p>Optional providers for payment collection, invoicing and tax.</p></div></div><div className="app-grid">{marketplace.data?.items.filter((listing) => !apps.data?.items.some((app) => app.type === listing.type)).map((listing) => <article className="app-card" key={listing.type}><div className={`app-logo ${listing.type}`}>{listing.type === 'stripe' ? <CreditCard size={24} /> : listing.type.includes('invoice') ? <ReceiptText size={24} /> : <Webhook size={24} />}</div><div className="app-card-title"><div><h3>{listing.name}</h3><span>{listing.type}</span></div><Badge tone="blue">available</Badge></div><p>{listing.description}</p><div className="capability-list">{listing.capabilities.slice(0, 3).map((capability) => <span key={capability.key}>{capability.name}</span>)}</div><footer>Configuration intentionally omitted from demo mode</footer></article>)}</div></section>
          {!apps.data?.items.length && <EmptyState icon={<Settings />} title="No billing apps installed" description="The quickstart normally provisions the Sandbox app automatically." />}
        </>
      )}
    </div>
  )
}
