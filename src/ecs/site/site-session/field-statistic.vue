<script lang="ts">
  export default {
    name: 'FieldStatistic',
  }
</script>
<script setup lang="ts">
  import { siteSearchType } from '@/types'

  import { siteAggregationsApi } from '~/src/api-ecs/site'

  import { useTableCopy } from '@/utils'

  import { formatTime } from '@/utils/time'
  import { ShortcutListType } from '@/types'
  import { getTableCopyData } from '~/src/utils/transition'
  const props = defineProps<{
    showFieldStatistic: boolean
    fieldList: string[]
    aggregationsPayload: siteSearchType
    title: string
    showContextMenu?: boolean
    tableColumn?: any
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'changeShortcut', data: ShortcutListType): void
  }>()

  const visible = ref(false) // 显隐

  const disabled = ref(true) // 显隐

  const topCountNmu = ref(10) // 前多少条

  const fieldLists = ref<string[]>([]) // 请求字段

  const aggregationsPayload1 = ref<siteSearchType>() // 回显数据

  const listLoading = ref(false) // 是否加载

  const option = reactive([
    { label: 'TOP5', value: 5 },
    { label: 'TOP10', value: 10 },
    { label: 'TOP20', value: 20 },
    { label: 'TOP50', value: 50 },
    { label: 'TOP100', value: 100 },
    { label: 'TOP1000', value: 1000 },
  ])

  // 表格数据
  const listDate = ref<object[]>([])

  // 表格数据
  // const layout = ref('total, sizes, prev, pager, next, jumper')

  onMounted(() => {
    initData()
  })

  const initData = async () => {
    visible.value = props.showFieldStatistic
    fieldLists.value = props.fieldList
    aggregationsPayload1.value = props.aggregationsPayload
    siteAggregations({ list: fieldLists.value, num: topCountNmu.value })
  }

  // 聚合检索
  const siteAggregations = async ({ list, num }: { list: string[]; num: number }) => {
    listLoading.value = true
    const {
      data: { aggObj },
    } =
      // @ts-ignore
      await siteAggregationsApi({
        ...aggregationsPayload1.value,
        aggregationFields: list,
        topCount: num,
      })
    const field = fieldLists.value[0]
    // @ts-ignore
    listDate.value = aggObj || []
    listLoading.value = false
  }

  // 检索
  const queryData = (e: any) => {
    topCountNmu.value = e
    siteAggregations({ list: fieldLists.value, num: topCountNmu.value })
  }

  const handleClose = () => {
    emit('on-closeEvent', false)
  }

  const tableColumn = [
    { fieldNameCn: '值', fieldNameEn: 'name' },
    { fieldNameCn: '次数', fieldNameEn: 'value' },
  ]

  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }

  function formatExcelData(val: any, key: string) {
    switch (key) {
      default:
        return val[key] || ''
    }
  }

  const handleDownloadExcel = () => {
    const tHeader: string[] = []
    const filterVal: any = []
    try {
      tableColumn.forEach(({ fieldNameCn, fieldNameEn }) => {
        tHeader.push(fieldNameCn)
        filterVal.push(fieldNameEn)
      })
      listLoading.value = true
      import('@/utils/excel').then((excel) => {
        const list = listDate.value
        const data = formatJson(filterVal, list)
        excel.export_json_to_excel({
          header: tHeader,
          data: data,
          filename: `${props.title}-${formatTime(new Date().getTime())}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
        listLoading.value = false
      })
    } catch (error) {
      console.log(error)
      listLoading.value = false
    }
  }

  const handleUpdateCallback = (res: ShortcutListType) => {
    emit('changeShortcut', res)
  }
  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    if (props.showContextMenu && column.label === '值') {
      const _column = { ...column, property: props.tableColumn?.fieldNameEn }
      const _row = { ...row, [_column.property]: row.name }
      const arr = getTableCopyData({
        tableColumn: [props.tableColumn],
        row: _row,
        column: _column,
        mothod: handleUpdateCallback,
      })
      useTableCopy(_row, _column, cell, event, arr)
    }
  }
</script>

<template>
  <div class="field-statistic">
    <el-dialog v-model="visible" :before-close="handleClose" :title="props.title" width="1170px">
      <vab-query-form style="width: 100%">
        <vab-query-form-left-panel :span="12">
          <el-button :disabled="listDate.length == 0" type="primary" @click="handleDownloadExcel">导出</el-button>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-form v-if="title != ('请求负载' || '响应负载')" inline @submit.prevent>
            <el-form-item>
              <el-select v-model="topCountNmu" class="m-2" :disabled="listLoading" @change="queryData">
                <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-form>
        </vab-query-form-right-panel>
      </vab-query-form>
      <!-- 表格 -->
      <el-table
        v-loading="listLoading"
        :border="true"
        :data="listDate"
        row-key="id"
        @cell-contextmenu="useTableCopyEvent"
      >
        <el-table-column :align="'center'" label="序号" width="100">
          <template #default="{ $index }">
            {{ $index + 1 }}
          </template>
        </el-table-column>
        <el-table-column :align="'center'" label="值" prop="name" show-overflow-tooltip />
        <el-table-column :align="'center'" label="次数" prop="value" show-overflow-tooltip width="200" />
        <template #empty>
          <el-empty class="vab-data-empty" description="暂无数据" />
        </template>
      </el-table>
      <!-- pagination -->
      <!-- <el-pagination background :current-page="aggregationsPayload1.pageNum" :layout="layout"
                                  :page-size="aggregationsPayload1.pageSize" :page-sizes="[10,20,50,100]" style="margin-bottom: 20px"
                                  :total="listDate.length" @current-change="handleCurrentChange" @size-change="handleSizeChange" /> -->
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .field-statistic {
    :deep() {
      .el-scrollbar {
        height: 480px;
      }
      .el-table__body-wrapper {
        border-bottom: var(--el-table-border);
      }
    }
  }
</style>
