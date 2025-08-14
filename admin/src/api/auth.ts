import http from './http'

export interface LoginReq {
  username: string
  password: string
}
export interface Profile {
  id: number
  username: string
  email?: string
  roles?: Array<{ id: string; name: string; code: string }>
}

export function login(data: LoginReq) {
  return http.post<{ token: string }>('/auth/login', data)
}

export function getProfile() {
  return http.get<Profile>('/auth/profile')
}
