import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authTokenKey, getData, postData } from '../api/client'
import type { User } from '../api/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(authTokenKey) ?? '')
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => Boolean(token.value))

  async function login(username: string, password: string) {
    const data = await postData<{ token: string; user: User }>('/auth/login', { username, password })
    token.value = data.token
    user.value = data.user
    localStorage.setItem(authTokenKey, data.token)
  }

  async function loadMe() {
    if (!token.value) return
    user.value = await getData<User>('/auth/me')
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem(authTokenKey)
  }

  return { token, user, isAuthenticated, login, loadMe, logout }
})
