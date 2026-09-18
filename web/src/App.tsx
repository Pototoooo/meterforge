import { Navigate, Route, Routes } from 'react-router-dom'
import { Layout } from './components/Layout'
import { OverviewPage } from './pages/OverviewPage'
import { MeteringPage } from './pages/MeteringPage'
import { CatalogPage } from './pages/CatalogPage'
import { BillingPage } from './pages/BillingPage'
import { CustomerDetailPage } from './pages/CustomerDetailPage'
import { SettingsPage } from './pages/SettingsPage'

export function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<OverviewPage />} />
        <Route path="metering" element={<MeteringPage />} />
        <Route path="catalog" element={<CatalogPage />} />
        <Route path="billing" element={<BillingPage />} />
        <Route path="billing/customers/:customerId" element={<CustomerDetailPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  )
}
