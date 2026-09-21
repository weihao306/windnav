<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Save, Trash2 } from 'lucide-vue-next'
import { computed, reactive, ref } from 'vue'
import { deleteData, getData, postData, putData } from '../../api/client'
import type { SearchEngine } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const success = ref('')
const saving = ref(false)
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
  success.value = ''
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
  success.value = ''
  saving.value = true
  const wasEditing = Boolean(editingId.value)
  try {
    if (wasEditing) {
      await putData(`/admin/search-engines/${editingId.value}`, form)
    } else {
      await postData('/admin/search-engines', form)
    }
    resetForm()
    await queryClient.invalidateQueries({ queryKey: ['admin-search-engines'] })
    await queryClient.invalidateQueries({ queryKey: ['public-search-engines'] })
    success.value = wasEditing ? '搜索引擎已更新' : '搜索引擎已添加'
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  } finally {
    saving.value = false
  }
}

async function removeEngine(id: number) {
  if (!window.confirm('确定删除这个搜索引擎吗？')) return
  await deleteData(`/admin/search-engines/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-search-engines'] })
  await queryClient.invalidateQueries({ queryKey: ['public-search-engines'] })
  success.value = '搜索引擎已删除'
}
</script>

<template>
  <div class="grid gap-5">
    <header class="admin-page-header">
      <div><h1>搜索引擎</h1><p>配置首页搜索使用的网络搜索服务与关键词模板。</p></div>
      <span class="admin-status">{{ engines.data.value?.length ?? 0 }} 个引擎</span>
    </header>

    <form class="admin-panel grid gap-4 p-5 lg:grid-cols-6" @submit.prevent="saveEngine">
      <input v-model="form.name" required placeholder="名称，如 Google" class="admin-input lg:col-span-2" />
      <input v-model="form.slug" required placeholder="标识，如 google" class="admin-input lg:col-span-2" />
      <input v-model="form.icon" placeholder="图标文字" class="admin-input lg:col-span-2" />
      <input v-model="form.searchUrl" required placeholder="https://www.google.com/search?q={query}" class="admin-input lg:col-span-4" />
      <input v-model.number="form.sortOrder" type="number" placeholder="排序" class="admin-input" />
      <label class="admin-checkbox"><input v-model="form.isDefault" type="checkbox" />设为默认</label>
      <label class="admin-checkbox"><input v-model="form.isVisible" type="checkbox" />公开显示</label>
      <div class="flex gap-2 lg:col-span-6">
        <button class="admin-button" type="submit" :disabled="!canSave || saving">
          <Save v-if="editingId" class="h-4 w-4" />
          <Plus v-else class="h-4 w-4" />
          {{ saving ? '正在保存...' : editingId ? '保存修改' : '新增搜索引擎' }}
        </button>
        <button v-if="editingId" type="button" class="admin-button-secondary" @click="resetForm">取消编辑</button>
      </div>
      <p v-if="error" class="admin-alert admin-alert-error lg:col-span-6">{{ error }}</p>
      <p v-if="success" class="admin-alert admin-alert-success lg:col-span-6">{{ success }}</p>
    </form>

    <section class="admin-panel overflow-hidden">
      <article v-for="engine in engines.data.value" :key="engine.id" class="admin-list-row lg:grid-cols-[1fr_160px_190px]">
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
           <button class="admin-button-secondary" type="button" @click="editEngine(engine)">编辑</button>
           <button class="admin-button-danger" type="button" @click="removeEngine(engine.id)"><Trash2 class="h-4 w-4" />删除</button>
        </div>
      </article>
      <div v-if="!engines.data.value?.length" class="admin-empty">还没有搜索引擎。</div>
    </section>
  </div>
</template>
