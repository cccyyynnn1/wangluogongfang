<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import { getServiceLinksApi } from '@/api-ecs/assets'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import { ServiceLink, ServiceNode } from '@/types/index'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import TraceabilityRelationsDatas from './traceability-relations-data.vue'
  import TraceabilityServiceVisit from './traceability-service-visit.vue'
  import dayjs from 'dayjs'
  import TraceabilityOnlineDecoding from './traceability-online-decoding.vue'
  const props = defineProps<{
    infoData: any
    alertData: any
    chartNodeData: any[]
  }>()

  const emits = defineEmits<{
    (e: 'putChart', data: { link: any; node: any }): void
  }>()
  const activeName = ref('relations')
  const { IpLocation: propsIpLocation, ip: propsIp, type: propsType, snat: clientIp, serverIp } = toRaw(props.infoData)
  // const propsIp = '10.99.19.33'
  const $baseMessage: any = inject('$baseMessage')
  const detailVisible = ref(false)
  const chartNodes = ref<ServiceNode[]>([])
  const chartLinks = ref<ServiceLink[]>([])
  const chartsOption = reactive({
    tooltip: {
      formatter: (params: any, ticket: string) => {
        const { dataType, data } = params
        if (dataType === 'node')
          return `<div style="white-space:pre-line;max-width:500px">${data?.IpLocation || data?.ip}</div>`
        const { alarm, serverPort } = data
        return `<div>
               ${Object.entries(serverPort).map(
                 ([key, val]: any) =>
                   `<p style="height:14px;margin:0 0 10px 0">${key}：${val.protocolStr}&nbsp;&nbsp;${val.appProtocol}</p>`
               )}
              ${alarm.length ? `<p style="height:14px;margin:0 0 10px 0">告警类型：${alarm}</p>` : ''}
            </div>`
      },
      backgroundColor: '#409EFF',
      textStyle: {
        color: '#fff',
      },
    },
    xAxis: {
      type: 'value',
      show: false,
    },
    yAxis: {
      type: 'value',
      show: false,
    },
    grid: { containLabel: false },
    layout: 'none',
    series: [
      {
        type: 'graph',
        layout: 'none',
        roam: true,
        draggable: true,
        label: {
          show: true,
          padding: [3, 10],
          height: 22,
          lineHeight: 22,
          backgroundColor: 'rgba(236, 244, 255, 1)',
          borderColor: '#4E7CBE',
          borderWidth: 1,
          position: 'inside',
          borderRadius: 2,
          distance: 20,
          color: '#4E7CBE',
          formatter: ({ data }: any) => data.name,
        },
        edgeLabel: {
          // 曲线文字标记
          show: false,
        },
        lineStyle: {
          color: 'rgba(92, 135, 187, 1)',
          curveness: 0.18,
        },
        nodeScaleRatio: 0.1,
        symbolSize: 0,
        scaleLimit: {
          min: 1,
        },
        emphasis: {
          scale: 1.4,
          focus: 'adjacency',
          disabled: false,
          lineStyle: {
            width: 2,
          },
        },
        data: [] as any[],
        links: [] as any[],
      },
    ],
  })
  const show = ref(false)
  const chartsLoading = ref(false)
  const showInfoData = reactive({
    data: {} as { id: string; x: number; y: number; ip: string },
    x: 0,
    y: 0,
  })
  // 自定义时间
  const timeDate = ref()
  const showTraceabilityBtn = ref(false)
  const queryData = reactive({
    page: 1,
    limit: 200,
    lastTime: '',
    ipAddr: propsType ? clientIp : propsIp,
    total: 0,
  })
  const formatNode = (nodes: ServiceNode[], links: ServiceLink[]) => {
    const top = 100
    const left = 200
    if (nodes.length === 0) {
      return [
        {
          id: propsIp || clientIp,
          ip: propsIp || clientIp,
          IpLocation: propsIpLocation || clientIp,
          x: left,
          y: top + 170,
        },
      ]
    }

    const targetNodekeys: string[] = []
    const sourceNodekeys: string[] = []
    const myTargetNodes: any[] = []
    const mySourceNodes: any[] = []
    let myCurNode = null
    const filterIp = propsType ? clientIp : propsIp
    for (let index = 0; index < links.length; index++) {
      const { clientIp, serverIp } = links[index]
      if (clientIp === filterIp) {
        targetNodekeys.push(serverIp)
        continue
      }
      if (serverIp === filterIp) {
        sourceNodekeys.push(clientIp)
        continue
      }
    }
    nodes.forEach((item) => {
      if (targetNodekeys.includes(item.ip)) {
        return myTargetNodes.push({ ...item, id: item.ip, x: left + 1820, y: (myTargetNodes.length + 1) * top })
      }
      if (sourceNodekeys.includes(item.ip)) {
        return mySourceNodes.push({
          ...item,
          id: item.ip,
          x: left - 1080,
          y: (mySourceNodes.length + 1) * top + 200,
        })
      }
      return (myCurNode = { ...item, id: item.ip, x: left, y: top + 170 })
    })
    chartNodes.value = [myCurNode, ...myTargetNodes, ...mySourceNodes]
    return chartNodes.value
  }
  const formatLinks = (links: ServiceLink[]) => {
    const newlinks = links.map((link, index: number) => {
      const hasErr = link.alarm.length > 0
      return {
        ...link,
        source: link.clientIp,
        target: link.serverIp,
        lineStyle: { color: hasErr ? '#F56C6C' : '#4E7CBE' },
        tooltip: { borderWidth: 0 },
      }
    })
    chartLinks.value = [...chartLinks.value, ...newlinks]
    return newlinks
  }
  async function getServiceAccessCarding() {
    // chartsLoading.value = true
    // const { lastTime, ipAddr } = queryData
    // const query = propsType ? { serverIp: serverIp, clientIp: clientIp, lastTime } : { ipAddr, lastTime }
    // try {
    //   const {
    //     data: { links, nodes },
    //   } = await getServiceLinksApi({
    //     query: JSON.stringify(query),
    //     limit: 100,
    //   })
    //   if (!nodes || nodes.length === 0) {
    //     chartsLoading.value = false
    //     chartsOption.series[0].data = []
    //     chartsOption.series[0].links = []
    //     return $baseMessage('暂无其他数据', 'warning', 'vab-hey-message-warning')
    //   }
    //   chartsOption.series[0].data = formatNode(nodes, links)
    //   chartsOption.series[0].links = formatLinks(links)
    // } catch (error) {
    //   console.error(error)
    //   chartsLoading.value = false
    // }
    // chartsLoading.value = false
  }

  function zrclick(event: any) {
    if (!event.target) {
      show.value = false
    }
  }

  function echartClicckHandler(params: any) {
    show.value = false
    detailVisible.value = true
    showInfoData.data = params.data
  }
  function echartContextmenuHandler(params: any) {
    const { event, data } = params
    event.event.preventDefault()
    if (data.ip === props.infoData.ip) return
    show.value = true
    showInfoData.data = data
    showInfoData.x = event.event.clientX - 50
    showInfoData.y = event.event.clientY - 30
  }
  // clientIp :左 ｜｜ serverIp：右。
  async function tipClick() {
    // chartsLoading.value = true
    const { data } = showInfoData
    const link = chartLinks.value.filter((item) => item.clientIp === data.ip || item.serverIp === data.ip)
    const chartHasNode = props.chartNodeData.some((i) => i.ip === data.ip)
    emits('putChart', { link: chartHasNode ? [] : link, node: chartHasNode ? [] : [data] })
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
      queryData.lastTime = `${formatTime(new Date(startDate).getTime())} - ${formatTime(new Date(endDate).getTime())}`
      getServiceAccessCarding()
    }
  })
</script>

<script lang="ts">
  export default {
    name: 'TraceabilityRelations',
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
      <vab-chart
        ref="serviceVisits_chart"
        class="charts"
        :click="echartClicckHandler"
        :contextmenu="echartContextmenuHandler"
        :option="chartsOption"
        theme="vab-echarts-theme"
        :zrclick="zrclick"
      />
    </div>
    <div v-if="show" class="chartsTip" :style="{ transform: `translate(${showInfoData.x}px,${showInfoData.y}px)` }">
      <el-button-group>
        <el-button size="small" type="primary" @click="tipClick">添加到溯源图</el-button>
      </el-button-group>
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
        <el-tab-pane label="数据" lazy name="relations">
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
</style>
