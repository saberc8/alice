import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useMenuStore } from '@/stores/menu'

const Login = () => import('@/views/login/Login.vue')
const Layout = () => import('@/layout/Layout.vue')
const Redirect = () => import('@/views/dashboard/Redirect.vue')

const routes: RouteRecordRaw[] = [
	{
		path: '/login',
		name: 'Login',
		component: Login,
		meta: { public: true, title: '登录' },
	},
	{
		path: '/',
		component: Layout,
		children: [
			{
				path: '',
				name: 'Home',
				component: Redirect,
				meta: { requiresAuth: true, title: '首页' },
			},
		],
	},
	{ path: '/:pathMatch(.*)*', redirect: '/' },
]

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

	// if logged in and not yet initialized, fetch profile and menus
	if (auth.token && !auth.initialized) {
		try {
			await auth.fetchProfile()
			await menu.fetchMenusForUser(auth.profile?.id)
		} catch (e) {
			// token invalid -> logout
			await auth.logout()
			return next('/login')
		}
	}

	next()
})

export default router
