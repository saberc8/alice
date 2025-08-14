<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const title = computed(() => (route.meta?.title as string) || '页面')
</script>

<template>
	<a-card :title="title">
		<div>这是占位页面（TitleOnly），用于演示基于菜单的动态路由。</div>
	</a-card>
</template>
