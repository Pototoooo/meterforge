export function formatNumber(value: number | string | undefined | null) {
  if (value === undefined || value === null || value === '') return '—'
  const number = typeof value === 'number' ? value : Number(value)
  return Number.isFinite(number)
    ? new Intl.NumberFormat('zh-CN', { maximumFractionDigits: 4 }).format(number)
    : String(value)
}

export function formatMoney(
  value: number | string | undefined | null,
  currency = 'USD',
) {
  const number = Number(value ?? 0)
  try {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
      maximumFractionDigits: 4,
    }).format(number)
  } catch {
    return `${currency} ${formatNumber(number)}`
  }
}

export function formatDate(value: Date | string | undefined | null) {
  if (!value) return '—'
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return new Intl.DateTimeFormat('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

export function shortId(value: string | undefined | null, length = 10) {
  if (!value) return '—'
  return value.length > length ? `${value.slice(0, length)}…` : value
}

export function slugify(value: string) {
  return value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '_')
    .replace(/^_+|_+$/g, '')
}

export function statusTone(status: string | undefined) {
  const value = status?.toLowerCase() ?? ''
  if (['active', 'paid', 'issued', 'ready', 'valid'].includes(value)) return 'green'
  if (['draft', 'gathering', 'pending', 'scheduled'].includes(value)) return 'amber'
  if (['failed', 'invalid', 'overdue', 'voided', 'canceled'].includes(value)) return 'red'
  return 'blue'
}
