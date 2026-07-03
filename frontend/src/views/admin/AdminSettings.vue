<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Save } from 'lucide-vue-next'
import { reactive, watch } from 'vue'
import { getData, putData } from '../../api/client'
import type { SettingMap } from '../../api/types'

const queryClient = useQueryClient()
const form = reactive({ site_title: '', site_subtitle: '', search_placeholder: '', default_theme: 'light' })
const settings = useQuery({ queryKey: ['admin-settings'], queryFn: () => getData<SettingMap>('/admin/settings') })

watch(() => settings.data.value, (value) => {
  if (!value) return
  form.site_title = value.site_title ?? 'WindNav'
  form.site_subtitle = value.site_subtitle ?? '简单轻快的自建导航页'
  form.search_placeholder = value.search_placeholder ?? '搜索站点、标签或描述'
  form.default_theme = value.default_theme ?? 'light'
}, { immediate: true })

async function saveSettings() {
  await putData('/admin/settings', { settings: form })
  await queryClient.invalidateQueries({ queryKey: ['admin-settings'] })
}
</script>

<template>
  <div class="grid gap-5">
    <header class="rounded-[8px] border border-slate-200 bg-white p-5">
      <h1 class="text-2xl font-semibold text-slate-950">站点设置</h1>
    </header>
    <form class="grid max-w-2xl gap-4 rounded-[8px] border border-slate-200 bg-white p-5" @submit.prevent="saveSettings">
      <label><span class="admin-label">标题</span><input v-model="form.site_title" class="admin-input" /></label>
      <label><span class="admin-label">副标题</span><input v-model="form.site_subtitle" class="admin-input" /></label>
      <label><span class="admin-label">搜索占位</span><input v-model="form.search_placeholder" class="admin-input" /></label>
      <label><span class="admin-label">主题</span><select v-model="form.default_theme" class="admin-input"><option value="light">浅色</option><option value="dark">深色</option></select></label>
      <button class="inline-flex h-11 w-fit items-center gap-2 rounded-[8px] bg-cyan-600 px-4 font-medium text-white hover:bg-cyan-700"><Save class="h-4 w-4" />保存设置</button>
    </form>
  </div>
</template>

<style scoped>
.admin-label { margin-bottom: 0.375rem; display: block; font-size: 0.875rem; font-weight: 500; color: rgb(71 85 105); }
.admin-input { height: 2.75rem; width: 100%; border-radius: 8px; border: 1px solid rgb(226 232 240); padding: 0 0.75rem; outline: none; }
.admin-input:focus { border-color: rgb(34 211 238); box-shadow: 0 0 0 4px rgb(207 250 254); }
</style>
