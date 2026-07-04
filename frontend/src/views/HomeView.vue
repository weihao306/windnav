<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { ExternalLink, Globe2, LayoutDashboard, Moon, PanelLeftClose, PanelLeftOpen, Search, SearchX, Sparkles, Sun } from 'lucide-vue-next'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { getData, postData } from '../api/client'
import type { Category, SearchEngine, SettingMap, Site } from '../api/types'

const searchText = ref('')
const activeCategory = ref('all')
const themeMode = ref<'dark' | 'light'>('light')
const systemTheme = ref<'dark' | 'light'>('light')
const selectedEngineSlug = ref('')
const sidebarCollapsed = ref(false)
const currentTime = ref(new Date())
const brokenIcons = ref(new Set<string>())
const searchInputRef = ref<HTMLInputElement | null>(null)
let clockTimer: number | undefined
let mediaQuery: MediaQueryList | undefined

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
const searchEnginesQuery = useQuery({
  queryKey: ['public-search-engines'],
  queryFn: () => getData<SearchEngine[]>('/public/search-engines'),
})

const summary = computed(() => summaryQuery.data.value ?? {})
const categories = computed(() => categoriesQuery.data.value ?? [])
const sites = computed(() => sitesQuery.data.value ?? [])
const searchEngines = computed(() => searchEnginesQuery.data.value ?? [])
const selectedEngine = computed(() => searchEngines.value.find((engine) => engine.slug === selectedEngineSlug.value) ?? searchEngines.value.find((engine) => engine.isDefault) ?? searchEngines.value[0])
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
const searchPlaceholder = computed(() => `搜索站点、标签或描述`)

function syncThemeFromSystem() {
  systemTheme.value = mediaQuery?.matches ? 'dark' : 'light'
  themeMode.value = systemTheme.value
}

watch(searchEngines, (engines) => {
  if (!engines.length || engines.some((engine) => engine.slug === selectedEngineSlug.value)) return
  selectedEngineSlug.value = engines.find((engine) => engine.isDefault)?.slug ?? engines[0].slug
}, { immediate: true })

onMounted(() => {
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  syncThemeFromSystem()
  mediaQuery.addEventListener('change', syncThemeFromSystem)

  clockTimer = window.setInterval(() => {
    currentTime.value = new Date()
  }, 1000)
  document.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  if (clockTimer) window.clearInterval(clockTimer)
  mediaQuery?.removeEventListener('change', syncThemeFromSystem)
  document.removeEventListener('keydown', handleGlobalKeydown)
})

function toggleTheme() {
  systemTheme.value = isDarkMode.value ? 'light' : 'dark'
  themeMode.value = systemTheme.value
}

function handleGlobalKeydown(e: KeyboardEvent) {
  if ((e.key === 'k' || e.key === 'K') && (e.ctrlKey || e.metaKey)) {
    e.preventDefault()
    searchInputRef.value?.focus()
    return
  }
  if (e.key === '/' && document.activeElement !== searchInputRef.value && !(document.activeElement instanceof HTMLInputElement) && !(document.activeElement instanceof HTMLTextAreaElement) && !(document.activeElement instanceof HTMLSelectElement)) {
    e.preventDefault()
    searchInputRef.value?.focus()
  }
}

function clearSearch() {
  searchText.value = ''
  activeCategory.value = 'all'
  nextTick(() => searchInputRef.value?.focus())
}

