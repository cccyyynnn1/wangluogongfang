<script lang="ts">
  export default {
    name: 'Operation',
  }
</script>

<script setup lang="ts">
  import { useTableCopy } from '@/utils'
  import { getAuditLogPageApi } from '@/api-ecs/system'
  import { GetAuditLogQuery, AuditSystemLog, tableSearch } from '@/types'
  import { formatNstime, formatTime } from '@/utils/time'
  import { useScroll } from '@vueuse/core'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  import dayjs from 'dayjs'
  const multipleSelection = ref<AuditSystemLog[]>([])
  const options = [
    {
      value: 0,
      label: '网络日志',
    },
    {
      value: 1,
      label: '系统日志',
    },

    {
      value: 2,
      label: '数据库日志',
    },
  ]

  const levelOptions = [
    {
      value: 0,
      label: '信息',
    },
    {
      value: 1,
      label: '告警',
    },

    {
      value: 2,
      label: '错误',
    },
  ]
  const protocol = [
    {
      value: 3,
      label: 'TCP',
    },
    {
      value: 4,
      label: 'UDP',
    },
    {
      value: 5,
      label: '文本',
    },
  ]
  const totals = ref(0)

  const listLoading = ref(false) // 是否加载
  // 表格数据
  const listDate = ref<AuditSystemLog[]>([])
  const timeDate = ref() // 创建时间
  const queryData = reactive<GetAuditLogQuery>({
    pageNum: 1,
    pageSize: 10,
    logType: undefined,
    startTime: '',
    endTime: '',
    orderField: 'timeStamp',
    scrollId: undefined,
    clientIp: '',
    apiPath: '',
    level: undefined,
    sysType: undefined,
    protocol: undefined,
    userName: '',
    searchStr: '',
  })
  // 多选项改变
  const setSelectRows = (val: any) => {
    multipleSelection.value = val
  }
  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }

  function formatExcelData(val: any, key: string) {
    switch (key) {
      case 'timeStamp':
        return formatNstime(val[key])
      default:
        return val[key] || ''
    }
  }
  const handleDownloadExcel = () => {
    const tHeader = ['登录名称', '登录IP', '请求URL', '日志类型', '系统类型', '操作内容', '协议', '创建时间']
    const filterVal = [
      'userName',
      'clientIp',
      'apiPath',
      'logTypeStr',
      'sysTypeStr',
      'sourceData',
      'protocolStr',
      'timeStamp',
    ]
    try {
      listLoading.value = true
      import('@/utils/excel').then((excel) => {
        const list = multipleSelection.value.length ? multipleSelection.value : listDate.value
        const data = formatJson(filterVal, list)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `系统日志-${formatTime(new Date().getTime())}`,
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
  provide(tableSearch, () => {
    queryData.pageNum = 1
    handleSearch(true)
  })
  // 查询
  const handleSearch = async (isClaer = false) => {
    if (isClaer) queryData.scrollId = undefined
    listLoading.value = true
    try {
      const {
        data: { resList, scrollId, total },
      } = await getAuditLogPageApi(queryData)
      totals.value = total
      listDate.value = isClaer ? resList : [...listDate.value, ...resList]
      queryData.scrollId = scrollId
    } finally {
      listLoading.value = false
    }
  }

  // 重置
  const handleReset = () => {
    timeDate.value = undefined
    queryData.logType = undefined
    queryData.startTime = ''
    queryData.endTime = ''
    queryData.scrollId = undefined
    queryData.clientIp = ''
    queryData.apiPath = ''
    queryData.sysType = undefined
    queryData.protocol = undefined
    queryData.userName = ''
    queryData.searchStr = ''
    queryData.level = undefined
    handleSearch(true)
  }

  watch(timeDate, () => {
    if (timeDate.value) {
      const [start, end] = timeDate.value
      queryData.startTime = dayjs(start).format('YYYY-MM-DD HH:mm:ss')
      queryData.endTime = dayjs(end).format('YYYY-MM-DD HH:mm:ss')
    } else {
      queryData.startTime = ''
      queryData.endTime = ''
    }
  })

  const data = ref()

  const el = ref<HTMLElement | null>(null)
  onMounted(() => {
    el.value = document.querySelector('.el-table__body-wrapper')
    queryData.pageSize = Math.round(el.value!.offsetHeight / 40) + 5
    const { arrivedState } = useScroll(el.value)
    data.value = arrivedState
    timeDate.value = [dayjs().startOf('day'), dayjs().endOf('day')]
    handleSearch()
  })

  watch(
    () => data.value,
    () => {
      if (data.value.bottom) {
        handleSearch(false)
      }
    },
    { deep: true }
  )
</script>

<template>
  <div class="operation-container">
    <vab-query-form>
      <vab-query-form-left-panel :span="12">
        <h3>系统日志</h3>
      </vab-query-form-left-panel>
      <vab-query-form-right-panel :span="12">
        共：{{ totals }}条&nbsp;&nbsp;&nbsp;&nbsp;
        <el-button plain style="margin-right: 0 !important" type="primary" @click="handleDownloadExcel">
          导出日志
        </el-button>
      </vab-query-form-right-panel>
    </vab-query-form>
    <el-table
      v-loading="listLoading"
      class="operationTable my-table"
      :data="listDate"
      row-key="id"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column fixed="left" show-overflow-tooltip type="selection" width="55" />
      <el-table-column label="名称" prop="userName" show-overflow-tooltip width="120px">
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.userName"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="源IP" prop="clientIp" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.clientIp"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="目的IP" prop="serverIp" show-overflow-tooltip />
      <el-table-column label="Url" prop="apiPath" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.apiPath"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="日志类型" prop="logTypeStr" show-overflow-tooltip width="120px">
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.logType"
            :column="column"
            :filter-option="options"
            filter-type="select"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="等级" prop="levelStr" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.level"
            :column="column"
            :filter-option="levelOptions"
            filter-type="select"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="系统类型" prop="sysTypeStr" show-overflow-tooltip width="100px" />
      <el-table-column label="操作内容" prop="sourceData" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.searchStr"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column label="协议" prop="protocolStr" show-overflow-tooltip width="80px">
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.protocol"
            :column="column"
            :filter-option="protocol"
            filter-type="select"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column
        :formatter="({ timeStamp }) => formatNstime(timeStamp)"
        label="创建时间"
        prop="timeStamp"
        show-overflow-tooltip
        width="180px"
      >
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="timeDate"
            :column="column"
            filter-type="date-range"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
  </div>
</template>

<style lang="scss" scoped>
  .operation-container {
    position: relative;

    .w12 {
      width: 150px;
    }
    .floatBox {
      float: right;
      margin-right: 0;
    }
    .btn {
      position: absolute;
      right: 10px;
      bottom: 0px;
    }
  }
  :deep() {
    .el-table__inner-wrapper {
      width: 100%;
    }
    .el-table__body-wrapper {
      position: initial !important;
      width: 100%;
    }
    .el-scrollbar {
      position: initial !important;
    }
    .el-scrollbar__bar.is-horizontal {
      position: absolute;
      display: block !important;
    }
  }
  :deep() {
    .el-table__body-wrapper {
      max-height: calc(100vh - 140px);
      min-height: calc(100vh - 140px);
      overflow-y: auto;
      .el-scrollbar__thumb {
        display: none !important;
      }
    }
  }
</style>
