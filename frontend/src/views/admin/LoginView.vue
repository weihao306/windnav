<script setup lang="ts">
import { LogIn } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const username = ref('admin')
const password = ref('admin123456')
const error = ref('')
const loading = ref(false)

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
  <main class="grid min-h-screen place-items-center bg-[#f6f8fb] px-5">
    <form class="w-full max-w-sm rounded-[8px] border border-slate-200 bg-white p-6 card-shadow" @submit.prevent="submit">
      <div class="mb-6">
        <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-[8px]">
          <svg width="42" height="42" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <use href="/icons.svg#windnav-icon"/>
          </svg>
        </div>
        <h1 class="text-2xl font-semibold text-slate-950">后台登录</h1>
      </div>
      <label class="mb-4 block">
        <span class="mb-1 block text-sm font-medium text-slate-600">用户名</span>
        <input v-model="username" class="h-11 w-full rounded-[8px] border border-slate-200 px-3 outline-none focus:border-cyan-400 focus:ring-4 focus:ring-cyan-100" />
      </label>
      <label class="mb-4 block">
        <span class="mb-1 block text-sm font-medium text-slate-600">密码</span>
        <input v-model="password" type="password" class="h-11 w-full rounded-[8px] border border-slate-200 px-3 outline-none focus:border-cyan-400 focus:ring-4 focus:ring-cyan-100" />
      </label>
      <p v-if="error" class="mb-4 rounded-[8px] bg-rose-50 px-3 py-2 text-sm text-rose-700">{{ error }}</p>
      <button class="inline-flex h-11 w-full items-center justify-center gap-2 rounded-[8px] bg-cyan-600 font-medium text-white hover:bg-cyan-700 disabled:opacity-60" :disabled="loading">
        <LogIn class="h-4 w-4" />
        {{ loading ? '登录中' : '登录' }}
      </button>
    </form>
  </main>
</template>
