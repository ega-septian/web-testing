<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'

const router = useRouter()
const auth = useAuthStore()
const cart = useCartStore()

const searchQuery = ref('')

// Below the `lg` breakpoint the nav links and search box collapse into this
// hamburger-triggered panel (see template) instead of disappearing outright.
const isMobileMenuOpen = ref(false)

function closeMobileMenu() {
  isMobileMenuOpen.value = false
}

// There's no separate authenticated page to send an already-logged-in user
// to — the whole site already reflects login state (this navbar included),
// so this CTA has nowhere left to go once authenticated.
function goToCta() {
  if (auth.isAuthenticated) return
  router.push({ name: 'login' })
}

// Sends the user to the Shop page with their term applied — the Shop page
// itself does the actual filtering/searching (see ShopView.vue), which also
// has its own search box for refining once already there.
function submitSearch() {
  const q = searchQuery.value.trim()
  if (!q) return
  router.push({ name: 'shop', query: { q } })
  searchQuery.value = ''
  closeMobileMenu()
}

// The cart is intentionally NOT cleared here — it's per-browser, not tied to
// the account (see stores/cart.js), so it should survive a logout the same
// way it survives a refresh.
function handleLogout() {
  auth.logout()
  router.push({ name: 'home' })
}
</script>

<template>
  <header data-testid="navbar" class="border-b border-slate-200">
    <div class="mx-auto flex max-w-7xl items-center gap-6 px-6 py-4">
      <RouterLink :to="{ name: 'home' }" data-testid="navbar-logo" class="flex items-center gap-1 text-xl font-black tracking-tight">
        SHOP<span class="text-slate-400">.CO</span>
      </RouterLink>

      <nav class="hidden items-center gap-6 text-sm font-medium text-slate-600 lg:flex">
        <RouterLink :to="{ name: 'shop' }" data-testid="navbar-link-shop" class="hover:text-black">Shop</RouterLink>
        <!-- These sections only exist on the homepage — routing there with a
             hash (plus the router's scrollBehavior) gets you to them from
             any page, not just when already on Home. -->
        <RouterLink :to="{ name: 'home', hash: '#sale' }" data-testid="navbar-link-sale" class="hover:text-black">On Sale</RouterLink>
        <RouterLink :to="{ name: 'home', hash: '#new' }" data-testid="navbar-link-new-arrivals" class="hover:text-black">New Arrivals</RouterLink>
        <RouterLink :to="{ name: 'home', hash: '#brands' }" data-testid="navbar-link-brands" class="hover:text-black">Brands</RouterLink>
      </nav>

      <div class="ml-auto flex flex-1 items-center gap-4 lg:flex-none">
        <div class="hidden flex-1 items-center rounded-full bg-slate-100 px-4 py-2 text-sm text-slate-500 sm:flex lg:w-72 lg:flex-none">
          <button type="button" data-testid="navbar-search-submit" class="mr-2" title="Cari" @click="submitSearch">🔍</button>
          <input
            v-model="searchQuery"
            type="text"
            data-testid="navbar-search-input"
            placeholder="Cari produk..."
            class="w-full bg-transparent outline-none placeholder:text-slate-400"
            @keyup.enter="submitSearch"
          />
        </div>

        <RouterLink :to="{ name: 'cart' }" data-testid="navbar-cart-button" class="relative text-xl" title="Keranjang">
          🛒
          <span
            v-if="cart.totalCount > 0"
            data-testid="navbar-cart-badge"
            class="absolute -right-2 -top-2 flex h-4 min-w-4 items-center justify-center rounded-full bg-black px-1 text-[10px] font-semibold text-white"
          >
            {{ cart.totalCount }}
          </span>
        </RouterLink>

        <button
          v-if="!auth.isAuthenticated"
          data-testid="navbar-cta"
          class="rounded-full bg-black px-5 py-2 text-sm font-semibold text-white transition hover:bg-slate-800"
          @click="goToCta"
        >
          Masuk
        </button>
        <!-- No dedicated account/profile page yet — order history is the one
             real destination behind this icon so far. -->
        <template v-else>
          <RouterLink :to="{ name: 'orders' }" data-testid="navbar-account-button" class="text-xl" title="Riwayat Pesanan">👤</RouterLink>
          <button
            data-testid="navbar-logout-button"
            title="Keluar"
            class="text-sm font-medium text-slate-500 hover:text-black"
            @click="handleLogout"
          >
            Keluar
          </button>
        </template>

        <!-- Toggles the collapsible panel below that carries the nav links
             and search box once they no longer fit in the header itself. -->
        <button
          type="button"
          data-testid="navbar-mobile-menu-toggle"
          class="text-xl lg:hidden"
          :aria-expanded="isMobileMenuOpen"
          title="Menu"
          @click="isMobileMenuOpen = !isMobileMenuOpen"
        >
          {{ isMobileMenuOpen ? '✕' : '☰' }}
        </button>
      </div>
    </div>

    <div v-if="isMobileMenuOpen" data-testid="navbar-mobile-menu" class="border-t border-slate-200 px-6 py-4 lg:hidden">
      <!-- Only shown below `sm`: from `sm` up, the header's own search box
           (just above) is already visible, so this would be a duplicate. -->
      <div class="mb-4 flex items-center rounded-full bg-slate-100 px-4 py-2 text-sm text-slate-500 sm:hidden">
        <button type="button" data-testid="navbar-mobile-search-submit" class="mr-2" title="Cari" @click="submitSearch">🔍</button>
        <input
          v-model="searchQuery"
          type="text"
          data-testid="navbar-mobile-search-input"
          placeholder="Cari produk..."
          class="w-full bg-transparent outline-none placeholder:text-slate-400"
          @keyup.enter="submitSearch"
        />
      </div>

      <nav class="flex flex-col gap-4 text-sm font-medium text-slate-600">
        <RouterLink :to="{ name: 'shop' }" data-testid="navbar-mobile-link-shop" class="hover:text-black" @click="closeMobileMenu">Shop</RouterLink>
        <RouterLink :to="{ name: 'home', hash: '#sale' }" data-testid="navbar-mobile-link-sale" class="hover:text-black" @click="closeMobileMenu">On Sale</RouterLink>
        <RouterLink :to="{ name: 'home', hash: '#new' }" data-testid="navbar-mobile-link-new-arrivals" class="hover:text-black" @click="closeMobileMenu">New Arrivals</RouterLink>
        <RouterLink :to="{ name: 'home', hash: '#brands' }" data-testid="navbar-mobile-link-brands" class="hover:text-black" @click="closeMobileMenu">Brands</RouterLink>
      </nav>
    </div>
  </header>
</template>
