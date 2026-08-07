<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { assetFileUrl, getAsset } from '../lib/api'

const router = useRouter()
const auth = useAuthStore()

const heroBannerUrl = ref(null)

onMounted(async () => {
  try {
    const asset = await getAsset('hero_banner')
    heroBannerUrl.value = assetFileUrl(asset)
  } catch {
    heroBannerUrl.value = null
  }
})

function goToCta() {
  router.push({ name: auth.isAuthenticated ? 'dashboard' : 'login' })
}

const newArrivals = [
  { name: 'T-shirt with Tape Details', icon: '👕', rating: 4.5, price: '$120' },
  { name: 'Skinny Fit Jeans', icon: '👖', rating: 3.5, price: '$240', oldPrice: '$260' },
  { name: 'Checkered Shirt', icon: '🧥', rating: 4.5, price: '$180' },
  { name: 'Sleeve Striped T-shirt', icon: '👕', rating: 4.5, price: '$130', oldPrice: '$160' },
]

const topSelling = [
  { name: 'Vertical Striped Shirt', icon: '🧥', rating: 5.0, price: '$212', oldPrice: '$232' },
  { name: 'Courage Graphic T-shirt', icon: '👕', rating: 4.0, price: '$145' },
  { name: 'Loose Fit Bermuda Shorts', icon: '🩳', rating: 3.0, price: '$80' },
  { name: 'Faded Skinny Jeans', icon: '👖', rating: 4.5, price: '$210' },
]

const dressStyles = [
  { name: 'Casual', color: 'bg-cream' },
  { name: 'Formal', color: 'bg-cream' },
  { name: 'Party', color: 'bg-cream' },
  { name: 'Gym', color: 'bg-cream' },
]

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

const brands = ['VERSACE', 'ZARA', 'GUCCI', 'PRADA', 'Calvin Klein']
</script>

