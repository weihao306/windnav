<script setup lang="ts">
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { Plus, Trash2 } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
import { deleteData, getData, postData } from '../../api/client'
import type { Tag } from '../../api/types'

const queryClient = useQueryClient()
const error = ref('')
const success = ref('')
const saving = ref(false)
const form = reactive({ name: '', slug: '', color: '#22c55e' })
const tags = useQuery({ queryKey: ['admin-tags'], queryFn: () => getData<Tag[]>('/admin/tags') })

async function createTag() {
  error.value = ''
  success.value = ''
  saving.value = true
  try {
    await postData('/admin/tags', form)
    form.name = ''
    form.slug = ''
    form.color = '#22c55e'
    await queryClient.invalidateQueries({ queryKey: ['admin-tags'] })
    success.value = '标签已添加'
  } catch (err) {
    error.value = err instanceof Error ? err.message : '保存失败'
  } finally {
    saving.value = false
  }
}

async function removeTag(id: number) {
  if (!window.confirm('确定删除这个标签吗？')) return
  await deleteData(`/admin/tags/${id}`)
  await queryClient.invalidateQueries({ queryKey: ['admin-tags'] })
  success.value = '标签已删除'
}
</script>

<template>
  <div class="grid gap-5">
    <header class="admin-page-header">
      <div><h1>标签管理</h1><p>用轻量标签补充站点的检索和识别信息。</p></div>
      <span class="admin-status">{{ tags.data.value?.length ?? 0 }} 个标签</span>
    </header>
    <form class="admin-panel grid gap-4 p-5 md:grid-cols-[1fr_1fr_80px_120px]" @submit.prevent="createTag">
      <input v-model="form.name" required placeholder="名称" class="admin-input" />
      <input v-model="form.slug" required placeholder="slug" class="admin-input" />
      <input v-model="form.color" type="color" class="h-11 rounded-[8px] border border-slate-200 bg-white p-1" />
      <button class="admin-button" type="submit" :disabled="saving"><Plus class="h-4 w-4" />{{ saving ? '正在添加...' : '新增标签' }}</button>
      <p v-if="error" class="admin-alert admin-alert-error md:col-span-4">{{ error }}</p>
      <p v-if="success" class="admin-alert admin-alert-success md:col-span-4">{{ success }}</p>
    </form>
    <section class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
      <article v-for="tag in tags.data.value" :key="tag.id" class="admin-panel flex items-center justify-between p-4">
        <div class="flex items-center gap-3"><span class="h-3 w-3 rounded-full" :style="{ background: tag.color || '#22c55e' }" /><div><div class="font-medium">{{ tag.name }}</div><div class="text-sm text-slate-500">{{ tag.slug }}</div></div></div>
        <button class="admin-button-icon" type="button" aria-label="删除标签" @click="removeTag(tag.id)"><Trash2 class="h-4 w-4" /></button>
      </article>
      <div v-if="!tags.data.value?.length" class="admin-empty sm:col-span-2 xl:col-span-3">还没有标签。</div>
    </section>
  </div>
</template>
