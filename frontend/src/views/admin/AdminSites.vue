<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Trash2 } from 'lucide-vue-next'
import { computed, reactive, ref } from 'vue'
import { deleteData, getData, postData } from '../../api/client'
import type { Category, Site, Tag } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const form = reactive({ categoryId: 0, title: '', url: '', description: '', iconUrl: '', fallbackIcon: '', sortOrder: 0, isPinned: false, isVisible: true, tagIds: [] as number[] })
const sites = useQuery({ queryKey: ['admin-sites'], queryFn: () => getData<Site[]>('/admin/sites', { page_size: 100 }) })
const categories = useQuery({ queryKey: ['admin-categories'], queryFn: () => getData<Category[]>('/admin/categories') })
const tags = useQuery({ queryKey: ['admin-tags'], queryFn: () => getData<Tag[]>('/admin/tags') })
const canCreate = computed(() => form.categoryId > 0 && form.title.trim() && form.url.trim())

function resetForm() {
  form.categoryId = categories.data.value?.[0]?.id ?? 0
  form.title = ''
  form.url = ''
  form.description = ''
  form.iconUrl = ''
  form.fallbackIcon = ''
  form.sortOrder = 0
  form.isPinned = false
  form.isVisible = true
  form.tagIds = []
}

async function createSite() {
  error.value = ''
  try {
    await postData('/admin/sites', form)
    resetForm()
    await queryClient.invalidateQueries({ queryKey: ['admin-sites'] })
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  }
}

async function removeSite(id: number) {
  await deleteData(`/admin/sites/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-sites'] })
}
</script>

<template>
  <div class="grid gap-5">
    <header class="rounded-[8px] border border-slate-200 bg-white p-5">
      <h1 class="text-2xl font-semibold text-slate-950">站点管理</h1>
    </header>

    <form class="grid gap-3 rounded-[8px] border border-slate-200 bg-white p-5 lg:grid-cols-6" @submit.prevent="createSite">
      <select v-model.number="form.categoryId" required class="admin-input lg:col-span-2">
        <option :value="0">选择分类</option>
        <option v-for="category in categories.data.value" :key="category.id" :value="category.id">{{ category.name }}</option>
      </select>
      <input v-model="form.title" required placeholder="站点名称" class="admin-input lg:col-span-2" />
      <input v-model="form.url" required placeholder="https://example.com" class="admin-input lg:col-span-2" />
      <input v-model="form.description" placeholder="描述" class="admin-input lg:col-span-3" />
      <input v-model="form.iconUrl" placeholder="图标 URL" class="admin-input lg:col-span-2" />
      <input v-model="form.fallbackIcon" placeholder="备用字母" class="admin-input" />
      <select v-model="form.tagIds" multiple class="min-h-24 rounded-[8px] border border-slate-200 px-3 py-2 outline-none focus:border-cyan-400 focus:ring-4 focus:ring-cyan-100 lg:col-span-3">
        <option v-for="tag in tags.data.value" :key="tag.id" :value="tag.id">{{ tag.name }}</option>
      </select>
      <input v-model.number="form.sortOrder" type="number" placeholder="排序" class="admin-input" />
      <label class="flex h-11 items-center gap-2 rounded-[8px] border border-slate-200 px-3 text-sm text-slate-600"><input v-model="form.isPinned" type="checkbox" />置顶</label>
      <label class="flex h-11 items-center gap-2 rounded-[8px] border border-slate-200 px-3 text-sm text-slate-600"><input v-model="form.isVisible" type="checkbox" />显示</label>
      <button class="inline-flex h-11 items-center justify-center gap-2 rounded-[8px] bg-cyan-600 px-4 font-medium text-white hover:bg-cyan-700 disabled:opacity-50" :disabled="!canCreate"><Plus class="h-4 w-4" />新增站点</button>
      <p v-if="error" class="text-sm text-rose-600 lg:col-span-6">{{ error }}</p>
    </form>

    <section class="grid gap-3">
      <article v-for="site in sites.data.value" :key="site.id" class="grid gap-3 rounded-[8px] border border-slate-200 bg-white p-4 lg:grid-cols-[1fr_160px_90px] lg:items-center">
        <div class="min-w-0"><div class="truncate font-medium text-slate-950">{{ site.title }}</div><div class="truncate text-sm text-slate-500">{{ site.url }}</div><div class="mt-1 text-xs text-slate-400">{{ site.category?.name ?? '未分类' }}</div></div>
        <div class="text-sm text-slate-500">点击 {{ site.clickCount }}</div>
        <button class="inline-flex h-9 items-center justify-center gap-2 rounded-[8px] bg-rose-50 px-3 text-sm font-medium text-rose-600 hover:bg-rose-100" @click="removeSite(site.id)"><Trash2 class="h-4 w-4" />删除</button>
      </article>
    </section>
  </div>
</template>

<style scoped>
.admin-input { height: 2.75rem; border-radius: 8px; border: 1px solid rgb(226 232 240); padding: 0 0.75rem; outline: none; }
.admin-input:focus { border-color: rgb(34 211 238); box-shadow: 0 0 0 4px rgb(207 250 254); }
</style>
