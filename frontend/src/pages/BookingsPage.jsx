import { useState } from 'react'
import ResourcePage from '../components/ResourcePage'
import { btnDanger, errorCls } from '../components/ui'
import { useAuth } from '../context/auth-context'
import { api } from '../lib/api'

// Tabel booking + tombol "Batalkan" (PATCH /cancel) untuk pemilik booking atau admin.
export default function BookingsPage({ resource }) {
  const { user } = useAuth()
  const [version, setVersion] = useState(0)
  const [error, setError] = useState('')

  async function cancel(item) {
    if (!window.confirm(`Batalkan booking #${item.id}?`)) return
    setError('')
    try {
      await api(`${resource.path}/${item.id}/cancel`, { method: 'PATCH' })
      setVersion((v) => v + 1)
    } catch (err) {
      setError(err.message)
    }
  }

  const actions = (item) =>
    item.status === 'confirmed' && (user.role === 'admin' || user.role === 'resident') ? (
      <button type="button" className={btnDanger} onClick={() => cancel(item)}>
        Batalkan
      </button>
    ) : null

  return (
    <>
      {error && (
        <p role="alert" className={`${errorCls} mb-4`}>
          {error}
        </p>
      )}
      <ResourcePage resource={resource} extraActions={actions} refreshKey={version} />
    </>
  )
}
