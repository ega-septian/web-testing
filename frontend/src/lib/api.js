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

// For endpoints that require login (checkout, order history). Reads the
// token straight from localStorage — the same key stores/auth.js writes —
// rather than importing that store here, which would create a circular
// import (stores/auth.js already imports this file).
function authedRequest(path, options = {}) {
  const token = localStorage.getItem('auth_token')
  return request(path, {
    ...options,
    headers: {
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  })
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

// Fetches every uploaded asset in one call — used for the homepage images
// (hero banner + each dress-style photo) instead of one request per key.
export function getAssets() {
  return request('/api/assets', { method: 'GET' })
}

// sort: "newest" (default, by created_at) or "best_selling" (by total units
// sold). brand/gender/category/subcategory/size are the Shop page's
// checklist filters — each accepts an array for multi-select (e.g. gender:
// ['Pria', 'Wanita']), sent as repeated query params; a plain string also
// works for a single value (e.g. the homepage's brand strip linking straight
// to one brand). q is free-text search (name/brand/description) — what the
// homepage search box and the Shop page's own search field use.
export function getProducts({ sort, limit, brand, gender, category, subcategory, size, q } = {}) {
  const params = new URLSearchParams()
  if (sort) params.set('sort', sort)
  if (limit) params.set('limit', limit)
  if (q) params.set('q', q)
  const appendAll = (key, value) => {
    for (const v of Array.isArray(value) ? value : [value]) {
      if (v) params.append(key, v)
    }
  }
  if (brand) appendAll('brand', brand)
  if (gender) appendAll('gender', gender)
  if (category) appendAll('category', category)
  if (subcategory) appendAll('subcategory', subcategory)
  if (size) appendAll('size', size)
  const query = params.toString()
  return request(`/api/products${query ? `?${query}` : ''}`, { method: 'GET' })
}

export function getProduct(id) {
  return request(`/api/products/${id}`, { method: 'GET' })
}

// The Shop page's filter sidebar options (Gender, Kategori, Sub Kategori,
// Ukuran) with live counts — see CatalogHandler.ListProductFilters for the
// "counts aren't narrowed by other active filters" trade-off.
export function getProductFilters() {
  return request('/api/products/filters', { method: 'GET' })
}

// No real checkout flow exists yet — this exists so a sale can be recorded
// manually (curl/admin), which is what feeds the "best_selling" sort above.
export function recordSale(productId, quantity) {
  return request(`/api/products/${productId}/sales`, {
    method: 'POST',
    body: JSON.stringify({ quantity }),
  })
}

// Places an order from the cart. items is [{ productId, size, quantity }] —
// only that (plus shipping info) is sent; the backend resolves each item's
// current name/brand/price and validates stock itself (see
// OrderRepo.Place), so a stale cart can't under/over-charge or oversell.
// Requires auth (Bearer token added automatically — see requireAuth below).
export function checkout({ recipientName, phone, address, items }) {
  return authedRequest('/api/orders', {
    method: 'POST',
    body: JSON.stringify({
      recipient_name: recipientName,
      phone,
      address,
      items: items.map((item) => ({ product_id: item.productId, size: item.size, quantity: item.quantity })),
    }),
  })
}

// The logged-in user's order history, newest first.
export function getOrders() {
  return authedRequest('/api/orders', { method: 'GET' })
}

export function getOrder(id) {
  return authedRequest(`/api/orders/${id}`, { method: 'GET' })
}

// Turns a relative path (e.g. "/uploads/x.jpg") into a full URL. Returns
// null as-is so templates can `v-if` on it directly.
export function fileUrl(path) {
  return path ? `${BASE_URL}${path}` : null
}

export function assetFileUrl(asset) {
  return fileUrl(asset.url)
}

// Formats a plain integer price (stored as whole Rupiah, no cents) as
// "Rp 81.100" — id-ID locale gives the "." thousands separator this store uses.
export function formatRupiah(amount) {
  return `Rp ${Number(amount).toLocaleString('id-ID')}`
}

// Derives the struck-through "before discount" price for display, from the
// actual charged price + discount percentage — there's no separate original
// price stored, only the discount rate, so this is a rounded approximation.
export function originalPrice(price, discount) {
  if (!discount) return null
  return Math.round(price / (1 - discount / 100))
}
