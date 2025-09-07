<script lang="ts">
  export default {
    name: 'AssetAccessInsights',
  }
</script>

<script setup lang="ts">
  import d3GraphForce2 from '@/components/d3-graph-force2.vue'
  import { getAssetVisitApi } from '@/api-ecs/retrieve'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import TraceabilityRelationsDatas from '@/ecs/alert/components/traceability-relations-data.vue'
  import TraceabilityServiceVisit from '@/ecs/alert/components/traceability-service-visit.vue'
  import TraceabilityOnlineDecoding from '@/ecs/alert/components/traceability-online-decoding.vue'
  import dayjs from 'dayjs'
  import { useUserStore } from '@/store/modules/user'
  import { timeDuratioOptions } from '@/data/constant'
  const { getTableColumns, setAllIndexType } = useUserStore()
  const $baseMessage: any = inject('$baseMessage')
  const formateTimeStr = 'YYYY-MM-DD HH:mm:ss'
  const chartLoading = ref(true)
  const route = useRoute()
  const graphRef = ref()
  const detailVisible = ref(false)
  const showInfoData = ref()

  const timeDuration = ref('30minute')
  const timeDate = ref()
  const searchIp = ref('')
  const queryData = reactive({
    endTime: '' as string,
    startTime: '' as string,
    searchSql: '',
    count: 100,
  })
  const searchChartsData = reactive<{
    nodes: any[]
    links: any[]
  }>({
    nodes: [],
    links: [],
  })
  const mergeSql = () => {
    let ipSql: string | string[] = searchIp.value
      .replace(/\s+/g, '')
      .replace('，', ',')
      .split(',')
      .filter(Boolean)
      .map((ip) => `源IP = ${ip} or 目的IP = ${ip}`)
      .join(' or ')
    return ipSql
  }
  // 获取资产访问数据
  const getAssetVisitsDate = async () => {
    chartLoading.value = true
    try {
      const {
        data: { links, nodes },
      } = await getAssetVisitApi({
        ...queryData,
        indexType: 11,
        isPathTrack: false,
        isTraceSource: false,
        isAssetVisit: true,
      })
      if (!links || links.length === 0) {
        searchChartsData.links = []
        searchChartsData.nodes = []
        chartLoading.value = false
        return $baseMessage('暂无其他数据', 'warning', 'vab-hey-message-warning')
      }
      searchChartsData.links = links.map((i) => ({ ...i, source: i.clientIp, target: i.serverIp }))
      searchChartsData.nodes = nodes.map((i) => ({ ...i, id: i.ip }))
      chartLoading.value = false
    } catch (error) {
      chartLoading.value = false
      searchChartsData.links = []
      searchChartsData.nodes = []
    }
  }
  const getChartData = () => {
    queryData.searchSql = mergeSql()
    getAssetVisitsDate()
  }
  function echartClicckHandler(params: any) {
    detailVisible.value = true
    showInfoData.value = params
  }
  watchEffect(() => {
    if (timeDuration.value === 'user-defined') return
    const curtimeDate = dayjs()
    const endDate = [
      '30minute',
      '1hours',
      '3hours',
      '24hours',
      'last-week',
      'last-two-week',
      'last-three-days',
      'last-one-month',
    ].includes(timeDuration.value)
      ? curtimeDate.add(10, 'minute').format(formateTimeStr)
      : curtimeDate.endOf('day').format(formateTimeStr)

    const cur = curtimeDate.startOf('day')
    let startDate = null
    switch (timeDuration.value) {
      case '30minute':
        startDate = curtimeDate.subtract(0.5, 'hour').format(formateTimeStr)
        break
      case '1hours':
        startDate = curtimeDate.subtract(1, 'hour').format(formateTimeStr)
        break
      case '3hours':
        startDate = curtimeDate.subtract(3, 'hour').format(formateTimeStr)
        break
      case '24hours':
        startDate = curtimeDate.subtract(24, 'hour').format(formateTimeStr)
        break
      case 'today':
        startDate = curtimeDate.startOf('day').format(formateTimeStr)
        break
      case 'yesterday':
        startDate = cur.subtract(1, 'day').format(formateTimeStr)
        break
      case 'last-three-days':
        startDate = curtimeDate.subtract(3, 'day').format(formateTimeStr)
        break
      case 'last-week':
        startDate = curtimeDate.subtract(1, 'week').format(formateTimeStr)
        break
      case 'last-two-week':
        startDate = curtimeDate.subtract(2, 'week').format(formateTimeStr)
        break
      case 'last-one-month':
        startDate = curtimeDate.subtract(1, 'month').format(formateTimeStr)
        break
      default:
        startDate = curtimeDate.startOf('day').format(formateTimeStr)
        break
    }
    queryData.endTime =
      timeDuration.value === 'yesterday' ? curtimeDate.endOf('day').subtract(1, 'day').format(formateTimeStr) : endDate
    queryData.startTime = startDate
  })
  watch(
    () => timeDuration.value,
    () => {
      if (timeDuration.value === 'user-defined') {
        timeDate.value = [dayjs(queryData.startTime), dayjs(queryData.endTime)]
      }
    }
  )
  watch(
    () => timeDate.value,
    (val) => {
      if (timeDate.value) {
        const [start, end] = val
        queryData.endTime = dayjs(end).format(formateTimeStr)
        queryData.startTime = dayjs(start).format(formateTimeStr)
      }
    }
  )
  onBeforeRouteUpdate((to) => {
    const { count, ip, startTime, endTime } = to.query
    queryData.count = Number(count) || 100
    searchIp.value = (ip as string) || ''
    if (startTime && endTime) {
      queryData.startTime = dayjs(new Date(+startTime)).format(formateTimeStr)
      queryData.endTime = dayjs(new Date(+endTime)).format(formateTimeStr)
      timeDuration.value = 'user-defined'
    }
    getChartData()
  })
  onMounted(() => {
    setAllIndexType()
    getTableColumns()
    const { count, ip, startTime, endTime } = route.query
    queryData.count = Number(count) || 100
    searchIp.value = (ip as string) || ''
    if (startTime && endTime) {
      queryData.startTime = dayjs(new Date(+startTime)).format(formateTimeStr)
      queryData.endTime = dayjs(new Date(+endTime)).format(formateTimeStr)
      timeDuration.value = 'user-defined'
    }
    getChartData()
  })
