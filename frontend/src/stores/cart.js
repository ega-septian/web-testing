import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

const STORAGE_KEY = 'cart_items'

function loadItems() {
  try {
    const parsed = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]')
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

// Client-side only, on purpose — there's no backend cart/order system yet
// (same "known gap" as checkout in general). Persisted to localStorage so it
// survives a refresh, but it's per-browser, not tied to the logged-in account.
export const useCartStore = defineStore('cart', () => {
  const items = ref(loadItems())

  watch(
    items,
    (value) => {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(value))
    },
    { deep: true },
  )

  const totalCount = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))
  const totalPrice = computed(() => items.value.reduce((sum, item) => sum + item.price * item.quantity, 0))

  // Same product+size combo bumps the existing line's quantity instead of
  // adding a duplicate row — "size" is part of the identity since each size
  // is really a distinct variant (own stock, same as the product detail page).
  function addItem({ productId, name, brand, imageUrl, price, size, quantity }) {
    const existing = items.value.find((item) => item.productId === productId && item.size === size)
    if (existing) {
      existing.quantity += quantity
    } else {
      items.value.push({ productId, name, brand, imageUrl, price, size, quantity })
    }
  }

  function removeItem(productId, size) {
    items.value = items.value.filter((item) => !(item.productId === productId && item.size === size))
  }

  function updateQuantity(productId, size, quantity) {
    const item = items.value.find((item) => item.productId === productId && item.size === size)
    if (item && quantity > 0) item.quantity = quantity
  }

  function clear() {
    items.value = []
  }

  return { items, totalCount, totalPrice, addItem, removeItem, updateQuantity, clear }
})
