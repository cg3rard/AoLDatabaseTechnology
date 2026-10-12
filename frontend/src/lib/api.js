// Pembungkus fetch: menambahkan token JWT, mengubah respons menjadi JSON,
// dan melempar ApiError (dengan pesan dari backend) jika permintaan gagal.

const API_URL = import.meta.env.VITE_API_URL ?? 'http://127.0.0.1:8080'
const TOKEN_KEY = 'aol_token'

export const getToken = () => localStorage.getItem(TOKEN_KEY)

export function setToken(token) {
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

export class ApiError extends Error {
  constructor(message, status) {
    super(message)
    this.status = status
  }
}

export async function api(path, { method = 'GET', body, params } = {}) {
  const url = new URL(path, API_URL)
  for (const [key, value] of Object.entries(params ?? {})) {
    if (value !== '' && value !== null && value !== undefined) {
      url.searchParams.set(key, value)
    }
  }

  const token = getToken()
  const headers = {}
  if (body !== undefined) headers['Content-Type'] = 'application/json'
  if (token) headers.Authorization = `Bearer ${token}`

  let res
  try {
    res = await fetch(url, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    })
  } catch {
    throw new ApiError('Tidak dapat terhubung ke server. Pastikan backend sedang berjalan.', 0)
  }

  if (res.status === 204) return null
  const data = await res.json().catch(() => null)

  if (!res.ok) {
    // Token kedaluwarsa/tidak valid: beri tahu AuthProvider agar logout.
    if (res.status === 401 && token) window.dispatchEvent(new Event('auth:expired'))
    throw new ApiError(data?.error ?? `Permintaan gagal (${res.status})`, res.status)
  }
  return data
}