</script>

<template>
  <div v-loading="chartLoading" class="asset-access-insights">
    <div class="search-box">
      <div class="search-left">
        <el-select v-model="timeDuration" style="margin-right: 15px">
          <el-option
            v-for="item in [
              {
                value: '30minute',
                label: '半小时',
              },
              ...timeDuratioOptions,
            ]"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
        <vab-date-time-picker v-if="timeDuration === 'user-defined'" v-model="timeDate" />
      </div>
      <div class="search-right">
        <el-input
          v-model="searchIp"
          placeholder="请输入IP地址，支持掩码，多个条件用英文逗号分开"
          style="width: calc(100% - 80px); margin-right: 15px"
        />
        <el-button class="myBtn" type="primary" @click="getChartData">检索</el-button>
      </div>
    </div>
    <div class="btnn">
      <el-tooltip content="重置位置" effect="dark" placement="right">
        <vab-icon class="icon" icon="loader-3-line" @click="graphRef.reset()" />
      </el-tooltip>
    </div>
    <d3-graph-force2 ref="graphRef" :graph-data="searchChartsData" @click="echartClicckHandler" />
    <vab-dialog
      v-model="detailVisible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      title="深度挖掘"
      width="1375px"
    >
      <el-tabs class="demo-tabs" model-value="relations">
        <el-tab-pane label="关系" lazy name="relations">
          <traceability-relations-datas
            :node-data="showInfoData"
            :time-range="[queryData.startTime, queryData.endTime]"
          />
        </el-tab-pane>
        <el-tab-pane label="服务访问" lazy name="service-visit">
          <traceability-service-visit
            :node-data="showInfoData"
            :time-range="[queryData.startTime, queryData.endTime]"
          />
        </el-tab-pane>
        <el-tab-pane label="数据包分析" lazy name="online-decoding">
          <traceability-online-decoding
            :node-data="showInfoData"
            :time-range="[queryData.startTime, queryData.endTime]"
          />
        </el-tab-pane>
      </el-tabs>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  .asset-access-insights {
    width: 100%;
    height: 100%;
    background-color: #f3f9ff;
    overflow: hidden;
    --el-color-primary: #6954f0;
    .search-box {
      .search-left {
        position: absolute;
        top: 20px;
        left: 20px;
        z-index: 10;
        width: 600px;
        display: flex;
      }
      .search-right {
        position: absolute;
        top: 20px;
        right: 20px;
        z-index: 10;
        box-sizing: border-box;
        width: 600px;
        height: auto;
        padding: 20px;
        background: #fff;
      }
    }
    .btnn {
      width: 60px;
      height: 64px;
      position: absolute;
      bottom: 60px;
      right: 20px;
      z-index: 1;
      display: flex;
      flex-direction: column;
      cursor: pointer;
      background-color: #25223f;
      color: #fff;
      border-radius: 10px;
      .icon {
        font-size: 25px;
        width: 100%;
        height: 64px;
        line-height: 64px;
      }
    }

    :deep() {
      .myBtn,
      .el-button.el-button--primary:not(.is-plain) {
        background-color: var(--el-color-primary);
        border-color: var(--el-color-primary);
      }
    }
  }
</style>
