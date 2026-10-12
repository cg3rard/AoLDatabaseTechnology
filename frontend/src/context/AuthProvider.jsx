import { useEffect, useState } from 'react'
import { api, getToken, setToken } from '../lib/api'
import { AuthContext } from './auth-context'

export default function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  // Kalau ada token tersimpan, kita harus memeriksanya dulu sebelum menampilkan halaman.
  const [loading, setLoading] = useState(() => Boolean(getToken()))

  useEffect(() => {
    if (!getToken()) return
    let cancelled = false
    api('/api/auth/me')
      .then((me) => {
        if (!cancelled) setUser(me)
      })
      .catch(() => setToken(null))
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [])

  // Dipicu api.js saat server menjawab 401 (token kedaluwarsa).
  useEffect(() => {
    const onExpired = () => {
      setToken(null)
      setUser(null)
    }
    window.addEventListener('auth:expired', onExpired)
    return () => window.removeEventListener('auth:expired', onExpired)
  }, [])

  async function login(email, password) {
    const data = await api('/api/auth/login', { method: 'POST', body: { email, password } })
    setToken(data.token)
    setUser(data.user)
  }

  function logout() {
    setToken(null)
    setUser(null)
  }

  return <AuthContext.Provider value={{ user, loading, login, logout }}>{children}</AuthContext.Provider>
}
