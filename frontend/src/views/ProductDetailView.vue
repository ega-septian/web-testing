<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { getProduct, assetFileUrl, formatRupiah, originalPrice } from '../lib/api'
import { useCartStore } from '../stores/cart'

const route = useRoute()
const cart = useCartStore()

const product = ref(null)
const loading = ref(true)
const notFound = ref(false)
const activeImageIndex = ref(0)
const selectedSize = ref(null)
const quantity = ref(1)
const activeTab = ref('details') // 'details' | 'faqs'
const justAdded = ref(false)
let justAddedTimeout = null

onMounted(async () => {
  try {
    product.value = await getProduct(route.params.id)
    // Default to the first size that's actually in stock, falling back to
    // the first size at all if every size is sold out.
    const sizes = product.value.sizes ?? []
    selectedSize.value = sizes.find((s) => s.stock > 0)?.size ?? sizes[0]?.size ?? null
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
})

// The struck-through "before discount" price is derived, not stored — only
// the discount percentage is (see api.js: originalPrice).
const oldPrice = computed(() => {
  const p = product.value
  return p ? originalPrice(p.price, p.discount) : null
})

// Stock is tracked per size, not one total for the whole product.
const selectedSizeStock = computed(() => {
  return product.value?.sizes?.find((s) => s.size === selectedSize.value)?.stock ?? 0
})

const mainImageUrl = computed(() => {
  const img = product.value?.images?.[activeImageIndex.value]
  return img ? assetFileUrl(img) : null
})

function selectImage(index) {
  activeImageIndex.value = index
}

// Switching size can lower the available stock below the current quantity
// (e.g. quantity 13 on a size that only has 5 left) — clamp it back down.
function selectSize(size) {
  selectedSize.value = size
  const stock = product.value.sizes.find((s) => s.size === size)?.stock ?? 0
  if (stock > 0 && quantity.value > stock) {
    quantity.value = stock
  }
}

function incrementQuantity() {
  if (quantity.value < selectedSizeStock.value) quantity.value++
}

function decrementQuantity() {
  if (quantity.value > 1) quantity.value--
}

// Adds the currently selected size + quantity as one line item (see
// stores/cart.js — client-side only, no backend order system yet). Clamped
// against what's already in the Cart for this same product/size, so adding
// again can't push the combined quantity past available stock (GASNTIN-48).
function addToCart() {
  if (selectedSizeStock.value === 0) return

  const alreadyInCart =
    cart.items.find((item) => item.productId === product.value.id && item.size === selectedSize.value)
      ?.quantity ?? 0
  const room = selectedSizeStock.value - alreadyInCart
  if (room <= 0) return

  cart.addItem({
    productId: product.value.id,
    name: product.value.name,
    brand: product.value.brand,
    imageUrl: mainImageUrl.value,
    price: product.value.price,
    size: selectedSize.value,
    quantity: Math.min(quantity.value, room),
  })

  // Brief "Ditambahkan ✓" confirmation on the button itself, since there's no
  // toast system on this site yet.
  justAdded.value = true
  clearTimeout(justAddedTimeout)
  justAddedTimeout = setTimeout(() => {
    justAdded.value = false
  }, 2000)
}

onUnmounted(() => clearTimeout(justAddedTimeout))
</script>

<template>
  <div data-testid="product-detail-page" class="min-h-screen bg-white font-sans text-slate-900">
    <div v-if="loading" class="flex min-h-screen items-center justify-center text-slate-400">Memuat...</div>

    <div v-else-if="notFound" data-testid="product-detail-not-found" class="flex min-h-screen flex-col items-center justify-center gap-2 text-center">
      <p class="text-2xl font-bold">Produk tidak ditemukan</p>
      <RouterLink :to="{ name: 'home' }" class="text-sm font-medium underline underline-offset-2">Kembali ke Homepage</RouterLink>
    </div>

    <div v-else class="mx-auto max-w-6xl px-6 py-12">
      <div class="grid gap-10 lg:grid-cols-2">
        <!-- Gallery -->
        <div>
          <div class="flex aspect-square items-center justify-center overflow-hidden rounded-2xl bg-cream">
            <img
              v-if="mainImageUrl"
              :src="mainImageUrl"
              :alt="product.name"
              data-testid="product-detail-main-image"
              class="h-full w-full object-cover"
            />
            <span v-else class="text-8xl text-slate-300">🖼️</span>
          </div>

          <div v-if="product.images.length > 1" class="mt-4 flex gap-3">
            <button
              v-for="(img, index) in product.images"
              :key="img.id"
              :data-testid="`product-detail-thumbnail-${index}`"
              class="h-20 w-20 overflow-hidden rounded-xl ring-2 transition"
              :class="index === activeImageIndex ? 'ring-black' : 'ring-transparent hover:ring-slate-300'"
              @click="selectImage(index)"
            >
              <img :src="assetFileUrl(img)" :alt="`${product.name} ${index + 1}`" class="h-full w-full object-cover" />
            </button>
          </div>
        </div>

        <!-- Info -->
        <div>
          <p data-testid="product-detail-brand" class="text-sm font-bold uppercase tracking-wide">{{ product.brand }}</p>

          <h1 data-testid="product-detail-name" class="mt-1 font-display text-3xl tracking-tight sm:text-4xl">
            {{ product.name }}
          </h1>

          <div class="mt-4 flex items-center gap-3">
            <span
              data-testid="product-detail-price"
              class="text-2xl font-bold"
              :class="product.discount ? 'text-red-500' : 'text-slate-900'"
            >
              {{ formatRupiah(product.price) }}
            </span>
            <span v-if="product.discount" data-testid="product-detail-old-price" class="text-lg text-slate-400 line-through">
              {{ formatRupiah(oldPrice) }}
            </span>
            <span
              v-if="product.discount"
              data-testid="product-detail-discount-badge"
              class="rounded-full bg-red-50 px-2.5 py-1 text-xs font-semibold text-red-500"
            >
              -{{ product.discount }}%
            </span>
          </div>

          <p data-testid="product-detail-description" class="mt-4 text-sm text-slate-600">
            {{ product.description }}
          </p>

          <div class="mt-6 border-t border-slate-200 pt-6">
            <p class="text-sm font-medium text-slate-700">Pilih Ukuran</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <button
                v-for="s in product.sizes"
                :key="s.size"
                :data-testid="`product-detail-size-${s.size}`"
                :disabled="s.stock === 0"
                class="rounded-full px-4 py-2 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-40 disabled:line-through"
                :class="s.size === selectedSize ? 'bg-black text-white' : 'bg-slate-100 text-slate-700 hover:bg-slate-200'"
                @click="selectSize(s.size)"
              >
                {{ s.size }}
              </button>
            </div>
          </div>

          <p class="mt-4 text-sm" :class="selectedSizeStock > 0 ? 'text-slate-500' : 'text-red-500'">
            {{ selectedSizeStock > 0 ? `Stok tersedia (ukuran ${selectedSize}): ${selectedSizeStock}` : `Stok habis untuk ukuran ${selectedSize}` }}
          </p>

          <div class="mt-6 flex items-center gap-3 border-t border-slate-200 pt-6">
            <div class="flex items-center rounded-full border border-slate-300">
              <button
                data-testid="product-detail-quantity-decrement"
                class="px-4 py-2.5 text-lg disabled:cursor-not-allowed disabled:opacity-40"
                :disabled="quantity <= 1"
                @click="decrementQuantity"
              >−</button>
              <span data-testid="product-detail-quantity-value" class="w-8 text-center text-sm font-medium">{{ quantity }}</span>
              <button
                data-testid="product-detail-quantity-increment"
                class="px-4 py-2.5 text-lg disabled:cursor-not-allowed disabled:opacity-40"
                :disabled="quantity >= selectedSizeStock"
                @click="incrementQuantity"
              >+</button>
            </div>
            <button
              data-testid="product-detail-add-to-cart"
              :disabled="selectedSizeStock === 0"
              class="flex-1 rounded-full px-6 py-3 text-sm font-semibold text-white transition disabled:cursor-not-allowed disabled:opacity-40"
              :class="justAdded ? 'bg-green-600' : 'bg-black hover:bg-slate-800'"
              @click="addToCart"
            >
              {{ justAdded ? 'Ditambahkan ✓' : 'Add to Cart' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="mt-16 border-t border-slate-200">
        <div class="flex gap-8">
          <button
            data-testid="product-detail-tab-details"
            class="border-b-2 py-4 text-sm font-medium"
            :class="activeTab === 'details' ? 'border-black text-black' : 'border-transparent text-slate-400'"
            @click="activeTab = 'details'"
          >
            Product Details
          </button>
          <button
            data-testid="product-detail-tab-faqs"
            class="border-b-2 py-4 text-sm font-medium"
            :class="activeTab === 'faqs' ? 'border-black text-black' : 'border-transparent text-slate-400'"
            @click="activeTab = 'faqs'"
          >
            FAQs
          </button>
        </div>

        <div class="py-8 text-sm text-slate-600">
          <p v-if="activeTab === 'details'" data-testid="product-detail-tab-panel-details">{{ product.description }}</p>
          <p v-else data-testid="product-detail-tab-panel-faqs">Belum ada FAQ untuk produk ini.</p>
        </div>
      </div>
    </div>
  </div>
</template>
