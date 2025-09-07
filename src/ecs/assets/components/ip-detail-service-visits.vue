<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { getServiceLinksApi } from '@/api-ecs/assets'
  import { AssetsIp_Detail } from '../type'
  import { injectStrict } from '@/utils/inject'
  import { formatTime, formatNstime } from '~/src/utils/time'
  import { ServiceLink, ServiceNode } from '@/types/index'
  import { timeDuratioOptions } from '@/data/constant'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import TraceabilityRelationsDatas from '@/ecs/alert/components/traceability-relations-data.vue'
  import TraceabilityServiceVisit from '@/ecs/alert/components/traceability-service-visit.vue'
  import TraceabilityOnlineDecoding from '@/ecs/alert/components/traceability-online-decoding.vue'
  import dayjs from 'dayjs'
  import d3GraphForce2 from '@/components/d3-graph-force2.vue'
  const { currentRow } = injectStrict(AssetsIp_Detail)
  const rowData = toRaw(currentRow?.value) as any
  // rowData.assetsIp = '10.99.19.33'
  const $baseMessage: any = inject('$baseMessage')

  const graphData = reactive<{ nodes: ServiceNode[]; links: ServiceLink[] }>({
    nodes: [],
    links: [],
  })

  const chartNodes = ref<ServiceNode[]>([])
  const chartLinks = ref<ServiceLink[]>([])
  const activeName = ref('relations')

  const show = ref(false)
  const chartsLoading = ref(false)
  const showInfoData = reactive({
    data: {} as { id: string; x: number; y: number; ip: string },
    x: 0,
    y: 0,
  })
  // 查询时间
  const timeDuration = ref('1hours')
  const detailVisible = ref(false)
  // 自定义时间
  const timeDate = ref()
  const gridData = reactive({
    list: [],
    targetList: [],
    sourceList: [],
  })
  const queryData = reactive({
    page: 1,
    limit: 200,
    startTime: '',
    endTime: '',
    lastTime: '',
    ipAddr: rowData.assetsIp,
    total: 0,
  })

  const contextMenu = [
    {
      title: '向前挖掘',
      action: (elm: any, d: any) => {
        showInfoData.data = d
        console.log(d, '向前挖掘')
        tipClick('clientIp')
      },
    },
    {
      title: '向后挖掘',
      action: (elm: any, d: any) => {
        console.log(d, '向后挖掘')
        showInfoData.data = d
        tipClick('serverIp')
      },
    },
  ]

  const formatNode = (nodes: ServiceNode[]) => {
    return nodes.map((i) => ({ ...i, id: i.ip }))
  }
  const formatLinks = (links: ServiceLink[]) => {
    const newlinks = links.map((link) => {
      const hasErr = link.alarm.length > 0
      return {
        ...link,
        source: link.clientIp,
        target: link.serverIp,
        lineStyle: { color: hasErr ? '#F56C6C' : '#4E7CBE' },
        tooltip: { borderWidth: 0 },
      }
    })
    return newlinks
  }
  async function getServiceAccessCarding() {
    chartsLoading.value = true
    const { startTime, ipAddr, endTime } = queryData
    try {
      const {
        data: { links, nodes },
      } = await getServiceLinksApi({
        // query: JSON.stringify({ searchSql: `源IP = ${ipAddr} or 目的IP = ${ipAddr}`, lastTime }),
        // limit: 100,
        indexType: 11,
        // tid: id,
        startTime: startTime,
        searchSql: `源IP = ${ipAddr} or 目的IP = ${ipAddr}`,
        endTime: endTime,
        count: 100,
      })
      if (!nodes || nodes.length === 0) {
        chartsLoading.value = false
        graphData.links = []
        graphData.nodes = []
        return $baseMessage('暂无其他数据', 'warning', 'vab-hey-message-warning')
      }
      graphData.links = formatLinks(links)
      graphData.nodes = formatNode(nodes)
    } catch (error) {
      console.error(error)
    }
    chartsLoading.value = false
  }
  function echartClicckHandler(params: any) {
    show.value = false
    detailVisible.value = true
    showInfoData.data = params
  }
  // clientIp :左 ｜｜ serverIp：右。
  async function tipClick(query: 'clientIp' | 'serverIp') {
    chartsLoading.value = true
    const { x: eventX, y: eventY, ip: eventIp, id: eventId } = showInfoData.data
    const { startTime, endTime } = queryData
    const queryLinkkey = query === 'serverIp' ? 'clientIp' : 'serverIp'
    const sameKye = chartLinks.value.filter((item) => item[queryLinkkey] === eventIp).map((i) => i[query])
    const sameNodeKye = chartNodes.value.map((i) => i?.ip)
    try {
      const {
        data: { links, nodes },
      } = await getServiceLinksApi({
        // query: JSON.stringify({ ipAddr: eventIp, lastTime }),
        // limit: 1,
        indexType: 11,
        // tid: id,
        startTime: startTime,
        searchSql: `源IP = ${eventIp} or 目的IP = ${eventIp}`,
        endTime: endTime,
        count: 100,
      })
      const myLinks = (links || [])
        .filter((item) => item[queryLinkkey] === eventIp && !sameKye.includes(item[query]))
        .map((link) => {
          const hasErr = link.alarm.length > 0
          return {
            ...link,
            source: link.clientIp,
            target: link.serverIp,
            lineStyle: { color: hasErr ? '#F56C6C' : '#4E7CBE' },
            tooltip: { borderWidth: 0 },
          }
        })
      const nodeKeys = myLinks.map((item) => item[query])
      const _x = query === 'clientIp' ? -500 : 500
      const _y = 100
      if (nodeKeys.length === 0) {
        $baseMessage('暂无其他数据', 'warning', 'vab-hey-message-warning')
        chartsLoading.value = false
        return
      }
      const myNodes = (nodes || [])
        .filter((node) => nodeKeys.includes(node.ip) && !sameNodeKye.includes(node.ip))
        .map((node, index) => ({
          ...node,
          id: node.ip,
          symbolSize: 0,
          x: eventX + _x,
          y: eventY + _y + index * 100,
        }))
      chartNodes.value = [...chartNodes.value, ...myNodes]
      chartLinks.value = [...chartLinks.value, ...myLinks]
    } catch (error) {
      console.error(error)
    }
    show.value = false
    chartsLoading.value = false
  }
  watchEffect(() => {
    const sources = [] as any
    const targets = [] as any
    for (let i = gridData.list.length; i--; ) {
      const { clientIp, serverIp } = gridData.list[i]
      if (clientIp === rowData.assetsIp) {
        sources.push(gridData.list[i])
        continue
      }
      if (serverIp === rowData.assetsIp) {
        targets.push(gridData.list[i])
        continue
      }
    }

    gridData.sourceList = sources
    gridData.targetList = targets
  })
  watchEffect(() => {
    if (timeDate.value && timeDuration.value === 'user-defined') {
      const [startDate, endDate] = timeDate.value
      queryData.lastTime = `${startDate} - ${endDate}`
      queryData.startTime = `${startDate}`
      queryData.endTime = `${endDate}`
      getServiceAccessCarding()
    }
  })
  watch(
    () => timeDuration.value,
    () => {
      const _timeDate = dayjs()
      const cur = _timeDate.startOf('day')
      const endDate =
        timeDuration.value === 'yesterday'
          ? _timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
          : _timeDate.format('YYYY-MM-DD HH:mm:ss')
      let startDate = null
      switch (timeDuration.value) {
        case '1hours':
          startDate = _timeDate.subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss')
          break
        case '3hours':
          startDate = _timeDate.subtract(3, 'hour').format('YYYY-MM-DD HH:mm:ss')
          break
        case '24hours':
          startDate = _timeDate.subtract(24, 'hour').format('YYYY-MM-DD HH:mm:ss')
          break
        case 'today':
          startDate = _timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
          break
        case 'yesterday':
          startDate = cur.subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
          break
        case 'last-three-days':
          startDate = _timeDate.subtract(3, 'day').format('YYYY-MM-DD HH:mm:ss')
          break
        case 'last-week':
          startDate = _timeDate.subtract(1, 'week').format('YYYY-MM-DD HH:mm:ss')
          break
        case 'last-two-week':
          startDate = _timeDate.subtract(2, 'week').format('YYYY-MM-DD HH:mm:ss')
          break
        default:
          return
      }
      timeDate.value = [startDate, endDate]
      queryData.lastTime = `${startDate} - ${endDate}`
      queryData.startTime = startDate
      queryData.endTime = endDate
      getServiceAccessCarding()
    },
    { immediate: true }
  )
