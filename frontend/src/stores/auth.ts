import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { api } from '../api/http'
import type { CurrentUser } from '../api/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<CurrentUser | null>(null)
  const loading = ref(false)
  const error = ref('')

  const isLoggedIn = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const profileComplete = computed(() => Boolean(user.value?.profileComplete))

  async function loadMe() {
    loading.value = true
    error.value = ''
    try {
      user.value = await api.me()
    } catch {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  async function login(email: string, password: string) {
    loading.value = true
    error.value = ''
    try {
      user.value = await api.login(email, password)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '登录失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function register(email: string, password: string, nickname: string) {
    loading.value = true
    error.value = ''
    try {
      user.value = await api.register(email, password, nickname)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '注册失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    await api.logout()
    user.value = null
  }

  return { user, loading, error, isLoggedIn, isAdmin, profileComplete, loadMe, login, register, logout }
})
