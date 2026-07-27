import type { ReactNode } from 'react'
import { AlertCircle, LoaderCircle, X } from 'lucide-react'
import { statusTone } from '../lib/format'

export function Button({
  children,
  variant = 'primary',
  className = '',
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
}) {
  return (
    <button className={`button button-${variant} ${className}`} {...props}>
      {children}
    </button>
  )
}

export function Badge({ children, tone }: { children: ReactNode; tone?: string }) {
  return <span className={`badge badge-${tone ?? statusTone(String(children))}`}>{children}</span>
}

export function PageHeader({
  eyebrow,
  title,
  description,
  actions,
}: {
  eyebrow?: string
  title: string
  description?: string
  actions?: ReactNode
}) {
  return (
    <div className="page-header">
      <div>
        {eyebrow && <div className="eyebrow">{eyebrow}</div>}
        <h1>{title}</h1>
        {description && <p>{description}</p>}
      </div>
      {actions && <div className="page-actions">{actions}</div>}
    </div>
  )
}

export function Tabs({
  items,
  value,
  onChange,
}: {
  items: Array<{ value: string; label: string; count?: number }>
  value: string
  onChange: (value: string) => void
}) {
  return (
    <div className="tabs">
      {items.map((item) => (
        <button
          className={item.value === value ? 'tab active' : 'tab'}
          key={item.value}
          onClick={() => onChange(item.value)}
        >
          {item.label}
          {item.count !== undefined && <span>{item.count}</span>}
        </button>
      ))}
    </div>
  )
}

export function EmptyState({
  icon,
  title,
  description,
  action,
}: {
  icon?: ReactNode
  title: string
  description: string
  action?: ReactNode
}) {
  return (
    <div className="empty-state">
      {icon && <div className="empty-icon">{icon}</div>}
      <h3>{title}</h3>
      <p>{description}</p>
      {action}
    </div>
  )
}

export function LoadingState({ label = '正在加载数据' }: { label?: string }) {
  return (
    <div className="loading-state">
      <LoaderCircle className="spin" size={22} /> {label}
    </div>
  )
}

export function ErrorState({ message }: { message: string }) {
  return (
    <div className="error-state">
      <AlertCircle size={20} />
      <div>
        <strong>请求失败</strong>
        <span>{message}</span>
      </div>
    </div>
  )
}

export function Modal({
  open,
  title,
  description,
  children,
  onClose,
}: {
  open: boolean
  title: string
  description?: string
  children: ReactNode
  onClose: () => void
}) {
  if (!open) return null
  return (
    <div className="modal-backdrop" role="presentation" onMouseDown={onClose}>
      <section
        className="modal"
        role="dialog"
        aria-modal="true"
        aria-label={title}
        onMouseDown={(event) => event.stopPropagation()}
      >
        <div className="modal-header">
          <div>
            <h2>{title}</h2>
            {description && <p>{description}</p>}
          </div>
          <button className="icon-button" aria-label="关闭" onClick={onClose}>
            <X size={18} />
          </button>
        </div>
        {children}
      </section>
    </div>
  )
}

export function Field({
  label,
  hint,
  children,
}: {
  label: string
  hint?: string
  children: ReactNode
}) {
  return (
    <label className="field">
      <span className="field-label">{label}</span>
      {children}
      {hint && <small>{hint}</small>}
    </label>
  )
}
