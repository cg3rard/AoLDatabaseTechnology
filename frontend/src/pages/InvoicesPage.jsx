import { useState } from 'react'
import Modal from '../components/Modal'
import ResourcePage from '../components/ResourcePage'
import { btnGhost, btnLink, btnPrimary, errorCls, inputCls } from '../components/ui'
import { useAuth } from '../context/auth-context'
import { api } from '../lib/api'

// Penghuni mengirim URL bukti bayar; status invoice menjadi "pending_verification".
function ProofModal({ invoice, onClose, onDone }) {
  const [url, setUrl] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  async function submit(e) {
    e.preventDefault()
    setBusy(true)
    setError('')
    try {
      await api(`/api/invoices/${invoice.id}/payment-proof`, { method: 'POST', body: { payment_proof_url: url.trim() } })
      onDone()
    } catch (err) {
      setError(err.message)
      setBusy(false)
    }
  }

  return (
    <Modal title={`Bukti bayar ${invoice.invoice_number}`} onClose={onClose}>
      <form onSubmit={submit} className="space-y-4">
        <div>
          <label htmlFor="proof-url" className="mb-1 block text-sm font-medium text-slate-700">
            URL gambar bukti bayar
          </label>
          <input
            id="proof-url"
            type="url"
            required
            placeholder="https://…"
            className={inputCls}
            value={url}
            onChange={(e) => setUrl(e.target.value)}
          />
          <p className="mt-1 text-xs text-slate-500">Harus diawali http:// atau https://. Admin akan memverifikasi pembayaran Anda.</p>
        </div>
        {error && (
          <p role="alert" className={errorCls}>
            {error}
          </p>
        )}
        <div className="flex justify-end gap-2">
          <button type="button" className={btnGhost} onClick={onClose}>
            Batal
          </button>
          <button type="submit" className={btnPrimary} disabled={busy}>
            {busy ? 'Mengirim…' : 'Kirim'}
          </button>
        </div>
      </form>
    </Modal>
  )
}

export default function InvoicesPage({ resource }) {
  const { user } = useAuth()
  const [proofFor, setProofFor] = useState(null)
  const [version, setVersion] = useState(0)

  const actions = (item) =>
    user.role === 'resident' && (item.payment_status === 'unpaid' || item.payment_status === 'pending_verification') ? (
      <button type="button" className={btnLink} onClick={() => setProofFor(item)}>
        Kirim bukti bayar
      </button>
    ) : null

  return (
    <>
      <ResourcePage resource={resource} extraActions={actions} refreshKey={version} />
      {proofFor && (
        <ProofModal
          invoice={proofFor}
          onClose={() => setProofFor(null)}
          onDone={() => {
            setProofFor(null)
            setVersion((v) => v + 1)
          }}
        />
      )}
    </>
  )
}
