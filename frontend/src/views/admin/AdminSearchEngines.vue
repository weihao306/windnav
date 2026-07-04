<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Save, Trash2 } from 'lucide-vue-next'
import { computed, reactive, ref } from 'vue'
import { deleteData, getData, postData, putData } from '../../api/client'
import type { SearchEngine } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const editingId = ref<number | null>(null)
const form = reactive({ name: '', slug: '', searchUrl: '', icon: '', sortOrder: 0, isDefault: false, isVisible: true })
const engines = useQuery({ queryKey: ['admin-search-engines'], queryFn: () => getData<SearchEngine[]>('/admin/search-engines') })
const canSave = computed(() => form.name.trim() && form.slug.trim() && form.searchUrl.trim())

function resetForm() {
  editingId.value = null
  form.name = ''
  form.slug = ''
  form.searchUrl = ''
  form.icon = ''
  form.sortOrder = 0
  form.isDefault = false
  form.isVisible = true
  error.value = ''
}

function editEngine(engine: SearchEngine) {
  editingId.value = engine.id
  form.name = engine.name
  form.slug = engine.slug
  form.searchUrl = engine.searchUrl
  form.icon = engine.icon
  form.sortOrder = engine.sortOrder
  form.isDefault = engine.isDefault
  form.isVisible = engine.isVisible
  error.value = ''
}

async function saveEngine() {
  error.value = ''
  try {
    if (editingId.value) {
      await putData(`/admin/search-engines/${editingId.value}`, form)
    } else {
      await postData('/admin/search-engines', form)
    }
    resetForm()
    await queryClient.invalidateQueries({ queryKey: ['admin-search-engines'] })
    await queryClient.invalidateQueries({ queryKey: ['public-search-engines'] })
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  }
}

async function removeEngine(id: number) {
  await deleteData(`/admin/search-engines/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-search-engines'] })
  await queryClient.invalidateQueries({ queryKey: ['public-search-engines'] })
}
</script>

<template>
  <div class="grid gap-5">
    <header class="rounded-[8px] border border-slate-200 bg-white p-5">
      <h1 class="text-2xl font-semibold text-slate-950">搜索引擎</h1>
      <p class="mt-2 text-sm text-slate-500">添加网络搜索引擎，搜索 URL 使用 <code>{query}</code> 作为关键词占位。</p>
    </header>

    <form class="grid gap-3 rounded-[8px] border border-slate-200 bg-white p-5 lg:grid-cols-6" @submit.prevent="saveEngine">
      <input v-model="form.name" required placeholder="名称，如 Google" class="admin-input lg:col-span-2" />
      <input v-model="form.slug" required placeholder="标识，如 google" class="admin-input lg:col-span-2" />
      <input v-model="form.icon" placeholder="图标文字" class="admin-input lg:col-span-2" />
      <input v-model="form.searchUrl" required placeholder="https://www.google.com/search?q={query}" class="admin-input lg:col-span-4" />
      <input v-model.number="form.sortOrder" type="number" placeholder="排序" class="admin-input" />
      <label class="flex h-11 items-center gap-2 rounded-[8px] border border-slate-200 px-3 text-sm text-slate-600"><input v-model="form.isDefault" type="checkbox" />默认</label>
      <label class="flex h-11 items-center gap-2 rounded-[8px] border border-slate-200 px-3 text-sm text-slate-600"><input v-model="form.isVisible" type="checkbox" />显示</label>
      <div class="flex gap-2 lg:col-span-6">
        <button class="inline-flex h-11 items-center justify-center gap-2 rounded-[8px] bg-cyan-600 px-4 font-medium text-white hover:bg-cyan-700 disabled:opacity-50" :disabled="!canSave">
          <Save v-if="editingId" class="h-4 w-4" />
          <Plus v-else class="h-4 w-4" />
          {{ editingId ? '保存修改' : '新增搜索引擎' }}
        </button>
        <button v-if="editingId" type="button" class="h-11 rounded-[8px] border border-slate-200 px-4 font-medium text-slate-600 hover:bg-slate-50" @click="resetForm">取消编辑</button>
      </div>
      <p v-if="error" class="text-sm text-rose-600 lg:col-span-6">{{ error }}</p>
    </form>

    <section class="grid gap-3">
      <article v-for="engine in engines.data.value" :key="engine.id" class="grid gap-3 rounded-[8px] border border-slate-200 bg-white p-4 lg:grid-cols-[1fr_160px_170px] lg:items-center">
        <div class="min-w-0">
          <div class="flex items-center gap-2 font-medium text-slate-950">
            <span class="inline-flex h-8 w-8 items-center justify-center rounded-[8px] bg-cyan-50 text-sm font-bold text-cyan-700">{{ engine.icon || engine.name.slice(0, 1) }}</span>
            {{ engine.name }}
            <span v-if="engine.isDefault" class="rounded-full bg-amber-50 px-2 py-0.5 text-xs font-medium text-amber-700">默认</span>
            <span v-if="!engine.isVisible" class="rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-500">隐藏</span>
          </div>
          <div class="mt-1 truncate text-sm text-slate-500">{{ engine.searchUrl }}</div>
        </div>
        <div class="text-sm text-slate-500">排序 {{ engine.sortOrder }}</div>
        <div class="flex gap-2">
          <button class="h-9 rounded-[8px] bg-slate-100 px-3 text-sm font-medium text-slate-600 hover:bg-slate-200" @click="editEngine(engine)">编辑</button>
          <button class="inline-flex h-9 items-center justify-center gap-2 rounded-[8px] bg-rose-50 px-3 text-sm font-medium text-rose-600 hover:bg-rose-100" @click="removeEngine(engine.id)"><Trash2 class="h-4 w-4" />删除</button>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped>
.admin-input { height: 2.75rem; border-radius: 8px; border: 1px solid rgb(226 232 240); padding: 0 0.75rem; outline: none; }
.admin-input:focus { border-color: rgb(34 211 238); box-shadow: 0 0 0 4px rgb(207 250 254); }
</style>
