import {
  Activity,
  Boxes,
  ChevronDown,
  CircleDollarSign,
  Command,
  Gauge,
  LayoutDashboard,
  Search,
  Settings,
  Sparkles,
  Waypoints,
} from 'lucide-react'
import { NavLink, Outlet } from 'react-router-dom'

const nav = [
  { to: '/', label: 'Overview', icon: LayoutDashboard, end: true },
  { to: '/metering', label: 'Metering', icon: Gauge },
  { to: '/catalog', label: 'Product Catalog', icon: Boxes },
  { to: '/billing', label: 'Billing', icon: CircleDollarSign },
  { to: '/settings', label: 'Settings', icon: Settings },
]

export function Layout() {
  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="brand-mark">
          <Waypoints size={22} />
        </div>
        <button className="workspace-switcher">
          <span className="workspace-avatar">L</span>
          <span>
            <small>Workspace</small>
            Local Demo
          </span>
          <ChevronDown size={15} />
        </button>
        <div className="topbar-spacer" />
        <span className="local-pill">
          <span /> Local stack
        </span>
        <button className="topbar-action">
          <Sparkles size={16} /> Demo assistant
        </button>
        <div className="profile-avatar">OM</div>
      </header>

      <aside className="sidebar">
        <button className="search-button">
          <Search size={17} />
          <span>Search</span>
          <kbd>
            <Command size={11} /> K
          </kbd>
        </button>

        <div className="nav-label">Metering & Billing</div>
        <nav>
          {nav.map(({ to, label, icon: Icon, end }) => (
            <NavLink key={to} to={to} end={end}>
              <Icon size={18} />
              {label}
            </NavLink>
          ))}
        </nav>

        <div className="sidebar-divider" />
        <div className="nav-label">System</div>
        <div className="system-item">
          <Activity size={17} />
          <span>MeterForge API</span>
          <i />
        </div>

        <div className="sidebar-footer">
          <div className="sandbox-card">
            <span className="sandbox-icon">
              <Sparkles size={16} />
            </span>
            <div>
              <strong>Sandbox billing</strong>
              <small>不会发起真实扣款</small>
            </div>
          </div>
        </div>
      </aside>

      <main className="content">
        <Outlet />
      </main>
    </div>
  )
}
