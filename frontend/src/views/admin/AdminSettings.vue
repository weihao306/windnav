<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Save } from 'lucide-vue-next'
import { reactive, ref, watch } from 'vue'
import { getData, putData } from '../../api/client'
import type { SettingMap } from '../../api/types'

const queryClient = useQueryClient()
const form = reactive({ site_title: '', site_subtitle: '', search_placeholder: '', default_theme: 'light' })
const saving = ref(false)
const saved = ref(false)
const settings = useQuery({ queryKey: ['admin-settings'], queryFn: () => getData<SettingMap>('/admin/settings') })

watch(() => settings.data.value, (value) => {
  if (!value) return
  form.site_title = value.site_title ?? 'WindNav'
  form.site_subtitle = value.site_subtitle ?? '简单轻快的自建导航页'
  form.search_placeholder = value.search_placeholder ?? '搜索站点、标签或描述'
  form.default_theme = value.default_theme ?? 'light'
}, { immediate: true })

async function saveSettings() {
  saving.value = true
  saved.value = false
  try {
    await putData('/admin/settings', { settings: form })
    await queryClient.invalidateQueries({ queryKey: ['admin-settings'] })
    saved.value = true
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="grid gap-5">
    <header class="admin-page-header">
      <div><h1>站点设置</h1><p>调整首页品牌文案、搜索提示和默认主题。</p></div>
    </header>
    <form class="admin-panel grid max-w-2xl gap-4 p-5" @submit.prevent="saveSettings">
      <label><span class="admin-label">标题</span><input v-model="form.site_title" class="admin-input" /></label>
      <label><span class="admin-label">副标题</span><input v-model="form.site_subtitle" class="admin-input" /></label>
      <label><span class="admin-label">搜索占位</span><input v-model="form.search_placeholder" class="admin-input" /></label>
      <label><span class="admin-label">主题</span><select v-model="form.default_theme" class="admin-input"><option value="light">浅色</option><option value="dark">深色</option></select></label>
      <div class="flex flex-wrap items-center gap-3">
        <button class="admin-button" type="submit" :disabled="saving"><Save class="h-4 w-4" />{{ saving ? '正在保存...' : '保存设置' }}</button>
        <span v-if="saved" class="admin-status admin-status-success">设置已保存</span>
      </div>
    </form>
  </div>
</template>
