import http from './http'

export interface PermissionItem {
  id: number
  name: string
  code: string
  resource: string
  action: string
  status: string
  description?: string | null
  menu_id?: number
  created_at?: string
  updated_at?: string
}

// ============ 菜单下权限 ============
export function listMenuPermissions(menuId: string | number) {
  return http.get<PermissionItem[]>(`/menus/${menuId}/permissions`)
}

export interface CreatePermissionPayload {
  name: string
  code: string
  resource: string
  action: string
  status?: string
  description?: string | null
}
export function createPermission(menuId: string | number, data: CreatePermissionPayload) {
  return http.post(`/menus/${menuId}/permissions`, data)
}

export interface UpdatePermissionPayload {
  name?: string
  code?: string
  resource?: string
  action?: string
  status?: string
  description?: string | null
}
export function updatePermission(id: number, data: UpdatePermissionPayload) {
  return http.put(`/permissions/${id}`, data)
}

export function deletePermission(id: number) {
  return http.delete(`/permissions/${id}`)
}

// ============ 全量权限（分页） ============
export interface ListPermissionsParams { page?: number; page_size?: number }
export interface ListPermissionsResult { items: PermissionItem[]; total: number; page: number; page_size: number }
export function listPermissions(params: ListPermissionsParams = { page: 1, page_size: 1000 }) {
  return http.get<ListPermissionsResult>('/permissions', { params })
}

// ============ 角色权限 ============
export function getRolePermissions(roleId: string | number) {
  return http.get<PermissionItem[]>(`/roles/${roleId}/permissions`)
}
export function assignPermissionsToRole(roleId: string | number, permissionIds: (string|number)[]) {
  return http.post(`/roles/${roleId}/permissions`, { permission_ids: permissionIds })
}
export function removePermissionsFromRole(roleId: string | number, permissionIds: (string|number)[]) {
  return http.delete(`/roles/${roleId}/permissions`, { data: { permission_ids: permissionIds } })
}
