<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { assetFileUrl, fileUrl, getAssets, getProducts, getProductFilters } from '../lib/api'
import ProductCard from '../components/ProductCard.vue'

const router = useRouter()
const auth = useAuthStore()

const heroBannerUrl = ref(null)
const newArrivals = ref([])
const topSelling = ref([])
const categories = ref([])

onMounted(async () => {
  try {
    const assets = await getAssets()
    const assetsByKey = Object.fromEntries(assets.map((asset) => [asset.key, asset]))
    heroBannerUrl.value = assetsByKey.hero_banner ? assetFileUrl(assetsByKey.hero_banner) : null
  } catch {
    heroBannerUrl.value = null
  }

  // Each catalog section fetches independently, so one failing (or the
  // backend being briefly unavailable) doesn't take the whole page down —
  // the section just renders empty instead of throwing.
  try {
    newArrivals.value = await getProducts({ sort: 'newest' })
  } catch {
    newArrivals.value = []
  }

  try {
    topSelling.value = await getProducts({ sort: 'best_selling' })
  } catch {
    topSelling.value = []
  }

  try {
    categories.value = await loadCategories()
  } catch {
    categories.value = []
  }
})

// "Shop by Category" replaces the old "Browse by Dress Style" section, which
// was disconnected from the real catalog (Casual/Formal/Party/Gym never
// matched any actual product's category) and its cards weren't even
// clickable. This is driven by the live category facet instead, so it's
// always in sync with whatever products actually exist — each card's photo
// is just the newest product in that category.
async function loadCategories() {
  const { category: options } = await getProductFilters()
  return Promise.all(
    options.map(async (opt) => {
      const [product] = await getProducts({ category: opt.value, sort: 'newest', limit: 1 })
      return { name: opt.value, count: opt.count, imageUrl: product?.image_url ? fileUrl(product.image_url) : null }
    }),
  )
}

// There's no separate authenticated page to send an already-logged-in user
// to anymore — the site itself reflects the logged-in state (see Navbar.vue),
// so this CTA has nowhere left to go once authenticated.
function goToCta() {
  if (auth.isAuthenticated) return
  router.push({ name: 'login' })
}

// "Shop Now" browses the catalog regardless of auth state — unlike the CTAs
// above, it was never really about login.
function goToShop() {
  router.push({ name: 'shop' })
}

// Clicking a brand in the strip below jumps straight to the Shop page
// pre-filtered to that brand (see ShopView.vue's brand checklist filter).
function goToBrand(brand) {
  router.push({ name: 'shop', query: { brand } })
}

// Same idea for a category card (see "Shop by Category" below).
function goToCategory(category) {
  router.push({ name: 'shop', query: { category } })
}

const testimonials = [
  {
    name: 'Sarah M.',
    quote:
      '"Belanja di sini jadi pengalaman yang menyenangkan — kualitas bahannya bagus dan potongannya selalu pas di badan."',
  },
  {
    name: 'Alex K.',
    quote:
      '"Pilihan gayanya lengkap banget, dari kasual sampai formal. Gampang nemuin outfit yang cocok buat acara apa pun."',
  },
  {
    name: 'James L.',
    quote:
      '"Desainnya selalu update sama tren terbaru tapi tetap nyaman dipakai sehari-hari. Jadi langganan tetap saya."',
  },
]

const brands = ['NEVADA', 'DISNEY', 'MARVEL', 'COLE', 'SUKO']
</script>

