import http from './http'

export interface MenuItem {
  id: string
  parent_id?: string | null
  name: string
  code?: string
  path?: string | null
  type?: number
  order?: number
  status?: string
  meta?: Record<string, any>
  children?: MenuItem[]
}

export function getUserMenusTree(userId: number | string) {
  return http.get<MenuItem[]>(`/users/${userId}/menus/tree`)
}

export function getRoleMenusTree(roleId: string) {
  return http.get<MenuItem[]>(`/roles/${roleId}/menus/tree`)
}

export function getRoleMenus(roleId: string) {
  return http.get<MenuItem[]>(`/roles/${roleId}/menus`)
}

// 角色菜单绑定
export function assignMenusToRole(roleId: string, menuIds: Array<string | number>) {
  return http.post(`/roles/${roleId}/menus`, { menu_ids: menuIds.map(id => Number(id)) })
}

export function removeMenusFromRole(roleId: string, menuIds: Array<string | number>) {
  return http.delete(`/roles/${roleId}/menus`, { data: { menu_ids: menuIds.map(id => Number(id)) } })
}

// ---- CRUD for menus ----
export function listMenus() {
  return http.get<MenuItem[]>(`/menus`)
}

export function getMenuTree() {
  return http.get<MenuItem[]>(`/menus/tree`)
}

export function createMenu(data: Partial<MenuItem>) {
  return http.post(`/menus`, data)
}

export function updateMenu(id: string, data: Partial<MenuItem>) {
  return http.put(`/menus/${id}`, data)
}

export function deleteMenu(id: string) {
  return http.delete(`/menus/${id}`)
}