function submitSearch() {
  const keyword = searchText.value.trim()
  if (!keyword || !selectedEngine.value) return
  const encodedKeyword = encodeURIComponent(keyword)
  const target = selectedEngine.value.searchUrl.includes('{query}')
    ? selectedEngine.value.searchUrl.replaceAll('{query}', encodedKeyword)
    : `${selectedEngine.value.searchUrl}${selectedEngine.value.searchUrl.includes('?') ? '&' : '?'}q=${encodedKeyword}`
  window.open(target, '_blank', 'noopener,noreferrer')
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

async function openSite(site: Site) {
  postData(`/public/sites/${site.id}/click`).catch(() => undefined)
  window.open(site.url, '_blank', 'noopener,noreferrer')
}

function onImgError(siteId: string) {
  brokenIcons.value.add(siteId)
}

function getFallbackIcon(site: Site) {
  return site.fallbackIcon || site.title.slice(0, 1).toUpperCase()
}
</script>

<template>
  <main class="home-layout min-h-screen" :class="{ 'sidebar-collapsed': sidebarCollapsed }" :data-theme="themeMode">
    <aside class="left-sidebar hidden xl:flex">
      <div class="sidebar-brand">
        <div class="brand-mark">
          <svg width="40" height="40" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <use href="/icons.svg#windnav-icon"/>
          </svg>
        </div>
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
        title="全部"
        @click="activeCategory = 'all'"
      >
        <span class="sidebar-icon">全</span>
        <span class="sidebar-text">全部</span>
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
            <span class="brand-mark brand-mark-sm">
              <svg width="30" height="30" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
                <use href="/icons.svg#windnav-icon"/>
              </svg>
            </span>
            {{ summary.site_title ?? 'WindNav' }}
          </div>
          <div class="hero-actions">
            <button class="theme-toggle" type="button" :aria-label="isDarkMode ? '当前跟随深色模式' : '当前跟随浅色模式'" @click="toggleTheme">
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

        <form class="search-box" role="search" @submit.prevent="submitSearch">
          <Search class="pointer-events-none h-5 w-5 text-slate-400" />
          <input
            ref="searchInputRef"
            v-model="searchText"
            :placeholder="searchPlaceholder"
            aria-label="搜索站点"
          />
          <select v-if="searchEngines.length" v-model="selectedEngineSlug" class="search-engine-select" aria-label="网络搜索引擎">
            <option v-for="engine in searchEngines" :key="engine.id" :value="engine.slug">{{ engine.name }}</option>
          </select>
          <button class="search-submit" type="submit" :disabled="!searchText.trim() || !selectedEngine">
            <Globe2 class="h-4 w-4" />
            搜索
          </button>
        </form>

        <div class="category-strip xl:hidden">
          <button
            class="category-pill"
            :class="activeCategory === 'all' ? 'category-pill-active' : ''"
            @click="activeCategory = 'all'"
          >
            全部（{{ sites.length }}）
          </button>
          <button
            v-for="category in categories"
            :key="category.id"
            class="category-pill"
            :class="activeCategory === category.slug ? 'category-pill-active' : ''"
            @click="activeCategory = category.slug"
          >
            {{ category.name }}（{{ sites.filter((site) => site.categoryId === category.id || site.category?.slug === category.slug).length }}）
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
              role="link"
              tabindex="0"
              :aria-label="'打开 ' + site.title"
              @click="openSite(site)"
              @keydown.enter="openSite(site)"
            >
              <div class="site-icon">
                <img v-if="site.iconUrl && !brokenIcons.has(String(site.id))" :src="site.iconUrl" alt="" @error="onImgError(String(site.id))" />
                <span v-else>{{ getFallbackIcon(site) }}</span>
              </div>
              <div class="min-w-0 flex-1">
                <h3>{{ site.title }}</h3>
                <p>{{ site.description || site.url }}</p>
                <div v-if="site.tags?.length" class="site-tags">
                  <span v-for="tag in site.tags" :key="tag.id" class="site-tag" :style="tag.color ? { background: tag.color + '22', color: tag.color, borderColor: tag.color + '44' } : {}">{{ tag.name }}</span>
                </div>
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
              role="link"
              tabindex="0"
              :aria-label="'打开 ' + site.title"
              @click="openSite(site)"
              @keydown.enter="openSite(site)"
            >
              <div class="site-icon">
                <img v-if="site.iconUrl && !brokenIcons.has(String(site.id))" :src="site.iconUrl" alt="" @error="onImgError(String(site.id))" />
                <span v-else>{{ getFallbackIcon(site) }}</span>
              </div>
              <div class="min-w-0 flex-1">
                <h3>{{ site.title }}</h3>
                <p>{{ site.description || site.url }}</p>
                <div v-if="site.tags?.length" class="site-tags">
                  <span v-for="tag in site.tags" :key="tag.id" class="site-tag" :style="tag.color ? { background: tag.color + '22', color: tag.color, borderColor: tag.color + '44' } : {}">{{ tag.name }}</span>
                </div>
              </div>
              <ExternalLink class="site-open-icon" />
            </article>
          </div>
        </div>
      </section>

      <div v-else class="empty-state">
        <SearchX class="empty-state-icon" />
        <p class="empty-state-title">没有找到匹配的站点</p>
        <p class="empty-state-hint">试试调整搜索关键词，或者</p>
        <button class="empty-state-clear" type="button" @click="clearSearch">清空搜索条件</button>
      </div>
      <footer class="home-footer">
        <p>Powered by <strong>WindNav</strong> · {{ new Date().getFullYear() }}</p>
      </footer>
    </section>
  </main>
</template>
