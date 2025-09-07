<script lang="ts">
  export default {
    name: 'FieldDeleteHistory',
  }
</script>

<script setup lang="ts">
  import { Back } from '@element-plus/icons-vue'
  import { ElDivider } from 'element-plus'
  import { formatNstime } from '@/utils/time'
  import { InitialisationItem, tableSearch } from '@/types/index'
  import {
    getDeleteFieldsTemplateApi,
    fieldsTemplateRecoverApi,
    fieldsTemplateHistoryDeleteApi,
  } from '~/src/api-ecs/custom-field'
  import TableFilterTool from '@/components/table-filter-tool.vue'

  const handleTableSearch = inject(tableSearch)
  const spacer = h(ElDivider, { direction: 'vertical' })
  const $baseMessage: any = inject('$baseMessage')
  const $baseConfirm: any = inject('$baseConfirm')

  // 获取表格序号
  const curIndex = computed(() => (queryPage.pageNum - 1) * queryPage.pageSize + 1)
  const queryPage = reactive({
    templateName: '',
    pageNum: 1,
    pageSize: 10,
  })
  const total = ref(0)
  const emits = defineEmits<{
    (e: 'back', show: boolean): void
    (e: 'show-info', showVal: any, display: 'add' | 'preview'): void
  }>()

  const listDate = ref<InitialisationItem[]>()
  const listLoading = ref(false) // 是否加载
  provide(tableSearch, reloadTemplate)
  function reloadTemplate() {
    queryPage.pageNum = 1
    getAllTemplate()
  }
  const getAllTemplate = async () => {
    listLoading.value = true
    try {
      const { data } = await getDeleteFieldsTemplateApi(queryPage)
      listDate.value = data.records || []
      total.value = data.total
    } finally {
      listLoading.value = false
    }
  }
  const handleEdit = (row: InitialisationItem) => {
    emits('show-info', row, 'preview')
  }

  const handleRecover = async (row: InitialisationItem) => {
    const { msg } = await fieldsTemplateRecoverApi({ ids: [row.id] })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    reloadTemplate()
    handleTableSearch?.()
  }
  const handleDelete = async (row: InitialisationItem) => {
    $baseConfirm('你确定要彻底删除当前项吗？', null, async () => {
      const { msg } = await fieldsTemplateHistoryDeleteApi({ ids: [row.id] })
      $baseMessage(msg, 'success', 'vab-hey-message-success')
      reloadTemplate()
    })
  }
  onMounted(() => {
    getAllTemplate()
  })
</script>

<template>
  <div>
    <el-space :spacer="spacer">
      <el-button :auto-insert-space="false" :icon="Back" plin @click="emits('back', false)">返回</el-button>
      <h3>删除历史</h3>
    </el-space>
    <el-table v-loading="listLoading" class="my-table" :data="listDate" row-key="id">
      <el-table-column align="left" :index="curIndex" label="序号" type="index" width="75" />
      <el-table-column label="名称" prop="templateName" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryPage.templateName"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" width="180">
        <template #default="{ row }">
          {{ formatNstime(row.createTime, false) }}
        </template>
      </el-table-column>
      <el-table-column label="修改时间" prop="updateTime" width="180">
        <template #default="{ row }">
          {{ formatNstime(row.updateTime, false) }}
        </template>
      </el-table-column>
      <el-table-column label="备注" prop="remark" show-overflow-tooltip />
      <el-table-column label="操作" width="220">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">详情</el-button>
          <el-button size="small" @click="handleRecover(row)">恢复</el-button>
          <el-button size="small" @click="handleDelete(row)">彻底删除</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <el-pagination
      v-model:current-page="queryPage.pageNum"
      v-model:page-size="queryPage.pageSize"
      background
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      @current-change="getAllTemplate"
      @size-change="getAllTemplate"
    />
  </div>
</template>

<style scoped lang="scss">
  h3 {
    margin-bottom: 0;
  }
  .my-table {
    margin-top: 16px;
    height: calc(100vh - 140px);
  }
</style>
