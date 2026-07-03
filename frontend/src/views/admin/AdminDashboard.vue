<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { Folder, Tags, Settings } from 'lucide-vue-next'
import { getData } from '../../api/client'
import type { Category, SettingMap, Site, Tag } from '../../api/types'

const categories = useQuery({ queryKey: ['admin-categories'], queryFn: () => getData<Category[]>('/admin/categories') })
const sites = useQuery({ queryKey: ['admin-sites'], queryFn: () => getData<Site[]>('/admin/sites', { page_size: 100 }) })
const tags = useQuery({ queryKey: ['admin-tags'], queryFn: () => getData<Tag[]>('/admin/tags') })
const settings = useQuery({ queryKey: ['admin-settings'], queryFn: () => getData<SettingMap>('/admin/settings') })
</script>

<template>
  <div class="grid gap-5">
    <header class="rounded-[8px] border border-slate-200 bg-white p-5">
      <h1 class="text-2xl font-semibold text-slate-950">管理概览</h1>
      <p class="mt-1 text-sm text-slate-500">{{ settings.data.value?.site_title ?? 'WindNav' }}</p>
    </header>
    <section class="grid gap-4 sm:grid-cols-3">
      <div class="rounded-[8px] border border-slate-200 bg-white p-5"><Folder class="mb-3 h-5 w-5 text-cyan-600" /><div class="text-2xl font-semibold">{{ sites.data.value?.length ?? 0 }}</div><p class="text-sm text-slate-500">站点</p></div>
      <div class="rounded-[8px] border border-slate-200 bg-white p-5"><Folder class="mb-3 h-5 w-5 text-emerald-600" /><div class="text-2xl font-semibold">{{ categories.data.value?.length ?? 0 }}</div><p class="text-sm text-slate-500">分类</p></div>
      <div class="rounded-[8px] border border-slate-200 bg-white p-5"><Tags class="mb-3 h-5 w-5 text-violet-600" /><div class="text-2xl font-semibold">{{ tags.data.value?.length ?? 0 }}</div><p class="text-sm text-slate-500">标签</p></div>
    </section>
    <RouterLink to="/admin/settings" class="inline-flex w-fit items-center gap-2 rounded-[8px] bg-slate-900 px-4 py-2 text-sm font-medium text-white"><Settings class="h-4 w-4" />站点设置</RouterLink>
  </div>
</template>