<template>
  <div data-testid="home-page" class="min-h-screen bg-white font-sans text-slate-900">
    <!-- Top promo bar -->
    <div data-testid="home-promo-bar" class="bg-black py-2 text-center text-xs text-white sm:text-sm">
      Daftar dan dapatkan diskon 20% untuk pembelian pertama —
      <button data-testid="home-promo-bar-cta" class="ml-1 font-semibold underline underline-offset-2" @click="goToCta">Daftar Sekarang</button>
    </div>

    <!-- Navbar -->
    <header data-testid="home-navbar" class="border-b border-slate-200">
      <div class="mx-auto flex max-w-7xl items-center gap-6 px-6 py-4">
        <div data-testid="home-navbar-logo" class="flex items-center gap-1 text-xl font-black tracking-tight">
          SHOP<span class="text-slate-400">.CO</span>
        </div>

        <nav class="hidden items-center gap-6 text-sm font-medium text-slate-600 lg:flex">
          <a href="#shop" data-testid="home-nav-link-shop" class="hover:text-black">Shop</a>
          <a href="#sale" data-testid="home-nav-link-sale" class="hover:text-black">On Sale</a>
          <a href="#new" data-testid="home-nav-link-new-arrivals" class="hover:text-black">New Arrivals</a>
          <a href="#brands" data-testid="home-nav-link-brands" class="hover:text-black">Brands</a>
        </nav>

        <div class="ml-auto flex flex-1 items-center gap-4 lg:flex-none">
          <div class="hidden flex-1 items-center rounded-full bg-slate-100 px-4 py-2 text-sm text-slate-500 sm:flex lg:w-72 lg:flex-none">
            <span class="mr-2">🔍</span>
            <input
              type="text"
              data-testid="home-search-input"
              placeholder="Cari produk..."
              class="w-full bg-transparent outline-none placeholder:text-slate-400"
            />
          </div>
          <button data-testid="home-cart-button" class="text-xl" title="Keranjang">🛒</button>
          <button
            data-testid="home-navbar-cta"
            class="rounded-full bg-black px-5 py-2 text-sm font-semibold text-white transition hover:bg-slate-800"
            @click="goToCta"
          >
            Masuk
          </button>
        </div>
      </div>
    </header>

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
            @click="goToCta"
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
        <span v-for="(brand, index) in brands" :key="brand" :data-testid="`home-brand-item-${index}`">{{ brand }}</span>
      </div>
    </section>

    <!-- New arrivals -->
    <section id="new" data-testid="home-new-arrivals-section" class="mx-auto max-w-7xl px-6 py-16">
      <h2 class="text-center font-display text-2xl tracking-tight sm:text-3xl">NEW ARRIVALS</h2>

      <div class="mt-10 grid grid-cols-2 gap-5 lg:grid-cols-4">
        <div v-for="(p, index) in newArrivals" :key="p.name" :data-testid="`home-new-arrival-card-${index}`" class="group">
          <div class="flex aspect-square items-center justify-center rounded-xl bg-cream text-5xl">
            {{ p.icon }}
          </div>
          <p data-testid="product-name" class="mt-3 text-sm font-medium">{{ p.name }}</p>
          <div class="mt-1 flex items-center gap-1 text-xs text-slate-500">
            <span class="text-amber-400">★</span>{{ p.rating }}/5
          </div>
          <div class="mt-1 flex items-center gap-2 text-sm">
            <span data-testid="product-price" class="font-semibold">{{ p.price }}</span>
            <span v-if="p.oldPrice" class="text-slate-400 line-through">{{ p.oldPrice }}</span>
          </div>
        </div>
      </div>

      <div class="mt-10 text-center">
        <button data-testid="home-new-arrivals-view-all" class="rounded-full border border-slate-300 px-6 py-2.5 text-sm font-medium transition hover:bg-slate-50">
          View All
        </button>
      </div>
    </section>

    <!-- Top selling -->
    <section id="shop" data-testid="home-top-selling-section" class="mx-auto max-w-7xl px-6 pb-16">
      <h2 class="text-center font-display text-2xl tracking-tight sm:text-3xl">TOP SELLING</h2>

      <div class="mt-10 grid grid-cols-2 gap-5 lg:grid-cols-4">
        <div v-for="(p, index) in topSelling" :key="p.name" :data-testid="`home-top-selling-card-${index}`" class="group">
          <div class="flex aspect-square items-center justify-center rounded-xl bg-cream text-5xl">
            {{ p.icon }}
          </div>
          <p data-testid="product-name" class="mt-3 text-sm font-medium">{{ p.name }}</p>
          <div class="mt-1 flex items-center gap-1 text-xs text-slate-500">
            <span class="text-amber-400">★</span>{{ p.rating }}/5
          </div>
          <div class="mt-1 flex items-center gap-2 text-sm">
            <span data-testid="product-price" class="font-semibold">{{ p.price }}</span>
            <span v-if="p.oldPrice" class="text-slate-400 line-through">{{ p.oldPrice }}</span>
          </div>
        </div>
      </div>

      <div class="mt-10 text-center">
        <button data-testid="home-top-selling-view-all" class="rounded-full border border-slate-300 px-6 py-2.5 text-sm font-medium transition hover:bg-slate-50">
          View All
        </button>
      </div>
    </section>

    <!-- Browse by dress style -->
    <section id="sale" data-testid="home-dress-style-section" class="mx-auto max-w-7xl px-6 pb-16">
      <div class="rounded-2xl bg-cream p-8">
        <h2 class="text-center font-display text-2xl tracking-tight sm:text-3xl">BROWSE BY DRESS STYLE</h2>

        <div class="mt-8 grid grid-cols-2 gap-4">
          <div
            v-for="(style, index) in dressStyles"
            :key="style.name"
            :data-testid="`home-dress-style-card-${index}`"
            class="flex h-32 items-center justify-center rounded-xl bg-white text-base font-semibold shadow-sm transition hover:shadow-md cursor-pointer sm:h-40"
          >
            {{ style.name }}
          </div>
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
