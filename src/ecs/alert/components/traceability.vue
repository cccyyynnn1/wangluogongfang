<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  // import TraceabilityRelations from './traceability-relations.vue'
  import TraceabilityRelations2 from './traceability-relations2.vue'
  import TraceabilityDatas from './traceability-data.vue'
  import TraceabilityServiceVisit from './traceability-service-visit.vue'
  import TraceabilityOnlineDecoding from './traceability-online-decoding.vue'
  import { chartImg } from '@/data/constant'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import { getServiceLinksApi, updateServiceLinksApi } from '@/api-ecs/assets'
  import { ServiceLink, ServiceNode } from '@/types/index'
  import { levelKey } from '../data/index'
  import { TableColumnItemType } from '/#/store'
  import { getAllDisPlaysFiledApi } from '~/src/api-ecs/public'
  import dayjs from 'dayjs'
  import { useUserStore } from '@/store/modules/user'
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    traceabilityVisible: boolean
    selectAlert: any
  }>()
  const alertIndexType = 10
  const userStore = useUserStore()
  const { getTableColumn } = userStore
  const tableColumn = ref<TableColumnItemType[]>([])
  // 自定义时间
  const dialogTableVisible = ref(false)
  const traceability_chart = ref<InstanceType<typeof VabChart>>()
  // 挖掘Dialog
  const diggingDialog = ref(false)
  const chartLoading = ref(false)
  const chartNodes = ref<ServiceNode[]>([])
  const chartLinks = ref<ServiceLink[]>([])
  const diggingNode = ref<ServiceNode>()

  const activeName = ref('relations-force')

  const emits = defineEmits<{
    (e: 'update:traceabilityVisible', traceabilityVisible: boolean): void
  }>()

  const chartsOption = reactive({
    tooltip: {
      formatter: (params: any, ticket: string) => {
        const { dataType, data } = params
        if (dataType === 'node')
          return `<div style="white-space:pre-line;max-width:500px">${data?.IpLocation || data?.ip}</div>`
        const { alarm, serverPort } = data
        return `<div>
               ${Object.entries(serverPort)
                 .map(
                   ([key, val]: any) =>
                     `<p style="height:14px;margin:0 0 10px 0">${key}：${val.protocolStr}&nbsp;&nbsp;${val.appProtocol}</p>`
                 )
                 .join('')}
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

    layout: 'none',
    series: [
      {
        type: 'graph',
        layout: 'none',
        roam: true,
        draggable: true,
        symbolSize: 80,
        label: {
          show: true,
          padding: [3, 10],
          height: 22,
          lineHeight: 22,
          backgroundColor: 'rgba(236, 244, 255, 1)',
          borderColor: '#4E7CBE',
          borderWidth: 1,
          position: 'insideBottom',
          borderRadius: 2,
          distance: -10,
          color: '#4E7CBE',
          formatter: ({ data }: any) => data?.name || data?.ip,
        },
        edgeLabel: {
          // 曲线文字标记
          show: false,
        },
        lineStyle: {
          curveness: 0.18,
        },
        nodeScaleRatio: 0.1,
        symbol: chartImg,
        scaleLimit: {
          min: 1,
        },
        edgeSymbol: ['none', 'arrow'],
        emphasis: {
          focus: 'adjacency',
          lineStyle: {
            width: 2,
          },
        },
        data: [] as any[],
        links: [] as any[],
      },
    ],
  })

  const queryData = reactive({
    page: 1,
    limit: 200,
    // lastTime: '',
    startTime: '',
    endTime: '',
    total: 0,
  })
  function formatDataTime(row: any, key: string) {
    const time = row[key] / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
  }
  async function formatColum() {
    const columnData = getTableColumn(alertIndexType)
    const { data: userColumnData } = await getAllDisPlaysFiledApi()
    const alertUserColumnData = userColumnData[10] || []
    if (alertUserColumnData?.length > 0) {
      tableColumn.value = alertUserColumnData.map((key: number) =>
        columnData.find((item) => item.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
    }
  }
  function changeVisible() {
    emits('update:traceabilityVisible', false)
  }
  const formatNode = (nodes: ServiceNode[]) => {
    const newNodes = nodes.map((node, index: number) => {
      const { x, y } = node
      return { ...node, id: node.ip, x: x || 50 * (index + 1), y: y || 100 }
    })
    chartNodes.value = newNodes
    return newNodes
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
    chartLinks.value = newlinks
    return newlinks
  }
  // 提交图表信息
  const submitChartsData = async (data?: { link: ServiceLink[]; node: ServiceNode[] }, done?: () => void) => {
    if (chartNodes.value.length === 0) return done && done()
    if (traceability_chart.value && traceability_chart.value.chart) {
      const { id } = props.selectAlert
      const { _itemLayouts } = traceability_chart.value.chart.getModel().getSeriesByIndex(0).getData()
      chartNodes.value.forEach((item, index) => {
        const [x, y] = _itemLayouts[index]
        item.x = x
        item.y = y
      })
      if (data) {
        chartLinks.value = [...chartLinks.value, ...data.link]
        chartNodes.value = [
          ...chartNodes.value,
          ...data.node.map((i) => {
            const { x, y } = i
            return {
              ...i,
              x: x ? x : (Math.random() + 30) * 100,
              y: y ? y : (Math.random() + 40) * 100,
            }
          }),
        ]
      }

      chartsOption.series[0].data = [...chartNodes.value]
      chartsOption.series[0].links = [...chartLinks.value]
      await updateServiceLinksApi([
        {
          tid: id,
          content: JSON.stringify({
            links: chartLinks.value,
            nodes: chartNodes.value,
          }),
        },
      ])
      diggingDialog.value = false
      if (done) done()
    }
  }
  const dialogCloseHandle = (done: () => void) => {
    submitChartsData(undefined, done)
  }
  const getServiceLinks = async () => {
    if (!props.selectAlert) return
    chartLoading.value = true
    const { endTime, startTime } = queryData
    const { clientIp, serverIp, id } = props.selectAlert
    const str = clientIp ? `源IP = ${clientIp}` : ''
    const str1 = clientIp && serverIp ? ' and ' : ''
    const str2 = serverIp ? `目的IP = ${serverIp}` : ''
    const searchSql = str + str1 + str2
    try {
      const {
        data: { links, nodes },
      } = await getServiceLinksApi({
        indexType: 11,
        // tid: id,
        tid: '7226171090309714293',
        startTime: startTime,
        searchSql: searchSql,
        endTime: endTime,
        count: 100,
      })
      if (!nodes || nodes.length === 0) {
        chartLoading.value = false
        chartsOption.series[0].data = formatNode([])
        chartsOption.series[0].links = formatLinks([])
        return $baseMessage('暂无其他数据', 'warning', 'vab-hey-message-warning')
      }
      chartsOption.series[0].data = formatNode(nodes)
      chartsOption.series[0].links = formatLinks(links)
      chartLoading.value = false
    } catch (error) {
      console.error(error)
    }
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
  const traceabilityClickHandle = (params: any) => {
    const { dataType, data } = params
    diggingNode.value = data
    diggingDialog.value = true
  }
  watchEffect(() => {
    if (props.selectAlert) {
      const { startTimeNs } = props.selectAlert
      const timeDate = dayjs(formatNstime(startTimeNs))
      // @ts-ignore
      queryData.startTime = formatTime(timeDate.subtract(5, 'minute'))
      // @ts-ignore
      queryData.endTime = formatTime(timeDate.add(5, 'minute'))
      // queryData.lastTime = `${formatTime(startDate)} - ${formatTime(endDate)}`
      if (props.traceabilityVisible) getServiceLinks()
    }
  })
  watch(
    () => props.traceabilityVisible,
    () => {
      dialogTableVisible.value = props.traceabilityVisible
      if (props.traceabilityVisible) formatColum()
    },
    { immediate: true }
  )
</script>

<script lang="ts">
  export default {
    name: 'Traceability',
  }
</script>

<template>
  <div class="traceability">
    <el-dialog
      v-model="dialogTableVisible"
      :before-close="dialogCloseHandle"
      :close-on-click-modal="false"
      :fullscreen="true"
      title="溯源"
      width="1155px"
      @close="changeVisible"
    >
      <el-table :border="true" class="alert-traceability-table" :data="[selectAlert]" style="width: 100%">
        <el-table-column
          v-for="item in tableColumn"
          :key="item.id"
          align="center"
          :label="item.fieldNameCn"
          :prop="item.fieldNameEn"
          :show-overflow-tooltip="{
            placement: 'left-start',
          }"
          :width="['threatName', 'warnTime', 'startTimeNs'].includes(item.fieldNameEn) ? 250 : 120"
        >
          <template
            #header
            v-if="
              item.fieldNameEn == 'startTimeNs' ||
              item.fieldNameEn == 'attackIp' ||
              item.fieldNameEn == 'xff' ||
              item.fieldNameEn == 'victimIp' ||
              item.fieldNameEn == 'threatType' ||
              item.fieldNameEn == 'threatName' ||
              item.fieldNameEn == 'threatLevel' ||
              item.fieldNameEn == 'attackResult'
            "
          >
            <span style="cursor: pointer">
              {{ item.fieldNameCn }}
            </span>
          </template>
          <template #default="{ row }" v-if="['startTimeNs', 'warnTime'].includes(item.fieldNameEn)">
            {{ formatDataTime(row, item.fieldNameEn) }}
          </template>
          <template #default="{ row }" v-else-if="item.fieldNameEn == 'threatLevel'">
            <span :class="['alert_tag', getThreatLevel(row.threatLevel)]">
              <el-icon><WarnTriangleFilled /></el-icon>
              {{ row.threatLevel }}
            </span>
          </template>
        </el-table-column>
      </el-table>
      <div v-loading="chartLoading" class="hostChart">
        <vab-chart
          v-if="dialogTableVisible"
          ref="traceability_chart"
          class="traceability-echart"
          :click="traceabilityClickHandle"
          :option="chartsOption"
          theme="vab-echarts-theme"
        />
      </div>
      <vab-dialog
        v-model="diggingDialog"
        :close-on-click-modal="false"
        destroy-on-close
        title="深度挖掘"
        width="1375px"
        @close="activeName = 'relations-force'"
      >
        <el-tabs v-model="activeName" class="demo-tabs">
          <!-- <el-tab-pane label="关系" name="relations">
            <traceability-relations
              :alert-data="props.selectAlert"
              :chart-node-data="chartNodes"
              :info-data="diggingNode"
              @putChart="submitChartsData"
            />
          </el-tab-pane> -->
          <el-tab-pane label="关系(引力图)" lazy name="relations-force">
            <traceability-relations2
              :alert-data="props.selectAlert"
              :chart-node-data="chartNodes"
              :info-data="diggingNode"
              @putChart="submitChartsData"
            />
          </el-tab-pane>
          <el-tab-pane label="数据" lazy name="data">
            <traceability-datas :alert-data="props.selectAlert" :node-data="diggingNode" />
          </el-tab-pane>
          <el-tab-pane label="服务访问" lazy name="service-visit">
            <traceability-service-visit :alert-data="props.selectAlert" :node-data="diggingNode" />
          </el-tab-pane>
          <el-tab-pane label="数据包分析" lazy name="online-decoding">
            <traceability-online-decoding :alert-data="props.selectAlert" :node-data="diggingNode" />
          </el-tab-pane>
        </el-tabs>
      </vab-dialog>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  $criticalColor: #ca0a08;
  $lowColor: #ffbe36;
  $midColor: #ff7212;
  $highColor: #ff2927;
  $defaultColor: #909399;
  .alert_tag {
    display: inline-block;
    width: 55px;
    padding: 2px 3px;
    border-radius: 5px;
    color: #fff;
    font-size: 14px;
    .el-icon {
      margin-right: -4px;
      vertical-align: -3px;
      font-size: 17px;
    }
    &.critical {
      background-color: rgba($criticalColor, 1);
    }
    &.low {
      background-color: rgba($lowColor, 1);
    }
    &.mid {
      background-color: rgba($midColor, 1);
    }
    &.high {
      background-color: rgba($highColor, 1);
    }
    &.default {
      background-color: rgba($defaultColor, 1);
      color: inherit;
    }
  }

  .hostChart {
    width: 100%;
    height: calc(100vh - 210px);
    border-radius: 4px;
    background: #f3f9ff;
    margin-top: 20px;
  }

  .hostTable {
    display: flex;

    .hostTable_item {
      flex: 1;
      overflow-x: auto;

      .hostTable_title {
        height: 14px;
        line-height: 12px;
        position: relative;
        text-indent: 0.8em;
        font-weight: 700;
        color: #303133;

        &::before {
          content: ' ';
          display: inline-block;
          position: absolute;
          width: 3px;
          height: 100%;
          left: 0;
          background: #0d88fe;
          margin-right: 5px;
        }
      }

      .hostTable_search {
        display: flex;
        justify-content: space-between;
        margin: 16px 0;

        .el-input {
          width: 88%;
        }
      }

      .alertLevel {
        margin-bottom: 20px;
      }
    }

    .host {
      width: 114px;
      height: 48px;
      margin: 20% 20px;
      background: #ecf4ff;
      border-radius: 4px;
      border: 1px solid #abcbff;
      line-height: 48px;
      text-align: center;
    }
  }

  .traceability-echart {
    width: 100%;
  }

  :deep() {
    .echarts {
      height: 100%;
    }
    .el-popper {
      max-width: 60% !important;
    }
    .alert-traceability-table {
      .el-table__body-wrapper {
        max-height: 50px;
        min-height: 50px;
      }
    }
  }

  .risk_status,
  .el-table__cell .el-tag {
    color: $defaultColor;

    &.critical {
      color: $criticalColor;
    }

    &.low {
      color: $lowColor;
    }

    &.mid {
      color: $midColor;
    }

    &.high {
      color: $highColor;
    }
  }

  .el-table__cell .el-tag {
    border: 0;
    background-color: rgba($criticalColor, 0.3);

    &.critical {
      background-color: rgba($criticalColor, 0.3);
    }

    &.low {
      background-color: rgba($lowColor, 0.3);
    }

    &.mid {
      background-color: rgba($midColor, 0.3);
    }

    &.high {
      background-color: rgba($highColor, 0.3);
    }
  }
</style>
