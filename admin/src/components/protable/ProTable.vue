<script lang="ts" setup>
import { ref, watch, computed } from 'vue'
import type { ProTableProps, ProColumn } from './types'

const props = defineProps<ProTableProps>()
const emits = defineEmits(['refresh'])

const loading = ref(false)
const data = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(props.pageSize || 10)
const searchModel = ref<Record<string, any>>({})

const rowKey = computed(() => props.rowKey || 'id')

async function fetch() {
  loading.value = true
  try {
    const params: any = { ...(searchModel.value || {}) }
    if (!props.isTree && props.pagination !== false) {
      params.page = page.value
      params.page_size = pageSize.value
    }
    const result = await props.request(params)
    if (Array.isArray(result)) {
      data.value = result
      if (!props.isTree && props.pagination !== false) total.value = result.length
    } else if (result && typeof result === 'object') {
      data.value = result.list
      total.value = result.total || result.list.length
    }
  } finally {
    loading.value = false
  }
}

watch(
  () => props.request,
  () => {
    fetch()
  },
  { immediate: true }
)

function get(obj: any, path: string) {
  return path.split('.').reduce((acc: any, k) => (acc ? acc[k] : undefined), obj)
}
function set(obj: any, path: string, value: any) {
  const keys = path.split('.')
  let cur = obj
  for (let i = 0; i < keys.length - 1; i++) {
    const k = keys[i]
    if (!cur[k] || typeof cur[k] !== 'object') cur[k] = {}
    cur = cur[k]
  }
  cur[keys[keys.length - 1]] = value
}

// 表单弹窗
const visible = ref(false)
const isEdit = ref(false)
const editingId = ref<string | number | null>(null)
const formModel = ref<Record<string, any>>({})

function openCreate() {
  isEdit.value = false
  editingId.value = null
  formModel.value = {}
  visible.value = true
}

function openEdit(record: any) {
  isEdit.value = true
  editingId.value = record[rowKey.value]
  formModel.value = JSON.parse(JSON.stringify(record))
  visible.value = true
}

async function onSubmit() {
  if (isEdit.value && props.onUpdate && editingId.value != null) {
    await props.onUpdate(editingId.value, formModel.value)
  } else if (!isEdit.value && props.onCreate) {
    await props.onCreate(formModel.value)
  }
  visible.value = false
  await fetch()
  emits('refresh')
}

async function onDelete(record: any) {
  if (!props.onDelete) return
  await props.onDelete(record[rowKey.value])
  await fetch()
}

// 生成搜索/表单字段
const safeColumns = computed(() => (Array.isArray(props.columns) ? props.columns.filter((c) => !!c) : []))
const searchableColumns = computed(() => safeColumns.value.filter((c) => (c as any).search))
const formColumns = computed(() => safeColumns.value.filter((c) => (c as any).form))

</script>

<template>
  <div class="pro-table">
    <!-- 搜索区域 -->
    <a-card :title="props.title" :bordered="false" class="mb-3">
      <a-form :model="searchModel" layout="inline" @submit.prevent>
        <template v-for="c in searchableColumns" :key="c.dataIndex">
          <a-form-item :label="c.title">
            <a-input v-model="searchModel[c.dataIndex]" :placeholder="c.searchPlaceholder || '请输入'" allow-clear />
          </a-form-item>
        </template>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="fetch">查询</a-button>
            <a-button @click="() => {searchModel.value = {}; page.value = 1; fetch()}" status="warning">重置</a-button>
          </a-space>
        </a-form-item>
        <a-form-item v-if="props.onCreate">
          <a-button type="primary" status="success" @click="openCreate">新建</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 表格区域 -->
    <a-table
      :data="data"
      :loading="loading"
      :row-key="rowKey"
      :pagination="false"
      :columns="[]"
      :bordered="false"
      :default-expand-all-rows="props.isTree"
      :row-class="() => 'pro-table-row'"
    >
      <template #columns>
        <template v-for="(c, idx) in safeColumns" :key="(c && c.dataIndex) ? c.dataIndex : idx">
          <a-table-column v-if="c && !c.hideInTable" :title="c.title" :data-index="c.dataIndex">
            <template #cell="{ record }">
              <slot :name="`col-${c.dataIndex}`" :record="record">
                {{ get(record, c.dataIndex) }}
              </slot>
            </template>
          </a-table-column>
        </template>
        <a-table-column title="操作" :width="200" v-if="props.onUpdate || props.onDelete">
          <template #cell="{ record }">
            <a-space>
              <a-button v-if="props.onUpdate" size="mini" @click="openEdit(record)">编辑</a-button>
              <a-popconfirm v-if="props.onDelete" content="确定删除该项？" type="warning" @ok="() => onDelete(record)">
                <a-button size="mini" status="danger">删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 分页 -->
    <div v-if="!props.isTree && props.pagination !== false" style="margin-top:12px; text-align:right;">
      <a-pagination
        :total="total"
        :current="page"
        :page-size="pageSize"
        show-total
        show-jumper
        show-page-size
        @change="(p:number)=>{ page.value = p; fetch() }"
        @page-size-change="(s:number)=>{ pageSize.value = s; page.value = 1; fetch() }"
      />
    </div>

    <!-- 创建/编辑弹窗 -->
    <a-modal v-model:visible="visible" :title="isEdit ? '编辑' : '新建'" @ok="onSubmit">
      <a-form :model="formModel" layout="vertical">
        <template v-for="c in formColumns" :key="c.dataIndex">
          <a-form-item :label="c.title" :required="typeof c.form === 'object' ? c.form.required : false">
            <!-- treeSelect 单独处理，以便传 data 与 field-names -->
            <a-tree-select
              v-if="typeof c.form === 'object' && c.form.widget === 'treeSelect'"
              :data="(c.form as any).treeOptions || []"
              :model-value="get(formModel, c.dataIndex)"
              @update:model-value="(val:any) => set(formModel, c.dataIndex, val)"
              :placeholder="(c.form as any).placeholder || ''"
              allow-search
              allow-clear
              :field-names="{ key: 'value', title: 'label', children: 'children' }"
              v-bind="(c.form as any).props || {}"
            />
            <!-- 其他控件使用通用 component 渲染 -->
            <component
              v-else
              :is="(typeof c.form === 'object' && c.form.widget === 'textarea') ? 'a-textarea' :
                    (typeof c.form === 'object' && c.form.widget === 'number') ? 'a-input-number' :
                    (typeof c.form === 'object' && c.form.widget === 'select') ? 'a-select' :
                    (typeof c.form === 'object' && c.form.widget === 'switch') ? 'a-switch' : 'a-input'"
              :model-value="get(formModel, c.dataIndex)"
              @update:model-value="(val:any) => set(formModel, c.dataIndex, val)"
              v-bind="typeof c.form === 'object' ? c.form.props : {}"
              :placeholder="typeof c.form === 'object' ? c.form.placeholder : ''"
              allow-clear
            >
              <template v-if="typeof c.form === 'object' && c.form.widget === 'select'">
                <a-option v-for="opt in c.form?.options || []" :key="String(opt.value)" :value="opt.value" :disabled="opt.disabled">{{ opt.label }}</a-option>
              </template>
            </component>
          </a-form-item>
        </template>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped>
.mb-3 { margin-bottom: 12px; }
</style>
