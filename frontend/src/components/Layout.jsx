import { useState } from 'react'
import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../context/auth-context'
import { ROLE_LABEL, navGroups } from '../lib/resources'

const linkCls = ({ isActive }) =>
  `block rounded-lg px-3 py-2 text-sm ${isActive ? 'bg-indigo-600 font-medium text-white' : 'text-slate-300 hover:bg-slate-800'}`

export default function Layout() {
  const { user, logout } = useAuth()
  const [open, setOpen] = useState(false)
  const groups = navGroups(user.role)
  const close = () => setOpen(false)

  return (
    <div className="flex min-h-screen bg-slate-50 text-slate-900">
      <aside
        className={`fixed inset-y-0 left-0 z-40 w-64 overflow-y-auto bg-slate-900 p-4 transition-transform md:static md:translate-x-0 ${open ? 'translate-x-0' : '-translate-x-full'}`}
        aria-label="Navigasi utama"
      >
        <p className="mb-6 px-3 text-lg font-bold text-white">Property Management</p>
        <nav className="space-y-5">
          <NavLink to="/" end className={linkCls} onClick={close}>
            Dashboard
          </NavLink>
          {groups.map((g) => (
            <div key={g.name}>
              <p className="mb-1 px-3 text-xs font-semibold uppercase tracking-wide text-slate-500">{g.name}</p>
              <div className="space-y-0.5">
                {g.items.map((r) => (
                  <NavLink key={r.key} to={`/${r.key}`} className={linkCls} onClick={close}>
                    {r.label}
                  </NavLink>
                ))}
              </div>
            </div>
          ))}
        </nav>
      </aside>

      {open && <div className="fixed inset-0 z-30 bg-slate-900/40 md:hidden" onClick={close} aria-hidden="true" />}

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex items-center justify-between gap-3 border-b border-slate-200 bg-white px-4 py-3">
          <button
            type="button"
            className="rounded-lg border border-slate-300 px-3 py-1.5 text-sm md:hidden"
            aria-label="Buka menu"
            onClick={() => setOpen(true)}
          >
            ☰
          </button>
          <div className="ml-auto flex items-center gap-3 text-sm">
            <span className="text-right">
              <span className="block font-medium">{user.full_name}</span>
              <span className="block text-xs text-slate-500">{ROLE_LABEL[user.role] ?? user.role}</span>
            </span>
            <button type="button" onClick={logout} className="rounded-lg border border-slate-300 px-3 py-1.5 hover:bg-slate-50">
              Keluar
            </button>
          </div>
        </header>
        <main className="flex-1 p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
