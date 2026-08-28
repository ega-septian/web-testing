<script setup>
import { fileUrl, formatRupiah, originalPrice } from '../lib/api'

defineProps({
  product: { type: Object, required: true },
  testid: { type: String, default: undefined },
})
</script>

<template>
  <RouterLink :to="{ name: 'product-detail', params: { id: product.id } }" :data-testid="testid" class="group block">
    <div class="flex aspect-square items-center justify-center overflow-hidden rounded-xl bg-cream">
      <img v-if="product.image_url" :src="fileUrl(product.image_url)" :alt="product.name" class="h-full w-full object-cover" />
      <span v-else class="text-4xl text-slate-300">🖼️</span>
    </div>
    <p data-testid="product-brand" class="mt-3 text-xs font-bold uppercase tracking-wide">{{ product.brand }}</p>
    <p data-testid="product-name" class="mt-1 text-sm font-medium text-slate-600">{{ product.name }}</p>
    <div class="mt-1 flex flex-wrap items-center gap-2 text-sm">
      <span data-testid="product-price" class="font-semibold" :class="product.discount ? 'text-red-500' : 'text-slate-900'">
        {{ formatRupiah(product.price) }}
      </span>
      <span v-if="product.discount" class="text-slate-400 line-through">{{ formatRupiah(originalPrice(product.price, product.discount)) }}</span>
      <span
        v-if="product.discount"
        class="rounded-full bg-red-50 px-2 py-0.5 text-xs font-semibold text-red-500"
      >
        -{{ product.discount }}%
      </span>
    </div>
  </RouterLink>
</template>
