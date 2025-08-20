import http from './http'

// 统一前缀（示例后端为 v1/storage/...）
const PFX = '/storage'

// Buckets
export function listBuckets() {
	return http.get<string[]>(`${PFX}/buckets`)
}

export function createBucket(bucket: string) {
	return http.post(`${PFX}/buckets/${encodeURIComponent(bucket)}`)
}

export function deleteBucket(bucket: string) {
	return http.delete(`${PFX}/buckets/${encodeURIComponent(bucket)}`)
}

// Objects
export interface ListObjectParams {
	prefix?: string
	recursive?: boolean
	limit?: number
}

export function listObjects(bucket: string, params: ListObjectParams = {}) {
	const query = new URLSearchParams()
	if (params.prefix) query.set('prefix', params.prefix)
	if (params.recursive) query.set('recursive', 'true')
	if (params.limit) query.set('limit', String(params.limit))
	const qs = query.toString()
	return http.get<string[]>(`${PFX}/buckets/${encodeURIComponent(bucket)}/objects${qs ? '?' + qs : ''}`)
}

export function uploadObject(bucket: string, file: File) {
	const form = new FormData()
	form.append('file', file)
	return http.post<{ url: string; object: string }>(`${PFX}/buckets/${encodeURIComponent(bucket)}/objects`, form, {
		headers: { 'Content-Type': 'multipart/form-data' },
	})
}

export function deleteObject(bucket: string, object: string) {
	return http.delete(`${PFX}/buckets/${encodeURIComponent(bucket)}/objects/${encodeURIComponent(object)}`)
}
