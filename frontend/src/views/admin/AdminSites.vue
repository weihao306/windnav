<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Trash2 } from 'lucide-vue-next'
import { computed, reactive, ref } from 'vue'
import { deleteData, getData, postData } from '../../api/client'
import type { Category, Site, Tag } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const success = ref('')
const saving = ref(false)
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
  success.value = ''
  saving.value = true
  try {
    await postData('/admin/sites', form)
    resetForm()
    await queryClient.invalidateQueries({ queryKey: ['admin-sites'] })
    success.value = '站点已添加'
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  } finally {
    saving.value = false
  }
}

async function removeSite(id: number) {
  if (!window.confirm('确定删除这个站点吗？删除后无法恢复。')) return
  error.value = ''
  success.value = ''
  await deleteData(`/admin/sites/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-sites'] })
  success.value = '站点已删除'
}
</script>

<template>
  <div class="grid gap-5">
    <header class="admin-page-header">
      <div><h1>站点管理</h1><p>维护首页展示的导航入口、分类与常用标签。</p></div>
      <span class="admin-status">{{ sites.data.value?.length ?? 0 }} 个站点</span>
    </header>

    <form class="admin-panel grid gap-4 p-5 lg:grid-cols-6" @submit.prevent="createSite">
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
      <label class="admin-checkbox"><input v-model="form.isPinned" type="checkbox" />置顶推荐</label>
      <label class="admin-checkbox"><input v-model="form.isVisible" type="checkbox" />公开显示</label>
      <button class="admin-button" type="submit" :disabled="!canCreate || saving"><Plus class="h-4 w-4" />{{ saving ? '正在添加...' : '新增站点' }}</button>
      <p v-if="error" class="admin-alert admin-alert-error lg:col-span-6">{{ error }}</p>
      <p v-if="success" class="admin-alert admin-alert-success lg:col-span-6">{{ success }}</p>
    </form>

    <section class="admin-panel overflow-hidden">
      <article v-for="site in sites.data.value" :key="site.id" class="admin-list-row lg:grid-cols-[1fr_160px_100px]">
        <div class="min-w-0"><div class="truncate font-medium text-slate-950">{{ site.title }}</div><div class="truncate text-sm text-slate-500">{{ site.url }}</div><div class="mt-1 text-xs text-slate-400">{{ site.category?.name ?? '未分类' }}</div></div>
        <div class="text-sm text-slate-500">点击 {{ site.clickCount }}</div>
        <div class="admin-list-actions"><button class="admin-button-danger" type="button" @click="removeSite(site.id)"><Trash2 class="h-4 w-4" />删除</button></div>
      </article>
      <div v-if="!sites.data.value?.length" class="admin-empty">还没有站点，使用上方表单添加第一个入口。</div>
    </section>
  </div>
</template>

<style scoped>
details[open] > summary.admin-input { border-color: rgb(34 211 238); box-shadow: 0 0 0 4px rgb(207 250 254); }
</style>
