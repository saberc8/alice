<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useMenuStore } from '@/stores/menu'
import MenuItemNode from '@/components/MenuItemNode.vue'

const auth = useAuthStore()
const menu = useMenuStore()
const router = useRouter()
const route = useRoute()

function onLogout() {
	auth.logout()
	menu.reset()
	window.location.href = '/login'
}

const selectedKeys = computed(() => [route.path])

function onMenuSelect(key: string) {
	if (key && key !== route.path) router.push(key)
}
</script>

<template>
	<a-layout style="min-height: 100vh">
		<a-layout-sider breakpoint="xl" collapsible>
			<div class="logo">Alice Admin</div>
					<a-menu :selected-keys="selectedKeys" @menu-item-click="onMenuSelect">
						<MenuItemNode v-for="n in menu.tree" :key="n.id" :node="n" />
					</a-menu>
		</a-layout-sider>
		<a-layout>
			<a-layout-header class="header">
				<div />
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
	</a-layout>
</template>

<style scoped lang="less">
.logo {
	height: 48px;
	display: flex;
	align-items: center;
	padding: 0 16px;
	font-weight: 600;
}
.header {
	display: flex;
	align-items: center;
	justify-content: space-between;
}
</style>
