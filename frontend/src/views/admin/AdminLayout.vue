<script setup lang="ts">
import { Folder, Home, LayoutDashboard, LogOut, Search, Settings, Tags } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.push('/admin/login')
}
</script>

<template>
  <main class="min-h-screen bg-[#f6f8fb] text-slate-900">
    <div class="mx-auto flex w-full max-w-7xl flex-col gap-5 px-5 py-5 lg:flex-row">
      <aside class="rounded-[8px] border border-slate-200 bg-white p-3 lg:sticky lg:top-5 lg:h-[calc(100vh-40px)] lg:w-60">
        <RouterLink to="/" class="mb-4 flex items-center justify-start gap-2 rounded-[8px] px-3 py-2 font-semibold text-slate-950">
          <Home class="h-5 w-5 text-cyan-600" />
          WindNav
        </RouterLink>
        <nav class="grid gap-1">
          <RouterLink to="/admin" class="admin-link"><LayoutDashboard class="h-4 w-4" /><span>概览</span></RouterLink>
          <RouterLink to="/admin/sites" class="admin-link"><Folder class="h-4 w-4" /><span>站点</span></RouterLink>
          <RouterLink to="/admin/categories" class="admin-link"><Folder class="h-4 w-4" /><span>分类</span></RouterLink>
          <RouterLink to="/admin/tags" class="admin-link"><Tags class="h-4 w-4" /><span>标签</span></RouterLink>
          <RouterLink to="/admin/search-engines" class="admin-link"><Search class="h-4 w-4" /><span>搜索引擎</span></RouterLink>
          <RouterLink to="/admin/settings" class="admin-link"><Settings class="h-4 w-4" /><span>设置</span></RouterLink>
        </nav>
        <button class="mt-4 flex w-full items-center justify-start gap-2 rounded-[8px] px-3 py-2 text-sm text-slate-500 hover:bg-rose-50 hover:text-rose-600" @click="logout">
          <LogOut class="h-4 w-4" />
          退出
        </button>
      </aside>
      <section class="min-w-0 flex-1">
        <RouterView />
      </section>
    </div>
  </main>
</template>

<style scoped>
.admin-link {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 0.5rem;
  border-radius: 8px;
  padding: 0.625rem 0.75rem;
  color: rgb(71 85 105);
  font-size: 0.925rem;
  font-weight: 500;
}
.admin-link:hover,
.admin-link.router-link-exact-active {
  background: rgb(236 254 255);
  color: rgb(14 116 144);
}
.admin-link span {
  min-width: 0;
}
</style>
