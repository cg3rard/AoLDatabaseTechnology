import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/auth-context'
import { api } from '../lib/api'
import { formatCell } from '../lib/format'
import { ROLE, ROLE_LABEL } from '../lib/resources'

const { admin, resident, laundry, security, technician } = ROLE

// Kartu ringkasan: jumlah data yang cocok dengan filter (diambil dari endpoint list).
const STATS = [
  { label: 'Unit terdaftar', key: 'units', roles: [admin, security, laundry, technician] },
  { label: 'Invoice menunggu verifikasi', key: 'invoices', params: { payment_status: 'pending_verification' }, roles: [admin] },
  { label: 'Invoice belum dibayar', key: 'invoices', params: { payment_status: 'unpaid' }, roles: [admin, resident] },
  { label: 'Tiket perbaikan terbuka', key: 'maintenance-tickets', params: { status: 'open' }, roles: [admin, technician, resident] },
  { label: 'Paket menunggu diambil', key: 'packages', params: { status: 'waiting_pickup' }, roles: [admin, security, resident] },
  { label: 'Permohonan kartu pending', key: 'card-requests', params: { status: 'pending' }, roles: [admin, security, resident] },
  { label: 'Booking terkonfirmasi', key: 'facility-bookings', params: { status: 'confirmed' }, roles: [admin, security, resident] },
  { label: 'Peringatan darurat aktif', key: 'emergency-alerts', params: { status: 'triggered' }, roles: [admin, security] },
]

export default function Dashboard() {
  const { user } = useAuth()
  const [cards, setCards] = useState(null)
  const [news, setNews] = useState([])

  useEffect(() => {
    let cancelled = false
    const list = STATS.filter((s) => s.roles.includes(user.role))

    Promise.allSettled(list.map((s) => api(`/api/${s.key}`, { params: { ...s.params, limit: 500 } }))).then((results) => {
      if (cancelled) return
      setCards(list.map((s, i) => ({ ...s, count: results[i].status === 'fulfilled' ? results[i].value.length : null })))
    })
    api('/api/announcements', { params: { limit: 3 } })
      .then((rows) => !cancelled && setNews(rows))
      .catch(() => {})

    return () => {
      cancelled = true
    }
  }, [user.role])

  return (
    <section className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold">Halo, {user.full_name}</h1>
        <p className="text-sm text-slate-500">Masuk sebagai {ROLE_LABEL[user.role] ?? user.role}</p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {cards === null && <p className="text-sm text-slate-500">Memuat ringkasan…</p>}
        {cards?.map((c) => (
          <Link
            key={c.label}
            to={`/${c.key}`}
            className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-200 hover:ring-indigo-400"
          >
            <p className="text-3xl font-bold">{c.count ?? '-'}</p>
            <p className="mt-1 text-sm text-slate-500">{c.label}</p>
          </Link>
        ))}
      </div>

      <div>
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-lg font-semibold">Pengumuman terbaru</h2>
          <Link to="/announcements" className="text-sm text-indigo-700 hover:underline">
            Lihat semua
          </Link>
        </div>
        {news.length === 0 ? (
          <p className="text-sm text-slate-500">Belum ada pengumuman.</p>
        ) : (
          <ul className="space-y-3">
            {news.map((n) => (
              <li key={n.id} className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-200">
                <p className="font-medium">{n.title}</p>
                <p className="mt-1 whitespace-pre-line text-sm text-slate-600">{n.content}</p>
                <p className="mt-2 text-xs text-slate-400">{formatCell('created_at', null, n.created_at)}</p>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  )
}
