<script setup lang="ts">
  import ApplicationConfiguration from '@/ecs/assets/components/application-configuration.vue'
  import { Setting } from '@element-plus/icons-vue'
  import DbDetail from '@/ecs/retrieve/components/detail-db.vue'
  import SmbDetail from '@/ecs/retrieve/components/detail-smb.vue'
  import DnsDetail from '@/ecs/retrieve/components/detail-dns.vue'
  import MailDetail from '@/ecs/retrieve/components/detail-mail.vue'
  import FtpDetail from '@/ecs/retrieve/components/detail-ftp.vue'
  import DetailOther from '@/ecs/retrieve/components/detail-other.vue'
  import DetailFile from '@/ecs/retrieve/components/detail-file.vue'
  import AlertDetail from '@/ecs/alert/components/alert-detail.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import SessionInfo from '@/ecs/site/site-session/session-info.vue'
  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'
  import { updateDisplayApi } from '@/api-ecs/custom-field'
  import { RetrieveIndexType, SearchBySqlParams, IndexTypeTpye } from '~/src/types'
  import { getSystemConfigApi } from '~/src/api-ecs/system'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import dayjs from 'dayjs'
  import { searchBySqlApi, getKeyVarApi } from '~/src/api-ecs/retrieve'
  import { useUserStore } from '@/store/modules/user'
  import { TableColumnItemType } from '~/types/store'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import { isIP } from '~/src/utils/validate'
  /* alertData 和 timeRange 必须有一个，alertData：{startTimeNs} */
  const props = defineProps<{
    nodeData: any
    startTime?: any
    timeRange?: [string, string]
    searchSql?: string
  }>()
  const moduleEnable = ref(false)
  const filterSql = ref(props.searchSql || '')
  const { getTableColumn, getAllIndexType } = useUserStore()
  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())
  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()
  const showTraceabilityField = (fieldCn: string, fieldEn: string) => {
    const { searchSql, indexType, startTime, endTime } = queryForm
    traceabilityFieldRef.value?.initData({
      title: fieldCn,
      query: { searchSql, indexType, startTime, endTime, aggregationFields: fieldEn, whiteType: 0 },
    })
  }
  const table_el = ref<HTMLElement | null>(null)
  const arrivedInfo = ref()
  // 检索返回的数据
  const queryPage = reactive({
    total: 0,
    title: '',
    listDate: [] as object[],
    queryTimes: '',
  })
  const detailVisible = ref(false)
  const sessionInfoVisible = ref(false)
  const detailVal = ref()
  const detailValIndex = ref(-1)
  const detailDom: { [key: number]: any } = {
    2: DnsDetail,
    3: DbDetail,
    4: FtpDetail,
    5: SmbDetail,
    6: MailDetail,
    10: AlertDetail,
    12: DetailOther,
    13: DetailOther,
    14: DetailOther,
    15: DetailOther,
    16: DetailOther,
    17: DetailOther,
    18: DetailOther,
    19: DetailOther,
    20: DetailOther,
    21: DetailOther,
    25: DetailOther,
    26: DetailOther,
    27: DetailFile,
    28: DetailOther,
    32: DetailOther,
  }
  const dataType = ref<RetrieveIndexType>(1)

  const listLoading = ref(false) // 是否加载

  // 自定义时间
  const timeDate = ref()

  const allDisPlaysFiled = ref<any[]>([]) // 用户所有字段

  // 检索信息表单数据
  const queryForm = reactive<SearchBySqlParams>({
    searchSql: '',
    indexType: 1,
    pageNum: 1,
    pageSize: 100,
    orderType: 'desc',
    orderField: 'requestTimeNs',
    startTime: '',
    endTime: '',
    scrollId: '',
  })

  // 表单数据
  const formData = reactive({
    clientIp: '',
    clientPort: '',
    serverIp: '',
    serverPort: '',
  })
  // 打开详情
  const handleInfo = (row: any, index: number) => {
    detailValIndex.value = index
    detailVal.value = { ...row, indexType: queryForm.indexType }
    if (queryForm.indexType === 1) return (sessionInfoVisible.value = true)
    detailVisible.value = true
  }
  const dialogVisible = ref(false)
  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)
  const tableColumn: { fieldNameCn: string; fieldNameEn: string }[] = [
    { fieldNameCn: '源IP', fieldNameEn: 'clientIp' },
    { fieldNameCn: '源端口', fieldNameEn: 'clientPort' },
    { fieldNameCn: '目的IP', fieldNameEn: 'serverIp' },
    { fieldNameCn: '目的端口', fieldNameEn: 'serverPort' },
  ]
  onMounted(() => {
    getData()
  })

  // 获取数据
  const getData = async (isClear = true) => {
    getAllDisPlaysFiled()
    jointSQL()
    if (isClear) {
      queryForm.scrollId = ''
      queryPage.listDate = []
    }
    listLoading.value = true
    try {
      const {
        data: { queryTime, resList, total, scrollId },
      } = await searchBySqlApi({ ...queryForm })
      const list = resList || []
      queryPage.listDate = [...queryPage.listDate, ...list]
      queryPage.total = total
      queryForm.scrollId = scrollId
      detailValIndex.value = -1
    } finally {
      listLoading.value = false
    }
  }

  // 切换类型
  function changeDataType(val: RetrieveIndexType) {
    queryForm.indexType = val
    dataType.value = val
    val == 1 ? (queryForm.orderField = 'requestTimeNs') : (queryForm.orderField = 'startTimeNs')
  }

  // 拼接sql
  const jointSQL = () => {
    // @ts-ignore
    const searchSql = props.nodeData.type ? '' : `( 源ip = "${props.nodeData.ip}" or 目的ip = "${props.nodeData.ip}" )`
    const serverIpValue = formData.serverIp ? `目的IP = "${formData.serverIp}"` : ''
    const serverPortValue = formData.serverPort ? `目的端口 = ${formData.serverPort}` : ''
    const clientIpValue = formData.clientIp ? `源IP = "${formData.clientIp}"` : ''
    const clientPortValue = formData.clientPort ? `源端口 = ${formData.clientPort}` : ''
    const arr = [clientIpValue, clientPortValue, serverIpValue, serverPortValue].filter(Boolean).join(' and ')
    let sqlStr = ''
    if (arr.length > 0 && searchSql) {
      sqlStr = `${searchSql} and ${arr}`
    } else if (arr.length > 0 && !searchSql) {
      sqlStr = arr
    } else {
      sqlStr = searchSql
    }
    const filterSqlVal = isIP(filterSql.value)
      ? `源IP = "${filterSql.value}" or 目的IP = "${filterSql.value}"`
      : filterSql.value
    queryForm.searchSql = filterSqlVal ? `${sqlStr} and ( ${filterSqlVal} )` : sqlStr
  }
  const getModelStatus = async () => {
    const {
      data: { module_config },
    } = await getSystemConfigApi({ keys: 'module_config' })
    moduleEnable.value = JSON.parse(module_config.value).isEnable || false
  }
  // 获取用户所有字段
  const getAllDisPlaysFiled = async () => {
    const { data } = await getAllDisPlaysFiledApi()
    const target = data[queryForm.indexType] as number[]
    const sourse = getTableColumn(queryForm.indexType)
    allDisPlaysFiled.value = target.map((key) => {
      return sourse.find((item) => {
        return item.id === key
      })
    })
  }

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryForm.pageSize = val
    getData()
  }

  // 页面改变
  const handleCurrentChange = (val: number) => {
    queryForm.pageNum = val
    getData()
  }

  function formatDate(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  function changeCellStyle(title: string) {
    if (['客户端口', '服务端口', '请求方式', '状态码'].includes(title)) {
      return '111'
    } else if (['客户端IP', '服务端IP', 'HOST', 'XFF'].includes(title)) {
      return '120'
    } else if (['请求时间', '响应时长'].includes(title)) {
      return '190'
    } else if (['请求负载', '响应负载'].includes(title)) {
      return '300'
    } else {
      return '150'
    }
  }

  // 弹窗更新页面
  const updateDisplay = async (showField: TableColumnItemType[]) => {
    const arr: number[] = []
    showField.forEach((item: TableColumnItemType) => {
      arr.push(item.id)
    })
    const { msg } = await updateDisplayApi({
      indexType: queryForm.indexType,
      displayIds: arr,
    })
    ElMessage({ message: msg, type: 'success' })
    getData()
  }
  const changeCurrentItemEvent = (val: boolean) => {
    if (detailValIndex.value >= queryPage.listDate.length - 5 && queryPage.listDate.length < queryPage.total) {
      getData(false)
    }
    const newIndex = val ? detailValIndex.value + 1 : detailValIndex.value - 1
    detailVal.value = { ...queryPage.listDate[newIndex], indexType: queryForm.indexType }
    detailValIndex.value = newIndex
  }
  // 得到自定义时间
  watchEffect(() => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      queryForm.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
      queryForm.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
    }
  })
  watchEffect(() => {
    if (props.startTime) {
      const newDate = dayjs(formatNstime(props.startTime))
      const startDate = newDate.subtract(5, 'minute')
      const endDate = newDate.add(5, 'minute')
      timeDate.value = [formatTime(startDate), formatTime(endDate)]
    }
  })
  watchEffect(() => {
    if (props.timeRange) {
      timeDate.value = props.timeRange
    }
  })
  watch(
    () => props.nodeData,
    () => {
      const { type, snat, serverIp } = props.nodeData
      // 当type存在的时候表示是连线
      if (type) {
        formData.clientIp = snat
        formData.serverIp = serverIp
      }
    },
    {
      deep: true,
      immediate: true,
    }
  )
  watch(
    () => arrivedInfo.value?.bottom,
    () => {
      if (arrivedInfo.value?.bottom) {
        getData(false)
      }
    }
  )
  onMounted(() => {
    getModelStatus()
    table_el.value = document.querySelector('.my-table.relations .el-table__body-wrapper')
    const { arrivedState } = useScroll(table_el.value)
    arrivedInfo.value = arrivedState
  })
