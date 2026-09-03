<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { checkout, formatRupiah, getProduct, ApiError } from '../lib/api'

const router = useRouter()
const cart = useCartStore()

const recipientName = ref('')
const phone = ref('')
const address = ref('')
const submitting = ref(false)
const errorMessage = ref('')

// Live stock per product, keyed by productId -> that product's `sizes`
// array, fetched fresh every time the Checkout page is opened — same
// approach as the Cart page (GASNTIN-42) — so a line item that became
// over-stock after being added to the cart is flagged here too, before
// the form is even touched (GASNTIN-49).
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

// Returns null while stock hasn't been fetched yet (or the product/size
// couldn't be resolved), so the template can tell "unknown" apart from "0".
function stockFor(item) {
  const sizes = stockByProduct.value[item.productId]
  if (!sizes) return null
  return sizes.find((s) => s.size === item.size)?.stock ?? 0
}

// True once real stock has dropped below what's already in the cart —
// surfaced as an inline warning naming the product and size (GASNTIN-49).
function isOverStock(item) {
  const stock = stockFor(item)
  return stock !== null && item.quantity > stock
}

// Landing here with nothing to check out (direct URL, or already checked
// out in another tab) has nowhere useful to go but back to the cart.
onMounted(() => {
  if (cart.items.length === 0) {
    router.replace({ name: 'cart' })
    return
  }
  loadStock()
})

async function handleSubmit() {
  submitting.value = true
  errorMessage.value = ''

  try {
    const order = await checkout({
      recipientName: recipientName.value,
      phone: phone.value,
      address: address.value,
      items: cart.items,
    })
    // Only clear the cart once the order actually exists — a failed
    // checkout (e.g. stock ran out) should leave the cart untouched.
    cart.clear()
    router.push({ name: 'order-detail', params: { id: order.id } })
  } catch (err) {
    errorMessage.value = err instanceof ApiError ? err.message : 'Gagal membuat pesanan, coba lagi'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div data-testid="checkout-page" class="min-h-screen bg-white font-sans text-slate-900">
    <div class="mx-auto max-w-5xl px-6 py-10">
      <h1 data-testid="checkout-heading" class="font-display text-3xl tracking-tight">Checkout</h1>

      <div class="mt-8 grid gap-10 lg:grid-cols-[1fr_360px]">
        <!-- Shipping form -->
        <form data-testid="checkout-form" class="space-y-5" @submit.prevent="handleSubmit">
          <div>
            <label for="recipient-name" class="block text-sm font-medium text-slate-700">Nama Penerima</label>
            <input
              id="recipient-name"
              v-model="recipientName"
              data-testid="checkout-recipient-name-input"
              type="text"
              required
              placeholder="Nama lengkap"
              class="mt-1.5 block w-full rounded-full border border-slate-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-black"
            />
          </div>

          <div>
            <label for="phone" class="block text-sm font-medium text-slate-700">Nomor HP</label>
            <input
              id="phone"
              v-model="phone"
              data-testid="checkout-phone-input"
              type="tel"
              required
              placeholder="08xxxxxxxxxx"
              class="mt-1.5 block w-full rounded-full border border-slate-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-black"
            />
          </div>

          <div>
            <label for="address" class="block text-sm font-medium text-slate-700">Alamat Pengiriman</label>
            <textarea
              id="address"
              v-model="address"
              data-testid="checkout-address-input"
              required
              rows="4"
              placeholder="Jalan, nomor rumah, kelurahan, kecamatan, kota, kode pos"
              class="mt-1.5 block w-full rounded-2xl border border-slate-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-black"
            ></textarea>
          </div>

          <div v-if="errorMessage" data-testid="checkout-error-message" class="rounded-lg bg-red-50 px-3.5 py-2.5 text-sm text-red-600 ring-1 ring-inset ring-red-200">
            {{ errorMessage }}
          </div>

          <button
            type="submit"
            data-testid="checkout-submit-button"
            :disabled="submitting"
            class="flex w-full items-center justify-center gap-2 rounded-full bg-black px-6 py-3 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <span v-if="submitting" class="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white"></span>
            {{ submitting ? 'Memproses...' : 'Buat Pesanan' }}
          </button>
        </form>

        <!-- Order summary -->
        <div class="h-fit rounded-2xl bg-slate-50 p-6">
          <p class="text-sm font-semibold">Ringkasan Pesanan</p>
          <ul class="mt-4 space-y-3">
            <li
              v-for="(item, index) in cart.items"
              :key="`${item.productId}-${item.size}`"
              :data-testid="`checkout-summary-item-${index}`"
            >
              <div class="flex items-center justify-between gap-3 text-sm">
                <span class="min-w-0 truncate text-slate-600">{{ item.name }} ({{ item.size }}) × {{ item.quantity }}</span>
                <span class="shrink-0 font-medium">{{ formatRupiah(item.price * item.quantity) }}</span>
              </div>
              <p
                v-if="isOverStock(item)"
                :data-testid="`checkout-summary-overstock-${index}`"
                class="mt-1 text-xs font-medium text-red-500"
              >
                {{ item.name }} (ukuran {{ item.size }}): Kuantitas melebihi stok yang tersedia.
              </p>
            </li>
          </ul>
          <div class="mt-4 flex items-center justify-between border-t border-slate-200 pt-4 text-sm">
            <span class="font-semibold">Total</span>
            <span data-testid="checkout-summary-total" class="text-base font-bold">{{ formatRupiah(cart.totalPrice) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
