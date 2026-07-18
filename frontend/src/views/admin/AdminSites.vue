<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Trash2 } from 'lucide-vue-next'
import { computed, reactive, ref } from 'vue'
import { deleteData, getData, postData } from '../../api/client'
import type { Category, Site, Tag } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const form = reactive({ categoryId: 0, title: '', url: '', description: '', iconUrl: '', fallbackIcon: '', sortOrder: 0, isPinned: false, isVisible: true, tagIds: [] as number[] })
const categorySearch = ref('')
const tagSearch = ref('')
const sites = useQuery({ queryKey: ['admin-sites'], queryFn: () => getData<Site[]>('/admin/sites', { page_size: 100 }) })
const categories = useQuery({ queryKey: ['admin-categories'], queryFn: () => getData<Category[]>('/admin/categories') })
const tags = useQuery({ queryKey: ['admin-tags'], queryFn: () => getData<Tag[]>('/admin/tags') })
const canCreate = computed(() => form.categoryId > 0 && form.title.trim() && form.url.trim())
const filteredCategories = computed(() => {
  const keyword = categorySearch.value.trim().toLowerCase()
  const list = categories.data.value ?? []
  if (!keyword) return list
  return list.filter((category) => category.name.toLowerCase().includes(keyword))
})
const selectedCategoryName = computed(() => categories.data.value?.find((category) => category.id === form.categoryId)?.name ?? '选择分类')
const filteredTags = computed(() => {
  const keyword = tagSearch.value.trim().toLowerCase()
  const list = tags.data.value ?? []
  if (!keyword) return list
  return list.filter((tag) => tag.name.toLowerCase().includes(keyword))
})
const selectedTagNames = computed(() => {
  const selected = tags.data.value?.filter((tag) => form.tagIds.includes(tag.id)).map((tag) => tag.name) ?? []
  return selected.length ? selected.join('、') : '选择标签'
})

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
  categorySearch.value = ''
  tagSearch.value = ''
}

function selectCategory(categoryId: number) {
  form.categoryId = categoryId
  categorySearch.value = ''
}

function toggleTag(tagId: number) {
  if (form.tagIds.includes(tagId)) {
    form.tagIds = form.tagIds.filter((id) => id !== tagId)
    return
  }
  form.tagIds = [...form.tagIds, tagId]
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
      <details class="relative lg:col-span-2">
        <summary class="admin-input flex cursor-pointer list-none items-center justify-between gap-3 text-sm text-slate-700 marker:hidden">
          <span class="truncate">{{ selectedCategoryName }}</span>
          <span class="text-xs text-slate-400">▼</span>
        </summary>
        <div class="absolute z-20 mt-2 grid max-h-72 w-full gap-2 overflow-hidden rounded-[8px] border border-slate-200 bg-white p-3 shadow-xl">
          <input v-model="categorySearch" placeholder="搜索分类" class="admin-input h-10" @click.stop />
          <div class="grid max-h-48 gap-1 overflow-y-auto pr-1">
            <button v-for="category in filteredCategories" :key="category.id" type="button" class="flex items-center justify-between gap-2 rounded-[6px] px-2 py-2 text-left text-sm text-slate-700 hover:bg-slate-50" @click="selectCategory(category.id)">
              <span class="truncate">{{ category.name }}</span>
              <span v-if="form.categoryId === category.id" class="text-xs text-cyan-600">已选</span>
            </button>
            <div v-if="!filteredCategories.length" class="px-2 py-4 text-center text-sm text-slate-400">没有匹配的分类</div>
          </div>
        </div>
      </details>
      <input v-model="form.title" required placeholder="站点名称" class="admin-input lg:col-span-2" />
      <input v-model="form.url" required placeholder="https://example.com" class="admin-input lg:col-span-2" />
      <input v-model="form.description" placeholder="描述" class="admin-input lg:col-span-3" />
      <input v-model="form.iconUrl" placeholder="图标 URL" class="admin-input lg:col-span-2" />
      <input v-model="form.fallbackIcon" placeholder="备用字母" class="admin-input" />
      <details class="relative lg:col-span-3">
        <summary class="admin-input flex cursor-pointer list-none items-center justify-between gap-3 text-sm text-slate-700 marker:hidden">
          <span class="truncate">{{ selectedTagNames }}</span>
          <span class="text-xs text-slate-400">▼</span>
        </summary>
        <div class="absolute z-20 mt-2 grid max-h-72 w-full gap-2 overflow-hidden rounded-[8px] border border-slate-200 bg-white p-3 shadow-xl">
          <input v-model="tagSearch" placeholder="搜索标签" class="admin-input h-10" @click.stop />
          <div class="grid max-h-48 gap-1 overflow-y-auto pr-1">
            <button v-for="tag in filteredTags" :key="tag.id" type="button" class="flex items-center gap-2 rounded-[6px] px-2 py-2 text-left text-sm text-slate-700 hover:bg-slate-50" @click="toggleTag(tag.id)">
              <input type="checkbox" :checked="form.tagIds.includes(tag.id)" class="pointer-events-none" />
              <span class="truncate">{{ tag.name }}</span>
            </button>
            <div v-if="!filteredTags.length" class="px-2 py-4 text-center text-sm text-slate-400">没有匹配的标签</div>
          </div>
        </div>
      </details>
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
details[open] > summary.admin-input { border-color: rgb(34 211 238); box-shadow: 0 0 0 4px rgb(207 250 254); }
</style>
