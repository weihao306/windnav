<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { ExternalLink, LayoutDashboard, Search, Sparkles } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { getData, postData } from '../api/client'
import type { Category, SettingMap, Site } from '../api/types'

const searchText = ref('')
const activeCategory = ref('all')

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

async function openSite(site: Site) {
  postData(`/public/sites/${site.id}/click`).catch(() => undefined)
  window.open(site.url, '_blank', 'noopener,noreferrer')
}
</script>

<template>
  <main class="min-h-screen bg-[#f6f8fb] text-slate-900">
    <section class="mx-auto flex w-full max-w-7xl flex-col gap-8 px-5 py-6 sm:px-8 lg:px-10">
      <header class="flex flex-col gap-5 rounded-[8px] border border-white/80 bg-white/85 p-5 card-shadow sm:p-7">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="min-w-0">
            <div class="mb-3 inline-flex items-center gap-2 rounded-full bg-cyan-50 px-3 py-1 text-sm font-medium text-cyan-700">
              <Sparkles class="h-4 w-4" />
              WindNav
            </div>
            <h1 class="text-3xl font-semibold tracking-normal text-slate-950 sm:text-4xl">
              {{ summary.site_title ?? 'WindNav' }}
            </h1>
            <p class="mt-2 max-w-2xl text-base leading-7 text-slate-600">
              {{ summary.site_subtitle ?? '简单轻快的自建导航页' }}
            </p>
          </div>
          <RouterLink to="/admin" class="inline-flex h-10 items-center justify-center gap-2 rounded-[8px] border border-slate-200 bg-white px-4 text-sm font-medium text-slate-700 transition hover:border-cyan-200 hover:text-cyan-700">
            <LayoutDashboard class="h-4 w-4" />
            管理
          </RouterLink>
        </div>

        <div class="flex flex-col gap-3 lg:flex-row lg:items-center">
          <label class="relative flex-1">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400" />
            <input v-model="searchText" class="h-12 w-full rounded-[8px] border border-slate-200 bg-white pl-10 pr-4 text-base outline-none transition focus:border-cyan-400 focus:ring-4 focus:ring-cyan-100" :placeholder="summary.search_placeholder ?? '搜索站点、标签或描述'" />
          </label>
          <div class="flex gap-2 overflow-x-auto pb-1 lg:max-w-xl">
            <button class="h-10 shrink-0 rounded-[8px] px-4 text-sm font-medium transition" :class="activeCategory === 'all' ? 'bg-slate-900 text-white' : 'bg-slate-100 text-slate-600 hover:bg-cyan-50 hover:text-cyan-700'" @click="activeCategory = 'all'">
              全部
            </button>
            <button v-for="category in categories" :key="category.id" class="h-10 shrink-0 rounded-[8px] px-4 text-sm font-medium transition" :class="activeCategory === category.slug ? 'bg-slate-900 text-white' : 'bg-slate-100 text-slate-600 hover:bg-cyan-50 hover:text-cyan-700'" @click="activeCategory = category.slug">
              {{ category.name }}
            </button>
          </div>
        </div>
      </header>

      <section v-if="sitesQuery.isLoading.value" class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        <div v-for="item in 6" :key="item" class="h-36 animate-pulse rounded-[8px] border border-slate-200 bg-white" />
      </section>

      <section v-else-if="filteredSites.length" class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        <article v-for="site in filteredSites" :key="site.id" class="group flex min-h-36 flex-col justify-between rounded-[8px] border border-slate-200 bg-white p-4 transition hover:-translate-y-0.5 hover:border-cyan-200 hover:shadow-lg">
          <div class="flex gap-3">
            <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-[8px] bg-gradient-to-br from-cyan-50 to-emerald-50 text-sm font-semibold text-cyan-700">
              <img v-if="site.iconUrl" :src="site.iconUrl" alt="" class="h-7 w-7 rounded-[6px] object-cover" />
              <span v-else>{{ site.fallbackIcon || site.title.slice(0, 1).toUpperCase() }}</span>
            </div>
            <div class="min-w-0 flex-1">
              <h2 class="truncate text-base font-semibold text-slate-950">{{ site.title }}</h2>
              <p class="mt-1 line-clamp-2 min-h-10 text-sm leading-5 text-slate-500">{{ site.description || site.url }}</p>
            </div>
          </div>
          <div class="mt-4 flex items-center justify-between gap-3">
            <div class="flex min-w-0 gap-1 overflow-hidden">
              <span v-for="tag in site.tags?.slice(0, 3)" :key="tag.id" class="shrink-0 rounded-[6px] bg-slate-100 px-2 py-1 text-xs text-slate-500">{{ tag.name }}</span>
            </div>
            <button class="inline-flex h-9 shrink-0 items-center gap-2 rounded-[8px] bg-cyan-600 px-3 text-sm font-medium text-white transition hover:bg-cyan-700" @click="openSite(site)">
              打开
              <ExternalLink class="h-4 w-4" />
            </button>
          </div>
        </article>
      </section>

      <section v-else class="rounded-[8px] border border-dashed border-slate-300 bg-white p-10 text-center text-slate-500">
        暂无匹配站点
      </section>
    </section>
  </main>
</template>
