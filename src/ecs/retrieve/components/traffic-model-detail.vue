<script lang="ts">
  export default {
    name: 'TrafficModelDetail',
  }
</script>
<script setup lang="ts">
  import { getRuleConfigsApi, searchBySqlApi } from '@/api-ecs/retrieve'
  import { PublishRuleParams } from '@/types/index'
  import numberFormatte from '@/utils/number'
  import EditDetail from '@/ecs/retrieve/components/edit-detail.vue'
  import VabChart from '@/plugins/VabChart/index.vue'
  import SessionInfo from '@/ecs/site/site-session/session-info.vue'
  import TraceabilityFieldStatistic from '@/ecs/alert/components/field-statistic.vue'
  import dayjs from 'dayjs'
  import { TableColumnItemType } from '~/types/store'
  import { useTableCopy } from '@/utils'
  import { useScroll } from '@vueuse/core'

  const data = ref()
  const router = useRouter()
  const el = ref<HTMLElement | null>(null)
  const $baseMessage: any = inject('$baseMessage')
  // 详情弹窗Visible 除HTTP检索使用
  const detailVisible = ref(false)
  const detailVal = ref()
  const { params } = useRoute()

  const traceabilityFieldRef = ref<InstanceType<typeof TraceabilityFieldStatistic>>()
  const showTraceabilityField = (fieldCn: string, fieldEn: string) => {
    const { searchSql, indexType, startTime, endTime } = queryForm
    traceabilityFieldRef.value?.initData({
      title: fieldCn,
      query: { searchSql, indexType, startTime, endTime, aggregationFields: fieldEn },
    })
  }
  // 检索信息表单数据
  const queryForm = reactive<PublishRuleParams>({
    searchSql: '',
    indexType: 1,
    pageNum: 1,
    pageSize: 100,
    orderType: 'desc',
    orderField: 'requestTimeNs',
    startTime: '',
    endTime: '',
    displayFields: [],
    ruleName: '',
    searchTime: 0,
  })

  // 获取表格序号
  const curIndex = computed(() => (queryForm.pageNum - 1) * queryForm.pageSize + 1)

  // 自动刷新间隔（秒）
  const timeDuration = ref(30)
  // 是否自动刷新
  const switchValue = ref(false)
  // 定时器ID
  const timerId = ref<NodeJS.Timer>()
  const timeDuratioOptions = [
    {
      value: 30,
      label: '30秒',
    },
    {
      value: 60,
      label: '1分钟',
    },
    {
      value: 120,
      label: '2分钟',
    },
    {
      value: 300,
      label: '5分钟',
    },
  ]

  // 图表配置
  const retrieveChartOption = reactive({
    xAxis: {
      type: 'category',
      data: null,
    },
    yAxis: {
      type: 'value',
    },
    series: [
      {
        data: null,
        type: 'bar',
        color: 'rgb(63, 157, 248)',
        backgroundStyle: {
          color: 'rgba(220, 220, 220, 0.8)',
        },
      },
    ],
    grid: {
      x: '5%',
      y: '6%',
      x2: '7%',
      y2: '40%',
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
  })
  const tableColumn = ref<TableColumnItemType[]>([])
  const listLoading = ref(false) // 是否加载
  const showDialog = ref(false) // 是否显示编辑页
  const isFullscreen = ref(false)
  // 表格数据
  const layout = ref('sizes, prev, pager, next, jumper')
  const queryPage = reactive({
    total: 0,
    title: '',
    listDate: [] as object[],
    scrollId: '',
    searchSql: '',
  })

  async function getRuleConfigsHandel() {
    const {
      data: { displayFields, indexType, pageSize, ruleName, searchSql, searchTime, id },
    } = await getRuleConfigsApi(+params.id)
    queryForm.indexType = +indexType || 1
    queryForm.pageSize = pageSize || 10
    queryForm.ruleName = ruleName || ''
    queryForm.searchSql = searchSql || ''
    queryForm.searchTime = searchTime || 1
    queryForm.id = id || undefined
    tableColumn.value = displayFields as TableColumnItemType[]
    await queryData()
  }

  const resetData = () => {
    queryPage.scrollId = ''
    queryPage.listDate = []
  }

  // 详情
  const handleInfo = (row: object) => {
    detailVal.value = { ...row, indexType: 1 }
    detailVisible.value = true
  }

  // 检索
  const queryData = async (remark = true) => {
    if (remark) {
      resetData()
    }
    listLoading.value = true
    const timeDate = dayjs()
    const endDate = timeDate.format('YYYY-MM-DD HH:mm:ss')
    queryForm.startTime = timeDate.subtract(queryForm.searchTime || 1, 'minute').format('YYYY-MM-DD HH:mm:ss')
    queryForm.endTime = endDate
    const sql = {
      searchSql: `( ${queryForm.searchSql} ) and ( ${queryPage.searchSql} )`,
    }
    const queryFormData = queryPage.searchSql ? { ...queryForm, ...sql } : queryForm
    try {
      const {
        data: { resList, sumaryMap, total, scrollId },
        code,
      } = await searchBySqlApi(queryFormData)
      queryPage.scrollId = scrollId
      if (resList.length > 0) {
        const id = resList[0].id
        const flag = queryPage.listDate.some((item: any) => {
          return item.id == id
        })
        if (flag) return false
        queryPage.listDate = [...queryPage.listDate, ...resList]
      }
      queryPage.total = total
      if (sumaryMap) {
        const { count, date } = JSON.parse(sumaryMap)
        retrieveChartOption.xAxis.data = date || []
        retrieveChartOption.series[0].data = count || []
      }
    } finally {
      listLoading.value = false
    }
  }

  // 多选项改变
  const setSelectRows = () => {}

  // 页容量改变
  const handleSizeChange = (val: number) => {
    queryForm.pageSize = val
    queryData()
  }

  // 页面改变
  const handleCurrentChange = () => {
    queryData()
  }

  // 展开编辑页
  const handleEdit = (val: boolean) => {
    showDialog.value = val
  }
  function formatDate(row: any) {
    const time = row.requestTimeNs / 100000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss.SSS')
  }

  function changeCellStyle(title: string) {
    if (['客户端口', '服务端口', '请求方式', '状态码'].includes(title)) {
      return '101'
    } else if (['客户端IP', '服务端IP', 'HOST', 'XFF'].includes(title)) {
      return '120'
    } else if (['请求时间', '响应时长'].includes(title)) {
      return '160'
    } else {
      return '150'
    }
  }
  // 全屏的开关
  const clickFullScreen = () => {
    isFullscreen.value = !isFullscreen.value
  }

  function configurationHandel(showField: TableColumnItemType[]) {
    tableColumn.value = showField
    getRuleConfigsHandel()
  }

  watchEffect(() => {
    if (switchValue.value) {
      if (timerId.value) clearInterval(timerId.value)
      timerId.value = setInterval(() => queryData(), timeDuration.value * 1000)
    } else {
      clearInterval(timerId.value)
    }
  })

  onMounted(async () => {
    getRuleConfigsHandel()
    el.value = document.querySelector('.el-table__body-wrapper')
    const { arrivedState } = useScroll(el.value)
    data.value = arrivedState
  })

  watch(
    () => data.value,
    () => {
      if (data.value.bottom) {
        queryData(false)
      }
    },
    { deep: true }
  )
  onUnmounted(() => {
    if (timerId.value) clearInterval(timerId.value)
  })
</script>

<template>
  <div class="traffic-model-detail-container" :class="{ 'vab-fullscreen': isFullscreen }">
    <!-- 头部检索栏 -->
    <vab-query-form style="width: 100%">
      <vab-query-form-right-panel :span="24">
        <div class="back">
          <el-button
            type="primary"
            @click="() => router.push({ name: 'retrieveIndex', query: { params: 'trafficModel' } })"
          >
            返回
          </el-button>
          <div class="name">{{ queryForm.ruleName }}</div>
        </div>
        <el-select v-model="timeDuration" class="m-2" :disabled="switchValue">
          <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-switch v-model="switchValue" inactive-text="自动刷新:" style="margin-left: 20px" />
        <div class="btn">
          <!-- <vab-icon
            class="icon"
            :icon="isFullscreen ? 'fullscreen-exit-fill' : 'fullscreen-fill'"
            @click="clickFullScreen"
          /> -->
          <vab-icon icon="edit-2-line" style="margin-left: 6px" @click="handleEdit(true)" />
        </div>
      </vab-query-form-right-panel>
      <vab-query-form-left-panel :span="24">
        <el-form class="my-form" inline @submit.prevent>
          <el-form-item class="my-input">
            <el-input v-model="queryPage.searchSql" clearable placeholder="请输入查询语句" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="() => queryData()">检索</el-button>
          </el-form-item>
        </el-form>
      </vab-query-form-left-panel>
    </vab-query-form>
    <!-- 图表数据 -->
    <div class="chart">
      <div class="panel">
        <div class="num">{{ numberFormatte.format(queryPage.total) }}</div>
        <div style="color: #a9acb3">查询结果总数</div>
      </div>
      <vab-chart class="target-echart2" :option="retrieveChartOption" theme="vab-echarts-theme" />
    </div>
    <!-- 表格 -->
    <el-table
      v-loading="listLoading"
      :border="true"
      class="my-table"
      :data="queryPage.listDate"
      style="margin-top: 20px"
      @cell-contextmenu="useTableCopy"
      @selection-change="setSelectRows"
    >
      <!-- <el-table-column align="center" type="selection" width="70" label="导出" fixed="left" /> -->
      <el-table-column
        :align="'center'"
        fixed="left"
        :index="(index) => curIndex + index"
        label="序号"
        type="index"
        width="55"
      />
      <template v-for="item in tableColumn" :key="item.id">
        <el-table-column
          align="center"
          :label="item.fieldNameCn"
          :min-width="changeCellStyle(item.fieldNameCn)"
          :prop="item.fieldNameEn"
          resizable
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
          <template v-if="item.fieldNameCn.includes('时间')" #default="{ row }">
            {{ formatDate(row) }}
          </template>
        </el-table-column>
      </template>
      <el-table-column :align="'center'" fixed="right" label="操作" width="100">
        <template #default="{ row }">
          <el-button size="small" @click="handleInfo(row)">详情</el-button>
        </template>
      </el-table-column>
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <edit-detail
      v-if="showDialog"
      :fields="tableColumn"
      :query-form="queryForm"
      :show-dialog="showDialog"
      @handleok="configurationHandel"
      @on-close-event="handleEdit"
    />
    <session-info
      v-if="detailVisible"
      :info-data="detailVal"
      :show-session-info="detailVisible"
      @on-close-event="detailVisible = false"
    />
    <traceability-field-statistic ref="traceabilityFieldRef" />
  </div>
</template>

<style scoped lang="scss">
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
  .traffic-model-detail-container {
    position: relative;
    .back {
      position: absolute;
      display: inline-flex;
      align-items: center;
      left: 0;
      .name {
        margin-left: 20px;
      }
    }
    .btn {
      font-size: 15px;
      color: #a9acb3;
      margin-left: 10px;
      &:hover {
        cursor: pointer;
      }
    }
    .my-form {
      width: 100%;
      display: flex;

      .my-input {
        flex: 1;
      }
    }
  }

  .chart {
    width: 100%;
    height: 6.5vw;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: center;

    .panel {
      width: 100px;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: #f7fbff;
      flex-direction: column;

      .num {
        font-size: 25px;
        font-weight: 500;
        color: #303133;
      }
    }

    .target-echart2 {
      flex: 1;
      height: 100%;
    }
  }
</style>
