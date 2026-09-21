<script setup lang="ts">
import { LogIn } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const usernameInput = ref<HTMLInputElement | null>(null)

onMounted(() => usernameInput.value?.focus())

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    await router.push(String(route.query.redirect ?? '/admin'))
  } catch (err) {
    error.value = err instanceof Error ? err.message : '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="admin-login">
    <form class="admin-login-card" @submit.prevent="submit">
      <RouterLink to="/" class="admin-login-brand" aria-label="返回首页">
        <svg width="42" height="42" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
          <use href="/icons.svg#windnav-icon" />
        </svg>
      </RouterLink>
      <h1>欢迎回来</h1>
      <p>登录 WindNav 管理工作台，维护你的导航内容与站点设置。</p>

      <div class="grid gap-4">
        <label>
          <span class="admin-label">用户名</span>
          <input ref="usernameInput" v-model="username" class="admin-input" autocomplete="username" required />
        </label>
        <label>
          <span class="admin-label">密码</span>
          <input v-model="password" type="password" class="admin-input" autocomplete="current-password" required />
        </label>
        <p v-if="error" class="admin-alert admin-alert-error" role="alert">{{ error }}</p>
        <button class="admin-button w-full" type="submit" :disabled="loading">
          <LogIn class="h-4 w-4" />
          {{ loading ? '正在登录...' : '登录管理台' }}
        </button>
      </div>
    </form>
  </main>
</template>
