import { defineStore } from 'pinia'
import { getUserMenusTree, type MenuItem } from '@/api/menu'

interface State {
	tree: MenuItem[]
	loaded: boolean
}

export const useMenuStore = defineStore('menu', {
	state: (): State => ({ tree: [], loaded: false }),
	actions: {
		async fetchMenusForUser(userId?: number) {
			if (!userId) return []
			const data = await getUserMenusTree(userId)
			this.tree = data || []
			this.loaded = true
			return this.tree
		},
		reset() {
			this.tree = []
			this.loaded = false
		},
	},
})
