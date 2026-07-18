<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { BarChart3, Folder, MousePointerClick, Settings, TrendingUp } from 'lucide-vue-next'
import { computed } from 'vue'
import { getData } from '../../api/client'
import type { DashboardStats, SettingMap } from '../../api/types'

const stats = useQuery({ queryKey: ['admin-dashboard-stats'], queryFn: () => getData<DashboardStats>('/admin/dashboard/stats') })
const settings = useQuery({ queryKey: ['admin-settings'], queryFn: () => getData<SettingMap>('/admin/settings') })

const statsData = computed(() => stats.data.value)

const siteCount = computed(() => statsData.value?.siteCount ?? 0)
const visibleSiteCount = computed(() => statsData.value?.visibleSiteCount ?? 0)
const totalClicks = computed(() => statsData.value?.totalClicks ?? 0)
const avgClicks = computed(() => {
  const count = siteCount.value
  return count > 0 ? Math.round(totalClicks.value / count) : 0
})
const pinnedCount = computed(() => statsData.value?.pinnedCount ?? 0)

const clickTrend = computed(() => {
  const list = statsData.value?.topClicked ?? []
  if (!list.length) return []

  const maxClicks = Math.max(...list.map((site) => site.clickCount), 1)

  return list.map((site) => ({
    id: site.id,
    title: site.title,
    shortTitle: site.title.length > 8 ? `${site.title.slice(0, 8)}…` : site.title,
    clicks: site.clickCount,
    percentage: Math.max((site.clickCount / maxClicks) * 100, 8),
  }))
})

const categoryTrend = computed(() => {
  const list = statsData.value?.categoryDist ?? []
  if (!list.length) return []

  const maxCount = Math.max(...list.map((item) => item.count), 1)

  return list.map((item) => ({
    id: item.id,
    name: item.name,
    count: item.count,
    percentage: Math.max((item.count / maxCount) * 100, item.count > 0 ? 12 : 0),
  }))
})

const summaryCards = computed(() => [
  { title: '站点', value: siteCount.value, note: `${visibleSiteCount.value} 个展示中`, icon: Folder, iconClass: 'text-cyan-600 bg-cyan-50' },
  { title: '分类', value: statsData.value?.categoryCount ?? 0, note: `${statsData.value?.tagCount ?? 0} 个标签已关联`, icon: BarChart3, iconClass: 'text-emerald-600 bg-emerald-50' },
  { title: '总点击', value: totalClicks.value, note: `单站平均 ${avgClicks.value}`, icon: MousePointerClick, iconClass: 'text-violet-600 bg-violet-50' },
  { title: '置顶推荐', value: pinnedCount.value, note: '可用于首页推荐位', icon: TrendingUp, iconClass: 'text-amber-600 bg-amber-50' },
])
</script>