<template>
  <div data-testid="home-page" class="min-h-screen bg-white font-sans text-slate-900">
    <!-- Top promo bar: only makes sense as an acquisition message for anonymous visitors -->
    <div v-if="!auth.isAuthenticated" data-testid="home-promo-bar" class="bg-black py-2 text-center text-xs text-white sm:text-sm">
      Daftar dan dapatkan diskon 20% untuk pembelian pertama —
      <button data-testid="home-promo-bar-cta" class="ml-1 font-semibold underline underline-offset-2" @click="goToCta">Daftar Sekarang</button>
    </div>

    <!-- Hero -->
    <section data-testid="home-hero-section" class="relative overflow-hidden bg-cream">
      <div class="mx-auto max-w-7xl px-6 py-16 lg:py-24">
        <div class="relative z-10 max-w-xl">
          <h1 data-testid="home-hero-heading" class="font-display text-4xl leading-[1.05] tracking-tight sm:text-5xl lg:text-6xl">
            FIND CLOTHES THAT MATCHES YOUR STYLE
          </h1>
          <p class="mt-6 max-w-md text-slate-600">
            Jelajahi koleksi busana pilihan yang dirancang dengan detail, dipilih untuk mencerminkan gaya
            dan kepribadian kamu.
          </p>
          <button
            data-testid="home-hero-cta"
            class="mt-8 rounded-full bg-black px-8 py-3.5 text-sm font-semibold text-white transition hover:bg-slate-800"
            @click="goToShop"
          >
            Shop Now
          </button>

          <div data-testid="home-hero-stats" class="mt-10 flex flex-wrap gap-6 border-t border-slate-300/70 pt-6 text-sm">
            <div>
              <p class="text-2xl font-black">200+</p>
              <p class="text-slate-500">Brand internasional</p>
            </div>
            <div class="border-l border-slate-300 pl-6">
              <p class="text-2xl font-black">2,000+</p>
              <p class="text-slate-500">Produk berkualitas</p>
            </div>
            <div class="border-l border-slate-300 pl-6">
              <p class="text-2xl font-black">30,000+</p>
              <p class="text-slate-500">Pelanggan puas</p>
            </div>
          </div>
        </div>

        <!-- Mobile/tablet: banner sits below the text, in normal flow -->
        <div data-testid="home-hero-banner-mobile" class="relative mt-10 aspect-[4/3] w-full overflow-hidden rounded-2xl lg:hidden">
          <img
            v-if="heroBannerUrl"
            :src="heroBannerUrl"
            alt="Hero banner"
            class="h-full w-full object-cover object-[80%_15%]"
          />
          <div v-else class="flex h-full items-center justify-center bg-gradient-to-br from-slate-200 to-slate-100">
            <span class="text-8xl">🧥👗</span>
          </div>
        </div>
      </div>

      <!-- Desktop: banner fills the section as a full-bleed right-anchored background -->
      <div data-testid="home-hero-banner" class="pointer-events-none absolute inset-y-0 right-0 hidden w-1/2 lg:block">
        <img
          v-if="heroBannerUrl"
          :src="heroBannerUrl"
          alt="Hero banner"
          data-testid="home-hero-banner-image"
          class="h-full w-full object-cover object-right-bottom"
        />
        <div v-else class="flex h-full items-center justify-center bg-gradient-to-br from-slate-200 to-slate-100">
          <span class="text-8xl">🧥👗</span>
        </div>
        <span class="absolute right-10 top-10 text-3xl">✦</span>
        <span class="absolute bottom-16 left-10 text-2xl">✦</span>
      </div>
    </section>

    <!-- Brand strip -->
    <section data-testid="home-brand-strip" class="bg-black py-6">
      <div class="mx-auto flex max-w-7xl flex-wrap items-center justify-center gap-x-10 gap-y-3 px-6 text-lg font-bold tracking-wide text-white sm:justify-between">
        <button
          v-for="(brand, index) in brands"
          :key="brand"
          type="button"
          :data-testid="`home-brand-item-${index}`"
          class="transition hover:text-slate-300"
          @click="goToBrand(brand)"
        >
          {{ brand }}
        </button>
      </div>
    </section>

    <!-- New arrivals -->
    <section id="new" data-testid="home-new-arrivals-section" class="mx-auto max-w-7xl px-6 py-16">
      <h2 class="text-center font-display text-2xl tracking-tight sm:text-3xl">NEW ARRIVALS</h2>

      <div class="mt-10 grid grid-cols-2 gap-5 lg:grid-cols-4">
        <ProductCard
          v-for="(p, index) in newArrivals"
          :key="p.id"
          :product="p"
          :testid="`home-new-arrival-card-${index}`"
        />
      </div>

      <div class="mt-10 text-center">
        <RouterLink :to="{ name: 'shop' }" data-testid="home-new-arrivals-view-all" class="rounded-full border border-slate-300 px-6 py-2.5 text-sm font-medium transition hover:bg-slate-50">
          View All
        </RouterLink>
      </div>
    </section>

    <!-- Top selling -->
    <section id="top-selling" data-testid="home-top-selling-section" class="mx-auto max-w-7xl px-6 pb-16">
      <h2 class="text-center font-display text-2xl tracking-tight sm:text-3xl">TOP SELLING</h2>

      <div class="mt-10 grid grid-cols-2 gap-5 lg:grid-cols-4">
        <ProductCard
          v-for="(p, index) in topSelling"
          :key="p.id"
          :product="p"
          :testid="`home-top-selling-card-${index}`"
        />
      </div>

      <div class="mt-10 text-center">
        <RouterLink :to="{ name: 'shop' }" data-testid="home-top-selling-view-all" class="rounded-full border border-slate-300 px-6 py-2.5 text-sm font-medium transition hover:bg-slate-50">
          View All
        </RouterLink>
      </div>
    </section>

    <!-- Shop by category -->
    <section id="sale" data-testid="home-category-section" class="mx-auto max-w-7xl px-6 pb-16">
      <div class="rounded-2xl bg-cream p-8">
        <h2 class="text-center font-display text-2xl tracking-tight sm:text-3xl">SHOP BY CATEGORY</h2>

        <div class="mt-8 grid grid-cols-2 gap-4">
          <button
            v-for="(cat, index) in categories"
            :key="cat.name"
            type="button"
            :data-testid="`home-category-card-${index}`"
            class="relative flex h-32 overflow-hidden rounded-xl bg-white text-left shadow-sm transition hover:shadow-md sm:h-40"
            :class="cat.imageUrl ? 'items-end' : 'items-center justify-center'"
            @click="goToCategory(cat.name)"
          >
            <img
              v-if="cat.imageUrl"
              :src="cat.imageUrl"
              :alt="cat.name"
              class="absolute inset-0 h-full w-full object-cover object-top"
            />
            <span
              class="relative z-10 text-base font-semibold"
              :class="cat.imageUrl ? 'w-full bg-gradient-to-t from-black/60 to-transparent px-4 py-3 text-white' : ''"
            >
              {{ cat.name }} <span class="font-normal opacity-80">({{ cat.count }})</span>
            </span>
          </button>
        </div>
      </div>
    </section>

    <!-- Testimonials -->
    <section id="brands" data-testid="home-testimonials-section" class="mx-auto max-w-7xl px-6 pb-16">
      <h2 class="font-display text-2xl tracking-tight sm:text-3xl">OUR HAPPY CUSTOMERS</h2>

      <div class="mt-8 grid gap-5 md:grid-cols-3">
        <div v-for="(t, index) in testimonials" :key="t.name" :data-testid="`home-testimonial-card-${index}`" class="rounded-xl border border-slate-200 p-6">
          <div class="text-amber-400">★★★★★</div>
          <p class="mt-3 text-sm text-slate-600">{{ t.quote }}</p>
          <p class="mt-4 text-sm font-semibold">{{ t.name }}</p>
        </div>
      </div>
    </section>

    <!-- Newsletter CTA -->
    <section data-testid="home-newsletter-section" class="mx-auto max-w-7xl px-6 pb-16">
      <div class="flex flex-col items-center justify-between gap-6 rounded-2xl bg-black px-8 py-10 sm:flex-row">
        <h3 class="max-w-sm text-center font-display text-xl leading-tight text-white sm:text-left sm:text-2xl">
          STAY UPTO DATE ABOUT OUR LATEST OFFERS
        </h3>
        <form data-testid="home-newsletter-form" class="flex w-full max-w-sm flex-col gap-3 sm:flex-row" @submit.prevent>
          <input
            type="email"
            data-testid="home-newsletter-email-input"
            placeholder="Masukkan email kamu"
            class="w-full rounded-full bg-white px-4 py-2.5 text-sm outline-none"
          />
          <button type="submit" data-testid="home-newsletter-submit" class="shrink-0 rounded-full bg-white px-5 py-2.5 text-sm font-semibold hover:bg-slate-100">
            Subscribe
          </button>
        </form>
      </div>
    </section>

    <!-- Footer -->
    <footer data-testid="home-footer" class="border-t border-slate-200">
      <div class="mx-auto grid max-w-7xl gap-10 px-6 py-12 sm:grid-cols-2 lg:grid-cols-5">
        <div class="lg:col-span-2">
          <div class="text-lg font-black">SHOP<span class="text-slate-400">.CO</span></div>
          <p class="mt-3 max-w-xs text-sm text-slate-500">
            Belanja busana dari brand pilihan, dirancang biar kamu tetap percaya diri dengan gaya sendiri.
          </p>
        </div>

        <div>
          <p class="text-sm font-semibold">Company</p>
          <ul class="mt-3 space-y-2 text-sm text-slate-500">
            <li>About</li>
            <li>Features</li>
            <li>Works</li>
            <li>Career</li>
          </ul>
        </div>

        <div>
          <p class="text-sm font-semibold">Help</p>
          <ul class="mt-3 space-y-2 text-sm text-slate-500">
            <li>Customer Support</li>
            <li>Delivery Details</li>
            <li>Terms &amp; Conditions</li>
            <li>Privacy Policy</li>
          </ul>
        </div>

        <div>
          <p class="text-sm font-semibold">FAQ</p>
          <ul class="mt-3 space-y-2 text-sm text-slate-500">
            <li>Account</li>
            <li>Manage Deliveries</li>
            <li>Orders</li>
            <li>Payments</li>
          </ul>
        </div>
      </div>
      <div class="border-t border-slate-200 py-6 text-center text-xs text-slate-400">
        Shop.co © 2026. Semua hak dilindungi.
      </div>
    </footer>
  </div>
</template>
