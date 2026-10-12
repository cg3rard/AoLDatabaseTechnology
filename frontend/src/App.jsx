import { BrowserRouter, Navigate, Outlet, Route, Routes } from 'react-router-dom'
import Layout from './components/Layout'
import ResourcePage from './components/ResourcePage'
import AuthProvider from './context/AuthProvider'
import { useAuth } from './context/auth-context'
import { RESOURCES, canRead } from './lib/resources'
import BookingsPage from './pages/BookingsPage'
import Dashboard from './pages/Dashboard'
import InvoicesPage from './pages/InvoicesPage'
import Login from './pages/Login'

function RequireAuth() {
  const { user, loading } = useAuth()
  if (loading) return <p className="p-8 text-sm text-slate-500">Memeriksa sesi…</p>
  if (!user) return <Navigate to="/login" replace />
  return <Outlet />
}

function ResourceRoute({ resource }) {
  const { user } = useAuth()
  if (!canRead(resource, user.role)) {
    return <p className="text-sm text-slate-600">Anda tidak punya akses ke halaman ini.</p>
  }
  if (resource.key === 'invoices') return <InvoicesPage resource={resource} />
  if (resource.key === 'facility-bookings') return <BookingsPage resource={resource} />
  return <ResourcePage key={resource.key} resource={resource} />
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route element={<RequireAuth />}>
            <Route element={<Layout />}>
              <Route index element={<Dashboard />} />
              {RESOURCES.map((r) => (
                <Route key={r.key} path={r.key} element={<ResourceRoute key={r.key} resource={r} />} />
              ))}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Route>
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}
