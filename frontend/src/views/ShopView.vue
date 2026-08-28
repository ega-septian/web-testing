<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getProducts, getProductFilters } from '../lib/api'
import ProductCard from '../components/ProductCard.vue'

const SHOP_PAGE_SIZE = 24

const route = useRoute()
const router = useRouter()

const filterOptions = ref({ gender: [], category: [], subcategory: [], size: [] })
// Each dimension is a checklist — an array of currently-checked values, not
// a single choice.
const selected = ref({ gender: [], category: [], subcategory: [], size: [] })
const products = ref([])
const loading = ref(true)

// Seeded from ?q= so arriving from the homepage search box lands here
// already searching — kept in sync both ways with the URL below.
const searchQuery = ref(typeof route.query.q === 'string' ? route.query.q : '')

onMounted(async () => {
  try {
    filterOptions.value = await getProductFilters()
  } catch {
    filterOptions.value = { gender: [], category: [], subcategory: [], size: [] }
  }
  await loadProducts()
})

// Navigating here again from the homepage search (while already on the Shop
// page) reuses this component instance instead of remounting it, so the
// query-string change needs its own watcher rather than relying on onMounted.
watch(
  () => route.query.q,
  (newQ) => {
    searchQuery.value = typeof newQ === 'string' ? newQ : ''
    loadProducts()
  },
)

async function loadProducts() {
  loading.value = true
  try {
    products.value = await getProducts({ ...selected.value, q: searchQuery.value, limit: SHOP_PAGE_SIZE })
  } catch {
    products.value = []
  } finally {
    loading.value = false
  }
}

// Updates the URL's ?q= (so the search term is shareable/refreshable) —
// the watcher above reacts to that change and actually reloads the results.
function submitSearch() {
  router.replace({ query: { ...route.query, q: searchQuery.value.trim() || undefined } })
}

function clearSearch() {
  searchQuery.value = ''
  submitSearch()
}

function toggleFilter(dimension, value) {
  const list = selected.value[dimension]
  const i = list.indexOf(value)
  if (i === -1) {
    list.push(value)
  } else {
    list.splice(i, 1)
  }
  loadProducts()
}

function clearFilters() {
  selected.value = { gender: [], category: [], subcategory: [], size: [] }
  loadProducts()
}

const hasActiveFilters = () => Object.values(selected.value).some((list) => list.length > 0)
</script>

