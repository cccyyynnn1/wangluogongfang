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
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import SessionInfo from '@/ecs/site/site-session/session-info.vue'
  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'
  import { updateDisplayApi } from '@/api-ecs/custom-field'
  import { RetrieveIndexType, SearchBySqlParams, IndexTypeTpye } from '~/src/types'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import { searchBySqlApi } from '~/src/api-ecs/retrieve'
  import { useUserStore } from '@/store/modules/user'
  import { TableColumnItemType } from '~/types/store'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import dayjs from 'dayjs'

  const props = defineProps<{
    nodeData: any
    alertData: any
  }>()
  const { getTableColumn, getAllIndexType } = useUserStore()

  const tagType = ref<IndexTypeTpye[]>(getAllIndexType())
  // 检索返回的数据
  const queryPage = reactive({
    total: 0,
    title: '',
    listDate: [] as object[],
    queryTimes: '',
    scrollId: '',
  })
  const detailVisible = ref(false)
  const sessionInfoVisible = ref(false)
  const detailVal = ref()
  const detailDom: { [key: number]: any } = {
    2: DnsDetail,
    3: DbDetail,
    4: FtpDetail,
    5: SmbDetail,
    6: MailDetail,
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
  const elmRef = ref<HTMLElement | null>(null)
  const arrivedInfo = ref()
  const listLoading = ref(false) // 是否加载

  const timeDate = ref() // 自定义时间

  const allDisPlaysFiled = ref<any[]>([]) // 用户所有字段
  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()
  const showTraceabilityField = (fieldCn: string, fieldEn: string) => {
    const { searchSql, indexType, startTime, endTime } = queryForm
    traceabilityFieldRef.value?.initData({
      title: fieldCn,
      query: { searchSql, indexType, startTime, endTime, aggregationFields: fieldEn, whiteType: 0 },
    })
  }
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
  })

  // 表单数据
  const formData = reactive({
    clientIp: '',
    clientPort: '',
    serverIp: '',
    serverPort: '',
  })
  // 打开详情
  const handleInfo = (row: any) => {
    detailVal.value = { ...row, indexType: queryForm.indexType }
    if (queryForm.indexType === 1) return (sessionInfoVisible.value = true)
    detailVisible.value = true
  }

  const dialogVisible = ref(false)
  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)

  onMounted(() => {
    getData()
  })

  // 获取数据
  const getData = async (isClear = true) => {
    getAllDisPlaysFiled()
    jointSQL()
    if (isClear) {
      queryPage.scrollId = ''
      queryPage.listDate = []
    }
    listLoading.value = true
    try {
      const {
        data: { queryTime, resList, total, scrollId },
      } = await searchBySqlApi({ ...queryForm, scrollId: queryPage.scrollId })
      queryPage.listDate = [...queryPage.listDate, ...resList]
      queryPage.queryTimes = (queryTime / 1000).toFixed(1)
      queryPage.total = total
      queryPage.scrollId = scrollId
    } finally {
      listLoading.value = false
    }
  }

  // 切换类型
  function changeDataType(val: RetrieveIndexType) {
    queryForm.indexType = val
    val == 1 ? (queryForm.orderField = 'requestTimeNs') : (queryForm.orderField = 'startTimeNs')
  }

  // 拼接sql
  const jointSQL = () => {
    // @ts-ignore
    const searchSql = props.nodeData.type ? '' : `( 源ip = ${props.nodeData.ip} or 目的ip = ${props.nodeData.ip} )`
    const serverIpValue = formData.serverIp ? `目的IP = ${formData.serverIp}` : ''
    const serverPortValue = formData.serverPort ? `目的端口 = ${formData.serverPort}` : ''
    const clientIpValue = formData.clientIp ? `源IP = ${formData.clientIp}` : ''
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
    queryForm.searchSql = sqlStr
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

  function formatDate(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  function changeCellStyle(title: string) {
    if (['客户端口', '服务端口', '请求方式', '状态码'].includes(title)) {
      return '121'
    } else if (['客户端IP', '服务端IP', 'HOST', 'XFF'].includes(title)) {
      return '120'
    } else if (['请求时间', '响应时长'].includes(title)) {
      return '160'
    } else if (['请求负载', '响应负载'].includes(title)) {
      return '300'
    } else {
      return '150'
    }
  }

  // 弹窗更新页面
  const updateDisplay = async (showField: TableColumnItemType[]) => {
    const { msg } = await updateDisplayApi({
      indexType: queryForm.indexType,
      displayIds: showField.map((i) => i.id),
    })
    ElMessage({ message: msg, type: 'success' })
    getData()
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
    if (props.alertData) {
      const { startTimeNs } = props.alertData
      const newDate = dayjs(formatNstime(startTimeNs))
      const startDate = newDate.subtract(5, 'minute')
      const endDate = newDate.add(5, 'minute')
      timeDate.value = [formatTime(startDate), formatTime(endDate)]
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
    elmRef.value = document.querySelector('.traceabilityTable .el-table__body-wrapper')!
    const { arrivedState } = useScroll(elmRef.value)
    arrivedInfo.value = arrivedState
  })
</script>

<script lang="ts">
  export default {
    name: 'TraceabilityDatas',
  }
</script>

<template>
  <el-form label-position="top" :model="formData">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-form-item label="源IP" prop="clientIp">
          <el-input v-model="formData.clientIp" clearable />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item label="源端口" prop="clientPort">
          <el-input v-model="formData.clientPort" clearable />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item label="目的IP" prop="serverIp">
          <el-input v-model="formData.serverIp" clearable />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item label="目的端口" prop="serverPort">
          <el-input v-model="formData.serverPort" clearable />
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
  <el-row :gutter="20">
    <el-col :span="6">
      <el-select v-model="queryForm.indexType" placeholder="请选择" style="width: 100%" @change="changeDataType">
        <el-option v-for="item in tagType" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </el-col>
    <el-col :offset="6" :span="8">
      <div class="timer">
        <vab-date-time-picker v-model="timeDate" />
      </div>
    </el-col>
    <el-col :span="4">
      <el-button :loading="listLoading" style="float: right; margin-left: 20px" type="primary" @click="getData()">
        检索
      </el-button>
      <el-button di :icon="Setting" style="float: right" @click="dialogVisible = true" />
    </el-col>
  </el-row>
  <el-table
    v-loading="listLoading"
    class="traceabilityTable"
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
            <el-image
              v-if="item.supportAgg"
              class="table-filter"
              :src="require('@/assets/tongji-3.svg')"
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
            <vab-icon
              v-if="item.supportAgg"
              icon="filter-line"
              style="font-size: 12px; margin-left: 6px; vertical-align: -1px !important"
              @click="showTraceabilityField(item.fieldNameCn, item.fieldNameEn)"
            />
          </span>
        </template>
      </el-table-column>
    </template>
    <el-table-column align="center" fixed="right" label="操作" width="80">
      <template #default="{ row }">
        <el-button class="row_action" size="small" @click="handleInfo(row)">详情</el-button>
      </template>
    </el-table-column>
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
    <component :is="detailDom[queryForm.indexType]" v-if="queryForm.indexType !== 1" :info-val="detailVal" />
  </vab-dialog>
  <session-info
    v-if="sessionInfoVisible"
    :info-data="detailVal"
    :show-session-info="sessionInfoVisible"
    @on-close-event="sessionInfoVisible = false"
  />
  <traceability-field-statistic ref="traceabilityFieldRef" />
</template>

<style scoped lang="scss">
  .timer {
    display: flex;
    align-items: center;
  }
  .m-2 {
    margin-right: 20px;
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
</style>
