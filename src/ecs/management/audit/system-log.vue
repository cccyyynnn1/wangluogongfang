<script lang="ts">
  export default {
    name: 'SystemLogs',
  }
</script>

<script setup lang="ts">
  import LogDetails from './log-details.vue'
  import { useTableCopy } from '@/utils'
  import { getAuditLogPageApi } from '@/api-ecs/system'
  import { GetAuditLogQuery, AuditSystemLog, tableSearch } from '@/types'
  import dayjs from 'dayjs'
  import { formatNstime, formatTime } from '@/utils/time'
  import { useScroll } from '@vueuse/core'
  import TableFilterTool from '@/components/table-filter-tool.vue'
  const multipleSelection = ref<AuditSystemLog[]>([])
  const scrollWatch = ref()
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
  const options = [
    // {
    //   value: 1,
    //   label: '进程日志',
    // },
    // {
    //   value: 2,
    //   label: '系统日志',
    // },
    {
      value: 3,
      label: '用户日志',
    },
  ]
  const listLoading = ref(false) // 是否加载
  // 表格数据
  const listDate = ref<AuditSystemLog[]>([])
  const totals = ref(0)
  const timeDate = ref('') // 创建时间
  const isShow = ref(true) // 是否加载
  const queryData = reactive<GetAuditLogQuery>({
    pageNum: 1,
    pageSize: 10,
    level: undefined,
    logType: 3,
    startTime: '',
    endTime: '',
    orderField: 'timeStamp',
    scrollId: undefined,
    searchStr: '',
  })
  const logInfo = ref()
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
      if (isClaer) el.value?.scrollTo({ top: 0 })
    }
  }
  provide(tableSearch, () => {
    queryData.pageNum = 1
    handleSearch(true)
  })
  // 重置
  const handleReset = () => {
    queryData.logType = 3
    queryData.level = undefined
    queryData.searchStr = ''
    timeDate.value = ''
    handleSearch(true)
  }

  // 多选项改变
  const setSelectRows = (val: any) => {
    multipleSelection.value = val
  }

  // 页面改变
  const handleInfo = (row: AuditSystemLog) => {
    logInfo.value = row
    isShowInfo(false)
  }

  // 二级页面的显隐
  const isShowInfo = (val: boolean) => {
    isShow.value = val
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
    const tHeader: string[] = ['源IP', 'URL', '请求类型', '等级', '日志类型', '日志摘要', '创建时间']
    const filterVal: any = ['clientIp', 'apiPath', 'reqType', 'levelStr', 'logTypeStr', 'sourceData', 'timeStamp']
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
  watch(timeDate, () => {
    if (timeDate.value) {
      const [start, end] = timeDate.value
      queryData.startTime = dayjs(start).startOf('day').format('YYYY-MM-DD HH:mm:ss')
      queryData.endTime = dayjs(end).endOf('day').format('YYYY-MM-DD HH:mm:ss')
    } else {
      queryData.startTime = ''
      queryData.endTime = ''
    }
  })
  watch(
    scrollWatch,
    () => {
      if (scrollWatch.value.bottom) handleSearch()
    },
    {
      deep: true,
    }
  )
  const data = ref()

  const el = ref<HTMLElement | null>(null)
  onMounted(() => {
    el.value = document.querySelector('.el-table__body-wrapper')
    queryData.pageSize = Math.round(el.value!.offsetHeight / 40) + 5
    const { arrivedState } = useScroll(el.value)
    data.value = arrivedState
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
        <h3>应用日志</h3>
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
      class="systemLogsTable"
      :data="listDate"
      row-key="id"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <el-table-column fixed="left" show-overflow-tooltip type="selection" width="55" />
      <el-table-column label="源IP" prop="clientIp" show-overflow-tooltip />
      <el-table-column label="URL" prop="apiPath" show-overflow-tooltip />
      <el-table-column label="请求类型" prop="reqType" show-overflow-tooltip />
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
      <el-table-column label="功能说明" prop="operate" show-overflow-tooltip />
      <el-table-column label="用户名" prop="userName" show-overflow-tooltip />
      <el-table-column label="操作内容" prop="function" show-overflow-tooltip />
      <el-table-column label="日志类型" prop="logTypeStr" show-overflow-tooltip />
      <el-table-column label="日志摘要" prop="sourceData" show-overflow-tooltip>
        <template #header="{ column }">
          <table-filter-tool
            v-model:filter-value="queryData.searchStr"
            :column="column"
            filter-type="text"
            :tools="['filter']"
          />
        </template>
      </el-table-column>
      <el-table-column
        :formatter="({ timeStamp }) => formatNstime(timeStamp)"
        label="创建时间"
        prop="timeStamp"
        show-overflow-tooltip
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
      <el-table-column label="操作" width="80">
        <template #default="{ row }">
          <el-button size="small" @click="handleInfo(row)">详情</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <div v-if="!isShow" class="log-detail-box">
      <log-details :info-val="logInfo" :is-show="isShow" @on-saveData="isShowInfo" />
    </div>
  </div>
</template>
<style lang="scss" scoped>
  :deep() {
    .floatBox {
      float: right;
      margin-right: 0;
    }
    // .el-scrollbar {
    //   height: calc(100vh - 230px) !important;
    // }
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
      max-height: calc(100vh - 130px);
      min-height: calc(100vh - 130px);
      overflow-y: auto;
      .el-scrollbar__thumb {
        display: none !important;
      }
    }
  }
  .operation-container {
    position: relative;
    .log-detail-box {
      position: absolute;
      inset: 10px;
      background-color: #fff;
      z-index: 9999;
    }
  }
</style>
