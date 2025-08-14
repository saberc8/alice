import { defineStore } from 'pinia'
import { getProfile, login, type Profile } from '@/api/auth'

interface State {
	token: string | null
	profile: Profile | null
	initialized: boolean
}

export const useAuthStore = defineStore('auth', {
	state: (): State => ({
		token: typeof window !== 'undefined' ? localStorage.getItem('token') : null,
		profile: null,
		initialized: false,
	}),
	actions: {
		async doLogin(username: string, password: string) {
			const data = await login({ username, password })
			this.token = data.token
			localStorage.setItem('token', data.token)
		},
		async fetchProfile() {
			const p = await getProfile()
			this.profile = p
		this.initialized = true
			return p
		},
		async logout() {
			this.token = null
			this.profile = null
			this.initialized = false
			localStorage.removeItem('token')
		},
	},
})
