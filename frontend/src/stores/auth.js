import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as api from '../lib/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('auth_token') || null)
  const user = ref(JSON.parse(localStorage.getItem('auth_user') || 'null'))
  const error = ref(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  function persist(data) {
    token.value = data.token
    user.value = data.user
    localStorage.setItem('auth_token', data.token)
    localStorage.setItem('auth_user', JSON.stringify(data.user))
  }

  async function login({ email, password }) {
    loading.value = true
    error.value = null

    try {
      const data = await api.login({ email, password })
      persist(data)
      return true
    } catch (err) {
      error.value = err.message || 'Login gagal'
      return false
    } finally {
      loading.value = false
    }
  }

  async function register({ email, password }) {
    loading.value = true
    error.value = null

    try {
      const data = await api.register({ email, password })
      persist(data)
      return true
    } catch (err) {
      error.value = err.message || 'Registrasi gagal'
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
  }

  return { token, user, error, loading, isAuthenticated, login, register, logout }
})
