<script setup lang="ts">
  import { getTabelListAPI, searchBySqlApi } from '~/src/api-ecs/retrieve'

  import { getBySiteIdApi, siteSearchApi } from '~/src/api-ecs/site'

  import { MonitoringItem, filedDataType } from '~/src/types'

  import { useUserStore } from '@/store/modules/user'

  import { useTableCopy, uuid } from '@/utils'

  import dayjs from 'dayjs'

  const { getTableColumn } = useUserStore()

  const listDate = ref<object[]>([]) // 表格数据

  const listLoading = ref(false) // 是否加载

  const currentRow = ref()

  const tableColumn = ref<filedDataType[]>([]) // 表头

  const timeSetting = ref()

  const props = defineProps<{
    currentItem: MonitoringItem
    timeSetting: any
  }>()

  const queryData = reactive({
    pageNum: 1,
    pageSize: 10,
  })

  const initData = async () => {
    listLoading.value = true
    currentRow.value = props.currentItem
    const res = currentRow.value.data
    queryData.pageNum = currentRow.value.data.pageNum
    queryData.pageSize = currentRow.value.data.pageSize
    timeSetting.value = props.timeSetting
    let allField: filedDataType[] = []
    if (res.indexType == 22) {
      allField = getTableColumn(1)
    } else {
      allField = getTableColumn(res.indexType)
      if (res.indexType == 1 && res.module == '0') {
        allField = getTableColumn(30)
        res.displayFields = []
        if (currentRow.value.hasError) return (listLoading.value = false)
        const {
          data: { displayFieldsArr },
        } = await getBySiteIdApi({ id: res.siteSessionId, apiId: res.siteApiId })
        displayFieldsArr.forEach((item: number) => {
          allField.forEach((td) => {
            if (item == td.id) {
              res.displayFields.push(td)
            }
          })
        })
      }
    }
    tableColumn.value = res.displayFields
      ?.map((i: any) => {
        return allField.find((item) => item.id === i.id)
      })
      .filter(Boolean)
    try {
      if (res.module == '0') {
        if (currentRow.value.hasError) return (listLoading.value = false)
        const {
          data: { resList },
        } = await siteSearchApi({
          indexType: res.indexType,
          orderField: res.orderField,
          orderType: res.orderType,
          pageNum: res.pageNum,
          pageSize: res.pageSize,
          startTime: timeSetting.value.startTime,
          endTime: timeSetting.value.endTime,
          siteSessionId: res.siteSessionId,
          siteApiId: res.siteApiId,
        })
        // queryPage.listDate = []
        listDate.value = resList || []
      } else if (res.indexType == 24) {
        const query = {
          topField: 'topField',
          timeRange: '',
          flowProbeIds: ['0'],
          searchTable: 'eventStat',
          searchModel: 'all',
          eventFilterList: [],
          keyword: undefined,
          ip: '',
          clientIp: '',
          serverIp: '',
          serverPort: '',
        }
        query.timeRange = `${timeSetting.value.startTime} - ${timeSetting.value.endTime}`
        query.flowProbeIds = res.flowProbeIds
        query.searchModel = res.searchModel
        query.ip = res.ip
        query.clientIp = res.clientIp
        query.serverIp = res.serverIp
        query.serverPort = res.serverPort
        const keyword = res.ip || res.clientIp || res.serverIp || res.serverPort || undefined
        query.keyword = keyword
        // query.eventFilterList =
        // remark == 'node' ? nodelist.value : eventFilterList.value && JSON.parse(JSON.stringify(eventFilterList.value))
        const res1 = JSON.stringify({ ...query })
        const result = await getTabelListAPI({
          top: res.topValue || '100',
          moduleType: 'eventStat',
          query: res1,
          groupBy: 'clientIp,serverIp,serverPort',
        })
        listDate.value = result.list
      } else {
        const {
          data: { resList },
        } = await searchBySqlApi({
          indexType: res.indexType,
          orderField: res.orderField,
          orderType: res.orderType,
          pageNum: res.pageNum,
          pageSize: res.pageSize,
          searchSql: res.searchSql,
          startTime: timeSetting.value.startTime,
          endTime: timeSetting.value.endTime,
        })
        listDate.value = resList || []
      }
    } finally {
      listLoading.value = false
    }
  }

  const useTableCopyEvent = (row: any, column: any, cell: any, event: any) => {
    const arr = getList(row, column)
    useTableCopy(row, column, cell, event, arr)
  }

  const params = ref()
  const router = useRouter()
  const handleUpdateCallback = () => {
    if (params.value.searchSql) {
      params.value.sql = `(${params.value.sql} and  ${params.value.searchSql})`
    }
    const resolveRouter = router.resolve({
      path: '/retrieve/index',
      query: { info: encodeURIComponent(JSON.stringify(params.value)) },
    })
    window.open(resolveRouter.href, '_blank')
  }

  const NOEXIST = 'not_exist'
  // 设置右键数组
  const getList = (row: any, column: any) => {
    const target = column.property
    let value = undefined
    let label = ''
    if (target) {
      value = row[target]
      const res = tableColumn.value.find((item: any) => {
        return item.fieldNameEn == target
      })
      label = res!.fieldNameCn
      params.value = {
        key: target,
        value: value ? value : NOEXIST,
        label: label,
        sql: value ? `${label} = ${value}` : `${label} ${NOEXIST}`,
        indexType: currentRow.value.data.indexType == 22 ? 1 : currentRow.value.data.indexType,
        searchSql: currentRow.value.data.searchSql,
      }
    }

    const arr = [
      {
        label: '详细查询',
        callback: handleUpdateCallback,
        value: '',
      },
    ]
    return arr
  }

  onMounted(() => {
    initData()
  })

  function formatDate(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  // 获取表格序号
  const curIndex = computed(() => (queryData.pageNum - 1) * queryData.pageSize + 1)

  const svg = `
        <path class="path" d="
          M 30 15
          L 28 17
          M 25.61 25.61
          A 15 15, 0, 0, 1, 15 30
          A 15 15, 0, 1, 1, 27.99 7.5
          L 15 15
        " style="stroke-width: 4px; fill: rgba(0, 0, 0, 0)"/>
      `

  defineExpose({
    initData,
  })
</script>

<script lang="ts">
  export default {
    name: 'ViewTable',
  }
</script>
<template>
  <el-table
    v-loading="listLoading"
    :border="true"
    :data="listDate"
    element-loading-background="rgba(122, 122, 122, 0.8)"
    :element-loading-spinner="svg"
    element-loading-svg-view-box="-10, -10, 50, 50"
    element-loading-text="加载中..."
    @cell-contextmenu="useTableCopyEvent"
  >
    <el-table-column v-if="tableColumn.length > 0" :align="'center'" label="序号" width="65">
      <template #default="{ $index }">
        <span>{{ curIndex + $index }}</span>
      </template>
    </el-table-column>
    <template v-for="item in tableColumn" :key="item?.id">
      <el-table-column
        v-if="item.fieldNameCn.includes('时间')"
        align="center"
        min-width="160"
        :prop="item.fieldNameEn"
        :resizable="true"
        show-overflow-tooltip
      >
        <template #header>
          {{ item.fieldNameCn }}
        </template>
        <template #default="{ row }">
          {{ formatDate(row, item.fieldNameEn) }}
        </template>
      </el-table-column>
      <el-table-column
        v-else
        align="center"
        :min-width="120"
        :prop="item.fieldNameEn"
        :resizable="true"
        show-overflow-tooltip
      >
        <template #header>
          {{ item.fieldNameCn }}
        </template>
      </el-table-column>
      <!-- <el-table-column v-else align="center" :label="item.fieldNameCn" min-width="100" :prop="item.fieldNameEn"
                                                                                                                                                                                    :resizable="false" show-overflow-tooltip /> -->
    </template>
    <template #empty>暂无数据</template>
  </el-table>
</template>

<style scoped lang="scss"></style>
