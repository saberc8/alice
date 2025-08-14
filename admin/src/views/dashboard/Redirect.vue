<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMenuStore } from '@/stores/menu'

const router = useRouter()
const menu = useMenuStore()

function findFirstPath(tree: any[]): string | null {
	const stack = [...tree]
	while (stack.length) {
		const n = stack.shift()
		if (n?.type === 2 && n?.path) return n.path
		if (n?.children?.length) stack.unshift(...n.children)
	}
	return null
}

onMounted(() => {
	const p = findFirstPath(menu.tree)
	if (p) router.replace(p)
})
</script>

<template>
	<a-spin dot :loading="true" style="width:100%;margin-top:20vh;display:flex;justify-content:center" />
</template>
