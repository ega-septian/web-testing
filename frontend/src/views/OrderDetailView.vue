<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getOrder, formatRupiah } from '../lib/api'

const route = useRoute()

const order = ref(null)
const loading = ref(true)
const notFound = ref(false)

onMounted(async () => {
  try {
    order.value = await getOrder(route.params.id)
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div data-testid="order-detail-page" class="min-h-screen bg-white font-sans text-slate-900">
    <div v-if="loading" class="flex min-h-screen items-center justify-center text-slate-400">Memuat...</div>

    <div v-else-if="notFound" data-testid="order-detail-not-found" class="flex min-h-screen flex-col items-center justify-center gap-2 text-center">
      <p class="text-2xl font-bold">Pesanan tidak ditemukan</p>
      <RouterLink :to="{ name: 'orders' }" class="text-sm font-medium underline underline-offset-2">Lihat Riwayat Pesanan</RouterLink>
    </div>

    <div v-else class="mx-auto max-w-3xl px-6 py-12">
      <div class="text-center">
        <p class="text-5xl">✅</p>
        <h1 data-testid="order-detail-heading" class="mt-4 font-display text-3xl tracking-tight">Pesanan Berhasil Dibuat</h1>
        <p class="mt-2 text-sm text-slate-500">
          Nomor Pesanan: <span data-testid="order-detail-id" class="font-mono font-medium text-slate-900">{{ order.id }}</span>
        </p>
      </div>

      <div class="mt-10 grid gap-6 sm:grid-cols-2">
        <div class="rounded-2xl bg-slate-50 p-6">
          <p class="text-sm font-semibold">Dikirim ke</p>
          <p data-testid="order-detail-recipient" class="mt-2 text-sm text-slate-700">{{ order.recipient_name }}</p>
          <p class="text-sm text-slate-700">{{ order.phone }}</p>
          <p data-testid="order-detail-address" class="mt-1 text-sm text-slate-500">{{ order.address }}</p>
        </div>
        <div class="rounded-2xl bg-slate-50 p-6">
          <p class="text-sm font-semibold">Status Pesanan</p>
          <p data-testid="order-detail-status" class="mt-2 inline-block rounded-full bg-green-100 px-3 py-1 text-xs font-semibold text-green-700">
            {{ order.status }}
          </p>
          <p class="mt-3 text-sm font-semibold">Dibuat pada</p>
          <p class="text-sm text-slate-500">{{ order.created_at }}</p>
        </div>
      </div>

      <div class="mt-6 rounded-2xl border border-slate-200">
        <ul class="divide-y divide-slate-200">
          <li
            v-for="(item, index) in order.items"
            :key="item.id"
            :data-testid="`order-detail-item-${index}`"
            class="flex items-center justify-between gap-3 px-6 py-4 text-sm"
          >
            <div class="min-w-0">
              <p class="text-xs font-bold uppercase tracking-wide text-slate-500">{{ item.brand }}</p>
              <p class="truncate font-medium">{{ item.product_name }}</p>
              <p class="text-xs text-slate-500">Ukuran: {{ item.size }} · {{ item.quantity }}x</p>
            </div>
            <p class="shrink-0 font-semibold">{{ formatRupiah(item.unit_price * item.quantity) }}</p>
          </li>
        </ul>
        <div class="flex items-center justify-between border-t border-slate-200 px-6 py-4">
          <span class="text-sm font-semibold">Total</span>
          <span data-testid="order-detail-total" class="text-base font-bold">{{ formatRupiah(order.total_amount) }}</span>
        </div>
      </div>

      <div class="mt-8 flex items-center justify-center gap-4">
        <RouterLink :to="{ name: 'shop' }" class="rounded-full bg-black px-6 py-3 text-sm font-semibold text-white hover:bg-slate-800">
          Lanjutkan Belanja
        </RouterLink>
        <RouterLink :to="{ name: 'orders' }" class="text-sm font-medium text-slate-500 underline underline-offset-2">
          Lihat Riwayat Pesanan
        </RouterLink>
      </div>
    </div>
  </div>
</template>
