<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { ExternalLink, Globe2, LayoutDashboard, Moon, PanelLeftClose, PanelLeftOpen, Search, SearchX, Sun } from 'lucide-vue-next'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { getData, postData } from '../api/client'
import type { Category, SearchEngine, SettingMap, Site } from '../api/types'

type HomeSection = {
  key: string
  name: string
  description: string
  caption: string
  sites: Site[]
}

type QuickFilter = {
  key: string
  name: string
}

const recentSiteStorageKey = 'windnav_recent_sites'

const searchText = ref('')
const activeCategory = ref('all')
const activeQuickFilter = ref('all')
const recentSiteIds = ref<number[]>([])
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
const uncategorizedSites = computed(() => filteredSites.value.filter((site) => !site.categoryId && !site.category?.slug))
const categoriesWithCounts = computed(() => categories.value.map((category) => ({
  ...category,
  count: sites.value.filter((site) => site.categoryId === category.id || site.category?.slug === category.slug).length,
})))
const quickFilters = computed<QuickFilter[]>(() => ([
  { key: 'all', name: '全部' },
  { key: 'recent', name: '最近访问' },
  { key: 'pinned', name: '常用站点' },
  ...categoriesWithCounts.value.map((category) => ({ key: category.slug, name: category.name })),
]))
const recentSites = computed(() => {
  const siteMap = new Map(filteredSites.value.map((site) => [site.id, site]))
  const fromHistory = recentSiteIds.value
    .map((id) => siteMap.get(id))
    .filter((site): site is Site => Boolean(site))

  if (fromHistory.length >= 6) return fromHistory.slice(0, 6)

  const supplemental = [...filteredSites.value]
    .filter((site) => !recentSiteIds.value.includes(site.id))
    .sort((a, b) => {
      const aTime = a.updatedAt ? new Date(a.updatedAt).getTime() : 0
      const bTime = b.updatedAt ? new Date(b.updatedAt).getTime() : 0
      return bTime - aTime
    })

  return [...fromHistory, ...supplemental].slice(0, 6)
})
const frequentSites = computed(() => {
  const ordered = [...filteredSites.value].sort((a, b) => {
    const aScore = (a.clickCount ?? 0) + (a.isPinned ? 1000 : 0)
    const bScore = (b.clickCount ?? 0) + (b.isPinned ? 1000 : 0)
    return bScore - aScore
  })
  return ordered.slice(0, 6)
})
const visibleSiteSections = computed(() => {
  if (activeQuickFilter.value === 'recent') {
    return [{ key: 'recent', name: '最近访问', description: '优先展示最近更新与访问痕迹较新的入口。', caption: 'RECENT', sites: recentSites.value }]
  }

  if (activeQuickFilter.value === 'pinned') {
    return [{ key: 'pinned', name: '常用站点', description: '基于置顶和访问频率整理的高频入口。', caption: 'PINNED', sites: frequentSites.value }]
  }

  return sectionBlocks.value
})
const sectionBlocks = computed<HomeSection[]>(() => {
  const sections = categories.value
    .map((category) => ({
      key: category.slug || String(category.id),
      name: category.name,
      description: category.description || '为这组常用站点提供更聚焦的入口体验。',
      caption: category.slug || `CAT-${category.id}`,
      sites: filteredSites.value.filter((site) => site.categoryId === category.id || site.category?.slug === category.slug),
    }))
    .filter((section) => section.sites.length > 0)

  if (activeCategory.value === 'all' && uncategorizedSites.value.length) {
    sections.push({
      key: 'uncategorized',
      name: '其他',
      description: '暂未归类但依旧值得保留的灵感与工具入口。',
      caption: 'UNSORTED',
      sites: uncategorizedSites.value,
    })
  }

  return sections
})
const isDarkMode = computed(() => themeMode.value === 'dark')
const displayTime = computed(() => currentTime.value.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }))
const searchPlaceholder = computed(() => '搜索站点、标签、描述或网址')
const totalSiteCount = computed(() => sites.value.length)
const totalTagCount = computed(() => {
  const tagIds = new Set<number>()
  for (const site of sites.value) {
    for (const tag of site.tags ?? []) {
      tagIds.add(tag.id)
    }
  }
  return tagIds.size
})
const filteredSiteCount = computed(() => filteredSites.value.length)
const launcherTitle = computed(() => summary.value.site_title ?? 'WindNav')
const launcherSubtitle = computed(() => summary.value.site_subtitle ?? '搜索优先的轻量导航启动器')

