import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useMenuStore } from '@/stores/menu'

// lazy components
const Login = () => import('@/views/login/Login.vue')
const Layout = () => import('@/layout/Layout.vue')
const Redirect = () => import('@/views/dashboard/Redirect.vue')
const TitleOnly = () => import('@/views/dynamic/TitleOnly.vue')

// map backend meta.component to SFC via Vite glob (using views convention)
const viewModules = import.meta.glob('/src/views/**/*.vue')

function resolveViewByMeta(metaComponent?: string) {
	if (!metaComponent) return null
	// normalize incoming strings like:
	// 'src/views/a/b.vue' | '/src/views/a/b.vue' | '/views/a/b' | 'views/a/b'
	let comp = metaComponent.trim()
	// strip leading slash
	if (comp.startsWith('/')) comp = comp.slice(1)
	// strip optional leading 'src/'
	if (comp.startsWith('src/')) comp = comp.slice(4)
	// accept 'pages/...' as alias of 'views/...'
	if (comp.startsWith('pages/')) comp = `views/${comp.slice(6)}`
	// now expect 'views/...'
	if (!comp.startsWith('views/')) return null
	// build candidate keys for import.meta.glob
	const base = `/src/${comp}`
	const candidates: string[] = []
	if (base.endsWith('.vue')) {
		candidates.push(base)
	} else {
		candidates.push(`${base}.vue`, `${base}/index.vue`)
	}
	for (const c of candidates) {
		const loader = (viewModules as any)[c]
		if (loader) return loader as () => Promise<any>
	}
	return null
}

// build routes from menu tree
function buildRoutesFromMenuTree(tree: any[]): RouteRecordRaw[] {
	const routes: RouteRecordRaw[] = []
	const walk = (nodes: any[]) => {
		nodes?.forEach((n) => {
			const type = n.type
			const path: string | null = n.path
			const name: string = n.code || n.id
			const title: string = n.name
			if (type === 2 && path) {
				const loader = resolveViewByMeta(n.meta?.component)
				routes.push({
					path,
					name,
					component: loader ? (loader as any) : TitleOnly,
					meta: { title },
				})
			}
			if (n.children?.length) walk(n.children)
		})
	}
	walk(tree)
	return routes
}

const routes: RouteRecordRaw[] = [
	{
		path: '/login',
		name: 'Login',
		component: Login,
		meta: { public: true, title: '登录' },
	},
	{
		path: '/',
	name: 'Root',
		component: Layout,
		children: [
			{
				path: '',
				name: 'Home',
				component: Redirect,
				meta: { requiresAuth: true, title: '首页' },
			},
		// 本地调试静态路由：对象存储浏览
		{
			path: 'minio',
			name: 'MinioLocal',
			component: () => import('@/views/minio/BucketBrowser.vue'),
			meta: { requiresAuth: true, title: '对象存储(卡片)' },
		},
		{
			path: 'minio-table',
			name: 'MinioTable',
			component: () => import('@/views/minio/MinioManager.vue'),
			meta: { requiresAuth: true, title: '对象存储(表格)' },
		},
		],
	},
	{ path: '/:pathMatch(.*)*', redirect: '/' },
]

// 记录首次进入时浏览器地址（用于动态路由注册后还原）
const initialFullPath = window.location.pathname + window.location.search + window.location.hash

const router = createRouter({
	history: createWebHistory(),
	routes,
})

router.beforeEach(async (to, _from, next) => {
	const auth = useAuthStore()
	const menu = useMenuStore()
	const isPublic = to.meta.public

	if (!isPublic && !auth.token) {
		return next({ path: '/login', query: { redirect: to.fullPath } })
	}

	// if logged in, ensure profile exists; then fetch menus by userId every refresh
	if (auth.token) {
		try {
			if (!auth.profile) await auth.fetchProfile()
			const userId = auth.profile?.id
			if (userId && !menu.loaded) {
				await menu.fetchMenusForUser(userId)
				// register dynamic routes (idempotent)
				if (!menu.routesRegistered) {
					const dynamicChildren = buildRoutesFromMenuTree(menu.tree)
					dynamicChildren.forEach((r) => router.addRoute('Root', r))
					menu.routesRegistered = true
				}
			}
			// 首次刷新：如果最初请求的路径在动态路由注册后已存在，且当前路由（可能已被重定向到默认页）不同，则还原
			if ((menu as any).initialNavigationDone !== true) {
				(menu as any).initialNavigationDone = true
				// 排除登录页与初始就是当前页的情况
				if (
					initialFullPath &&
					initialFullPath !== '/login' &&
					initialFullPath !== router.currentRoute.value.fullPath
				) {
					const resolved = router.resolve(initialFullPath)
					if (resolved.matched.length) {
						return next({ path: initialFullPath, replace: true })
					}
				}
			}
		} catch (e) {
			// token invalid -> logout
			await auth.logout()
			return next('/login')
		}
	}

	next()
})

export default router
