<script setup lang="ts">
defineOptions({ name: 'MenuItemNode' })
interface MenuNode {
  id: string
  code?: string
  name: string
  path?: string | null
  type?: number
  children?: MenuNode[]
}
const props = defineProps<{ node: MenuNode }>()
</script>

<template>
  <a-sub-menu v-if="props.node.children && props.node.children.length" :key="props.node.code || props.node.id">
    <template #title>{{ props.node.name }}</template>
    <MenuItemNode v-for="c in props.node.children" :key="c.id" :node="c" />
  </a-sub-menu>
  <a-menu-item v-else-if="props.node.type === 2 && props.node.path" :key="props.node.path">
    {{ props.node.name }}
  </a-menu-item>
</template>
