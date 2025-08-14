<script setup lang="ts">
import { ref, computed } from 'vue'
import ProTable from '@/components/protable/ProTable.vue'
import type { ProColumn } from '@/components/protable/types'
import { getMenuTree, createMenu, updateMenu, deleteMenu, type MenuItem } from '@/api/menu'

// 将后端菜单类型映射为中文
const typeOptions = [
  { label: '分组', value: 0 },
  { label: '目录', value: 1 },
  { label: '菜单', value: 2 },
  { label: '按钮', value: 3 },
]

const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
]

// 列定义：同时驱动搜索与表单
const columns = ref<ProColumn<MenuItem>[]>([
  { title: '名称', dataIndex: 'name', search: true, form: { required: true, placeholder: '请输入名称' } },
  { title: '编码', dataIndex: 'code', search: true, form: { required: true, placeholder: '唯一编码' } },
  { title: '路径', dataIndex: 'path', form: { placeholder: '/path' } },
  {
    title: '类型',
    dataIndex: 'type',
    form: { widget: 'select', options: typeOptions, required: true },
    render: (r) => typeOptions.find((o) => o.value === r.type)?.label || r.type,
  },
  { title: '排序', dataIndex: 'order', form: { widget: 'number', props: { min: 0, step: 1 } } },
  {
    title: '状态',
    dataIndex: 'status',
    search: true,
    form: { widget: 'select', options: statusOptions, required: true },
    render: (r) => statusOptions.find((o) => o.value === r.status)?.label || r.status,
  },
  { title: '图标', dataIndex: 'meta.icon', form: { placeholder: 'local:icon' } },
  { title: '组件', dataIndex: 'meta.component', form: { placeholder: '/views/xxx' } },
])

// 父级选择（用树选择），需要把树形菜单转为 treeOptions
const treeOptions = ref<any[]>([])
function toTreeOptions(list: MenuItem[]): any[] {
  return (list || []).map((n) => ({
    label: n.name,
    value: n.id,
    children: n.children ? toTreeOptions(n.children) : undefined,
  }))
}

columns.value.unshift({
  title: '父级',
  dataIndex: 'parent_id',
  form: { widget: 'treeSelect', treeOptions: treeOptions.value, props: { allowClear: true, treeCheckable: false } },
})

// 请求函数：返回树结构
async function request() {
  const data = await getMenuTree()
  // 更新父级选择树
  treeOptions.value = toTreeOptions(data)
  // 将 options 写回到列
  const parentCol = columns.value.find((c) => c.dataIndex === 'parent_id')
  if (parentCol && typeof parentCol.form === 'object') parentCol.form.treeOptions = treeOptions.value
  return data
}

// CRUD 行为
async function onCreate(values: Partial<MenuItem>) {
  // 处理 meta 嵌套
  const payload: any = { ...values }
  // 如果 meta 不是对象，组装
  if (!payload.meta) {
    payload.meta = {}
  }
  return createMenu(payload)
}

async function onUpdate(id: string, values: Partial<MenuItem>) {
  const payload: any = { ...values }
  if (!payload.meta) payload.meta = {}
  return updateMenu(id, payload)
}

async function onDeleteRow(id: string) {
  return deleteMenu(id)
}

const title = computed(() => '菜单管理')
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
    :is-tree="true"
  />
</template>
