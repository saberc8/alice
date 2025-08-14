import { defineStore } from 'pinia'
import { getUserMenusTree, type MenuItem } from '@/api/menu'

interface State {
	tree: MenuItem[]
	loaded: boolean
	routesRegistered: boolean
}

export const useMenuStore = defineStore('menu', {
	state: (): State => ({ tree: [], loaded: false, routesRegistered: false }),
	actions: {
			async fetchMenusForUser(userId?: number | string) {
				if (!userId) return []
				const data = await getUserMenusTree(userId)
				this.tree = (data as any) || []
				this.loaded = true
				return this.tree
			},
		reset() {
			this.tree = []
			this.loaded = false
			this.routesRegistered = false
		},
	},
})