<template>
  <div data-testid="shop-page" class="min-h-screen bg-white font-sans text-slate-900">
    <div class="mx-auto max-w-7xl px-6 py-10">
      <h1 data-testid="shop-heading" class="font-display text-3xl tracking-tight">Shop</h1>

      <div class="mt-6 flex items-center gap-3">
        <div class="flex flex-1 items-center rounded-full bg-slate-100 px-4 py-2 text-sm text-slate-500 sm:max-w-sm">
          <button type="button" data-testid="shop-search-submit" class="mr-2" title="Cari" @click="submitSearch">🔍</button>
          <input
            v-model="searchQuery"
            type="text"
            data-testid="shop-search-input"
            placeholder="Cari produk..."
            class="w-full bg-transparent outline-none placeholder:text-slate-400"
            @keyup.enter="submitSearch"
          />
        </div>
        <p v-if="route.query.q" data-testid="shop-search-summary" class="text-sm text-slate-500">
          Hasil pencarian untuk "<span class="font-medium text-slate-900">{{ route.query.q }}</span>"
          <button data-testid="shop-search-clear" class="ml-1 underline underline-offset-2 hover:text-black" @click="clearSearch">
            Hapus
          </button>
        </p>
      </div>

      <div class="mt-8 grid gap-10 lg:grid-cols-[240px_1fr]">
        <!-- Filter sidebar -->
        <aside data-testid="shop-filter-sidebar" class="space-y-8">
          <div v-if="filterOptions.gender.length" data-testid="shop-filter-gender">
            <p class="text-sm font-semibold">Gender</p>
            <ul class="mt-3 space-y-2">
              <li v-for="opt in filterOptions.gender" :key="opt.value">
                <label class="flex cursor-pointer items-center justify-between gap-2 text-sm text-slate-600 hover:text-black">
                  <span class="flex items-center gap-2">
                    <input
                      type="checkbox"
                      :data-testid="`shop-filter-gender-${opt.value}`"
                      class="h-4 w-4 rounded border-slate-300 accent-black"
                      :checked="selected.gender.includes(opt.value)"
                      @change="toggleFilter('gender', opt.value)"
                    />
                    {{ opt.value }}
                  </span>
                  <span class="text-slate-400">({{ opt.count }})</span>
                </label>
              </li>
            </ul>
          </div>

          <div v-if="filterOptions.category.length" data-testid="shop-filter-category" class="border-t border-slate-200 pt-6">
            <p class="text-sm font-semibold">Kategori</p>
            <ul class="mt-3 space-y-2">
              <li v-for="opt in filterOptions.category" :key="opt.value">
                <label class="flex cursor-pointer items-center justify-between gap-2 text-sm text-slate-600 hover:text-black">
                  <span class="flex items-center gap-2">
                    <input
                      type="checkbox"
                      :data-testid="`shop-filter-category-${opt.value}`"
                      class="h-4 w-4 rounded border-slate-300 accent-black"
                      :checked="selected.category.includes(opt.value)"
                      @change="toggleFilter('category', opt.value)"
                    />
                    {{ opt.value }}
                  </span>
                  <span class="text-slate-400">({{ opt.count }})</span>
                </label>
              </li>
            </ul>
          </div>

          <div v-if="filterOptions.subcategory.length" data-testid="shop-filter-subcategory" class="border-t border-slate-200 pt-6">
            <p class="text-sm font-semibold">Sub Kategori</p>
            <ul class="mt-3 space-y-2">
              <li v-for="opt in filterOptions.subcategory" :key="opt.value">
                <label class="flex cursor-pointer items-center justify-between gap-2 text-sm text-slate-600 hover:text-black">
                  <span class="flex items-center gap-2">
                    <input
                      type="checkbox"
                      :data-testid="`shop-filter-subcategory-${opt.value}`"
                      class="h-4 w-4 rounded border-slate-300 accent-black"
                      :checked="selected.subcategory.includes(opt.value)"
                      @change="toggleFilter('subcategory', opt.value)"
                    />
                    {{ opt.value }}
                  </span>
                  <span class="text-slate-400">({{ opt.count }})</span>
                </label>
              </li>
            </ul>
          </div>

          <div v-if="filterOptions.size.length" data-testid="shop-filter-size" class="border-t border-slate-200 pt-6">
            <p class="text-sm font-semibold">Ukuran</p>
            <ul class="mt-3 space-y-2">
              <li v-for="opt in filterOptions.size" :key="opt.value">
                <label class="flex cursor-pointer items-center justify-between gap-2 text-sm text-slate-600 hover:text-black">
                  <span class="flex items-center gap-2">
                    <input
                      type="checkbox"
                      :data-testid="`shop-filter-size-${opt.value}`"
                      class="h-4 w-4 rounded border-slate-300 accent-black"
                      :checked="selected.size.includes(opt.value)"
                      @change="toggleFilter('size', opt.value)"
                    />
                    {{ opt.value }}
                  </span>
                  <span class="text-slate-400">({{ opt.count }})</span>
                </label>
              </li>
            </ul>
          </div>

          <button
            v-if="hasActiveFilters()"
            data-testid="shop-filter-clear"
            class="text-sm font-medium text-slate-500 underline underline-offset-2 hover:text-black"
            @click="clearFilters"
          >
            Hapus semua filter
          </button>
        </aside>

        <!-- Product grid -->
        <div>
          <div v-if="loading" class="py-16 text-center text-slate-400">Memuat...</div>
          <div v-else-if="products.length === 0" data-testid="shop-empty-state" class="py-16 text-center text-slate-400">
            Tidak ada produk yang cocok dengan filter ini.
          </div>
          <div v-else class="grid grid-cols-2 gap-5 sm:grid-cols-3 lg:grid-cols-4">
            <ProductCard
              v-for="(p, index) in products"
              :key="p.id"
              :product="p"
              :testid="`shop-product-card-${index}`"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
