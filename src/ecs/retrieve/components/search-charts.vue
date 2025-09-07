<script lang="ts">
  export default {
    name: 'SearchCharts',
  }
</script>

<script setup lang="ts">
  import { timeDuratioOptions } from '@/data/constant'
  import SearchSql from '@/components/search-sql/index.vue'
  import d3GraphForce2 from '@/components/d3-graph-force2.vue'
  import { getAttackPerspectiveApi } from '@/api-ecs/alert'
  // import { getServiceLinksApi } from '@/api-ecs/assets'
  import { getAssetVisitApi, getPathTracingApi } from '@/api-ecs/retrieve'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import TraceabilityRelationsDatas from '@/ecs/alert/components/traceability-relations-data.vue'
  import TraceabilityServiceVisit from '@/ecs/alert/components/traceability-service-visit.vue'
  import TraceabilityOnlineDecoding from '@/ecs/alert/components/traceability-online-decoding.vue'
  import AttackAlerts from '../../index/components/attack-alert.vue'
  import dayjs from 'dayjs'
  const $baseMessage: any = inject('$baseMessage')

  const props = defineProps<{
    chartType: string
  }>()
  const bool = inject('isFullscreen') ? inject('isFullscreen') : false
  const isFullscreen = ref(bool)
  const graphRef = ref()
  const activeName = ref(props.chartType === 'attack-perspective' ? 'alert' : 'relations')
  const detailVisible = ref(false)

  const showInfoData = reactive({
    data: {} as { id: string; x: number; y: number; ip: string },
  })

  const formateTimeStr = 'YYYY-MM-DD HH:mm:ss'

  const _timeDuratioOptions = [...timeDuratioOptions]
  _timeDuratioOptions.splice(-1, 0, { value: 'last-one-month', label: '最近一月' })
  const timeOption =
    props.chartType === 'attack-perspective'
      ? timeDuratioOptions.filter((i) => !i.value.includes('hours') && i.value !== 'user-defined')
      : _timeDuratioOptions
  const tableColumn: { fieldNameCn: string; fieldNameEn: string }[] = [
    { fieldNameCn: '源IP', fieldNameEn: 'clientIp' },
    { fieldNameCn: '源端口', fieldNameEn: 'clientPort' },
    { fieldNameCn: '目的IP', fieldNameEn: 'serverIp' },
    { fieldNameCn: '目的端口', fieldNameEn: 'serverPort' },
  ]
  // 搜索dom
  const searchBoxHidden = ref(true)
  const showMoreSearch = ref(false)
  const chartLoading = ref(true)

  const md5Data: any = inject('md5Data')
  const sqlComponentsRef = ref<InstanceType<typeof SearchSql>>()
  // 查询时间
  const timeDuration = ref(props.chartType === 'attack-perspective' ? 'today' : '1hours')
  // 自定义时间
  const timeDate = ref()

  const queryData = reactive({
    endDate: '',
    startDate: '',
    searchSql: '',
    ip: '',
  })
  //表单数据
  const queryForm = reactive({
    assetsIp: '',
    filter: '',
    searchSql: '',
    attackLimit: 100,
    topCount: 200,
  })

  const searchChartsData = reactive<{
    nodes: any[]
    links: any[]
  }>({
    nodes: [],
    links: [],
  })
  const mergeSql = () => {
    let ipSql: string | string[] = queryForm.assetsIp
      .replace(/\s+/g, '')
      .replace('，', ',')
      .split(',')
      .filter(Boolean)
      .map((ip) => `源IP = ${ip} or 目的IP = ${ip}`)
    ipSql = ipSql.join(' or ')
    return [ipSql, queryForm.filter.trim()]
      .filter(Boolean)
      .map((i) => `( ${i} )`)
      .join(' and ')
  }
  // 获取资产访问数据
  const getAssetVisitsDate = async (query: any) => {
    try {
      const {
        data: { links, nodes },
      } = await getAssetVisitApi({
        indexType: 11,
        startTime: queryData.startDate,
        endTime: queryData.endDate,
        count: 100,
        searchSql: query.searchSql || null,
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
  // 获取攻击透视数据
  const getAttackPerspectiveDate = async (query: { ip: string; startDate: string; endDate: string; limit: number }) => {
    try {
      const {
        data: { links, nodes },
      } = await getAttackPerspectiveApi(query)
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
      searchChartsData.links = []
      searchChartsData.nodes = []
      chartLoading.value = false
    }
  }
  // 获取路径追踪数据
  const getPathTracingData = async (query: {
    searchSql: string
    indexType: number
    startTime: string
    endTime: string
    aggregationFields: string[]
    count: number
    isPathTrack: boolean
    isTraceSource: boolean
    isAssetVisit: boolean
  }) => {
    try {
      const {
        data: { links, nodes },
      } = await getPathTracingApi(query)
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
      searchChartsData.links = []
      searchChartsData.nodes = []
      chartLoading.value = false
    }
  }
  const getChartData = () => {
    if (sqlComponentsRef.value && !sqlComponentsRef.value?.getValidateSql()) return
    chartLoading.value = true
    if (props.chartType === 'asset-visits') queryData.searchSql = mergeSql()
    if (props.chartType === 'path-tracing') queryData.searchSql = queryForm.searchSql.trim()
    const { endDate, startDate, ip, searchSql } = queryData
    if (searchSql) {
      queryData.searchSql =
        props.chartType === 'path-tracing' ? searchSql : `源IP = ${searchSql} or 目的IP = ${searchSql}`
    } else {
      if (ip) {
        queryData.searchSql = `源IP = ${ip} or 目的IP = ${ip}`
      } else {
        queryData.searchSql = ''
      }
    }
    if (!sqlComponentsRef.value?.isValidate && props.chartType === 'path-tracing') {
      chartLoading.value = false
      return
    }
    const { attackLimit, topCount } = queryForm
    const FetchHandle: {
      [key: string]: () => void
    } = {
      'asset-visits': () => getAssetVisitsDate({ searchSql }),
      'attack-perspective': () =>
        getAttackPerspectiveDate({
          endDate,
          startDate,
          ip,
          limit: ip ? (attackLimit === 10000 ? 0 : attackLimit) : 20,
        }),
      'path-tracing': () =>
        getPathTracingData({
          searchSql,
          indexType: md5Data?.value ? 23 : 1,
          count: topCount,
          aggregationFields: ['clientIp', 'serverIp', 'serverPort'],
          startTime: startDate,
          endTime: endDate,
          isPathTrack: true,
          isTraceSource: false,
          isAssetVisit: false,
        }),
    }
    if (props.chartType === 'path-tracing') {
      sqlComponentsRef.value?.changeHistories(searchSql)
    }
    FetchHandle[props.chartType]()
  }
  function echartClicckHandler(params: any) {
    detailVisible.value = true
    searchBoxHidden.value = true
    showInfoData.data = params
  }
  const searchhandle = (val: string) => {
    getChartData()
  }
  watchEffect(() => {
    const curtimeDate = dayjs()
    const endDate = [
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
    queryData.endDate =
      timeDuration.value === 'yesterday' ? curtimeDate.endOf('day').subtract(1, 'day').format(formateTimeStr) : endDate
    queryData.startDate = startDate
  })
  watchEffect(() => {
    if (queryForm.attackLimit > 10000) {
      return (queryForm.attackLimit = 10000)
    }
    // if (queryForm.attackLimit < 1) {
    //   return (queryForm.attackLimit = 0)
    // }
  })

  watchEffect(() => {
    if (timeDuration.value === 'user-defined') {
      timeDate.value = [dayjs(queryData.startDate), dayjs(queryData.endDate)]
    }
  })
  watchEffect(() => {
    if (timeDate.value) {
      const [start, end] = timeDate.value
      queryData.endDate = dayjs(end).format(formateTimeStr)
      queryData.startDate = dayjs(start).format(formateTimeStr)
    }
  })
  watch(
    () => md5Data?.value,
    () => {
      if (md5Data?.value) {
        queryForm.searchSql = `用户代理 like "Tuhuan-${md5Data.value}"`
        queryData.startDate = '2023-01-21 00:00:00'
        getChartData()
      } else {
        getChartData()
      }
    }
  )

  onMounted(() => {
    getChartData()
  })
</script>

<template>
  <div v-loading="chartLoading" class="charts-content" :class="{ 'vab-fullscreen': isFullscreen }">
    <div class="search-box" :class="{ isHidden: searchBoxHidden }" @click="searchBoxHidden = false">
      <div class="search-left">
        <el-select v-if="!md5Data" v-model="timeDuration" style="margin-right: 15px">
          <el-option v-for="item in timeOption" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <vab-date-time-picker v-if="timeDuration === 'user-defined'" v-model="timeDate" />
      </div>
      <div class="search-right">
        <template v-if="chartType === 'asset-visits'">
          <el-input
            v-model="queryForm.assetsIp"
            placeholder="请输入IP地址，支持掩码，多个条件用英文逗号分开"
            style="width: calc(100% - 127px); margin-right: 15px"
          />
          <el-button type="primary" @click="getChartData">检索</el-button>
          <div class="card-header-tag">
            <img
              class="form_icon"
              :class="{ upward: showMoreSearch }"
              src="@/assets/alert_images/down.png"
              @click="showMoreSearch = !showMoreSearch"
            />
          </div>
          <div v-if="showMoreSearch" class="search-filter">
            <p>过滤条件：</p>
            <el-input v-model="queryForm.filter" :autosize="{ minRows: 4, maxRows: 4 }" resize="none" type="textarea" />
          </div>
        </template>
        <template v-else-if="chartType === 'attack-perspective'">
          <el-input
            v-model="queryData.ip"
            placeholder="请输入IP，多个IP使用英文逗号分割"
            style="width: calc(100% - 127px); margin-right: 15px"
          />
          <el-button type="primary" @click="getChartData">检索</el-button>
          <!-- 折叠按钮 -->
          <div class="card-header-tag">
            <img
              class="form_icon"
              :class="{ upward: !showMoreSearch }"
              src="@/assets/alert_images/down.png"
              @click="showMoreSearch = !showMoreSearch"
            />
          </div>
          <div v-if="showMoreSearch" class="search-filter">
            <p>查询条数：</p>
            <el-input v-model="queryForm.attackLimit" :min="0" placeholder="最大查询条数：1000" />
          </div>
        </template>
        <template v-else-if="chartType === 'path-tracing'">
          <!-- <el-input v-model="queryData.searchSql" placeholder="请输入需要追踪的特征条件" /> -->
          <SearchSql
            ref="sqlComponentsRef"
            :index-type="md5Data ? 23 : 1"
            :model-value="queryForm.searchSql"
            placeholder="请输入需要追踪的特征条件"
            style="display: inline-block; width: calc(100% - 127px); margin-right: 15px"
            @on-change="(str) => (queryForm.searchSql = str)"
            @onSearch="searchhandle"
          />
          <el-button type="primary" @click="getChartData">检索</el-button>
          <!-- 折叠按钮 -->
          <div class="card-header-tag">
            <img
              class="form_icon"
              :class="{ upward: !showMoreSearch }"
              src="@/assets/alert_images/down.png"
              @click="showMoreSearch = !showMoreSearch"
            />
          </div>
          <div v-if="showMoreSearch" class="search-filter">
            <p>查询条数：</p>
            <el-input v-model.number="queryForm.topCount" placeholder="最大查询条数：1000" />
          </div>
        </template>
      </div>
    </div>
    <div class="btnn">
      <el-tooltip content="重置位置" effect="dark" placement="right">
        <vab-icon class="icon" icon="loader-3-line" @click="graphRef.reset()" />
      </el-tooltip>
      <el-tooltip :content="isFullscreen ? '关闭全屏' : '全屏'" effect="dark" placement="right">
        <vab-icon
          class="icon fullscreen"
          :icon="isFullscreen ? 'fullscreen-exit-fill' : 'fullscreen-fill'"
          @click="isFullscreen = !isFullscreen"
        />
      </el-tooltip>
    </div>
    <d3-graph-force2
      ref="graphRef"
      :graph-data="searchChartsData"
      @click="echartClicckHandler"
      @graph-move="() => (searchBoxHidden = true)"
    />

    <vab-dialog
      v-model="detailVisible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      title="深度挖掘"
      width="1375px"
      @close="activeName = chartType === 'attack-perspective' ? 'alert' : 'relations'"
    >
      <el-tabs v-model="activeName" class="demo-tabs">
        <el-tab-pane v-if="chartType === 'attack-perspective'" label="告警" lazy name="alert">
          <attack-alerts :node-data="showInfoData.data" :time-range="[queryData.startDate, queryData.endDate]" />
        </el-tab-pane>
        <el-tab-pane label="关系" lazy name="relations">
          <traceability-relations-datas
            :node-data="showInfoData.data"
            :search-sql="chartType === 'path-tracing' ? queryData.searchSql : ''"
            :time-range="[queryData.startDate, queryData.endDate]"
          />
        </el-tab-pane>
        <el-tab-pane label="服务访问" lazy name="service-visit">
          <traceability-service-visit
            :node-data="showInfoData.data"
            :time-range="[queryData.startDate, queryData.endDate]"
          />
        </el-tab-pane>
        <el-tab-pane label="数据包分析" lazy name="online-decoding">
          <traceability-online-decoding
            :node-data="showInfoData.data"
            :time-range="[queryData.startDate, queryData.endDate]"
          />
        </el-tab-pane>
      </el-tabs>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  .my-form {
    display: flex;
    width: 100%;
    .my-input {
      flex: 1;
      :deep() {
        .el-row {
          width: 100%;
        }
        .el-form-item {
          width: 100%;
        }
      }
    }
  }
  .charts-content {
    position: relative;
    width: 100%;
    height: calc(100vh - 50px);
    background-color: #f3f9ff;
    .btnn {
      width: 60px;
      height: 108px;
      position: absolute;
      bottom: 88px;
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
        height: 54px;
        line-height: 64px;
        &.fullscreen {
          line-height: 44px;
        }
      }
    }
    .search-box {
      &.isHidden {
        opacity: 1;
      }
      .search-left {
        position: absolute;
        top: 20px;
        left: 20px;
        z-index: 10;
        width: 600px;
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
        .search-filter {
          width: 100%;
          margin-top: 15px;
          overflow: hidden;
        }
        .card-header-tag {
          display: inline-flex;
          align-items: center;
          width: 32px;
          height: 32px;
          margin-left: 15px;
          vertical-align: middle;
          cursor: pointer;
          border: 1px solid #e4e7ed;
          border-radius: 2px;
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
      }
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
</style>
