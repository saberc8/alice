<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
const auth = useAuthStore()

function onLogout() {
	auth.logout()
	window.location.href = '/login'
}
</script>

<template>
	<a-layout style="min-height: 100vh">
		<a-layout-header style="display:flex;align-items:center;justify-content:space-between">
			<div style="font-weight:600">Alice Admin</div>
			<div>
				<a-dropdown trigger="click">
					<a-button type="text">{{ auth.profile?.username || '用户' }}</a-button>
					<template #content>
						<a-doption @click="onLogout">退出登录</a-doption>
					</template>
				</a-dropdown>
			</div>
		</a-layout-header>
		<a-layout-content style="padding: 16px">
			<router-view />
		</a-layout-content>
	</a-layout>
</template>

<style scoped lang="less">
</style>
