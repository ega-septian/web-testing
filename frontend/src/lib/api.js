const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081'

export class ApiError extends Error {
  constructor(message, status) {
    super(message)
    this.status = status
  }
}

async function request(path, options = {}) {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })

  const data = await res.json().catch(() => null)

  if (!res.ok) {
    throw new ApiError(data?.error || 'Terjadi kesalahan, coba lagi', res.status)
  }

  return data
}

export function login({ email, password }) {
  return request('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function register({ email, password }) {
  return request('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export function getAsset(key) {
  return request(`/api/assets/${key}`, { method: 'GET' })
}

export function assetFileUrl(asset) {
  return `${BASE_URL}${asset.url}`
}
