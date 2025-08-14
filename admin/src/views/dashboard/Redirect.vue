<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useMenuStore } from '@/stores/menu'
const auth = useAuthStore()
const menu = useMenuStore()
</script>

<template>
	<a-card>
		<template #title>登录成功</template>
		<div>欢迎，{{ auth.profile?.username }}。</div>
		<div style="margin-top:8px;">已获取菜单项数：{{ menu.tree.length }}</div>
	</a-card>
</template>
