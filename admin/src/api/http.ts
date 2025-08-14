import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE || '/api'
const timeout = Number(import.meta.env.VITE_REQUEST_TIMEOUT || 15000)

export const http = axios.create({ baseURL, timeout })

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (res) => {
    // normalize common response shape { code, message, data }
    const payload = res.data
    if (payload && typeof payload === 'object' && 'code' in payload) {
      if (payload.code === 200) return payload.data
      // throw error with message
      const err = new Error(payload.message || 'Request error') as any
      err.code = payload.code
      throw err
    }
    return res.data
  },
  (err) => {
    return Promise.reject(err)
  }
)

export default http
