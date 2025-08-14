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
