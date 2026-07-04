<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { ExternalLink, LayoutDashboard, Moon, PanelLeftClose, PanelLeftOpen, Search, Sparkles, Sun } from 'lucide-vue-next'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { getData, postData } from '../api/client'
import type { Category, SettingMap, Site } from '../api/types'

const searchText = ref('')
const activeCategory = ref('all')
const themeMode = ref<'dark' | 'light'>('dark')
const sidebarCollapsed = ref(false)
const currentTime = ref(new Date())
let clockTimer: number | undefined

const summaryQuery = useQuery({
  queryKey: ['public-summary'],
  queryFn: () => getData<SettingMap>('/public/summary'),
})
const categoriesQuery = useQuery({
  queryKey: ['public-categories'],
  queryFn: () => getData<Category[]>('/public/categories'),
})
const sitesQuery = useQuery({
  queryKey: ['public-sites'],
  queryFn: () => getData<Site[]>('/public/sites', { page_size: 100 }),
})

const summary = computed(() => summaryQuery.data.value ?? {})
const categories = computed(() => categoriesQuery.data.value ?? [])
const sites = computed(() => sitesQuery.data.value ?? [])
const filteredSites = computed(() => {
  const keyword = searchText.value.trim().toLowerCase()
  return sites.value.filter((site) => {
    const categoryMatched = activeCategory.value === 'all' || site.category?.slug === activeCategory.value || String(site.categoryId) === activeCategory.value
    const text = `${site.title} ${site.description} ${site.url} ${site.tags?.map((tag) => tag.name).join(' ') ?? ''}`.toLowerCase()
    return categoryMatched && (!keyword || text.includes(keyword))
  })
})
const groupedSections = computed(() => categories.value
  .map((category) => ({
    category,
    sites: filteredSites.value.filter((site) => site.categoryId === category.id || site.category?.slug === category.slug),
  }))
  .filter((section) => section.sites.length > 0))
const uncategorizedSites = computed(() => filteredSites.value.filter((site) => !site.categoryId && !site.category?.slug))
const isDarkMode = computed(() => themeMode.value === 'dark')
const displayTime = computed(() => currentTime.value.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }))
const displayDate = computed(() => currentTime.value.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }))

onMounted(() => {
  clockTimer = window.setInterval(() => {
    currentTime.value = new Date()
  }, 1000)
})

onUnmounted(() => {
  if (clockTimer) window.clearInterval(clockTimer)
})

