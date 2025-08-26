<script setup lang="ts">
import { ref, computed } from 'vue'
import ProTable from '@/components/protable/ProTable.vue'
import type { ProColumn } from '@/components/protable/types'
import { getMenuTree, createMenu, updateMenu, deleteMenu, type MenuItem } from '@/api/menu'
import { listMenuPermissions, createPermission, updatePermission, deletePermission, getRolePermissions, assignPermissionsToRole, removePermissionsFromRole, type PermissionItem, type CreatePermissionPayload, type UpdatePermissionPayload } from '@/api/permission'
import { listRoles, type RoleListResult } from '@/api/role'

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
  { title: '权限数', dataIndex: 'meta.perms', render: (r:any)=> r.meta?.perms?.length || 0 },
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

// ====== 权限管理（抽屉） ======
const permDrawerVisible = ref(false)
const currentMenu = ref<MenuItem | null>(null)
const permissions = ref<PermissionItem[]>([])
const permLoading = ref(false)
const permFormVisible = ref(false)
const isPermEdit = ref(false)
const editingPermId = ref<number | null>(null)
const permForm = ref<Partial<CreatePermissionPayload & UpdatePermissionPayload>>({})

// 角色权限分配（在菜单权限抽屉内）
const roleOptions = ref<{label:string; value:number}[]>([])
const selectedRoleId = ref<number | null>(null)
const rolePermissionIdsForMenu = ref<number[]>([])
const selectedPermissionIdsForRole = ref<number[]>([])
const rolePermSaving = ref(false)

async function loadRolesForPermission() {
  const res: RoleListResult = await listRoles({ page:1, page_size: 100 })
  roleOptions.value = res.items.map(r => ({ label: r.name, value: r.id }))
}

async function onSelectRoleLoadPermissions() {
  if (!selectedRoleId.value || !currentMenu.value) return
  // 获取该角色的全部权限
  const rolePerms = await getRolePermissions(String(selectedRoleId.value))
  // 当前菜单下的权限ID集合
  const menuPermIds = new Set(permissions.value.map(p => p.id))
  rolePermissionIdsForMenu.value = rolePerms.filter(p => menuPermIds.has(p.id)).map(p => p.id)
  selectedPermissionIdsForRole.value = [...rolePermissionIdsForMenu.value]
}

async function saveRoleMenuPermissions() {
  if (!selectedRoleId.value || !currentMenu.value) return
  rolePermSaving.value = true
  try {
    const original = new Set(rolePermissionIdsForMenu.value)
    const current = new Set(selectedPermissionIdsForRole.value)
    const toAdd:number[] = []
    const toRemove:number[] = []
    current.forEach(id => { if(!original.has(id)) toAdd.push(id) })
    original.forEach(id => { if(!current.has(id)) toRemove.push(id) })
    if (toAdd.length) await assignPermissionsToRole(String(selectedRoleId.value), toAdd)
    if (toRemove.length) await removePermissionsFromRole(String(selectedRoleId.value), toRemove)
    // 更新原始集合
    rolePermissionIdsForMenu.value = [...selectedPermissionIdsForRole.value]
  } finally {
    rolePermSaving.value = false
  }
}

async function openPermissionDrawer(menu: MenuItem) {
  currentMenu.value = menu
  permDrawerVisible.value = true
  await loadPermissions()
  if (!roleOptions.value.length) loadRolesForPermission()
}

async function loadPermissions() {
  if (!currentMenu.value) return
  permLoading.value = true
  try {
    permissions.value = await listMenuPermissions(currentMenu.value.id)
  } finally {
    permLoading.value = false
  }
}

function openCreatePermission() {
  isPermEdit.value = false
  editingPermId.value = null
  permForm.value = { status: 'active' }
  permFormVisible.value = true
}

function openEditPermission(p: PermissionItem) {
  isPermEdit.value = true
  editingPermId.value = p.id
  permForm.value = { name: p.name, code: p.code, resource: p.resource, action: p.action, status: p.status, description: p.description || undefined }
  permFormVisible.value = true
}

async function submitPermission() {
  if (!currentMenu.value) return
  if (!permForm.value.name || !permForm.value.code || !permForm.value.resource || !permForm.value.action) return
  if (isPermEdit.value && editingPermId.value != null) {
    await updatePermission(editingPermId.value, permForm.value as UpdatePermissionPayload)
  } else {
    await createPermission(currentMenu.value.id, permForm.value as CreatePermissionPayload)
  }
  permFormVisible.value = false
  await loadPermissions()
  // 刷新主菜单树以更新 perms 数量
  // 通过触发 ProTable 内部 refresh：方案——修改 request 引用或简单再次调用 request（这里直接调用 request() 不会刷新组件内部数据）。
  // 所以我们设置一个 reloadKey 绑定在 ProTable :key 来强制重渲染。
  reloadKey.value++
}