function syncThemeFromSystem() {
  systemTheme.value = mediaQuery?.matches ? 'dark' : 'light'
  themeMode.value = systemTheme.value
}

watch(searchEngines, (engines) => {
  if (!engines.length || engines.some((engine) => engine.slug === selectedEngineSlug.value)) return
  selectedEngineSlug.value = engines.find((engine) => engine.isDefault)?.slug ?? engines[0].slug
}, { immediate: true })

onMounted(() => {
  loadRecentSites()
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
  activeQuickFilter.value = 'all'
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

function selectQuickFilter(key: string) {
  activeQuickFilter.value = key
  if (categoriesWithCounts.value.some((category) => category.slug === key)) {
    activeCategory.value = key
  }
  else {
    activeCategory.value = 'all'
  }
}

function selectCategory(key: string) {
  activeCategory.value = key
  activeQuickFilter.value = key
}

function loadRecentSites() {
  try {
    const raw = localStorage.getItem(recentSiteStorageKey)
    if (!raw) {
      recentSiteIds.value = []
      return
    }

    const parsed = JSON.parse(raw)
    recentSiteIds.value = Array.isArray(parsed)
      ? parsed.map((item) => Number(item)).filter((item) => Number.isInteger(item) && item > 0).slice(0, 12)
      : []
  }
  catch {
    recentSiteIds.value = []
  }
}

function rememberRecentSite(siteId: number) {
  const nextIds = [siteId, ...recentSiteIds.value.filter((id) => id !== siteId)].slice(0, 12)
  recentSiteIds.value = nextIds
  localStorage.setItem(recentSiteStorageKey, JSON.stringify(nextIds))
}

function openSite(site: Site) {
  rememberRecentSite(site.id)
  postData(`/public/sites/${site.id}/click`).catch(() => undefined)
  window.open(site.url, '_blank', 'noopener,noreferrer')
}

function onImgError(siteId: string) {
  brokenIcons.value.add(siteId)
}

function getFallbackIcon(site: Site) {
  return site.fallbackIcon || site.title.slice(0, 1).toUpperCase()
}

function getSiteHost(url: string) {
  try {
    return new URL(url).hostname.replace(/^www\./, '')
  }
  catch {
    return url
  }
}

function getCategoryLabelForSite(site: Site) {
  return site.category?.name
    ?? categories.value.find((category) => category.id === site.categoryId || category.slug === site.category?.slug)?.name
    ?? '未分类'
}

function getTagStyle(color?: string) {
  return color
    ? {
        background: `${color}22`,
        color,
        borderColor: `${color}44`,
      }
    : {}
}
</script>

<template>
  <main class="home-layout" :class="{ 'sidebar-collapsed': sidebarCollapsed }" :data-theme="themeMode">
    <div class="home-backdrop" aria-hidden="true">
      <span class="backdrop-orb backdrop-orb-cyan" />
      <span class="backdrop-orb backdrop-orb-violet" />
      <span class="backdrop-orb backdrop-orb-blue" />
      <span class="backdrop-grid" />
      <span class="backdrop-noise" />
    </div>

    <aside class="left-sidebar hidden xl:flex">
      <div class="sidebar-brand">
        <div class="brand-mark">
          <svg width="40" height="40" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <use href="/icons.svg#windnav-icon" />
          </svg>
        </div>
        <div class="sidebar-brand-text">
          <p>{{ summary.site_title ?? 'WindNav' }}</p>
          <span>Glass Navigation</span>
        </div>
        <button class="sidebar-toggle" type="button" :aria-label="sidebarCollapsed ? '展开侧边栏' : '收起侧边栏'" @click="toggleSidebar">
          <PanelLeftOpen v-if="sidebarCollapsed" class="h-4 w-4" />
          <PanelLeftClose v-else class="h-4 w-4" />
        </button>
      </div>

      <section class="sidebar-overview glass-surface">
        <p class="sidebar-eyebrow">Control Center</p>
        <div class="sidebar-metrics">
          <div class="sidebar-metric">
            <span>站点</span>
            <strong>{{ totalSiteCount }}</strong>
          </div>
          <div class="sidebar-metric">
            <span>分类</span>
            <strong>{{ categories.length }}</strong>
          </div>
          <div class="sidebar-metric">
            <span>标签</span>
            <strong>{{ totalTagCount }}</strong>
          </div>
        </div>
      </section>

      <div class="sidebar-section-header">
        <span>导航分类</span>
      </div>

      <button
        class="sidebar-link"
        :class="activeCategory === 'all' ? 'sidebar-link-active' : ''"
        title="全部"
        @click="selectCategory('all')"
      >
        <span class="sidebar-icon">全</span>
        <span class="sidebar-text">全部站点</span>
        <span class="sidebar-count">{{ totalSiteCount }}</span>
      </button>

      <button
        v-for="category in categoriesWithCounts"
        :key="category.id"
        class="sidebar-link"
        :class="activeCategory === category.slug ? 'sidebar-link-active' : ''"
        :title="category.name"
        @click="selectCategory(category.slug)"
      >
        <span class="sidebar-icon">{{ category.name.slice(0, 1) }}</span>
        <span class="sidebar-text">{{ category.name }}</span>
        <span class="sidebar-count">{{ category.count }}</span>
      </button>

      <div class="sidebar-footer glass-surface">
        <p>快捷键</p>
        <span>Ctrl / ⌘ + K</span>
      </div>
    </aside>

    <section class="main-content">
      <header class="top-toolbar">
        <div class="toolbar-chip glass-chip">
          <span class="toolbar-dot" />
          <span>{{ summary.site_title ?? 'WindNav' }}</span>
        </div>

        <div class="toolbar-actions">
          <div class="toolbar-chip glass-chip toolbar-shortcut">/ 聚焦搜索</div>
          <button class="theme-toggle" type="button" :aria-label="isDarkMode ? '切换到浅色模式' : '切换到深色模式'" @click="toggleTheme">
            <Sun v-if="isDarkMode" class="h-4 w-4" />
            <Moon v-else class="h-4 w-4" />
          </button>
          <RouterLink to="/admin" class="admin-link">
            <LayoutDashboard class="h-4 w-4" />
            管理台
          </RouterLink>
        </div>
      </header>

      <section class="launcher-shell glass-surface">
        <div class="launcher-copy">
          <h1>{{ launcherTitle }}</h1>
          <p>{{ launcherSubtitle }}</p>
        </div>

        <form class="search-panel glass-panel" role="search" @submit.prevent="submitSearch">
          <div class="search-panel-main">
            <Search class="search-panel-icon" />
            <input
              ref="searchInputRef"
              v-model="searchText"
              :placeholder="searchPlaceholder"
              aria-label="搜索站点"
            />
          </div>

          <div class="search-panel-actions">
            <select v-if="searchEngines.length" v-model="selectedEngineSlug" class="search-engine-select" aria-label="网络搜索引擎">
              <option v-for="engine in searchEngines" :key="engine.id" :value="engine.slug">{{ engine.name }}</option>
            </select>
            <button v-if="searchText.trim()" class="search-clear" type="button" @click="clearSearch">清空</button>
            <button class="search-submit" type="submit" :disabled="!searchText.trim() || !selectedEngine">
              <Globe2 class="h-4 w-4" />
              搜索
            </button>
          </div>
        </form>

        <div class="category-rail">
          <button
            v-for="filter in quickFilters"
            :key="filter.key"
            class="category-pill"
            :class="activeQuickFilter === filter.key ? 'category-pill-active' : ''"
            @click="selectQuickFilter(filter.key)"
          >
            {{ filter.name }}
            <span>{{ filter.key === 'all' ? totalSiteCount : filter.key === 'recent' ? recentSites.length : filter.key === 'pinned' ? frequentSites.length : (categoriesWithCounts.find((category) => category.slug === filter.key)?.count ?? 0) }}</span>
          </button>
        </div>
        <div class="launcher-meta">
          <div class="launcher-meta-card glass-chip">
            <span>当前结果</span>
            <strong>{{ filteredSiteCount }}</strong>
          </div>
          <div class="launcher-meta-card glass-chip">
            <span>当前时间</span>
            <strong>{{ displayTime }}</strong>
          </div>
        </div>
      </section>

      <section class="content-stack home-dashboard">
        <div class="quick-row">
          <section class="quick-strip glass-surface">
            <div class="quick-strip-header">
              <div>
                <p class="section-kicker">RECENT</p>
                <h2 class="quick-strip-title">最近访问</h2>
              </div>
              <span>{{ recentSites.length }} 个入口</span>
            </div>

            <div class="quick-strip-grid">
              <button v-for="site in recentSites" :key="site.id" class="quick-pill" type="button" @click="openSite(site)">
                <span class="quick-pill-title">{{ site.title }}</span>
                <span class="quick-pill-host">{{ getSiteHost(site.url) }}</span>
              </button>
            </div>
          </section>

          <section class="quick-strip glass-surface">
            <div class="quick-strip-header">
              <div>
                <p class="section-kicker">PINNED</p>
                <h2 class="quick-strip-title">常用站点</h2>
              </div>
              <span>{{ frequentSites.length }} 个入口</span>
            </div>

            <div class="quick-strip-grid">
              <button v-for="site in frequentSites" :key="site.id" class="quick-pill" type="button" @click="openSite(site)">
                <span class="quick-pill-title">{{ site.title }}</span>
                <span class="quick-pill-host">{{ getSiteHost(site.url) }}</span>
              </button>
            </div>
          </section>
        </div>
      </section>

      <section v-if="sitesQuery.isLoading.value" class="content-stack">
        <div class="site-section glass-surface loading-section">
          <div class="section-heading">
            <div>
              <p class="section-kicker">LOADING</p>
              <h2 class="section-title">正在加载导航内容</h2>
            </div>
          </div>
          <div class="site-grid">
            <div v-for="item in 6" :key="item" class="skeleton-card" />
          </div>
        </div>
      </section>

      <section v-else-if="filteredSites.length" class="content-stack">
        <div v-for="section in visibleSiteSections" :key="section.key" class="site-section glass-surface">
          <div class="section-heading">
            <div>
              <p class="section-kicker">{{ section.caption }}</p>
              <h2 class="section-title">{{ section.name }}</h2>
              <p class="section-description">{{ section.description }}</p>
            </div>
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
              <div class="site-card-head">
                <div class="site-icon">
                  <img v-if="site.iconUrl && !brokenIcons.has(String(site.id))" :src="site.iconUrl" alt="" @error="onImgError(String(site.id))" />
                  <span v-else>{{ getFallbackIcon(site) }}</span>
                </div>

                <div class="site-card-title-wrap">
                  <h3>{{ site.title }}</h3>
                  <span class="site-host">{{ getSiteHost(site.url) }}</span>
                </div>

                <ExternalLink class="site-open-icon" />
              </div>

              <p class="site-description">{{ site.description || site.url }}</p>

              <div class="site-card-footer">
                <div v-if="site.tags?.length" class="site-tags">
                  <span
                    v-for="tag in site.tags"
                    :key="tag.id"
                    class="site-tag"
                    :style="getTagStyle(tag.color)"
                  >
                    {{ tag.name }}
                  </span>
                </div>

                <span class="site-card-badge">{{ getCategoryLabelForSite(site) }}</span>
              </div>
            </article>
          </div>
        </div>
      </section>

      <div v-else class="empty-state glass-surface">
        <div class="empty-state-orb" />
        <SearchX class="empty-state-icon" />
        <p class="empty-state-title">没有找到匹配的站点</p>
        <p class="empty-state-hint">试试调整关键词、切换分类，或者直接清空当前筛选。</p>
        <button class="empty-state-clear" type="button" @click="clearSearch">恢复默认视图</button>
      </div>

      <footer class="home-footer">
        <p>Powered by <strong>WindNav</strong> · {{ new Date().getFullYear() }}</p>
      </footer>
    </section>
  </main>
</template>
