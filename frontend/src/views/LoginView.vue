<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const showPassword = ref(false)

const highlights = [
  { name: 'New Arrivals', desc: 'Koleksi terbaru tiap minggu', icon: '👕' },
  { name: 'Top Selling', desc: 'Favorit pelanggan kami', icon: '🔥' },
  { name: 'Dress Style', desc: 'Casual, formal, party, gym', icon: '👗' },
]

async function handleSubmit() {
  const ok = await auth.login({ email: email.value, password: password.value })
  if (ok) {
    // Sent here by the router guard (e.g. from checkout) with where to go
    // back to — falls back to the homepage otherwise.
    router.push(route.query.redirect || { name: 'home' })
  }
}
</script>

<template>
  <div data-testid="login-page" class="min-h-screen grid bg-white font-sans text-slate-900 lg:grid-cols-2">
    <!-- Store showcase panel -->
    <div class="relative hidden flex-col justify-between overflow-hidden bg-cream p-12 lg:flex">
      <div>
        <div class="flex items-center gap-1 text-2xl font-black tracking-tight">
          SHOP<span class="text-slate-400">.CO</span>
        </div>
        <p class="mt-2 text-sm text-slate-500">Belanja fashion favorit kamu, kapan saja.</p>
      </div>

      <div class="space-y-6">
        <h1 class="font-display text-4xl leading-[1.05] tracking-tight">
          TEMUKAN GAYA YANG MENCERMINKAN DIRIMU
        </h1>
        <p class="max-w-md text-slate-600">
          Masuk untuk menjelajahi koleksi terbaru, penawaran eksklusif, dan rekomendasi outfit pilihan.
        </p>

        <div data-testid="login-highlights-list" class="space-y-3 pt-2">
          <div
            v-for="(h, index) in highlights"
            :key="h.name"
            :data-testid="`login-highlight-item-${index}`"
            class="flex items-center gap-4 rounded-xl bg-white p-4 shadow-sm ring-1 ring-slate-200/70 transition hover:shadow-md"
          >
            <span class="text-2xl">{{ h.icon }}</span>
            <div>
              <p class="font-medium">{{ h.name }}</p>
              <p class="text-sm text-slate-500">{{ h.desc }}</p>
            </div>
          </div>
        </div>
      </div>

      <p class="text-xs text-slate-400">© 2026 Shop.co. Semua hak dilindungi.</p>
    </div>

    <!-- Login form panel -->
    <div class="flex items-center justify-center p-6 sm:p-12">
      <div class="w-full max-w-sm">
        <div class="mb-8 flex items-center gap-1 text-xl font-black tracking-tight lg:hidden">
          SHOP<span class="text-slate-400">.CO</span>
        </div>

        <h2 data-testid="login-heading" class="font-display text-2xl tracking-tight">MASUK KE AKUN KAMU</h2>
        <p class="mt-2 text-sm text-slate-500">
          Belum punya akun?
          <a href="#" data-testid="login-register-link" class="font-semibold text-black underline underline-offset-2">Daftar gratis</a>
        </p>

        <form data-testid="login-form" class="mt-8 space-y-5" @submit.prevent="handleSubmit">
          <div>
            <label for="email" class="block text-sm font-medium text-slate-700">Email</label>
            <input
              id="email"
              data-testid="login-email-input"
              v-model="email"
              type="email"
              autocomplete="email"
              required
              placeholder="nama@perusahaan.com"
              class="mt-1.5 block w-full rounded-full border border-slate-300 bg-white px-4 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 outline-none transition focus:border-black"
            />
          </div>

          <div>
            <div class="flex items-center justify-between">
              <label for="password" class="block text-sm font-medium text-slate-700">Password</label>
              <a href="#" data-testid="login-forgot-password-link" class="text-sm font-medium text-slate-500 hover:text-black">Lupa password?</a>
            </div>
            <div class="relative mt-1.5">
              <input
                id="password"
                data-testid="login-password-input"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                placeholder="••••••••"
                class="block w-full rounded-full border border-slate-300 bg-white px-4 py-2.5 pr-11 text-sm text-slate-900 placeholder:text-slate-400 outline-none transition focus:border-black"
              />
              <button
                type="button"
                data-testid="login-toggle-password-visibility"
                class="absolute inset-y-0 right-1 flex items-center px-3 text-slate-400 hover:text-slate-700"
                @click="showPassword = !showPassword"
              >
                {{ showPassword ? '🙈' : '👁️' }}
              </button>
            </div>
          </div>

          <div v-if="auth.error" data-testid="login-error-message" class="rounded-lg bg-red-50 px-3.5 py-2.5 text-sm text-red-600 ring-1 ring-inset ring-red-200">
            {{ auth.error }}
          </div>

          <button
            type="submit"
            data-testid="login-submit-button"
            :disabled="auth.loading"
            class="flex w-full items-center justify-center gap-2 rounded-full bg-black px-3.5 py-2.5 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <span v-if="auth.loading" class="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white"></span>
            {{ auth.loading ? 'Memproses...' : 'Masuk' }}
          </button>
        </form>

        <div class="mt-6 flex items-center gap-3">
          <div class="h-px flex-1 bg-slate-200"></div>
          <span class="text-xs text-slate-400">atau lanjutkan dengan</span>
          <div class="h-px flex-1 bg-slate-200"></div>
        </div>

        <button
          type="button"
          data-testid="login-google-button"
          class="mt-6 flex w-full items-center justify-center gap-2 rounded-full border border-slate-300 px-3.5 py-2.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50"
        >
          <span>🔵</span> Masuk dengan Google
        </button>
      </div>
    </div>
  </div>
</template>
