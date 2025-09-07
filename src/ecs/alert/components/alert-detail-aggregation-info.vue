<script lang="ts">
  export default {
    name: 'AlertAggregationDetailInfo',
  }
  type InfoType = 'basic-info' | 'raw-data' | 'decoding-stream' | 'decoding-packet'
</script>

<script setup lang="ts">
  import { AlertItem } from '~/src/types'
  import BasicInfo from './basic-info.vue'
  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'
  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'
  import JsonPreview from '@/components/json-preview.vue'
  import Whitelist from './white-list/index.vue'
  import { PacketDecodeQuery } from '~/src/types'
  import { TableColumnItemType } from '/#/store'
  import { formatNstime, formatTime } from '@/utils/time'
  import { downloadLogPacket } from '@/utils/download'
  import AlertAiAnalyse from './alert-ai-analyse.vue'

  import dayjs from 'dayjs'
  const router = useRouter()
  const props = defineProps<{
    alertInfo: AlertItem
    infoIndex: number
    tableColumn: TableColumnItemType[]
    detailTotalIndex: number
    moduleEnable: boolean
  }>()
  const emits = defineEmits<{
    (e: 'next', row: AlertItem, index: number): void
    (e: 'prev', row: AlertItem, index: number): void
  }>()
  const aiRef = ref<HTMLDivElement>()
  const infoRef = ref<HTMLDivElement>()
  const collapse = ref(false)
  const infoType = ref<InfoType>('basic-info')
  const decodingTime = ref('')
  const decodingData = ref<PacketDecodeQuery>()

  const whitelistShow = ref(false)
  const isNeedToFormet = ref(props.alertInfo.url ? true : false)
  provide('isNeedToFormet', isNeedToFormet)

  const decodingHandle = () => {
    const {
      attackIp: clientIp,
      sourcePort: clientPort,
      victimIp: serverIp,
      targetPort: serverPort,
      startTimeNs,
    } = props.alertInfo as any
    const timeDate = dayjs(formatNstime(startTimeNs))
    const startDate = timeDate.subtract(1, 'minute')
    const endDate = timeDate.add(1, 'minute')
    decodingTime.value = `${formatTime(startDate)} - ${formatTime(endDate)}`
    decodingData.value = {
      clientIp,
      clientPort,
      serverIp,
      serverPort,
    }
  }

  const formatSourceData = (val: string | undefined) => {
    if (!val) return ''
    let _val = ''
    try {
      _val = JSON.parse(val)
    } catch (error) {
      _val = val
    }
    return _val
  }
  async function download() {
    try {
      const {
        attackIp: clientIp,
        sourcePort: clientPort,
        victimIp: serverIp,
        targetPort: serverPort,
        startTimeNs,
      } = props.alertInfo as any
      downloadLogPacket({ clientIp, clientPort, serverIp, serverPort, requestTimeNs: startTimeNs }, 2)
    } catch (error) {
      console.error(error)
    }
  }
  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }
  function formatExcelData(val: any, key: string) {
    switch (key) {
      case 'startTimeNs':
      case 'warnTime':
        return formatNstime(val[key])
      default:
        return val[key] || ''
    }
  }
  const handleDownloadExcel = () => {
    const tHeader: string[] = []
    const filterVal: any = []
    try {
      props.tableColumn.forEach(({ fieldNameCn, fieldNameEn }) => {
        tHeader.push(fieldNameCn)
        filterVal.push(fieldNameEn)
      })
      import('@/utils/excel').then((excel) => {
        const data = formatJson(filterVal, [props.alertInfo])
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `告警日志-${formatNstime(+props.alertInfo?.startTimeNs)}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
      })
    } catch (error) {
      console.log(error)
    }
  }
  const addWhiteList = async () => {
    whitelistShow.value = true
  }
  const fullFlowSurveyHandle = async () => {
    const { clientIp, serverIp, sourcePort, targetPort, startTimeNs } = props.alertInfo as any
    const time = startTimeNs / 1000000
    const sqlStr = `源ip = "${clientIp}" and 源端口 = "${sourcePort}" and 目的ip = "${serverIp}" and 目的端口 = "${targetPort}"`
    const indexType = 1
    const start = formatTime(dayjs(time).subtract(15, 'minute'))
    const end = formatTime(dayjs(time).add(15, 'minute'))
    const resolveRouter = router.resolve({
      path: '/retrieve/index',
      query: {
        info: encodeURIComponent(
          JSON.stringify({
            sql: sqlStr,
            indexType,
            timeRanges: [start, end],
          })
        ),
      },
    })
    window.open(resolveRouter.href, '_blank')
  }
  const handleClick = (event: MouseEvent) => {
    aiRef.value?.classList.toggle('active')
    infoRef.value?.classList.toggle('active')
    collapse.value = !collapse.value
  }
  onMounted(() => {
    decodingHandle()
  })
</script>

<template>
  <div class="alert-aggregation-detail-info">
    <div ref="aiRef" class="alert-aggregation-detail-info-ai">
      <alert-ai-analyse :alarm-data="alertInfo" :collapse="collapse" :module-enable="moduleEnable" />
      <div v-if="moduleEnable" class="collapse" @click="handleClick">
        <el-icon><ArrowDownBold /></el-icon>
      </div>
    </div>
    <div ref="infoRef" class="alert-aggregation-detail-info-content active">
      <div class="alert-aggregation-detail-info-header">
        <el-button-group>
          <el-button :type="infoType === 'basic-info' ? 'primary' : 'default'" @click="() => (infoType = 'basic-info')">
            基本信息
          </el-button>
          <el-button :type="infoType === 'raw-data' ? 'primary' : 'default'" @click="() => (infoType = 'raw-data')">
            原始数据
          </el-button>
          <el-button
            :type="infoType === 'decoding-stream' ? 'primary' : 'default'"
            @click="() => (infoType = 'decoding-stream')"
          >
            数据流
          </el-button>
          <el-button
            :type="infoType === 'decoding-packet' ? 'primary' : 'default'"
            @click="() => (infoType = 'decoding-packet')"
          >
            数据包
          </el-button>
        </el-button-group>
        <div>
          <el-button @click="fullFlowSurveyHandle">全流量调查</el-button>
          <el-button plain @click="addWhiteList">加入白名单</el-button>
          <el-button plain @click="download">PCAP保存</el-button>
          <el-button plain @click="handleDownloadExcel">下载告警详情</el-button>
          <el-button :disabled="infoIndex === 0" plain @click="() => emits('prev', props.alertInfo, props.infoIndex)">
            上一条
          </el-button>
          <el-button
            :disabled="detailTotalIndex === infoIndex"
            plain
            @click="() => emits('next', props.alertInfo, props.infoIndex)"
          >
            下一条
          </el-button>
        </div>
      </div>
      <basic-info v-if="infoType === 'basic-info'" :info-val="alertInfo" type="alert-aggregation" />
      <decoding-stream
        v-else-if="infoType === 'decoding-stream'"
        :current-data="decodingData"
        :time-date="decodingTime"
      />
      <decoding-packet
        v-else-if="infoType === 'decoding-packet'"
        :current-data="decodingData"
        :time-date="decodingTime"
      />
      <div v-else class="mask">
        <json-preview :json-value="JSON.stringify(formatSourceData(props.alertInfo.sourceData), null, 4)" />
      </div>
      <Whitelist
        v-if="whitelistShow"
        :cur-data="alertInfo"
        :is-show="whitelistShow"
        @on-closeEvent="whitelistShow = false"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
  .alert-aggregation-detail-info {
    padding: 0 30px 25px;
    height: 790px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    .alert-aggregation-detail-info-ai,
    .alert-aggregation-detail-info-content {
      &.active {
        flex: 4;
        .collapse .el-icon {
          transform: rotateX(180deg);
        }
      }
      transition: all 0.3s ease-out;
      flex: 1;
      overflow: auto;
    }
    .alert-aggregation-detail-info-ai {
      position: relative;
      border-radius: 20px;
      border: 1px solid #f1f0ff;
      overflow: hidden;
      min-height: 190px;
      .collapse {
        position: absolute;
        bottom: 27px;
        left: 50%;
        margin-left: -18px;
        width: 36px;
        height: 14px;
        background: #cac1ff;
        border-radius: 2px 2px 0px 0px;
        text-align: center;
        color: #ffffff;
        cursor: pointer;
        .el-icon {
          vertical-align: text-top;
        }
      }
    }
    .alert-aggregation-detail-info-content {
      border-radius: 20px;
      border: 1px solid #f1f0ff;
      padding: 0 20px 24px;
      margin-top: -28px;
      background-color: #fff;
      z-index: 9;
      .alert-aggregation-detail-info-header {
        display: flex;
        justify-content: space-between;
        background: #fff;
        padding: 20px 0 15px;
        background: #fff;
        position: sticky;
        top: 0;
        z-index: 999999999;
      }
    }
    .mask {
      height: calc(100% - 60px);
      overflow-y: auto;
      border: 1px solid rgb(230 231 240);
      background-color: #f8f7ff;
      border-radius: 8px;
    }
  }
</style>
