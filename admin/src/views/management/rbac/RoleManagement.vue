<script setup lang="ts">
import { ref, computed, h } from 'vue'
import ProTable from '@/components/protable/ProTable.vue'
import type { ProColumn } from '@/components/protable/types'
import { listRoles, createRole, updateRole, deleteRole, type RoleItem } from '@/api/role'
import { listPermissions, getRolePermissions, assignPermissionsToRole, removePermissionsFromRole, type PermissionItem } from '@/api/permission'
import { getRoleMenusTree, type MenuItem, getMenuTree, getRoleMenus, assignMenusToRole, removeMenusFromRole } from '@/api/menu'

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
  const { page = 1, page_size = 10, name, code, status } = params
  const res = await listRoles({ page, page_size, name, code, status })
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

async function onUpdate(id: number, values: Partial<RoleItem>) {
  if (!values.name || !values.code) throw new Error('缺少必填项')
  return updateRole(id, {
    name: values.name,
    code: values.code,
    description: values.description || undefined,
    status: values.status || 'active',
  })
}

async function onDeleteRow(id: number) {
  return deleteRole(id)
}

const title = computed(() => '角色管理')

// ===== 权限分配抽屉状态 =====
const permDrawerVisible = ref(false)
const currentRole = ref<RoleItem | null>(null)
const allPermissions = ref<PermissionItem[]>([])
const rolePermissionIds = ref<number[]>([])
const selectedPermissionIds = ref<number[]>([])
const permLoading = ref(false)
const saving = ref(false)

async function openPermDrawer(record: RoleItem) {
  currentRole.value = record
  permDrawerVisible.value = true
  permLoading.value = true
  try {
    const [all, owned] = await Promise.all([
      listPermissions({ page: 1, page_size: 2000 }),
      getRolePermissions(String(record.id)),
    ])
    allPermissions.value = all.items
    rolePermissionIds.value = owned.map(p => p.id)
    selectedPermissionIds.value = [...rolePermissionIds.value]
  } finally {
    permLoading.value = false
  }
}

async function saveRolePermissions() {
  if (!currentRole.value) return
  saving.value = true
  try {
    const original = new Set(rolePermissionIds.value)
    const current = new Set(selectedPermissionIds.value)
    const toAdd: number[] = []
    const toRemove: number[] = []
    current.forEach(id => { if (!original.has(id)) toAdd.push(id) })
    original.forEach(id => { if (!current.has(id)) toRemove.push(id) })
    if (toAdd.length) await assignPermissionsToRole(String(currentRole.value.id), toAdd)
    if (toRemove.length) await removePermissionsFromRole(String(currentRole.value.id), toRemove)
    permDrawerVisible.value = false
  } finally {
    saving.value = false
  }
}

function filterPermissions(keyword: string) {
  if (!keyword) return allPermissions.value
  const k = keyword.toLowerCase()
  return allPermissions.value.filter(p => p.name.toLowerCase().includes(k) || p.code.toLowerCase().includes(k))
}

const searchPermKey = ref('')
const displayedPermissions = computed(() => filterPermissions(searchPermKey.value))

// ===== 角色菜单树展示 =====
const menuTreeDrawerVisible = ref(false)
const roleMenuTree = ref<MenuItem[]>([])
const loadingMenuTree = ref(false)

async function openRoleMenuTree(record: RoleItem) {
  currentRole.value = record
  menuTreeDrawerVisible.value = true
  loadingMenuTree.value = true
  try {
    roleMenuTree.value = await getRoleMenusTree(String(record.id))
  } finally {
    loadingMenuTree.value = false
  }
}

function renderMenuTree(nodes: MenuItem[]): any[] {
  return (nodes || []).map(n => ({
    key: n.id,
    title: () => h('div', { style: 'display:flex;align-items:center;gap:6px;' }, [
      h('span', n.name),
      (n.meta && n.meta.perms && n.meta.perms.length) ? h('span', { style: 'display:flex;gap:4px;flex-wrap:wrap;' }, n.meta.perms.slice(0,3).map(p=>h('span',{ class:'perm-chip'}, p))) : null,
      (n.meta && n.meta.perms && n.meta.perms.length > 3) ? h('span', { class:'perm-more' }, '+'+(n.meta.perms.length-3)) : null
    ]),
    children: n.children ? renderMenuTree(n.children) : undefined
  }))
}

// ===== 菜单绑定抽屉状态 =====
const menuBindDrawerVisible = ref(false)
const menuTreeData = ref<MenuItem[]>([])
const checkedMenuIds = ref<(string|number)[]>([])
const originalMenuIds = ref<Set<number>>(new Set())
const loadingMenuBind = ref(false)

