import { useEffect, useState } from 'react'
import { useAuth } from '../context/auth-context'
import { api } from '../lib/api'
import { formatCell } from '../lib/format'
import { canCreate, canWrite } from '../lib/resources'
import FormModal from './FormModal'
import { btnDanger, btnGhost, btnLink, btnPrimary, errorCls, inputCls } from './ui'

const GOOD = ['active', 'paid', 'completed', 'confirmed', 'resolved', 'success', 'collected', 'approved', 'closed', 'claimed', 'deposit_refunded']
const BAD = ['blocked', 'lost', 'cancelled', 'rejected', 'failed', 'urgent', 'triggered', 'expired', 'false_alarm', 'high']

function Badge({ value, children }) {
  const tone = GOOD.includes(value)
    ? 'bg-emerald-50 text-emerald-700 ring-emerald-200'
    : BAD.includes(value)
      ? 'bg-red-50 text-red-700 ring-red-200'
      : 'bg-amber-50 text-amber-800 ring-amber-200'
  return <span className={`inline-block whitespace-nowrap rounded-full px-2 py-0.5 text-xs font-medium ring-1 ${tone}`}>{children}</span>
}

// Halaman daftar + tambah/ubah/hapus yang dibangun dari konfigurasi di lib/resources.js.
//   extraActions(item) -> tombol tambahan per baris (opsional)
//   refreshKey         -> ubah nilainya untuk memuat ulang data dari luar
export default function ResourcePage({ resource, extraActions, refreshKey = 0 }) {
  const { user } = useAuth()
  const [filters, setFilters] = useState({})
  const [version, setVersion] = useState(0)
  const [modal, setModal] = useState(null) // null = tertutup, { item: null } = tambah, { item } = ubah
  const [actionError, setActionError] = useState('')
  const [result, setResult] = useState({ key: '', items: [], error: '' })

  const filterQuery = new URLSearchParams(Object.entries(filters).filter(([, v]) => v)).toString()
  const reqKey = `${filterQuery}|${version}|${refreshKey}`
  const loading = result.key !== reqKey

  useEffect(() => {
    let cancelled = false
    const params = Object.fromEntries(new URLSearchParams(filterQuery))
    api(resource.path, { params: { ...params, limit: 500 } })
      .then((items) => !cancelled && setResult({ key: reqKey, items, error: '' }))
      .catch((err) => !cancelled && setResult({ key: reqKey, items: [], error: err.message }))
    return () => {
      cancelled = true
    }
  }, [resource.path, filterQuery, reqKey])

  const fieldByName = Object.fromEntries(resource.fields.map((f) => [f.name, f]))
  const filterFields = resource.fields.filter((f) => f.filter)
  const writable = canWrite(resource, user.role)
  const showActions = writable || Boolean(extraActions)

  async function remove(item) {
    if (!window.confirm(`Hapus ${resource.label} #${item.id}?`)) return
    setActionError('')
    try {
      await api(`${resource.path}/${item.id}`, { method: 'DELETE' })
      setVersion((v) => v + 1)
    } catch (err) {
      setActionError(err.message)
    }
  }

  function saved() {
    setModal(null)
    setVersion((v) => v + 1)
  }

  return (
    <section>
      <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-2xl font-bold text-slate-900">{resource.label}</h1>
        {canCreate(resource, user.role) && (
          <button type="button" className={btnPrimary} onClick={() => setModal({ item: null })}>
            + Tambah
          </button>
        )}
      </div>

      {filterFields.length > 0 && (
        <div className="mb-4 flex flex-wrap gap-3">
          {filterFields.map((f) => (
            <label key={f.name} className="flex items-center gap-2 text-sm text-slate-600">
              {f.label}
              <select
                className={`${inputCls} w-auto`}
                value={filters[f.name] ?? ''}
                onChange={(e) => setFilters((prev) => ({ ...prev, [f.name]: e.target.value }))}
              >
                <option value="">Semua</option>
                {f.options.map((o) => (
                  <option key={o.value} value={o.value}>
                    {o.label}
                  </option>
                ))}
              </select>
            </label>
          ))}
          {filterQuery && (
            <button type="button" className={btnGhost} onClick={() => setFilters({})}>
              Reset filter
            </button>
          )}
        </div>
      )}

      {(result.error || actionError) && (
        <p role="alert" className={`${errorCls} mb-4`}>
          {result.error || actionError}
        </p>
      )}

      <div className="overflow-x-auto rounded-2xl bg-white shadow-sm ring-1 ring-slate-200">
        <table className="min-w-full text-left text-sm">
          <thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              {resource.columns.map((name) => (
                <th key={name} scope="col" className="whitespace-nowrap px-4 py-3 font-semibold">
                  {fieldByName[name]?.label ?? (name === 'id' ? 'ID' : name.replace(/_/g, ' '))}
                </th>
              ))}
              {showActions && (
                <th scope="col" className="px-4 py-3 text-right font-semibold">
                  Aksi
                </th>
              )}
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {loading && (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={resource.columns.length + 1}>
                  Memuat…
                </td>
              </tr>
            )}
            {!loading && result.items.length === 0 && !result.error && (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={resource.columns.length + 1}>
                  Belum ada data.
                </td>
              </tr>
            )}
            {!loading &&
              result.items.map((item) => (
                <tr key={item.id} className="hover:bg-slate-50">
                  {resource.columns.map((name) => {
                    const text = formatCell(name, fieldByName[name], item[name])
                    return (
                      <td key={name} className="px-4 py-3 align-top">
                        {name.endsWith('status') && item[name] ? <Badge value={item[name]}>{text}</Badge> : text}
                      </td>
                    )
                  })}
                  {showActions && (
                    <td className="whitespace-nowrap px-4 py-3 text-right">
                      {extraActions?.(item)}
                      {writable && (
                        <>
                          <button type="button" className={btnLink} onClick={() => setModal({ item })}>
                            Ubah
                          </button>
                          <button type="button" className={btnDanger} onClick={() => remove(item)}>
                            Hapus
                          </button>
                        </>
                      )}
                    </td>
                  )}
                </tr>
              ))}
          </tbody>
        </table>
      </div>
      {!loading && result.items.length > 0 && (
        <p className="mt-2 text-xs text-slate-500">{result.items.length} data</p>
      )}

      {modal && <FormModal resource={resource} user={user} item={modal.item} onClose={() => setModal(null)} onSaved={saved} />}
    </section>
  )
}
