<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

function handleLogout() {
  auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <div data-testid="dashboard-page" class="min-h-screen bg-slate-950 text-white">
    <header class="flex items-center justify-between border-b border-slate-800 px-6 py-4">
      <div class="flex items-center gap-1 text-lg font-bold">
        SHOP<span class="text-slate-500">.CO</span>
      </div>
      <button
        data-testid="dashboard-logout-button"
        class="rounded-lg bg-slate-900 px-3.5 py-2 text-sm font-medium text-slate-200 ring-1 ring-inset ring-slate-800 hover:bg-slate-800"
        @click="handleLogout"
      >
        Keluar
      </button>
    </header>

    <main class="mx-auto max-w-2xl px-6 py-16 text-center">
      <p class="text-sm text-brand-400">Login berhasil</p>
      <h1 data-testid="dashboard-welcome-heading" class="mt-2 text-2xl font-semibold">
        Selamat datang, {{ auth.user?.email }}
      </h1>
      <p class="mt-3 text-slate-400">
        Halaman dashboard akun kamu akan dibangun di sini.
      </p>
    </main>
  </div>
</template>
