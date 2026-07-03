<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Trash2 } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
import { deleteData, getData, postData } from '../../api/client'
import type { Tag } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const form = reactive({ name: '', slug: '', color: '#22c55e' })
const tags = useQuery({ queryKey: ['admin-tags'], queryFn: () => getData<Tag[]>('/admin/tags') })

async function createTag() {
  error.value = ''
  try {
    await postData('/admin/tags', form)
    form.name = ''
    form.slug = ''
    form.color = '#22c55e'
    await queryClient.invalidateQueries({ queryKey: ['admin-tags'] })
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  }
}

async function removeTag(id: number) {
  await deleteData(`/admin/tags/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-tags'] })
}
</script>

<template>
  <div class="grid gap-5">
    <header class="rounded-[8px] border border-slate-200 bg-white p-5">
      <h1 class="text-2xl font-semibold text-slate-950">标签管理</h1>
    </header>
    <form class="grid gap-3 rounded-[8px] border border-slate-200 bg-white p-5 md:grid-cols-[1fr_1fr_80px_120px]" @submit.prevent="createTag">
      <input v-model="form.name" required placeholder="名称" class="admin-input" />
      <input v-model="form.slug" required placeholder="slug" class="admin-input" />
      <input v-model="form.color" type="color" class="h-11 rounded-[8px] border border-slate-200 bg-white p-1" />
      <button class="inline-flex h-11 items-center justify-center gap-2 rounded-[8px] bg-cyan-600 px-4 font-medium text-white hover:bg-cyan-700"><Plus class="h-4 w-4" />新增</button>
      <p v-if="error" class="text-sm text-rose-600 md:col-span-4">{{ error }}</p>
    </form>
    <section class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
      <article v-for="tag in tags.data.value" :key="tag.id" class="flex items-center justify-between rounded-[8px] border border-slate-200 bg-white p-4">
        <div class="flex items-center gap-3"><span class="h-3 w-3 rounded-full" :style="{ background: tag.color || '#22c55e' }" /><div><div class="font-medium">{{ tag.name }}</div><div class="text-sm text-slate-500">{{ tag.slug }}</div></div></div>
        <button class="rounded-[8px] p-2 text-rose-600 hover:bg-rose-50" @click="removeTag(tag.id)"><Trash2 class="h-4 w-4" /></button>
      </article>
    </section>
  </div>
</template>

<style scoped>
.admin-input { height: 2.75rem; border-radius: 8px; border: 1px solid rgb(226 232 240); padding: 0 0.75rem; outline: none; }
.admin-input:focus { border-color: rgb(34 211 238); box-shadow: 0 0 0 4px rgb(207 250 254); }
</style>
