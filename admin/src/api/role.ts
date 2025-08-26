import http from './http'

// ===== 类型定义 =====
export interface RoleItem {
  id: string
  name: string
  code: string
  description?: string | null
  status: string
  created_at?: string
  updated_at?: string
}

export interface RoleListResult {
  items: RoleItem[]
  total: number
  page: number
  page_size: number
}

// ===== 列表（分页） =====
export interface ListRoleParams {
  page?: number
  page_size?: number
}

export function listRoles(params: ListRoleParams) {
  return http.get<RoleListResult>('/roles', { params })
}

// ===== 创建 =====
export interface CreateRolePayload {
  name: string
  code: string
  description?: string | null
  status?: string
}
export function createRole(data: CreateRolePayload) {
  return http.post<RoleItem>('/roles', data)
}

// ===== 详情 =====
export function getRole(id: string) {
  return http.get<RoleItem>(`/roles/${id}`)
}

// ===== 更新 =====
export interface UpdateRolePayload {
  name: string
  code: string
  description?: string | null
  status?: string
}
export function updateRole(id: string, data: UpdateRolePayload) {
  return http.put(`/roles/${id}`, data)
}

// ===== 删除 =====
export function deleteRole(id: string) {
  return http.delete(`/roles/${id}`)
}

// ===== 用户角色关联 =====
export function getUserRoles(userId: string) {
  return http.get<RoleItem[]>(`/users/${userId}/roles`)
}

export function assignRolesToUser(userId: string, roleIds: string[]) {
  return http.post(`/users/${userId}/roles`, { role_ids: roleIds })
}

export function removeRolesFromUser(userId: string, roleIds: string[]) {
  return http.delete(`/users/${userId}/roles`, { data: { role_ids: roleIds } })
}
