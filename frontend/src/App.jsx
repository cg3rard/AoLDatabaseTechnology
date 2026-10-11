import { useEffect, useState } from 'react'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

function App() {
  // 'loading' | 'ok' | 'error'
  const [status, setStatus] = useState('loading')

  useEffect(() => {
    const controller = new AbortController()
    fetch(`${API_URL}/api/health`, { signal: controller.signal })
      .then((res) => setStatus(res.ok ? 'ok' : 'error'))
      .catch((err) => {
        if (err.name !== 'AbortError') setStatus('error')
      })
    return () => controller.abort()
  }, [])

  const label = {
    loading: 'Memeriksa koneksi…',
    ok: 'Backend & database terhubung',
    error: 'Backend atau database tidak terjangkau',
  }[status]

  const color = {
    loading: 'bg-slate-400',
    ok: 'bg-emerald-500',
    error: 'bg-red-500',
  }[status]

  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 p-6 text-slate-900">
      <section className="w-full max-w-md rounded-2xl bg-white p-8 shadow-sm ring-1 ring-slate-200">
        <h1 className="text-2xl font-bold">Property Management</h1>
        <p className="mt-1 text-sm text-slate-500">AoL Database Technology</p>

        <div className="mt-6 flex items-center gap-3" role="status" aria-live="polite">
          <span className={`h-3 w-3 rounded-full ${color}`} aria-hidden="true" />
          <span className="text-sm">{label}</span>
        </div>
      </section>
    </main>
  )
}

export default App
