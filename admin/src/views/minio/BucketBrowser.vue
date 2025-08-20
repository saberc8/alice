<template>
  <div class="bucket-browser">
    <a-page-header title="对象存储 (本地调试)" :breadcrumb="[{ path: '/', label: '首页' }, { label: '对象存储' }]" />

    <div class="buckets-section">
      <div class="section-header">
        <h3>Buckets</h3>
        <div class="actions">
          <a-input v-model="newBucket" placeholder="新建 bucket" style="width:200px" />
          <a-button type="primary" :disabled="!newBucket" @click="createNewBucket">创建</a-button>
          <a-button @click="loadBuckets" :loading="loadingBuckets">刷新</a-button>
        </div>
      </div>
      <div class="bucket-grid" v-if="!loadingBuckets">
        <div
          v-for="b in buckets"
          :key="b"
          class="bucket-item"
          :class="{ active: b === currentBucket }"
          @click="selectBucket(b)"
        >
          <div class="icon">📁</div>
          <div class="name" :title="b">{{ b }}</div>
          <div class="ops" @click.stop>
            <a-popconfirm content="确认删除该 bucket? (对象需已清空)" @ok="() => onDeleteBucket(b)">
              <a-button size="mini" status="danger" type="text">删</a-button>
            </a-popconfirm>
          </div>
        </div>
        <div v-if="!buckets.length" class="empty-tip">暂无 bucket</div>
      </div>
      <div v-else class="loading-tip">加载中...</div>
    </div>

    <div class="objects-section" v-if="currentBucket">
      <div class="section-header">
        <h3>对象 - {{ currentBucket }}</h3>
        <div class="actions">
          <input type="file" ref="fileInput" style="display:none" @change="onFileChange" />
          <a-button type="primary" @click="() => fileInput?.click()">上传文件</a-button>
          <a-input v-model="prefix" placeholder="前缀过滤" style="width:160px" @input="debouncedFetchObjects" />
          <a-switch v-model="recursive" size="small" @change="fetchObjects">递归</a-switch>
          <a-button @click="fetchObjects" :loading="loadingObjects">刷新</a-button>
        </div>
      </div>
  <a-table :data="objectRows" row-key="o" :loading="loadingObjects" :pagination="false" size="small">
        <a-table-column title="对象名" data-index="o">
          <template #cell="{ record }">
    <span class="obj" @click="copy(record.o)" :title="record.o">{{ record.o }}</span>
          </template>
        </a-table-column>
  <a-table-column title="操作">
          <template #cell="{ record }">
    <a-popconfirm content="确认删除该对象?" @ok="() => onDeleteObject(record.o)">
              <a-button size="mini" status="danger">删除</a-button>
            </a-popconfirm>
          </template>
        </a-table-column>
      <template #empty>
        <div style="padding:12px; color:var(--color-text-3)">
          {{ loadingObjects ? '加载中...' : '无对象 (或前缀过滤结果为空)' }}
        </div>
      </template>
      </a-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import { listBuckets, createBucket, deleteBucket, listObjects, uploadObject, deleteObject } from '@/api/minio'

const buckets = ref<string[]>([])
const loadingBuckets = ref(false)
const newBucket = ref('')
const currentBucket = ref('')

// 原始对象名数组
const objects = ref<string[]>([])
// 表格行（Arco Table 需要对象而不是原始字符串）
const objectRows = computed(() => objects.value.map((o) => ({ o })))
const loadingObjects = ref(false)
const prefix = ref('')
const recursive = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
let debounceTimer: number | null = null

function loadBuckets() {
  loadingBuckets.value = true
  listBuckets()
    .then((data) => {
      console.log('[Minio] buckets resp:', data)
      buckets.value = data || []
    })
    .catch((e) => Message.error(e.message || '加载失败'))
    .finally(() => (loadingBuckets.value = false))
}

function createNewBucket() {
  if (!newBucket.value) return
  createBucket(newBucket.value)
    .then(() => {
      Message.success('创建成功')
      newBucket.value = ''
      loadBuckets()
    })
    .catch((e) => Message.error(e.message || '创建失败'))
}

function onDeleteBucket(name: string) {
  deleteBucket(name)
    .then(() => {
      Message.success('删除成功')
      if (currentBucket.value === name) {
        currentBucket.value = ''
        objects.value = []
      }
      loadBuckets()
    })
    .catch((e) => Message.error(e.message || '删除失败'))
}

function selectBucket(name: string) {
  currentBucket.value = name
  fetchObjects()
}

function fetchObjects() {
  if (!currentBucket.value) return
  loadingObjects.value = true
  listObjects(currentBucket.value, { prefix: prefix.value, recursive: recursive.value })
    .then((list) => {
      console.log('[Minio] objects resp:', list)
      objects.value = list || []
    })
    .catch((e) => Message.error(e.message || '加载对象失败'))
    .finally(() => (loadingObjects.value = false))
}

function debouncedFetchObjects() {
  if (debounceTimer) window.clearTimeout(debounceTimer)
  debounceTimer = window.setTimeout(fetchObjects, 300)
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length || !currentBucket.value) return
  const file = input.files[0]
  uploadObject(currentBucket.value, file)
    .then(() => {
      Message.success('上传成功')
      fetchObjects()
    })
    .catch((e) => Message.error(e.message || '上传失败'))
    .finally(() => {
      if (fileInput.value) fileInput.value.value = ''
    })
}

function onDeleteObject(name: string) {
  if (!currentBucket.value) return
  deleteObject(currentBucket.value, name)
    .then(() => {
      Message.success('删除成功')
      fetchObjects()
    })
    .catch((e) => Message.error(e.message || '删除失败'))
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  Message.success('已复制')
}

onMounted(() => {
  loadBuckets()
})
</script>

<style scoped>
.bucket-browser { display: flex; flex-direction: column; gap: 20px; }
.buckets-section, .objects-section { background: var(--color-bg-2); padding: 16px; border-radius: 6px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; flex-wrap: wrap; gap: 8px; }
.bucket-grid { display: flex; flex-wrap: wrap; gap: 12px; }
.bucket-item { width: 140px; border: 1px solid var(--color-border-2); padding: 10px 8px 8px; border-radius: 6px; cursor: pointer; position: relative; background: var(--color-bg-1); display: flex; flex-direction: column; align-items: center; gap: 4px; }
.bucket-item.active { outline: 2px solid var(--color-primary-6); }
.bucket-item:hover { background: var(--color-fill-2); }
.bucket-item .icon { font-size: 32px; line-height: 1; }
.bucket-item .name { font-size: 14px; width: 100%; text-align: center; word-break: break-all; }
.bucket-item .ops { position: absolute; top: 4px; right: 4px; }
.empty-tip, .loading-tip { color: var(--color-text-3); padding: 8px; }
.obj { cursor: pointer; }
.obj:hover { text-decoration: underline; }
</style>