<template>
  <div class="grid gap-5">
    <header class="rounded-[8px] border border-slate-200 bg-white p-5">
      <div class="flex flex-col gap-2 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-slate-950">管理概览</h1>
          <p class="mt-1 text-sm text-slate-500">{{ settings.data.value?.site_title ?? 'WindNav' }}</p>
        </div>
        <div class="rounded-[8px] bg-slate-50 px-4 py-3 text-sm text-slate-500">
          当前已收录 <span class="font-semibold text-slate-900">{{ siteCount }}</span> 个站点，累计点击 <span class="font-semibold text-slate-900">{{ totalClicks }}</span>
        </div>
      </div>
    </header>

    <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <article v-for="card in summaryCards" :key="card.title" class="rounded-[8px] border border-slate-200 bg-white p-5">
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="text-sm text-slate-500">{{ card.title }}</p>
            <div class="mt-2 text-3xl font-semibold text-slate-950">{{ card.value }}</div>
            <p class="mt-2 text-sm text-slate-400">{{ card.note }}</p>
          </div>
          <div :class="['flex h-11 w-11 items-center justify-center rounded-[8px]', card.iconClass]">
            <component :is="card.icon" class="h-5 w-5" />
          </div>
        </div>
      </article>
    </section>

    <section class="grid gap-5 xl:grid-cols-[1.35fr_1fr]">
      <article class="rounded-[8px] border border-slate-200 bg-white p-5">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-slate-950">访问点击趋势</h2>
            <p class="mt-1 text-sm text-slate-500">按站点点击量排序，观察当前热门内容</p>
          </div>
          <div class="rounded-full bg-cyan-50 px-3 py-1 text-xs font-medium text-cyan-700">TOP {{ clickTrend.length }}</div>
        </div>

        <div v-if="clickTrend.length" class="mt-5 space-y-4">
          <div v-for="site in clickTrend" :key="site.id" class="space-y-2">
            <div class="flex items-center justify-between gap-3 text-sm">
              <span class="truncate font-medium text-slate-700">{{ site.shortTitle }}</span>
              <span class="text-slate-500">{{ site.clicks }} 次</span>
            </div>
            <div class="h-3 overflow-hidden rounded-full bg-slate-100">
              <div class="admin-chart-bar h-full rounded-full bg-gradient-to-r from-cyan-500 via-sky-500 to-violet-500" :style="{ width: `${site.percentage}%` }"></div>
            </div>
          </div>
        </div>
        <div v-else class="mt-5 rounded-[8px] border border-dashed border-slate-200 bg-slate-50 px-4 py-8 text-center text-sm text-slate-500">
          暂无点击数据，新增站点并产生访问后会显示趋势图
        </div>
      </article>

      <article class="rounded-[8px] border border-slate-200 bg-white p-5">
        <div>
          <h2 class="text-lg font-semibold text-slate-950">分类收录分布</h2>
          <p class="mt-1 text-sm text-slate-500">展示各分类下站点数量，便于识别内容重心</p>
        </div>

        <div v-if="categoryTrend.length" class="mt-5 space-y-4">
          <div v-for="category in categoryTrend" :key="category.id" class="rounded-[8px] border border-slate-100 bg-slate-50 p-4">
            <div class="flex items-center justify-between gap-3">
              <div class="font-medium text-slate-800">{{ category.name }}</div>
              <div class="text-sm text-slate-500">{{ category.count }} 个站点</div>
            </div>
            <div class="mt-3 h-2 overflow-hidden rounded-full bg-white">
              <div class="h-full rounded-full bg-gradient-to-r from-emerald-500 to-cyan-500" :style="{ width: `${category.percentage}%` }"></div>
            </div>
          </div>
        </div>
        <div v-else class="mt-5 rounded-[8px] border border-dashed border-slate-200 bg-slate-50 px-4 py-8 text-center text-sm text-slate-500">
          暂无分类数据
        </div>
      </article>
    </section>

    <section class="grid gap-5 lg:grid-cols-[1.1fr_0.9fr]">
      <article class="rounded-[8px] border border-slate-200 bg-white p-5">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-slate-950">运营建议</h2>
            <p class="mt-1 text-sm text-slate-500">基于现有数据的简要提示</p>
          </div>
          <TrendingUp class="h-5 w-5 text-amber-500" />
        </div>
        <ul class="mt-4 grid gap-3 text-sm text-slate-600">
          <li class="rounded-[8px] bg-slate-50 px-4 py-3">热门站点集中在前 {{ clickTrend.length }} 个项目，可优先优化这些站点的描述与排序。</li>
          <li class="rounded-[8px] bg-slate-50 px-4 py-3">当前展示中站点占比 {{ siteCount ? Math.round((visibleSiteCount / siteCount) * 100) : 0 }}%，可检查隐藏内容是否需要重新开放。</li>
          <li class="rounded-[8px] bg-slate-50 px-4 py-3">置顶推荐位共有 {{ pinnedCount }} 个，适合承载高点击或重点导航内容。</li>
        </ul>
      </article>

      <article class="rounded-[8px] border border-slate-200 bg-white p-5">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-slate-950">快捷入口</h2>
            <p class="mt-1 text-sm text-slate-500">常用后台操作</p>
          </div>
          <Settings class="h-5 w-5 text-slate-400" />
        </div>
        <div class="mt-4 flex flex-wrap gap-3">
          <RouterLink to="/admin/sites" class="inline-flex items-center gap-2 rounded-[8px] bg-cyan-600 px-4 py-2 text-sm font-medium text-white hover:bg-cyan-700">
            <Folder class="h-4 w-4" />站点管理
          </RouterLink>
          <RouterLink to="/admin/settings" class="inline-flex items-center gap-2 rounded-[8px] bg-slate-900 px-4 py-2 text-sm font-medium text-white">
            <Settings class="h-4 w-4" />站点设置
          </RouterLink>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.admin-chart-bar {
  transition: width 0.5s ease;
}
</style>
