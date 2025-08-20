<template>
	<div class="minio-page">
		<a-card title="Buckets" class="mb16">
			<div class="toolbar">
				<a-input
					v-model="newBucket"
					placeholder="新建 bucket 名称"
					style="width: 220px"
				/>
				<a-button
					type="primary"
					:disabled="!newBucket"
					@click="handleCreateBucket"
					>创建</a-button
				>
				<a-button @click="fetchBuckets" :loading="loadingBuckets"
					>刷新</a-button
				>
			</div>
            <a-table
				:data="buckets"
				:columns="bucketColumns"
				:loading="loadingBuckets"
				row-key="name"
				:pagination="false"
				@row-click="(row) => selectBucket(row.name)"
			>
				<template #bucketActions="{ record }">
					<a-popconfirm
						content="确认删除该 bucket? (其对象需手动清理)"
						@ok="() => handleDeleteBucket(record.name)"
					>
						<a-button status="danger" size="mini">删除</a-button>
					</a-popconfirm>
				</template>
				<template #empty>
					<div style="padding: 12px; color: var(--color-text-3)">
						暂无 Bucket 数据（调试：buckets.length = {{ buckets.length }}）
					</div>
				</template>
			</a-table>
			<div
				style="margin-top: 6px; display: flex; align-items: center; gap: 8px"
			>
				<a-switch v-model="__debug" size="small" />
				<span style="font-size: 12px; color: var(--color-text-3)"
					>调试开关</span
				>
			</div>
			<pre
				v-if="__debug"
				style="
					margin-top: 4px;
					max-height: 140px;
					overflow: auto;
					background: #f6f6f6;
					padding: 6px;
					font-size: 12px;
				"
				>{{ buckets }}</pre
			>
		</a-card>

		<a-card v-if="currentBucket" :title="'Objects in ' + currentBucket">
			<div class="toolbar">
				<input
					type="file"
					ref="fileInput"
					style="display: none"
					@change="onFileChange"
				/>
				<a-button type="primary" @click="() => fileInput?.click()"
					>上传文件</a-button
				>
				<a-input
					v-model="prefix"
					placeholder="过滤前缀"
					style="width: 200px"
					@change="fetchObjects"
				/>
				<a-switch v-model="recursive" @change="fetchObjects">递归</a-switch>
				<a-button @click="fetchObjects" :loading="loadingObjects"
					>刷新</a-button
				>
			</div>
			<a-table
				:data="objectRows"
				:columns="objectColumns"
				:loading="loadingObjects"
				row-key="o"
				:pagination="false"
			>
				<template #objectName="{ record }">
					<span class="obj-name" @click="copyName(record.o)" :title="record.o">{{
						record.o
					}}</span>
				</template>
				<template #objectActions="{ record }">
					<a-popconfirm
						content="确认删除该对象?"
						@ok="() => handleDeleteObject(record.o)"
					>
						<a-button size="mini" status="danger">删除</a-button>
					</a-popconfirm>
				</template>
			</a-table>
		</a-card>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
	listBuckets,
	createBucket,
	deleteBucket,
	listObjects,
	uploadObject,
	deleteObject,
} from '@/api/minio'

interface BucketRow {
	name: string
}

const buckets = ref<BucketRow[]>([])
const loadingBuckets = ref(false)
const newBucket = ref('')
const currentBucket = ref('')

const objects = ref<string[]>([])
const loadingObjects = ref(false)
const prefix = ref('')
const recursive = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
// 调试
const __debug = ref(false)

// table 列定义（采用 columns 方式，避免 a-table-column 渲染异常）
const bucketColumns = [
	{ title: '名称', dataIndex: 'name' },
	{ title: '操作', slotName: 'bucketActions', width: 120 },
]

// objects 需要转换为行对象
const objectRows = computed(() => objects.value.map((o) => ({ o })))
const objectColumns = [
	{ title: '对象名', dataIndex: 'o', slotName: 'objectName' },
	{ title: '操作', slotName: 'objectActions', width: 120 },
]

function fetchBuckets() {
	loadingBuckets.value = true
	listBuckets()
		.then((list) => {
			console.log(list)
			if (__debug.value) console.log('[MinioManager] listBuckets 返回:', list)
			buckets.value = list.map((n) => ({ name: n }))
		})
		.catch((e) => Message.error(e.message || '加载失败'))
		.finally(() => (loadingBuckets.value = false))
}

function handleCreateBucket() {
	if (!newBucket.value) return
	createBucket(newBucket.value)
		.then(() => {
			Message.success('创建成功')
			fetchBuckets()
			newBucket.value = ''
		})
		.catch((e) => Message.error(e.message || '创建失败'))
}

function handleDeleteBucket(name: string) {
	deleteBucket(name)
		.then(() => {
			Message.success('删除成功')
			if (currentBucket.value === name) {
				currentBucket.value = ''
				objects.value = []
			}
			fetchBuckets()
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
	listObjects(currentBucket.value, {
		prefix: prefix.value,
		recursive: recursive.value,
	})
		.then((list) => {
			objects.value = list
		})
		.catch((e) => Message.error(e.message || '加载对象失败'))
		.finally(() => (loadingObjects.value = false))
}

function onFileChange(e: Event) {
	const target = e.target as HTMLInputElement
	if (!target.files?.length || !currentBucket.value) return
	const file = target.files[0]
	uploadObject(currentBucket.value, file)
		.then(() => {
			Message.success('上传成功')
			fetchObjects()
		})
		.catch((err) => Message.error(err.message || '上传失败'))
		.finally(() => {
			if (fileInput.value) fileInput.value.value = ''
		})
}

function handleDeleteObject(name: string) {
	if (!currentBucket.value) return
	deleteObject(currentBucket.value, name)
		.then(() => {
			Message.success('删除成功')
			fetchObjects()
		})
		.catch((e) => Message.error(e.message || '删除失败'))
}

function copyName(name: string) {
	navigator.clipboard.writeText(name)
	Message.success('已复制对象名')
}

onMounted(() => {
	fetchBuckets()
})
</script>

<style scoped>
.minio-page {
	display: flex;
	flex-direction: column;
	gap: 16px;
	padding: 8px;
}
.toolbar {
	display: flex;
	gap: 8px;
	align-items: center;
	margin-bottom: 12px;
	flex-wrap: wrap;
}
.obj-name {
	cursor: pointer;
	color: var(--color-text-1);
}
.obj-name:hover {
	text-decoration: underline;
}
.mb16 {
	margin-bottom: 16px;
}
</style>
