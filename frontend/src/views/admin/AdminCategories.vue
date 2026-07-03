<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Trash2 } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
import { deleteData, getData, postData } from '../../api/client'
import type { Category } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
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
  try {
    await postData('/admin/categories', form)
    resetForm()
    await queryClient.invalidateQueries({ queryKey: ['admin-categories'] })
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  }
}

async function removeCategory(id: number) {
  await deleteData(`/admin/categories/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-categories'] })
}
</script>

<template>
  <div class="grid gap-5">
    <header class="rounded-[8px] border border-slate-200 bg-white p-5">
      <h1 class="text-2xl font-semibold text-slate-950">分类管理</h1>
    </header>

    <form class="grid gap-3 rounded-[8px] border border-slate-200 bg-white p-5 md:grid-cols-6" @submit.prevent="createCategory">
      <input v-model="form.name" required placeholder="名称" class="admin-input md:col-span-2" />
      <input v-model="form.slug" required placeholder="slug" class="admin-input md:col-span-2" />
      <input v-model.number="form.sortOrder" type="number" placeholder="排序" class="admin-input" />
      <label class="flex h-11 items-center gap-2 rounded-[8px] border border-slate-200 px-3 text-sm text-slate-600"><input v-model="form.isVisible" type="checkbox" />显示</label>
      <input v-model="form.description" placeholder="描述" class="admin-input md:col-span-3" />
      <input v-model="form.icon" placeholder="图标名称" class="admin-input" />
      <input v-model="form.color" type="color" class="h-11 rounded-[8px] border border-slate-200 bg-white p-1" />
      <button class="inline-flex h-11 items-center justify-center gap-2 rounded-[8px] bg-cyan-600 px-4 font-medium text-white hover:bg-cyan-700"><Plus class="h-4 w-4" />新增</button>
      <p v-if="error" class="text-sm text-rose-600 md:col-span-6">{{ error }}</p>
    </form>

    <section class="overflow-hidden rounded-[8px] border border-slate-200 bg-white">
      <div v-for="category in categories.data.value" :key="category.id" class="grid gap-3 border-b border-slate-100 p-4 last:border-b-0 md:grid-cols-[1fr_1fr_120px_80px] md:items-center">
        <div><div class="font-medium text-slate-950">{{ category.name }}</div><div class="text-sm text-slate-500">{{ category.description || category.slug }}</div></div>
        <div class="text-sm text-slate-500">{{ category.slug }}</div>
        <div class="text-sm text-slate-500">排序 {{ category.sortOrder }}</div>
        <button class="inline-flex h-9 items-center justify-center gap-2 rounded-[8px] bg-rose-50 px-3 text-sm font-medium text-rose-600 hover:bg-rose-100" @click="removeCategory(category.id)"><Trash2 class="h-4 w-4" />删除</button>
      </div>
    </section>
  </div>
</template>

<style scoped>
.admin-input { height: 2.75rem; border-radius: 8px; border: 1px solid rgb(226 232 240); padding: 0 0.75rem; outline: none; }
.admin-input:focus { border-color: rgb(34 211 238); box-shadow: 0 0 0 4px rgb(207 250 254); }
</style>
