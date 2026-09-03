<script setup>
import { ref, onMounted } from 'vue'
import { getOrders, formatRupiah } from '../lib/api'

const orders = ref([])
const loading = ref(true)
const loadError = ref(false)

onMounted(async () => {
  try {
    orders.value = await getOrders()
  } catch {
    loadError.value = true
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div data-testid="order-history-page" class="min-h-screen bg-white font-sans text-slate-900">
    <div class="mx-auto max-w-3xl px-6 py-10">
      <h1 data-testid="order-history-heading" class="font-display text-3xl tracking-tight">Riwayat Pesanan</h1>

      <div v-if="loading" class="py-16 text-center text-slate-400">Memuat...</div>

      <div v-else-if="loadError" data-testid="order-history-error-state" class="py-24 text-center text-slate-400">
        <p class="text-6xl">⚠️</p>
        <p class="mt-4 text-sm">Gagal memuat riwayat pesanan. Silakan coba lagi.</p>
      </div>

      <div v-else-if="orders.length === 0" data-testid="order-history-empty-state" class="py-24 text-center text-slate-400">
        <p class="text-6xl">📦</p>
        <p class="mt-4 text-sm">Belum ada pesanan.</p>
        <RouterLink :to="{ name: 'shop' }" class="mt-4 inline-block rounded-full bg-black px-6 py-3 text-sm font-semibold text-white hover:bg-slate-800">
          Mulai Belanja
        </RouterLink>
      </div>

      <ul v-else class="mt-8 space-y-4">
        <li v-for="(order, index) in orders" :key="order.id" :data-testid="`order-history-item-${index}`">
          <RouterLink
            :to="{ name: 'order-detail', params: { id: order.id } }"
            class="block rounded-2xl border border-slate-200 p-5 transition hover:border-slate-300 hover:shadow-sm"
          >
            <div class="flex items-center justify-between gap-3">
              <div class="min-w-0">
                <p class="font-mono text-xs text-slate-400">{{ order.id }}</p>
                <p class="mt-1 text-sm text-slate-500">{{ order.created_at }}</p>
              </div>
              <span class="shrink-0 rounded-full bg-green-100 px-3 py-1 text-xs font-semibold text-green-700">{{ order.status }}</span>
            </div>
            <p class="mt-3 truncate text-sm text-slate-600">
              {{ order.items.map((item) => `${item.product_name} (${item.size}) ×${item.quantity}`).join(', ') }}
            </p>
            <p class="mt-2 text-sm font-semibold">{{ formatRupiah(order.total_amount) }}</p>
          </RouterLink>
        </li>
      </ul>
    </div>
  </div>
</template>
