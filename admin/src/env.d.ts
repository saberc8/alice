/// <reference types="vite/client" />

interface ImportMetaEnv {
	readonly VITE_PORT?: string
	readonly VITE_API_BASE?: string
	readonly VITE_REQUEST_TIMEOUT?: string
	readonly VITE_API_PROXY_TARGET?: string
}
interface ImportMeta {
	readonly env: ImportMetaEnv
}