</script>

<script lang="ts">
  export default {
    name: 'TraceabilityRelationsDatas',
  }
</script>

<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <div class="timer">
        <el-select
          v-model="queryForm.indexType"
          placeholder="请选择"
          style="width: 250px; margin-right: 15px"
          @change="changeDataType"
        >
          <el-option v-for="item in tagType" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-input v-model="filterSql" clearable placeholder="请输入过滤条件" style="width: 45%; margin-right: 15px" />
        <vab-date-time-picker v-model="timeDate" />
        <el-button :icon="Setting" style="margin: 0 3px 0 15px" @click="dialogVisible = true" />
        <el-button :loading="listLoading" type="primary" @click="() => getData()">检索</el-button>
      </div>
    </el-col>
  </el-row>
  <el-table
    v-loading="listLoading"
    class="my-table relations"
    :data="queryPage.listDate"
    style="margin-top: 20px; width: 100%"
  >
    <el-table-column :index="curIndex" label="序号" type="index" width="55px" />
    <template v-for="item in allDisPlaysFiled" :key="item.id">
      <el-table-column
        v-if="item.fieldNameCn.includes('时间')"
        align="center"
        :label="item.fieldNameCn"
        :prop="item.fieldNameEn"
        :resizable="false"
        show-overflow-tooltip
        width="160px"
      >
        <template #header>
          <span style="cursor: pointer">
            {{ item.fieldNameCn }}
            <vab-icon
              v-if="item.supportAgg"
              icon="filter-line"
              style="font-size: 12px; margin-left: 6px; vertical-align: -1px !important"
              @click="showTraceabilityField(item.fieldNameCn, item.fieldNameEn)"
            />
          </span>
        </template>
        <template #default="{ row }">
          {{ formatDate(row, item.fieldNameEn) }}
        </template>
      </el-table-column>
      <el-table-column
        v-else
        :label="item.fieldNameCn"
        :min-width="changeCellStyle(item.fieldNameCn)"
        :prop="item.fieldNameEn"
        show-overflow-tooltip
      >
        <template #header>
          <span style="cursor: pointer">
            {{ item.fieldNameCn }}
            <el-image
              v-if="item.supportAgg"
              class="table-filter"
              :src="require('@/assets/tongji-3.svg')"
              @click="showTraceabilityField(item.fieldNameCn, item.fieldNameEn)"
            />
          </span>
        </template>
      </el-table-column>
    </template>
    <el-table-column align="center" fixed="right" label="操作" width="80">
      <template #default="{ row, $index }">
        <el-button class="row_action" size="small" @click="handleInfo(row, $index)">详情</el-button>
      </template>
    </el-table-column>
    <template #empty>
      <el-empty class="vab-data-empty" description="暂无数据" />
    </template>
  </el-table>
  <application-configuration
    v-model="dialogVisible"
    :fields="allDisPlaysFiled"
    :retrieve-index-type="dataType"
    @handleok="updateDisplay"
  />

  <vab-dialog
    v-model="detailVisible"
    align-center
    :close-on-click-modal="false"
    destroy-on-close
    title="详情"
    width="1175px"
  >
    <component
      :is="detailDom[queryForm.indexType]"
      v-if="queryForm.indexType !== 1"
      v-model:alert-detail-visible="detailVisible"
      :info-val="detailVal"
      :module-enable="moduleEnable"
      :next="detailValIndex < queryPage.total - 1"
      :prev="detailValIndex > 0"
      :select-alert="detailVal"
      @on-skip-event="changeCurrentItemEvent"
    />
  </vab-dialog>
  <session-info
    v-if="sessionInfoVisible"
    :info-data="detailVal"
    :next="detailValIndex < queryPage.total - 1"
    :prev="detailValIndex > 0"
    :show-session-info="sessionInfoVisible"
    @on-close-event="sessionInfoVisible = false"
    @on-skip-event="changeCurrentItemEvent"
  />
  <traceability-field-statistic ref="traceabilityFieldRef" />
</template>

<style scoped lang="scss">
  .my-table {
    position: relative;
  }
  :deep() {
    .el-table__body-wrapper {
      position: initial !important;
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
      max-height: 409px;
      min-height: 409px;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0 !important;
        height: 0;
      }
    }
  }
  .timer {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .m-2 {
    margin-right: 20px;
  }
</style>