async function openMenuBindDrawer(record: RoleItem) {
  currentRole.value = record
  menuBindDrawerVisible.value = true
  loadingMenuBind.value = true
  try {
    const [allTree, ownedFlat] = await Promise.all([
      getMenuTree(),
      getRoleMenus(String(record.id)),
    ])
    menuTreeData.value = allTree
    originalMenuIds.value = new Set(ownedFlat.map(m => Number(m.id)))
    checkedMenuIds.value = [...originalMenuIds.value]
  } finally {
    loadingMenuBind.value = false
  }
}

function extractAllIds(nodes: MenuItem[]): number[] {
  const res: number[] = []
  const walk = (arr: MenuItem[]) => {
    arr.forEach(n => {
      res.push(Number(n.id))
      if (n.children && n.children.length) walk(n.children)
    })
  }
  walk(nodes)
  return res
}

async function saveRoleMenus() {
  if (!currentRole.value) return
  loadingMenuBind.value = true
  try {
    const current = new Set(checkedMenuIds.value.map(id => Number(id)))
    const toAdd: number[] = []
    const toRemove: number[] = []
    current.forEach(id => { if (!originalMenuIds.value.has(id)) toAdd.push(id) })
    originalMenuIds.value.forEach(id => { if (!current.has(id)) toRemove.push(id) })
    if (toAdd.length) await assignMenusToRole(String(currentRole.value.id), toAdd)
    if (toRemove.length) await removeMenusFromRole(String(currentRole.value.id), toRemove)
    menuBindDrawerVisible.value = false
  } finally {
    loadingMenuBind.value = false
  }
}
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
  >
    <template #actions="{ record }">
      <a-space>
        <a-button size="mini" @click="openPermDrawer(record)">权限</a-button>
        <a-button size="mini" type="outline" @click="openMenuBindDrawer(record)">菜单</a-button>
      </a-space>
    </template>
  </ProTable>

  <!-- 权限分配抽屉 -->
  <a-drawer :visible="permDrawerVisible" :width="640" :title="(currentRole?.name || '') + ' - 权限分配'" @cancel="permDrawerVisible=false">
    <template #footer>
      <a-space>
        <a-button @click="permDrawerVisible=false" :disabled="saving">取消</a-button>
        <a-button type="primary" @click="saveRolePermissions" :loading="saving">保存</a-button>
      </a-space>
    </template>
    <a-spin :loading="permLoading">
      <div style="margin-bottom:12px;display:flex;gap:8px;align-items:center;">
        <a-button size="mini" status="warning" @click="openRoleMenuTree(currentRole!)" :disabled="!currentRole">菜单树</a-button>
        <a-input v-model="searchPermKey" placeholder="搜索名称或代码" allow-clear style="flex:1;" />
        <a-tag color="arcoblue">共 {{ allPermissions.length }} 条</a-tag>
      </div>
      <div style="max-height:480px;overflow:auto;border:1px solid var(--color-border-2);padding:12px;border-radius:6px;">
        <a-checkbox-group v-model="selectedPermissionIds">
          <a-space direction="vertical" :size="4" fill>
            <template v-for="p in displayedPermissions" :key="p.id">
              <a-checkbox :value="p.id">
                <a-space size="mini">
                  <span style="font-weight:500">{{ p.name }}</span>
                  <a-tag size="small" color="gray">{{ p.code }}</a-tag>
                  <span style="color:var(--color-text-3);font-size:12px;">{{ p.resource }}/{{ p.action }}</span>
                </a-space>
              </a-checkbox>
            </template>
          </a-space>
        </a-checkbox-group>
      </div>
    </a-spin>
  </a-drawer>

  <!-- 菜单绑定抽屉 -->
  <a-drawer :visible="menuBindDrawerVisible" :width="680" :title="(currentRole?.name||'') + ' - 菜单绑定'" @cancel="menuBindDrawerVisible=false">
    <template #footer>
      <a-space>
        <a-button @click="menuBindDrawerVisible=false" :disabled="loadingMenuBind">取消</a-button>
        <a-button type="primary" @click="saveRoleMenus" :loading="loadingMenuBind">保存</a-button>
      </a-space>
    </template>
    <a-spin :loading="loadingMenuBind">
      <a-space style="margin-bottom:8px;" wrap>
  <a-button size="mini" @click="() => { checkedMenuIds.value = extractAllIds(menuTreeData.value) }">全选</a-button>
  <a-button size="mini" @click="() => { checkedMenuIds.value = [] }">清空</a-button>
      </a-space>
      <a-tree
        :data="menuTreeData"
        checkable
        v-model:checked-keys="checkedMenuIds"
        :field-names="{ key: 'id', title: 'name', children: 'children' }"
        :default-expand-all="true"
        style="max-height:70vh;overflow:auto;border:1px solid var(--color-border-2);padding:8px;border-radius:6px;"
      />
    </a-spin>
  </a-drawer>
</template>
