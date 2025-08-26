import http from './http'

// ===== 类型定义 =====
export interface UserItem {
  id: number
  username: string
  email: string
  status: string
  created_at?: string
  updated_at?: string
}

export interface UserListResult {
  items: UserItem[]
  total: number
  page: number
  page_size: number
}

// ===== 列表（分页） =====
export interface ListUserParams {
  page?: number
  page_size?: number
}

export function listUsers(params: ListUserParams) {
  return http.get<UserListResult>('/users', { params })
}

// ===== 创建 =====
export interface CreateUserPayload {
  username: string
  password: string
  email: string
  status?: string
}
export function createUser(data: CreateUserPayload) {
  return http.post<UserItem>('/users', data)
}

// ===== 详情 =====
export function getUser(id: number | string) {
  return http.get<UserItem>(`/users/${id}`)
}

// ===== 更新 =====
export interface UpdateUserPayload {
  email?: string
  status?: string
  password?: string
}
export function updateUser(id: number | string, data: UpdateUserPayload) {
  return http.put<UserItem>(`/users/${id}` , data)
}

// ===== 删除 =====
export function deleteUser(id: number | string) {
  return http.delete(`/users/${id}`)
}