</script>

<script lang="ts">
  export default {
    name: 'ServiceVisits',
  }
</script>

<template>
  <el-row>
    <el-col :offset="10" :span="14" style="text-align: right">
      <el-select v-model="timeDuration" placeholder="Select" style="margin-right: 20px">
        <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <vab-date-time-picker v-if="timeDuration === 'user-defined'" v-model="timeDate" />
    </el-col>
  </el-row>
  <div v-loading="chartsLoading" class="serviceVisitsChart">
    <d3-graph-force2 :context-menu="[]" :graph-data="graphData" @click="echartClicckHandler" />
  </div>
  <vab-dialog
    v-model="detailVisible"
    align-center
    :close-on-click-modal="false"
    destroy-on-close
    title="深度挖掘"
    width="1375px"
    @close="activeName = 'relations'"
  >
    <el-tabs v-model="activeName" class="demo-tabs">
      <el-tab-pane label="关系" name="relations">
        <traceability-relations-datas :node-data="showInfoData.data" :time-range="timeDate" />
      </el-tab-pane>
      <el-tab-pane label="服务访问" lazy name="service-visit">
        <traceability-service-visit :node-data="showInfoData.data" :time-range="timeDate" />
      </el-tab-pane>
      <el-tab-pane label="数据包分析" lazy name="online-decoding">
        <traceability-online-decoding :node-data="showInfoData.data" :time-range="timeDate" />
      </el-tab-pane>
    </el-tabs>
  </vab-dialog>
</template>

<style scoped lang="scss">
  .serviceVisitsChart {
    width: 100%;
    height: calc(100vh - 222px);
    background: #f3f9ff;
    margin-top: 20px;
  }
  :deep() {
    .el-select {
      vertical-align: bottom;
    }
    .el-date-editor {
      margin-left: 10px;
    }
    .demo-tabs {
      min-height: 600px;
    }
  }
</style>
