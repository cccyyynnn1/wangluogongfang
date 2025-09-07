<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import Traceability from '@/ecs/alert/components/traceability.vue'
  import AlertDetail from '@/ecs/alert/components/alert-detail.vue'
  import dayjs from 'dayjs'
  import { level_options, levelKey } from '@/ecs/alert/data/index'
  import { timeDuratioOptions as date_options } from '@/data/constant'
  import { AlertItem } from '@/types'
  import { getAlertApi, updateAlertStatusApi } from '@/api-ecs/alert'
  import { formatNstime, formatTime } from '@/utils/time'
  import numberFormatte from '@/utils/number'
  import { Algorithms_alerts } from '../type'
  import { injectStrict } from '@/utils/inject'
  import { useTableCopy } from '@/utils'
  const { isAdvanced, autoRefresh, autoRefreshDuration } = injectStrict(Algorithms_alerts)

  const multipleSelection = ref<AlertItem[]>([])
  const infoVal = ref()
  const $baseMessage: any = inject('$baseMessage')
  // 查询时间
  const timeDuration = ref('1hours')
  // 自定义时间
  const timeDate = ref()
  // 定时器ID
  const timerId = ref<NodeJS.Timer>()
  // 快速筛选折叠
  const loading = ref(false)
  const chartDetail = ref(false)
  const activeName = ref('overview')

  // 快速筛选折叠
  const fold = ref(true)
  const traceabilityVisible = ref(false)
  const alertDetailVisible = ref(false)
  const must_parames = reactive({
    /** 排序字段 */
    orderField: 'startTimeNs',
    /** 排序方式 */
    orderType: 'decs',
    /** 索引类型 */
    indexType: 10,
    /** 当前页码 */
    pageNum: 1,
    /** 当前页数 */
    pageSize: 10,
    /** 开始时间 */
    startTime: '',
    /** 结束时间 */
    endTime: '',
  })

  const alertInfo = reactive({
    list: [] as AlertItem[],
    total: 0,
    pieData: [] as { name: string; value: string | number }[],
    lineData: {
      date: [] as any[],
      count: [] as any[],
    },
    threat_level: [] as any[],
  })

  // 告警检索条件
  const queryForm = reactive({
    /** 查询SQL */
    searchSql: '',
    /** 攻击IP */
    attackIp: '',
    /** 源端口 */
    sourcePort: null,
    xff: '',
    url: '',
    host: '',
    /** 受害IP */
    victimIp: '',
    /** 目的端口 */
    targetPort: null,
    /** 威胁类型 */
    threatType: '',
    /** 威胁名称 */
    threatName: '',
    /** 威胁等级 */
    threatLevel: '',
    /** 攻击结果 */
    attackResult: '',
    /** 攻击次数 */
    attackCount: '',
    /** 告警设备IP */
    warnDeviceIp: '',
    /** 规则ID */
    ruleId: '',
    /** 攻击阶段 */
    attackStage: '',
    /** 状态 */
    readStatus: '',
  })

  // 柱状图配置
  const chartOption = reactive({
    xAxis: {
      type: 'category',
      data: [] as string[],
    },
    yAxis: {
      type: 'value',
      splitNumber: 5,
      axisLabel: {
        formatter: function (value: number) {
          return numberFormatte.format(+value)
        },
      },
    },
    series: {
      data: [] as number[],
      type: 'bar',
      color: 'rgb(63, 157, 248)',
      backgroundStyle: {
        color: 'rgba(220, 220, 220, 0.8)',
      },
    },
    grid: {
      top: 10,
      left: '2%',
      right: '2%',
      bottom: 10,
      height: 90,
      containLabel: true,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
    },
  })
  // 面积图配置
  const trendEchartOption = reactive({
    title: {
      text: '告警趋势图',
    },
    tooltip: {
      trigger: 'axis',
      extraCssText: 'z-index:1',
    },
    grid: {
      containLabel: true,
      height: 200,
    },
    xAxis: {
      type: 'category',
      data: [] as any[],
      boundaryGap: false,
    },
    yAxis: [
      {
        type: 'value',
      },
    ],
    series: {
      type: 'line',
      data: [] as any[],
      symbol: 'circle',
      smooth: true,
      yAxisIndex: 0,
      showSymbol: false,
      areaStyle: {
        opacity: 0.2,
      },
    },
  })

  // 饼图配置
  const branchEchartOption = reactive({
    title: {
      text: '告警统计图',
    },
    tooltip: {
      trigger: 'item',
    },
    legend: {
      bottom: '0',
      left: 'center',
    },
    grig: {
      top: 0,
      bottom: 10,
    },
    series: {
      type: 'pie',
      radius: ['35%', '60%'],
      center: ['50%', '35%'],
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2,
      },
      label: {
        show: false,
        position: 'center',
      },
      labelLine: {
        show: false,
      },
      emphasis: {
        label: {
          show: false,
        },
      },
      data: [] as any[],
    },
  })

  /**
   * type: 标记已读1 标记未读0
   */
  const updateStatus = async (status: '1' | '0') => {
    const ids = multipleSelection.value.map((i) => i.id)
    const { msg } = await updateAlertStatusApi({ ...must_parames, readStatus: status, ids })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    multipleSelection.value = []
  }

  /**
   * type: 标记忽略1 取消忽略0
   */
  const updateIgnoreStatus = async (status: '1' | '0') => {
    const ids = multipleSelection.value.map((i) => i.id)
    const { msg } = await updateAlertStatusApi({ ...must_parames, readStatus: status, ids })
    $baseMessage(msg, 'success', 'vab-hey-message-success')
    multipleSelection.value = []
  }

  // 查询告警日志
  async function queryData() {
    loading.value = true
    try {
      const { searchSql } = queryForm
      const queryData = isAdvanced.value
        ? { ...must_parames, searchSql }
        : { ...queryForm, searchSql: '', ...must_parames }
      const {
        data: { total, sumaryMap, resList },
      } = await getAlertApi(queryData)
      const { pillar, pie, threat_level } = JSON.parse(sumaryMap)
      alertInfo.list = resList
      alertInfo.total = total
      alertInfo.pieData = pie || []
      alertInfo.lineData = {
        date: pillar.date,
        count: pillar.count,
      }
      alertInfo.threat_level = threat_level || []
      loading.value = false
    } catch (error) {
      loading.value = false
      console.error(error)
    }
  }

  const obj = { 低危: 1, 中危: 2, 高危: 3, 危急: 4 }
  function getLevel(str: string) {
    // @ts-ignore
    const level = obj[str]
    return levelKey[level]
  }
  const getThreatLevel = (type: string) => {
    const obj: {
      [key: string]: string
    } = {
      高危: 'high',
      中危: 'mid',
      低危: 'low',
      危急: 'critical',
    }
    return obj[type] || 'default'
  }
  const handleDownloadExcel = () => {
    loading.value = true
    import('@/utils/excel').then((excel) => {
      const tHeader = ['告警时间', '受害IP', '攻击IP', 'XFF', '威胁类型', '威胁名称', '威胁等级', '攻击结果']
      const filterVal = [
        'startTimeNs',
        'victimIp',
        'attackIp',
        'xff',
        'threatType',
        'threatName',
        'threatLevel',
        'result',
      ]
      const list = multipleSelection.value.length ? multipleSelection.value : alertInfo.list
      const data = formatJson(filterVal, list)
      excel.export_json_to_excel({
        header: tHeader,
        data,
        filename: `告警列表-${formatTime(new Date().getTime())}`,
        autoWidth: true,
        bookType: 'xlsx',
      })
      loading.value = false
    })
  }
  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }
  function download() {
    console.log('download')
  }
  function formatExcelData(val: any, key: string) {
    switch (key) {
      case 'startTimeNs':
        return formatNstime(val[key])
      default:
        return val[key] || ''
    }
  }
  function showTraceability(row: any) {
    console.log(row, 'showTraceability_row')
    infoVal.value = row
    traceabilityVisible.value = true
  }
  function showAlertDetail(row: any) {
    infoVal.value = row
    alertDetailVisible.value = true
  }
  // 多选表格
  function handleSelectionChange(val: AlertItem[]) {
    multipleSelection.value = val
  }
  watch(
    () => alertInfo.lineData,
    () => {
      const { date, count } = alertInfo.lineData
      chartOption.xAxis.data = date
      chartOption.series.data = count
      trendEchartOption.xAxis.data = date.slice(-10)
      trendEchartOption.series.data = count.slice(-10)
    }
  )
  watch(
    () => alertInfo.pieData,
    () => {
      branchEchartOption.series.data = alertInfo.pieData
    }
  )
  watch(
    () => isAdvanced.value,
    () => {
      must_parames.pageNum = 1
      must_parames.pageSize = 10
    }
  )
  watchEffect(() => {
    if (timeDuration.value === 'user-defined') {
      formatUserTime()
      return
    }
    const timeDate = dayjs()
    const endDate = timeDate.format('YYYY-MM-DD HH:mm:ss')
    const cur = timeDate.startOf('day')
    let startDate = null
    switch (timeDuration.value) {
      case '1hours':
        startDate = timeDate.subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '3hours':
        startDate = timeDate.subtract(3, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '6hours':
        startDate = timeDate.subtract(6, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '24hours':
        startDate = timeDate.subtract(24, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'today':
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'yesterday':
        startDate = cur.subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-three-days':
        startDate = timeDate.subtract(3, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-week':
        startDate = timeDate.subtract(1, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-two-week':
        startDate = timeDate.subtract(2, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      default:
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
    }
    must_parames.endTime =
      timeDuration.value === 'yesterday'
        ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        : endDate
    must_parames.startTime = startDate
  })

  const formatUserTime = () => {
    const [startDate, endDate] = timeDate.value
    must_parames.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
    must_parames.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
  }
  watchEffect(() => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      must_parames.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
      must_parames.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
    }
  })
  watchEffect(() => {
    if (autoRefresh.value) {
      if (timerId.value) clearInterval(timerId.value)
      timerId.value = setInterval(() => queryData(), autoRefreshDuration.value * 1000)
    } else {
      clearInterval(timerId.value)
    }
  })
  onUnmounted(() => {
    if (timerId.value) clearInterval(timerId.value)
  })
  onMounted(() => {
    queryData()
  })
</script>

<script lang="ts">
  export default {
    name: 'AlgorithmsAlerts',
  }
</script>

<template>
  <!-- 查询表单 -->
  <el-form inline label-width="110px">
    <el-row>
      <el-col :span="6">
        <el-form-item label="时间">
          <el-select v-model="timeDuration" :style="{ width: timeDuration === 'user-defined' ? '90%' : '100%' }">
            <el-option v-for="item in date_options" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-col>
      <el-col v-if="timeDuration === 'user-defined'" :span="6">
        <el-form-item>
          <vab-date-time-picker v-model="timeDate" />
        </el-form-item>
      </el-col>
      <el-col v-if="isAdvanced" :span="isAdvanced && timeDuration === 'user-defined' ? 9 : 15">
        <el-form-item label="查询SQL" prop="searchSql">
          <el-input
            v-model="queryForm.searchSql"
            placeholder="src_ip = '192.168.1.1' AND phase in ( 'recon','control' )"
          />
        </el-form-item>
      </el-col>
      <template v-else>
        <el-col :span="6">
          <el-form-item label="攻击IP" prop="attackIp">
            <el-input v-model="queryForm.attackIp" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="攻击端口" prop="targetPort">
            <el-input v-model.number="queryForm.targetPort" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="受害IP" prop="victimIp">
            <el-input v-model="queryForm.victimIp" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="源端口" prop="sourcePort">
            <el-input v-model.number="queryForm.sourcePort" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item label="Host" prop="host">
            <el-input v-model="queryForm.host" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item label="威胁类型" prop="threatType">
            <!-- <el-select v-model="queryForm.threatType">
                <el-option label="攻击" value="攻击" />
                <el-option label="风险" value="2" />
              </el-select> -->
            <el-input v-model="queryForm.threatType" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item label="威胁名称" prop="threatName">
            <el-input v-model="queryForm.threatName" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item label="威胁等级" prop="threatLevel">
            <!-- <el-input v-model="queryForm.threatLevel" placeholder="进行模糊匹配" /> -->
            <el-select v-model="queryForm.threatLevel">
              <el-option v-for="item in level_options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item v-show="!fold" label="攻击结果" prop="attackResult">
            <el-input v-model="queryForm.attackResult" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item label="XFF" prop="xff">
            <el-input v-model="queryForm.xff" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item label="URL" prop="url">
            <el-input v-model="queryForm.url" placeholder="进行模糊匹配" />
          </el-form-item>
        </el-col>
        <el-col v-if="!fold" :span="6">
          <el-form-item label="已读状态" prop="readStatus">
            <el-select v-model="queryForm.readStatus">
              <el-option label="已读" value="1" />
              <el-option label="未读" value="2" />
            </el-select>
          </el-form-item>
        </el-col>
      </template>
      <el-col :span="isAdvanced ? 3 : timeDuration === 'user-defined' ? 12 : 18">
        <el-form-item class="submit_item" label=" ">
          <el-button type="primary" @click="queryData">查询</el-button>
          <template v-if="!isAdvanced">
            <div class="card-header-tag">
              <img
                class="form_icon"
                :class="{ upward: !fold }"
                src="@/assets/alert_images/down.png"
                @click="fold = !fold"
              />
            </div>
            <!-- <div class="card-header-tag">
                <img class="form_icon" src="@/assets/alert_images/time.png" />
              </div> -->
          </template>
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
  <!-- 告警日志图表 -->
  <vab-card class="access" shadow="hover" skeleton>
    <template #header>
      <p>
        共计 - {{ alertInfo.total }}条
        <span v-for="alert of alertInfo.threat_level" :key="alert.name" class="log">
          <span class="risk_status" :class="[getLevel(alert.name)?.class]">{{ getLevel(alert.name)?.label }}</span>
          <span>&nbsp;-&nbsp;{{ alert.value }}条</span>
        </span>
        <!-- <span class="log">
            <span class="risk_status height">高危</span>
            &nbsp;&nbsp;-&nbsp;&nbsp;16.3K条
          </span>
          <span class="log">
            <span class="risk_status mid">中危</span>
            &nbsp;&nbsp;-&nbsp;&nbsp;16.3K条
          </span>
          <span class="log">
            <span class="risk_status low">低危</span>
            &nbsp;&nbsp;-&nbsp;&nbsp;16.3K条
          </span> -->
      </p>
      <div class="card-header-tag">
        <img
          class="form_icon"
          :class="{ upward: chartDetail }"
          src="@/assets/alert_images/down.png"
          @click="chartDetail = !chartDetail"
        />
      </div>
    </template>

    <template v-if="chartDetail">
      <vab-chart class="trend-echart" :option="branchEchartOption" theme="vab-echarts-theme" />
      <vab-chart :option="trendEchartOption" theme="vab-echarts-theme" />
    </template>
    <template v-else>
      <vab-chart class="small" :option="chartOption" theme="vab-echarts-theme" />
    </template>
  </vab-card>
  <!-- 告警日志操作 -->
  <div class="table_action">
    <el-dropdown trigger="click">
      <el-button type="primary">
        更多操作
        <el-icon class="el-icon--right"><arrow-down /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item :disabled="multipleSelection.length === 0" @click="download">下载Pcap文件</el-dropdown-item>
          <el-dropdown-item @click="handleDownloadExcel">倒出告警日志</el-dropdown-item>
          <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateStatus('1')">
            标记已读
          </el-dropdown-item>
          <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateStatus('0')">
            标记未读
          </el-dropdown-item>
          <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateIgnoreStatus('1')">
            标记忽略
          </el-dropdown-item>
          <el-dropdown-item :disabled="multipleSelection.length === 0" @click="updateIgnoreStatus('0')">
            取消忽略
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
  <!-- 告警日志列表 -->
  <el-table
    align="center"
    border
    :data="alertInfo.list"
    style="width: 100%"
    @cell-contextmenu="useTableCopy"
    @selection-change="handleSelectionChange"
  >
    <el-table-column align="center" type="selection" width="55" />
    <el-table-column
      align="center"
      :formatter="({ startTimeNs }) => formatNstime(startTimeNs)"
      label="告警时间"
      prop="startTimeNs"
      width="180"
    />
    <el-table-column align="center" label="受害IP" prop="victimIp" width="180" />
    <el-table-column align="center" label="攻击IP" prop="attackIp" />
    <el-table-column align="center" label="XFF" prop="xff" />
    <el-table-column align="center" label="威胁类型" prop="threatType" />
    <el-table-column align="center" label="威胁名称" prop="threatName" />
    <el-table-column align="center" label="威胁等级" prop="threatLevel" width="100">
      <template #default="{ row }">
        <span :class="['alert_tag', getThreatLevel(row.threatLevel)]">{{ row.threatLevel }}</span>
      </template>
    </el-table-column>
    <el-table-column align="center" label="攻击结果" prop="result" width="90" />
    <el-table-column align="center" fixed="right" label="操作" width="150">
      <template #default="{ row }">
        <el-button class="row_action" size="small" @click="showAlertDetail(row)">详情</el-button>
        <el-button class="row_action" size="small" @click="showTraceability(row)">溯源</el-button>
        <!-- <el-button class="row_action" link size="small" type="primary">
                                                                                                                        加白
                                                                                                             </el-button> -->
      </template>
    </el-table-column>
  </el-table>
  <el-pagination
    v-model:current-page="must_parames.pageNum"
    v-model:page-size="must_parames.pageSize"
    background
    layout="total, sizes, prev, pager, next, jumper"
    :page-sizes="[10, 20, 30]"
    :total="alertInfo.total"
    @current-change="queryData"
    @size-change="queryData"
  />

  <!-- 告警详情详情 -->
  <alert-detail v-model:alert-detail-visible="alertDetailVisible" :next="false" :prev="false" :select-alert="infoVal" />
  <traceability v-model:traceabilityVisible="traceabilityVisible" :select-alert="infoVal" />
</template>

<style scoped lang="scss">
  $criticalColor: #ff0202;
  $lowColor: #1b81fe;
  $midColor: #f1b04d;
  $highColor: #fa6d15;
  $defaultColor: #909399;
  .alert_tag {
    display: inline-block;
    width: 55px;
    padding: 2px 3px;
    border-radius: 5px;
    color: #fff;
    &.critical {
      background-color: rgba($criticalColor, 0.8);
    }
    &.low {
      background-color: rgba($lowColor, 0.8);
    }
    &.mid {
      background-color: rgba($midColor, 0.8);
    }
    &.high {
      background-color: rgba($highColor, 0.8);
    }
    &.default {
      background-color: rgba($defaultColor, 0.8);
      color: inherit;
    }
  }
  :deep() {
    .el-card__header {
      padding: 20px 20 14px;

      p {
        margin: 0;
      }
    }

    .el-card__body {
      padding-top: 0;
      display: flex;

      .echarts {
        flex: 1;
        width: auto;
        height: 275px;
        &.small {
          height: 100px;
        }
      }
    }
    .el-form-item {
      width: 100%;
      .el-form-item__content {
        .el-input,
        .el-select {
          width: 100%;
        }
      }
    }
    .submit_item {
      .el-form-item__content {
        justify-content: flex-end;
      }
    }
  }

  .table_action {
    margin-bottom: 20px;
    display: flex;
    justify-content: flex-end;

    .action {
      margin-left: 30px;
      display: flex;
      align-items: center;
    }

    .el-button {
      margin-left: 0;
    }

    .el-dropdown {
      float: right;
    }
  }
  .log {
    margin: 0 10px;
    color: #6c6d6e;
    font-size: 14px;
  }

  .card-header-tag {
    width: 32px;
    height: 32px;
    border: 1px solid #e4e7ed;
    border-radius: 2px;
    display: flex;
    align-items: center;
    top: 10px !important;
    margin-left: 10px;
    cursor: pointer;
  }

  .form_icon {
    width: 20px;
    height: 20px;
    margin: 0 auto;
    transition: transform 0.2s;

    &.upward {
      transform: rotate(180deg);
    }
  }

  .chart_detail {
    width: 50%;
    display: inline-block;

    h5 {
      margin: 10px 0 10px;
      font-weight: 500;
      color: #303133;
    }
  }
</style>
