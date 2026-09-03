<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { formatRupiah, getProduct } from '../lib/api'

const router = useRouter()
const cart = useCartStore()

// Live stock per product, keyed by productId -> that product's `sizes`
// array (same shape the Product Detail page reads), fetched fresh every
// time the Cart page is opened so it never shows a value cached from when
// the item was added (GASNTIN-42).
const stockByProduct = ref({})

async function loadStock() {
  const productIds = [...new Set(cart.items.map((item) => item.productId))]
  const products = await Promise.all(productIds.map((id) => getProduct(id).catch(() => null)))
  const map = {}
  productIds.forEach((id, index) => {
    if (products[index]) map[id] = products[index].sizes ?? []
  })
  stockByProduct.value = map
}

onMounted(loadStock)

// Returns null while stock hasn't been fetched yet (or the product/size
// couldn't be resolved), so the template can tell "unknown" apart from "0".
function stockFor(item) {
  const sizes = stockByProduct.value[item.productId]
  if (!sizes) return null
  return sizes.find((s) => s.size === item.size)?.stock ?? 0
}

// True once real stock has dropped below what's already in the cart (e.g.
// another order consumed it after this item was added) — surfaced as an
// inline warning that blocks checkout until resolved (GASNTIN-46).
function isOverStock(item) {
  const stock = stockFor(item)
  return stock !== null && item.quantity > stock
}

const hasOverStockItem = () => cart.items.some((item) => isOverStock(item))

// The router guard on /checkout (meta.requiresAuth) sends a guest to /login
// first (with a redirect back here) — this button doesn't need to check
// auth itself.
function goToCheckout() {
  router.push({ name: 'checkout' })
}

// Clamped to at least 1 — going to 0 should be done via the remove button,
// not by decrementing the quantity field to nothing.
function decrement(item) {
  if (item.quantity > 1) {
    cart.updateQuantity(item.productId, item.size, item.quantity - 1)
  }
}

// Clamped to the item's current available stock, so quantity in the Cart
// can never be pushed past what Checkout would actually allow (GASNTIN-43).
function increment(item) {
  const stock = stockFor(item)
  if (stock !== null && item.quantity >= stock) return
  cart.updateQuantity(item.productId, item.size, item.quantity + 1)
}
</script>

<template>
  <div data-testid="cart-page" class="min-h-screen bg-white font-sans text-slate-900">
    <div class="mx-auto max-w-5xl px-6 py-10">
      <div class="flex items-center justify-between">
        <h1 data-testid="cart-heading" class="font-display text-3xl tracking-tight">Keranjang</h1>
        <RouterLink :to="{ name: 'shop' }" class="text-sm font-medium text-slate-500 underline underline-offset-2">
          Lanjutkan Belanja
        </RouterLink>
      </div>

      <div v-if="cart.items.length === 0" data-testid="cart-empty-state" class="py-24 text-center text-slate-400">
        <p class="text-6xl">🛒</p>
        <p class="mt-4 text-sm">Keranjang kamu masih kosong.</p>
        <RouterLink :to="{ name: 'shop' }" class="mt-4 inline-block rounded-full bg-black px-6 py-3 text-sm font-semibold text-white hover:bg-slate-800">
          Mulai Belanja
        </RouterLink>
      </div>

      <div v-else class="mt-8 grid gap-10 lg:grid-cols-[1fr_320px]">
        <!-- Line items -->
        <ul class="divide-y divide-slate-200">
          <li
            v-for="(item, index) in cart.items"
            :key="`${item.productId}-${item.size}`"
            :data-testid="`cart-item-${index}`"
            class="flex items-center gap-4 py-5"
          >
            <div class="flex h-20 w-20 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-cream">
              <img v-if="item.imageUrl" :src="item.imageUrl" :alt="item.name" class="h-full w-full object-cover" />
              <span v-else class="text-3xl text-slate-300">🖼️</span>
            </div>

            <div class="min-w-0 flex-1">
              <p class="text-xs font-bold uppercase tracking-wide text-slate-500">{{ item.brand }}</p>
              <p data-testid="cart-item-name" class="truncate text-sm font-medium">{{ item.name }}</p>
              <p class="mt-1 text-xs text-slate-500">Ukuran: {{ item.size }}</p>
              <p v-if="stockFor(item) !== null" :data-testid="`cart-item-stock-${index}`" class="mt-1 text-xs text-slate-500">
                Stok tersedia (ukuran {{ item.size }}): {{ stockFor(item) }}
              </p>
              <p v-if="isOverStock(item)" :data-testid="`cart-item-overstock-${index}`" class="mt-1 text-xs font-medium text-red-500">
                Kuantitas melebihi stok yang tersedia.
              </p>
            </div>

            <div class="flex items-center rounded-full border border-slate-300">
              <button
                :data-testid="`cart-item-decrement-${index}`"
                class="px-3 py-1.5 text-sm disabled:cursor-not-allowed disabled:opacity-40"
                :disabled="item.quantity <= 1"
                @click="decrement(item)"
              >−</button>
              <span :data-testid="`cart-item-quantity-${index}`" class="w-6 text-center text-sm font-medium">{{ item.quantity }}</span>
              <button
                :data-testid="`cart-item-increment-${index}`"
                class="px-3 py-1.5 text-sm"
                @click="increment(item)"
              >+</button>
            </div>

            <p :data-testid="`cart-item-subtotal-${index}`" class="w-28 shrink-0 text-right text-sm font-semibold">
              {{ formatRupiah(item.price * item.quantity) }}
            </p>

            <button
              :data-testid="`cart-item-remove-${index}`"
              title="Hapus"
              class="shrink-0 text-slate-400 hover:text-red-500"
              @click="cart.removeItem(item.productId, item.size)"
            >
              ✕
            </button>
          </li>
        </ul>

        <!-- Summary -->
        <div class="h-fit rounded-2xl bg-slate-50 p-6">
          <p class="text-sm font-semibold">Ringkasan Belanja</p>
          <div class="mt-4 flex items-center justify-between text-sm text-slate-600">
            <span>Subtotal ({{ cart.totalCount }} item)</span>
            <span data-testid="cart-summary-subtotal" class="font-medium text-slate-900">{{ formatRupiah(cart.totalPrice) }}</span>
          </div>
          <p class="mt-2 text-xs text-slate-400">Ongkos kirim dihitung saat checkout.</p>

          <button
            data-testid="cart-checkout-button"
            class="mt-6 w-full rounded-full bg-black px-6 py-3 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="hasOverStockItem()"
            @click="goToCheckout"
          >
            Checkout
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
