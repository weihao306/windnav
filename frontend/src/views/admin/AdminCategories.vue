<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Trash2 } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
import { deleteData, getData, postData } from '../../api/client'
import type { Category } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const success = ref('')
const saving = ref(false)
const form = reactive({ name: '', slug: '', description: '', icon: '', color: '#06b6d4', sortOrder: 0, isVisible: true })
const categories = useQuery({ queryKey: ['admin-categories'], queryFn: () => getData<Category[]>('/admin/categories') })

function resetForm() {
  form.name = ''
  form.slug = ''
  form.description = ''
  form.icon = ''
  form.color = '#06b6d4'
  form.sortOrder = 0
  form.isVisible = true
}

async function createCategory() {
  error.value = ''
  success.value = ''
  saving.value = true
  try {
    await postData('/admin/categories', form)
    resetForm()
    await queryClient.invalidateQueries({ queryKey: ['admin-categories'] })
    success.value = '分类已添加'
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  } finally {
    saving.value = false
  }
}

async function removeCategory(id: number) {
  if (!window.confirm('确定删除这个分类吗？')) return
  await deleteData(`/admin/categories/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-categories'] })
  success.value = '分类已删除'
}
</script>

<template>
  <div class="grid gap-5">
    <header class="admin-page-header">
      <div><h1>分类管理</h1><p>用清晰的内容分组组织首页导航入口。</p></div>
      <span class="admin-status">{{ categories.data.value?.length ?? 0 }} 个分类</span>
    </header>

    <form class="admin-panel grid gap-4 p-5 md:grid-cols-6" @submit.prevent="createCategory">
      <input v-model="form.name" required placeholder="名称" class="admin-input md:col-span-2" />
      <input v-model="form.slug" required placeholder="slug" class="admin-input md:col-span-2" />
      <input v-model.number="form.sortOrder" type="number" placeholder="排序" class="admin-input" />
      <label class="admin-checkbox"><input v-model="form.isVisible" type="checkbox" />公开显示</label>
      <input v-model="form.description" placeholder="描述" class="admin-input md:col-span-3" />
      <input v-model="form.icon" placeholder="图标名称" class="admin-input" />
      <input v-model="form.color" type="color" class="h-11 rounded-[8px] border border-slate-200 bg-white p-1" />
      <button class="admin-button" type="submit" :disabled="saving"><Plus class="h-4 w-4" />{{ saving ? '正在添加...' : '新增分类' }}</button>
      <p v-if="error" class="admin-alert admin-alert-error md:col-span-6">{{ error }}</p>
      <p v-if="success" class="admin-alert admin-alert-success md:col-span-6">{{ success }}</p>
    </form>

    <section class="admin-panel overflow-hidden">
      <div v-for="category in categories.data.value" :key="category.id" class="admin-list-row md:grid-cols-[1fr_1fr_120px_80px]">
        <div><div class="font-medium text-slate-950">{{ category.name }}</div><div class="text-sm text-slate-500">{{ category.description || category.slug }}</div></div>
        <div class="text-sm text-slate-500">{{ category.slug }}</div>
        <div class="text-sm text-slate-500">排序 {{ category.sortOrder }}</div>
        <button class="admin-button-danger" type="button" @click="removeCategory(category.id)"><Trash2 class="h-4 w-4" />删除</button>
      </div>
      <div v-if="!categories.data.value?.length" class="admin-empty">还没有分类。</div>
    </section>
  </div>
</template>
