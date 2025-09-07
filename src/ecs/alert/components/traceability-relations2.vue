<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { getServiceLinksApi } from '@/api-ecs/assets'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import { ServiceLink, ServiceNode } from '@/types/index'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import TraceabilityRelationsDatas from './traceability-relations-data.vue'
  import TraceabilityServiceVisit from './traceability-service-visit.vue'
  import dayjs from 'dayjs'
  import d3GraphForce2 from '@/components/d3-graph-force2.vue'
  import TraceabilityOnlineDecoding from './traceability-online-decoding.vue'
  const props = defineProps<{
    infoData: any
    alertData: any
    chartNodeData: any[]
  }>()

  const emits = defineEmits<{
    (e: 'putChart', data: { link: any; node: any }): void
  }>()
  const { IpLocation: propsIpLocation, ip: propsIp, type: propsType, snat: clientIp, serverIp } = toRaw(props.infoData)
  // const propsIp = '10.99.19.33'
  const activeName = ref('relations')
  const $baseMessage: any = inject('$baseMessage')
  const detailVisible = ref(false)
  const chartNodes = ref<ServiceNode[]>([])
  const chartLinks = ref<ServiceLink[]>([])
  const graphData = reactive<{ nodes: ServiceNode[]; links: ServiceLink[] }>({
    nodes: [],
    links: [],
  })
  const contextMenu = [
    {
      title: '添加到溯源图',
      action: (elm: any, d: any) => {
        if (props.infoData.ip === d.ip) return
        const { IpLocation, coordinate, ip, id } = d
        tipClick({ IpLocation, coordinate, ip, id })
      },
    },
  ]
  const show = ref(false)
  const chartsLoading = ref(false)
  const showInfoData = reactive({
    data: {} as { id: string; x: number; y: number; ip: string },
    x: 0,
    y: 0,
  })
  // 自定义时间
  const timeDate = ref()
  const queryData = reactive({
    page: 1,
    limit: 200,
    // lastTime: '',
    startTime: '',
    endTime: '',
    ipAddr: propsType ? clientIp : propsIp,
    total: 0,
  })
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
    const { startTime, endTime, ipAddr } = queryData
    const { id } = props.infoData
    const query = propsType ? { serverIp: serverIp, clientIp: clientIp } : { ipAddr }
    const str = clientIp ? `源IP = ${clientIp}` : ''
    const str1 = clientIp && serverIp ? ' and ' : ''
    const str2 = serverIp ? `目的IP = ${serverIp}` : ''
    const searchSql = propsType ? str + str1 + str2 : `源IP = ${ipAddr} or 目的IP = ${ipAddr}`
    try {
      const {
        data: { links, nodes },
      } = await getServiceLinksApi({
        indexType: 11,
        // tid: id,
        startTime: startTime,
        searchSql: searchSql,
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
      chartsLoading.value = false
    }
    chartsLoading.value = false
  }

  function zrclick(event: any) {
    if (!event.target) {
      show.value = false
    }
  }
  function echartMouseUpHandler(params: any) {
    if (chartNodes.value.length < 5) return
    const { event, dataIndex } = params
    chartNodes.value[dataIndex].x = event.offsetX
    chartNodes.value[dataIndex].y = event.offsetY

    chartNodes.value[dataIndex].fixed = true
  }
  function echartClicckHandler(params: any) {
    show.value = false
    detailVisible.value = true
    showInfoData.data = params
  }
  async function tipClick(nade_data: any) {
    const link = graphData.links.filter((item) => item.clientIp === nade_data.ip || item.serverIp === nade_data.ip)
    const chartHasNode = props.chartNodeData.some((i) => i.ip === nade_data.ip)
    emits('putChart', { link: chartHasNode ? [] : link, node: chartHasNode ? [] : [nade_data] })
  }
  watchEffect(() => {
    if (props.alertData) {
      const { startTimeNs } = props.alertData
      const newDate = dayjs(formatNstime(startTimeNs))
      const startDate = newDate.subtract(5, 'minute')
      const endDate = newDate.add(5, 'minute')
      timeDate.value = [formatTime(startDate), formatTime(endDate)]
    }
  })
  watchEffect(() => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      queryData.endTime = `${formatTime(new Date(endDate).getTime())}`
      queryData.startTime = `${formatTime(new Date(startDate).getTime())}`
      getServiceAccessCarding()
    }
  })
</script>

<script lang="ts">
  export default {
    name: 'TraceabilityRelations2',
  }
</script>

<template>
  <div v-loading="chartsLoading">
    <el-row>
      <el-col :offset="10" :span="14" style="text-align: right">
        <vab-date-time-picker v-model="timeDate" />
      </el-col>
    </el-row>
    <div class="serviceVisitsChart">
      <d3-graph-force2 :context-menu="contextMenu" :graph-data="graphData" @click="echartClicckHandler" />
    </div>
    <vab-dialog
      v-model="detailVisible"
      align-center
      :close-on-click-modal="false"
      destroy-on-close
      title="深度挖掘"
      width="1175px"
      @close="activeName = 'relations'"
    >
      <el-tabs v-model="activeName" class="demo-tabs">
        <el-tab-pane label="关系" lazy name="relations">
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
  </div>
</template>

<style scoped lang="scss">
  .serviceVisitsChart {
    width: 100%;
    height: 648px;
    background: #f3f9ff;
    margin-top: 20px;

    .charts {
      width: 100%;
      height: 100%;
    }
  }

  .chartsTip {
    position: fixed;
    top: 0;
    left: 0;
  }

  .serviceVisits-IpQuery {
    display: flex;
    margin-top: 24px;

    .serviceVisits-source,
    .serviceVisits-target {
      flex: 1;
      overflow-x: auto;

      .el-form {
        width: 100%;
        display: flex;

        .el-form-item {
          flex: 1;
          margin-right: 20px;

          &:last-child {
            width: 65px;
            flex: none;
            margin-right: 0;
          }
        }
      }

      .el-row {
        margin-right: 0 !important;
      }
    }

    .serviceVisits-host {
      width: 114px;
      height: 68px;
      background: #ecf4ff;
      border-radius: 4px;
      border: 1px solid #abcbff;
      flex: none;
      margin: 10% 20px;
      color: #4e7cbe;
      text-align: center;
      display: flex;
      line-height: 22px;
      font-size: 14px;
      align-items: center;
      word-break: keep-all;
      word-wrap: break-word;
    }
  }
  :deep() {
    .el-select {
      vertical-align: bottom;
    }
    .el-date-editor {
      margin-left: 10px;
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
