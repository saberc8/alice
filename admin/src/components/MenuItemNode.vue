<script setup lang="ts">
import { Icon } from '@iconify/vue'
defineOptions({ name: 'MenuItemNode' })
interface MenuNode {
  id: string
  code?: string
  name: string
  path?: string | null
  type?: number
  meta?: Record<string, any>
  children?: MenuNode[]
}
const props = defineProps<{ node: MenuNode }>()
const iconName = props.node.meta?.icon as string | undefined
</script>

<template>
  <a-sub-menu v-if="props.node.children && props.node.children.length" :key="props.node.code || props.node.id">
    <template #title>
      <span class="menu-title">
        <Icon v-if="iconName" :icon="iconName" class="menu-icon" />
        <span>{{ props.node.name }}</span>
      </span>
    </template>
    <MenuItemNode v-for="c in props.node.children" :key="c.id" :node="c" />
  </a-sub-menu>
  <a-menu-item v-else-if="props.node.type === 2 && props.node.path" :key="props.node.path">
    <span class="menu-title">
      <Icon v-if="iconName" :icon="iconName" class="menu-icon" />
      <span>{{ props.node.name }}</span>
    </span>
  </a-menu-item>
</template>

<style scoped>
.menu-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.menu-icon {
  font-size: 16px;
  width: 16px;
  height: 16px;
}
</style>