async function deletePermissionItem(id: number) {
  await deletePermission(id)
  await loadPermissions()
  reloadKey.value++
}

// 强制重载 ProTable 的 key
const reloadKey = ref(0)
</script>

<template>
  <ProTable
    :key="reloadKey"
    :title="title"
    row-key="id"
    :columns="columns"
    :request="request"
    :on-create="onCreate"
    :on-update="onUpdate"
    :on-delete="onDeleteRow"
    :is-tree="true"
  >
    <template #col-name="{ record }">
      <span>{{ record.name }}</span>
      <a-tag v-if="record.meta?.perms?.length" color="arcoblue" size="small" style="margin-left:4px">{{ record.meta.perms.length }}</a-tag>
    </template>
    <template #actions="{ record }">
      <a-space>
        <a-button size="mini" @click="openPermissionDrawer(record)">权限</a-button>
        <!-- 透传原来的编辑/删除由 ProTable 内部操作列处理，这里不重复 -->
      </a-space>
    </template>
  </ProTable>

  <!-- 权限抽屉 -->
  <a-drawer :visible="permDrawerVisible" @cancel="permDrawerVisible=false" :title="currentMenu?.name + ' - 权限管理'" width="520">
    <template #footer>
      <a-space>
        <a-button @click="permDrawerVisible=false">关闭</a-button>
        <a-button type="primary" @click="openCreatePermission">新增权限</a-button>
      </a-space>
    </template>
    <a-spin :loading="permLoading">
      <!-- 角色勾选区 -->
      <div style="margin-bottom:12px; padding:8px; border:1px solid var(--color-border-2); border-radius:6px;">
        <a-space direction="vertical" fill :size="8">
          <a-space wrap>
            <a-select v-model="selectedRoleId" placeholder="选择角色查看/分配此菜单权限" style="min-width:200px" allow-clear @change="onSelectRoleLoadPermissions">
              <a-option v-for="r in roleOptions" :key="r.value" :value="r.value">{{ r.label }}</a-option>
            </a-select>
            <a-button size="mini" type="primary" :disabled="!selectedRoleId" :loading="rolePermSaving" @click="saveRoleMenuPermissions">保存角色权限</a-button>
          </a-space>
          <div v-if="selectedRoleId">
            <a-checkbox-group v-model="selectedPermissionIdsForRole">
              <a-space wrap>
                <a-checkbox v-for="p in permissions" :key="p.id" :value="p.id">{{ p.name }}</a-checkbox>
              </a-space>
            </a-checkbox-group>
          </div>
        </a-space>
      </div>
      <a-table :data="permissions" :pagination="false" size="small" row-key="id">
        <a-table-column title="名称" data-index="name" />
        <a-table-column title="代码" data-index="code" />
        <a-table-column title="资源" data-index="resource" />
        <a-table-column title="动作" data-index="action" />
        <a-table-column title="状态" data-index="status" />
        <a-table-column title="操作" :width="140">
          <template #cell="{ record }">
            <a-space>
              <a-button size="mini" @click="openEditPermission(record)">编辑</a-button>
              <a-popconfirm content="确认删除?" type="warning" @ok="() => deletePermissionItem(record.id)">
                <a-button size="mini" status="danger">删</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table-column>
      </a-table>
    </a-spin>
  </a-drawer>

  <!-- 权限表单弹窗 -->
  <a-modal v-model:visible="permFormVisible" :title="isPermEdit ? '编辑权限' : '新增权限'" @ok="submitPermission">
    <a-form :model="permForm" layout="vertical">
      <a-form-item label="名称" required>
        <a-input v-model="permForm.name" placeholder="名称" />
      </a-form-item>
      <a-form-item label="代码" required>
        <a-input v-model="permForm.code" placeholder="如 system:menu:create" />
      </a-form-item>
      <a-form-item label="资源" required>
        <a-input v-model="permForm.resource" placeholder="资源标识，如 menu" />
      </a-form-item>
      <a-form-item label="动作" required>
        <a-input v-model="permForm.action" placeholder="动作，如 create" />
      </a-form-item>
      <a-form-item label="状态">
        <a-select v-model="permForm.status" placeholder="状态">
          <a-option value="active">启用</a-option>
          <a-option value="inactive">禁用</a-option>
        </a-select>
      </a-form-item>
      <a-form-item label="描述">
        <a-textarea v-model="permForm.description" placeholder="描述" allow-clear :auto-size="{minRows:2,maxRows:4}" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>