function toggleTheme() {
  themeMode.value = isDarkMode.value ? 'light' : 'dark'
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

async function openSite(site: Site) {
  postData(`/public/sites/${site.id}/click`).catch(() => undefined)
  window.open(site.url, '_blank', 'noopener,noreferrer')
}
</script>

<template>
  <main class="home-layout min-h-screen" :class="{ 'sidebar-collapsed': sidebarCollapsed }" :data-theme="themeMode">
    <aside class="left-sidebar hidden xl:flex">
      <div class="sidebar-brand">
        <div class="brand-mark">W</div>
        <div class="sidebar-brand-text">
          <p class="text-base font-bold text-white">{{ summary.site_title ?? 'WindNav' }}</p>
          <p class="text-xs text-slate-400">Navigation</p>
        </div>
        <button class="sidebar-toggle" type="button" :aria-label="sidebarCollapsed ? '展开侧边栏' : '收起侧边栏'" @click="toggleSidebar">
          <PanelLeftOpen v-if="sidebarCollapsed" class="h-4 w-4" />
          <PanelLeftClose v-else class="h-4 w-4" />
        </button>
      </div>

      <div class="sidebar-section">导航分类</div>
      <button
        class="sidebar-link"
        :class="activeCategory === 'all' ? 'sidebar-link-active' : ''"
        title="常用工具"
        @click="activeCategory = 'all'"
      >
        <span class="sidebar-icon">常</span>
        <span class="sidebar-text">常用工具</span>
        <span class="sidebar-count">{{ sites.length }}</span>
      </button>
      <button
        v-for="category in categories"
        :key="category.id"
        class="sidebar-link"
        :class="activeCategory === category.slug ? 'sidebar-link-active' : ''"
        :title="category.name"
        @click="activeCategory = category.slug"
      >
        <span class="sidebar-icon">{{ category.name.slice(0, 1) }}</span>
        <span class="sidebar-text">{{ category.name }}</span>
        <span class="sidebar-count">{{ sites.filter((site) => site.categoryId === category.id || site.category?.slug === category.slug).length }}</span>
      </button>
    </aside>

    <section class="main-content">
      <header class="hero-panel">
        <nav class="hero-nav">
          <div class="inline-flex items-center gap-2 text-sm font-medium text-slate-300 xl:hidden">
            <span class="brand-mark brand-mark-sm">W</span>
            {{ summary.site_title ?? 'WindNav' }}
          </div>
          <div class="hero-actions">
            <button class="theme-toggle" type="button" :aria-label="isDarkMode ? '切换浅色模式' : '切换深色模式'" @click="toggleTheme">
              <Sun v-if="isDarkMode" class="h-4 w-4" />
              <Moon v-else class="h-4 w-4" />
            </button>
            <RouterLink to="/admin" class="admin-link">
              <LayoutDashboard class="h-4 w-4" />
              管理
            </RouterLink>
          </div>
        </nav>

        <div class="hero-body">
          <div class="hero-copy">
            <div class="hero-badge">
              <Sparkles class="h-4 w-4" />
              WindNav Start Page
            </div>
            <h1>{{ summary.site_title ?? 'WindNav' }}</h1>
            <p>{{ summary.site_subtitle ?? '简单轻快的自建导航页' }}</p>
          </div>

          <div class="clock-card" aria-label="当前日期时间">
            <div class="clock-time">{{ displayTime }}</div>
            <div class="clock-date">{{ displayDate }}</div>
          </div>
        </div>

        <label class="search-box">
          <Search class="pointer-events-none h-5 w-5 text-slate-400" />
          <input
            v-model="searchText"
            :placeholder="summary.search_placeholder ?? '搜索站点、标签或描述'"
          />
        </label>

        <div class="category-strip xl:hidden">
          <button
            class="category-pill"
            :class="activeCategory === 'all' ? 'category-pill-active' : ''"
            @click="activeCategory = 'all'"
          >
            全部
          </button>
          <button
            v-for="category in categories"
            :key="category.id"
            class="category-pill"
            :class="activeCategory === category.slug ? 'category-pill-active' : ''"
            @click="activeCategory = category.slug"
          >
            {{ category.name }}
          </button>
        </div>
      </header>

      <section v-if="sitesQuery.isLoading.value" class="site-grid mt-8">
        <div v-for="item in 9" :key="item" class="skeleton-card" />
      </section>

      <section v-else-if="filteredSites.length" class="content-stack">
        <div v-for="section in groupedSections" :key="section.category.id" class="site-section">
          <div class="section-heading">
            <h2 class="section-title">{{ section.category.name }}</h2>
            <span>{{ section.sites.length }} 个站点</span>
          </div>
          <div class="site-grid">
            <article
              v-for="site in section.sites"
              :key="site.id"
              class="site-card"
              @click="openSite(site)"
            >
              <div class="site-icon">
                <img v-if="site.iconUrl" :src="site.iconUrl" alt="" />
                <span v-else>{{ site.fallbackIcon || site.title.slice(0, 1).toUpperCase() }}</span>
              </div>
              <div class="min-w-0 flex-1">
                <h3>{{ site.title }}</h3>
                <p>{{ site.description || site.url }}</p>
              </div>
              <ExternalLink class="site-open-icon" />
            </article>
          </div>
        </div>

        <div v-if="activeCategory === 'all' && uncategorizedSites.length" class="site-section">
          <div class="section-heading">
            <h2 class="section-title">其他</h2>
            <span>{{ uncategorizedSites.length }} 个站点</span>
          </div>
          <div class="site-grid">
            <article
              v-for="site in uncategorizedSites"
              :key="site.id"
              class="site-card"
              @click="openSite(site)"
            >
              <div class="site-icon">
                <img v-if="site.iconUrl" :src="site.iconUrl" alt="" />
                <span v-else>{{ site.fallbackIcon || site.title.slice(0, 1).toUpperCase() }}</span>
              </div>
              <div class="min-w-0 flex-1">
                <h3>{{ site.title }}</h3>
                <p>{{ site.description || site.url }}</p>
              </div>
              <ExternalLink class="site-open-icon" />
            </article>
          </div>
        </div>
      </section>

      <section v-else class="empty-state">
        暂无匹配站点
      </section>
    </section>
  </main>
</template>
