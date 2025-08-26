<script setup lang="ts">
import { ref, computed } from 'vue'
import ProTable from '@/components/protable/ProTable.vue'
import type { ProColumn } from '@/components/protable/types'
import { listRoles, createRole, updateRole, deleteRole, type RoleItem } from '@/api/role'

// 状态选项
const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
]

// 列定义
const columns = ref<ProColumn<RoleItem>[]>([
  { title: 'ID', dataIndex: 'id', search: false },
  { title: '名称', dataIndex: 'name', search: true, form: { required: true, placeholder: '角色名称' } },
  { title: '标识码', dataIndex: 'code', search: true, form: { required: true, placeholder: '唯一代码' } },
  { title: '描述', dataIndex: 'description', search: false, form: { widget: 'textarea', placeholder: '描述(可选)' } },
  {
    title: '状态',
    dataIndex: 'status',
    search: true,
    form: { widget: 'select', options: statusOptions, required: true },
    render: (r) => statusOptions.find((o) => o.value === r.status)?.label || r.status,
  },
  { title: '创建时间', dataIndex: 'created_at', search: false },
])

// request 适配 ProTable 返回 { list, total }
async function request(params: any) {
  const { page = 1, page_size = 10 } = params
  const res = await listRoles({ page, page_size })
  return { list: res.items, total: res.total }
}

// onCreate / onUpdate / onDelete
async function onCreate(values: Partial<RoleItem>) {
  if (!values.name || !values.code) throw new Error('缺少必填项')
  return createRole({
    name: values.name,
    code: values.code,
    description: values.description || undefined,
    status: values.status || 'active',
  })
}

async function onUpdate(id: string, values: Partial<RoleItem>) {
  if (!values.name || !values.code) throw new Error('缺少必填项')
  return updateRole(id, {
    name: values.name,
    code: values.code,
    description: values.description || undefined,
    status: values.status || 'active',
  })
}

async function onDeleteRow(id: string) {
  return deleteRole(id)
}

const title = computed(() => '角色管理')
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
