<script setup lang="ts">
import { Folder, LayoutDashboard, LogOut, Search, Settings, Tags } from 'lucide-vue-next'
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
  <main class="admin-page">
    <div class="admin-shell">
      <aside class="admin-sidebar">
        <RouterLink to="/" class="admin-brand">
          <span class="admin-brand-mark">
            <svg width="34" height="34" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
              <use href="/icons.svg#windnav-icon" />
            </svg>
          </span>
          <span class="admin-brand-copy">
            <strong>WindNav</strong>
            <span>Admin workspace</span>
          </span>
        </RouterLink>

        <p class="admin-nav-label">工作台</p>
        <nav class="admin-nav" aria-label="后台导航">
          <RouterLink to="/admin" class="admin-nav-link"><LayoutDashboard class="h-4 w-4" /><span>概览</span></RouterLink>
          <RouterLink to="/admin/sites" class="admin-nav-link"><Folder class="h-4 w-4" /><span>站点</span></RouterLink>
          <RouterLink to="/admin/categories" class="admin-nav-link"><Folder class="h-4 w-4" /><span>分类</span></RouterLink>
          <RouterLink to="/admin/tags" class="admin-nav-link"><Tags class="h-4 w-4" /><span>标签</span></RouterLink>
          <RouterLink to="/admin/search-engines" class="admin-nav-link"><Search class="h-4 w-4" /><span>搜索引擎</span></RouterLink>
          <RouterLink to="/admin/settings" class="admin-nav-link"><Settings class="h-4 w-4" /><span>设置</span></RouterLink>
        </nav>

        <button class="admin-logout" type="button" aria-label="退出登录" @click="logout">
          <LogOut class="h-4 w-4" />
          <span>退出登录</span>
        </button>
      </aside>

      <section class="admin-main">
        <RouterView />
      </section>
    </div>
  </main>
</template>
