<script lang="ts">
  export default {
    name: 'TraceabilityFieldStatistic',
  }
</script>
<script setup lang="ts">
  import { EventExportAggAPI, getTopFieldApi } from '@/api-ecs/retrieve'

  import { PROTOCOLDICT } from '@/data/constant'
  import { ShortcutListType } from '@/types'
  import { formatTime } from '@/utils/time'
  import { TableColumnItemType } from '/#/store'
  import { useTableCopy } from '@/utils'
  import { getTableCopyData } from '~/src/utils/transition'
  import { downloadExcel } from '~/src/utils/download'

  const props = defineProps<{
    dnsDict?: any
    showContextMenu?: boolean
  }>()

  const emits = defineEmits<{
    (e: 'changeShortcut', data: ShortcutListType): void
  }>()

  const option = [
    { label: 'TOP5', value: 5 },
    { label: 'TOP10', value: 10 },
    { label: 'TOP20', value: 20 },
    { label: 'TOP50', value: 50 },
    { label: 'TOP100', value: 100 },
    { label: 'TOP1000', value: 1000 },
  ]

  const visible = ref(false) // 显隐

  const queryData = reactive({
    title: '',
    listDate: [] as { name: string; value: number }[],
    topCount: 10,
  })
  const statisticType = ref()
  const myQuery = ref()

  const listLoading = ref(false)

  const disabled = ref(false)
  let my_column: TableColumnItemType | undefined = undefined
  const searchHandle = async (query: {
    indexType: number
    startTime: string
    endTime: string
    aggregationFields: string
    searchSql: string
    whiteType?: number
  }) => {
    listLoading.value = true
    const { topCount } = queryData
    const {
      data: { aggObj },
    } = await getTopFieldApi({
      ...query,
      topCount: topCount,
    })
    queryData.listDate = aggObj || []
    if (query.indexType == 2) {
      if (query.aggregationFields == 'requestType' || query.aggregationFields == 'requestClass') {
        changeDNSData(query.aggregationFields)
      }
    }
    //  else if (query.aggregationFields == 'transProtocol') {
    //   formatTransProtocol()
    // } else if (query.aggregationFields == 'appProtocol') {
    //   formatAppProtocol()
    // }
    listLoading.value = false
  }
  const initData = async (option: {
    type?: 'alert'
    title: string
    query: {
      indexType: number
      startTime: string
      endTime: string
      aggregationFields: string
      searchSql: string
      whiteType?: number
    }
    column?: TableColumnItemType
  }) => {
    const { title, query, column, type } = option
    visible.value = true
    queryData.title = title
    myQuery.value = query
    statisticType.value = type
    my_column = column
    searchHandle(query)
  }
  const topCountChange = (topCount: number) => {
    searchHandle(myQuery.value)
  }

  const changeDNSData = (val: string) => {
    queryData.listDate.forEach((item: any) => {
      const res = item.name
      item['name'] = props.dnsDict[val][res] || item.name
    })
  }

  const formatAppProtocol = () => {
    queryData.listDate.forEach((item: any) => {
      const res = item.name
      // @ts-ignore
      item['name'] = PROTOCOLDICT[res]
    })
  }

  const formatTransProtocol = () => {
    queryData.listDate.forEach((item: any) => {
      const res = item.name
      item['name'] = formatTcp(res)
    })
  }

  const formatTcp = (val: string) => {
    let value = undefined
    if (val == '1') {
      value = 'TCP'
    } else if (val == '2') {
      value = 'UDP'
    }
    return value
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
  const handleUpdateCallback = (res: ShortcutListType) => {
    emits('changeShortcut', res)
  }
  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    if (props.showContextMenu && column.label === '值') {
      const _column = { ...column, property: my_column?.fieldNameEn }
      const _row = { ...row, [_column.property]: row.name }
      const arr = getTableCopyData({
        tableColumn: [my_column],
        row: _row,
        column: _column,
        mothod: handleUpdateCallback,
      })
      useTableCopy(_row, _column, cell, event, arr)
    }
  }
  const handleDownloadExcel = async () => {
    try {
      disabled.value = true
      const res = await EventExportAggAPI(queryData.listDate)
      const filename = `${queryData.title}-${formatTime(new Date().getTime())}`
      downloadExcel(res, filename)
    } finally {
      disabled.value = false
    }
    // const blob = new Blob([res]) // 把得到的结果用流对象转一下
    // var a = document.createElement('a') //创建一个<a></a>标签
    // a.href = URL.createObjectURL(blob) // 将流文件写入a标签的href属性值
    // a.download = '基础工艺数据.xlsx' //设置文件名
    // a.style.display = 'none' // 障眼法藏起来a标签
    // document.body.appendChild(a) // 将a标签追加到文档对象中
    // a.click() // 模拟点击了a标签，会触发a标签的href的读取，浏览器就会自动下载了
    // a.remove()
    // const tHeader: string[] = []
    // const filterVal: any = []
    // try {
    //   tableColumn.forEach(({ fieldNameCn, fieldNameEn }) => {
    //     tHeader.push(fieldNameCn)
    //     filterVal.push(fieldNameEn)
    //   })
    //   listLoading.value = true
    //   import('@/utils/excel').then((excel) => {
    //     const list = queryData.listDate
    //     const data = formatJson(filterVal, list)
    //     excel.export_json_to_excel({
    //       header: tHeader,
    //       data: data,
    //       ,
    //       autoWidth: true,
    //       bookType: 'xlsx',
    //     })
    //     listLoading.value = false
    //   })
    // } catch (error) {
    //   console.log(error)
    //   listLoading.value = false
    // }
  }

  defineExpose({
    queryData: queryData,
    initData,
  })
</script>

<template>
  <div class="field-statistic">
    <el-dialog v-model="visible" :title="queryData.title" width="1170px">
      <vab-query-form style="width: 100%">
        <vab-query-form-left-panel :span="12">
          <el-button :disabled="queryData.listDate.length == 0 || disabled" type="primary" @click="handleDownloadExcel">
            导出
          </el-button>
        </vab-query-form-left-panel>
        <vab-query-form-right-panel :span="12">
          <el-form inline @submit.prevent>
            <el-form-item>
              <el-select v-model="queryData.topCount" @change="topCountChange">
                <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-form>
        </vab-query-form-right-panel>
      </vab-query-form>
      <!-- 表格 -->
      <el-table
        v-loading="listLoading"
        class="field"
        :data="queryData.listDate"
        row-key="id"
        @cell-contextmenu="useTableCopyEvent"
      >
        <el-table-column :align="'center'" label="序号" width="100">
          <template #default="{ $index }">
            {{ $index + 1 }}
          </template>
        </el-table-column>
        <el-table-column :align="'center'" label="值" prop="name" show-overflow-tooltip />
        <el-table-column
          :align="'center'"
          :label="statisticType === 'alert' ? '总攻击次数' : '次数'"
          prop="value"
          show-overflow-tooltip
          width="200"
        />
        <template #empty>
          <div style="height: 300px; line-height: 300px">暂无其他数据</div>
        </template>
      </el-table>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .field {
    :deep() {
      .el-scrollbar {
        height: 480px;
      }
      .el-table__body-wrapper {
        border-bottom: var(--el-table-border);
      }
      .el-scrollbar__bar.is-vertical > div {
        margin-top: 40px;
      }
    }
  }
</style>
