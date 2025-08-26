<script setup lang="ts">
import { ref, computed } from 'vue'
import ProTable from '@/components/protable/ProTable.vue'
import type { ProColumn } from '@/components/protable/types'
import { listUsers, createUser, updateUser, deleteUser, type UserItem } from '@/api/user'

// 状态选项
const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
  { label: '封禁', value: 'banned' },
]

// 列定义
const columns = ref<ProColumn<UserItem>[]>([
  { title: 'ID', dataIndex: 'id', search: false },
  { title: '用户名', dataIndex: 'username', search: true, form: { required: true, placeholder: '用户名' } },
  { title: '密码', dataIndex: 'password', hideInTable: true, form: { placeholder: '密码（创建必填，修改留空不变）', props: { type: 'password' } } },
  { title: '邮箱', dataIndex: 'email', search: true, form: { required: true, placeholder: '邮箱' } },
  {
    title: '状态',
    dataIndex: 'status',
    search: true,
    form: { widget: 'select', options: statusOptions, required: true },
    render: (r) => statusOptions.find((o) => o.value === r.status)?.label || r.status,
  },
])

// request: 需要适配 ProTable 的返回结构 { list, total }
async function request(params: any) {
  const { page = 1, page_size = 10 } = params
  const res = await listUsers({ page, page_size })
  return { list: res.items, total: res.total }
}

// onCreate / onUpdate / onDelete
async function onCreate(values: Partial<UserItem & { password?: string }>) {
  if (!values.username || !values.password || !values.email) throw new Error('缺少必填项')
  return createUser({
    username: values.username,
    password: values.password,
    email: values.email,
    status: values.status || 'active',
  })
}

async function onUpdate(id: string | number, values: Partial<UserItem & { password?: string }>) {
  const payload: any = {}
  if (values.email) payload.email = values.email
  if (values.status) payload.status = values.status
  if ((values as any).password) payload.password = (values as any).password
  return updateUser(id, payload)
}

async function onDeleteRow(id: string | number) {
  return deleteUser(id)
}

const title = computed(() => '用户管理')
</script>

<template>
  <ProTable
    :title="title"
    row-key="id"
    :columns="columns"
    :request="request"
    :on-create="onCreate"
    :on-update="onUpdate"
    :on-delete="onDeleteRow"
    :pagination="true"
  />
</template>
