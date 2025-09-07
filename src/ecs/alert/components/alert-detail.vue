<script setup lang="ts">
  import Download from '~/library/components/VabColumnBar/download.vue'
  import ToolBox from '~/library/components/VabColumnBar/toolbox.vue'
  import BasicInfo from './basic-info.vue'
  import { PacketDecodeQuery } from '~/src/types'
  import dayjs from 'dayjs'
  import { formatNstime, formatTime } from '@/utils/time'
  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'
  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'
  import { useUserStore } from '@/store/modules/user'
  import { getAllDisPlaysFiledApi } from '@/api-ecs/public'
  import { TableColumnItemType } from '/#/store'
  import { downloadLogPacket } from '~/src/utils/download'
  import { useCopy } from '@/utils'
  import JsonPreview from '@/components/json-preview.vue'
  import Whitelist from './white-list/index.vue'
  import AlertAiAnalyse from './alert-ai-analyse.vue'
  import { getAlarmDescriptionApi } from '~/src/api-ecs/alert'
  const alertIndexType = 10
  const router = useRouter()
  const alarmIcon = require('@/assets/alarm-icon.svg')
  const props = defineProps<{
    alertDetailVisible: boolean
    selectAlert: any
    moduleEnable?: boolean
    prev?: boolean
    next?: boolean
  }>()
  const emits = defineEmits<{
    (e: 'update:alertDetailVisible', traceabilityVisible: boolean): void
    (e: 'on-skipEvent', val: boolean): void
  }>()
  const show = ref(false)
  const aiRef = ref<HTMLDivElement>()
  const infoRef = ref<HTMLDivElement>()
  const collapse = ref(false)
  const infoType = ref<InfoType>('basic-info')
  const userStore = useUserStore()
  const { getTableColumn } = userStore
  const alertColumns = getTableColumn(alertIndexType)
  const tableColumn = ref<TableColumnItemType[]>([])
  const decodingTime = ref('')
  const decodingData = ref<PacketDecodeQuery>()
  const detailVisible = ref(false)
  const activeName = ref('basic-info')
  const jsonVal = ref()
  const showPrevBtn = ref(false)
  const attackResultType = ['成功', '失败', '未知', '企图']
  const isLoading = ref(false)
  const showNextBtn = ref(false)
  const decodingHandle = () => {
    const {
      attackIp: clientIp,
      sourcePort: clientPort,
      victimIp: serverIp,
      targetPort: serverPort,
      startTimeNs,
    } = props.selectAlert
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
  function changeVisible() {
    activeName.value = 'basic-info'
    window.removeEventListener('keyup', handleKeyDown, true)
    emits('update:alertDetailVisible', false)
  }
  const formatSourceData = (val: string) => {
    if (!val) return ''
    let _val = ''
    try {
      _val = JSON.parse(val)
    } catch (error) {
      _val = val
    }
    return _val
  }
  const isNeedToFormet = ref(true)

  provide('isNeedToFormet', isNeedToFormet)

  watchEffect(() => {
    detailVisible.value = props.alertDetailVisible
    if (props.alertDetailVisible) {
      isNeedToFormet.value = props.selectAlert.url ? true : false
      showPrevBtn.value = props.prev
      showNextBtn.value = props.next
      window.addEventListener('keyup', handleKeyDown, true)
      const _sourceData = formatSourceData(props.selectAlert.sourceData) as unknown as any
      jsonVal.value = JSON.stringify(_sourceData, null, 4)
      decodingHandle()
    } else {
      collapse.value = false
    }
  })
  const basicInfo = ref()
  const handlePrev = () => {
    emits('on-skipEvent', false)
    isLoading.value = true
    setTimeout(() => {
      basicInfo.value?.getAssetName()
      isLoading.value = false
    }, 400)
  }
  const handleNext = () => {
    emits('on-skipEvent', true)
    isLoading.value = true
    setTimeout(() => {
      basicInfo.value?.getAssetName()
      isLoading.value = false
    }, 400)
  }
  const addWhiteList = async () => {
    show.value = true
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
  function formatJson(filterVal: any, jsonData: any) {
    return jsonData.map((v: any) =>
      filterVal.map((j: any) => {
        return formatExcelData(v, j)
      })
    )
  }
  const handleDownloadExcel = () => {
    const tHeader: string[] = []
    const filterVal: any = []
    try {
      tableColumn.value.forEach(({ fieldNameCn, fieldNameEn }) => {
        tHeader.push(fieldNameCn)
        filterVal.push(fieldNameEn)
      })
      import('@/utils/excel').then((excel) => {
        const data = formatJson(filterVal, [props.selectAlert])
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: `告警日志-${formatNstime(props.selectAlert?.startTimeNs)}`,
          autoWidth: true,
          bookType: 'xlsx',
        })
      })
    } catch (error) {
      console.log(error)
    }
  }
  const handleKeyDown = (e: any) => {
    if (e.target.nodeName === 'TEXTAREA') return
    if (e.keyCode == 37 && showPrevBtn.value) {
      handlePrev()
    } else if (e.keyCode == 39 && showNextBtn.value) {
      handleNext()
    }
  }
  async function download() {
    try {
      const {
        attackIp: clientIp,
        sourcePort: clientPort,
        victimIp: serverIp,
        targetPort: serverPort,
        startTimeNs,
      } = props.selectAlert
      downloadLogPacket({ clientIp, clientPort, serverIp, serverPort, requestTimeNs: startTimeNs }, 2)
    } catch (error) {
      console.error(error)
    }
  }
  const fullFlowSurveyHandle = async () => {
    const { clientIp, serverIp, sourcePort, targetPort, startTimeNs } = props.selectAlert as any
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
  async function formatColum() {
    const { data: userDisPlaysFiled } = await getAllDisPlaysFiledApi()
    const userColumnData = userDisPlaysFiled[alertIndexType] as number[]
    if (userColumnData?.length > 0) {
      tableColumn.value = userColumnData.map((key) =>
        alertColumns.find((item) => item.id === key)
      ) as TableColumnItemType[]
    } else {
      tableColumn.value = []
    }
  }
  onMounted(() => {
    formatColum()
  })
</script>

<script lang="ts">
  export default {
    name: 'AlertDetail',
  }
  type InfoType = 'basic-info' | 'raw-data' | 'decoding-stream' | 'decoding-packet'
</script>

<template>
  <div class="alert-detail">
    <el-drawer v-model="detailVisible" destroy-on-close size="85%" @close="changeVisible">
      <template #header="{ titleId, titleClass }">
        <h2 :id="titleId" :class="titleClass">详情</h2>
      </template>
      <template #default>
        <div v-loading="isLoading" class="alarm-aggregation-detail-content">
          <div class="alarm-aggregation-detail-left">
            <dl class="alarm-aggregation-detail-header">
              <dt>
                <el-image :src="alarmIcon" style="width: 42px; height: 52px; margin-bottom: 4px" />
                <a class="alarm-btn">{{ selectAlert?.threatLevel || '-' }}</a>
              </dt>
              <dd>
                <label>告警时间</label>
                <div @click="useCopy(formatNstime(selectAlert?.startTimeNs))">
                  {{ formatNstime(selectAlert?.startTimeNs) }}
                </div>
              </dd>
              <dd>
                <label>威胁名称</label>
                <div @click="useCopy(selectAlert?.threatName)">{{ selectAlert?.threatName || '-' }}</div>
              </dd>
              <dd>
                <label>威胁类型</label>
                <div @click="useCopy(selectAlert?.threatType)">{{ selectAlert?.threatType || '-' }}</div>
              </dd>
              <dd>
                <label>源IP</label>
                <div style="display: ruby" @click="useCopy(selectAlert?.clientIp)">
                  <span>{{ selectAlert?.clientIp || '-' }}</span>
                  <span v-if="selectAlert?.client_country" style="white-space: nowrap">
                    <span>（</span>
                    <span>{{ selectAlert?.client_country }}</span>
                    <img
                      v-if="selectAlert?.client_country_code"
                      :src="require(`@/assets/flag/${selectAlert?.client_country_code.toLowerCase()}.png`)"
                      style="width: 20px; height: 15px; margin: 0 0 3px 4px"
                    />
                    <span>）</span>
                  </span>
                </div>
              </dd>
              <dd>
                <label>源端口</label>
                <div @click="useCopy(selectAlert?.sourcePort)">{{ selectAlert?.sourcePort || '-' }}</div>
              </dd>
              <dd>
                <label>目的IP</label>
                <div style="display: ruby" @click="useCopy(selectAlert?.victimIp)">
                  <span>{{ selectAlert?.victimIp || '-' }}</span>
                  <span v-if="selectAlert?.server_country" style="white-space: nowrap">
                    <span>（</span>
                    <span>{{ selectAlert?.server_country }}</span>
                    <img
                      v-if="selectAlert?.server_country_code"
                      :src="require(`@/assets/flag/${selectAlert?.server_country_code.toLowerCase()}.png`)"
                      style="width: 20px; height: 15px; margin: 0 0 3px 4px"
                    />
                    <span>）</span>
                  </span>
                </div>
              </dd>
              <dd>
                <label>目的端口</label>
                <div @click="useCopy(selectAlert?.targetPort)">{{ selectAlert?.targetPort || '-' }}</div>
              </dd>
              <dd>
                <label class="copyPath" @click="useCopy(selectAlert?.xff)">XFF</label>
                <div>{{ selectAlert?.xff || '-' }}</div>
              </dd>

              <dd>
                <label>告警设备IP</label>
                <div @click="useCopy(selectAlert?.warnDeviceIp)">{{ selectAlert?.warnDeviceIp || '-' }}</div>
              </dd>
              <dd>
                <label>攻击阶段</label>
                <div @click="useCopy(selectAlert?.killchain)">{{ selectAlert?.killchain || '-' }}</div>
              </dd>
              <dd>
                <label>攻击结果</label>
                <div @click="useCopy(selectAlert?.attackResult)">
                  <span
                    :class="[
                      'attackResult',
                      `attackResult-${attackResultType.findIndex((type) => type === selectAlert?.attackResult)}`,
                    ]"
                  >
                    <el-icon><WarnTriangleFilled /></el-icon>
                    {{ selectAlert.attackResult }}
                  </span>
                </div>
              </dd>
              <dd v-if="selectAlert.vulnHarm">
                <label>风险危害</label>
                <el-tooltip
                  :content="selectAlert.vulnHarm"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.vulnHarm)">
                    {{ selectAlert.vulnHarm }}
                  </div>
                </el-tooltip>
              </dd>
              <dd v-if="selectAlert.detailInfo">
                <label>威胁详情</label>
                <el-tooltip
                  :content="selectAlert.detailInfo"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.detailInfo)">
                    {{ selectAlert.detailInfo }}
                  </div>
                </el-tooltip>
              </dd>
              <dd v-if="selectAlert.attackDesc">
                <label>威胁描述</label>
                <el-tooltip
                  :content="selectAlert.attackDesc"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.attackDesc)">{{ selectAlert.attackDesc }}</div>
                </el-tooltip>
              </dd>
              <dd v-if="selectAlert.bulletin">
                <label>解决方案</label>
                <el-tooltip
                  :content="selectAlert.bulletin"
                  placement="top"
                  popper-class="selectAlert-attackDesc"
                  :show-after="1000"
                >
                  <div class="word-content" @click="useCopy(selectAlert.bulletin)">{{ selectAlert.bulletin }}</div>
                </el-tooltip>
              </dd>
            </dl>
          </div>
          <div class="alarm-aggregation-detail-center">
            <div class="alert-aggregation-detail-info" style="width: calc(100%)">
              <div ref="aiRef" class="alert-aggregation-detail-info-ai">
                <alert-ai-analyse :alarm-data="selectAlert" :collapse="collapse" :module-enable="moduleEnable" />
                <div v-if="moduleEnable" class="collapse" @click="handleClick">
                  <el-icon><ArrowDownBold /></el-icon>
                </div>
              </div>
              <div ref="infoRef" class="alert-aggregation-detail-info-content active">
                <div class="alert-aggregation-detail-info-header">
                  <el-button-group>
                    <el-button
                      :type="infoType === 'basic-info' ? 'primary' : 'default'"
                      @click="() => (infoType = 'basic-info')"
                    >
                      基本信息
                    </el-button>
                    <el-button
                      :type="infoType === 'raw-data' ? 'primary' : 'default'"
                      @click="() => (infoType = 'raw-data')"
                    >
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
                </div>
                <div style="width: 100%">
                  <basic-info
                    v-if="infoType === 'basic-info'"
                    ref="basicInfo"
                    :info-val="selectAlert"
                    type="alert-aggregation"
                  />
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
                  <div v-else-if="infoType === 'raw-data'" class="mask">
                    <json-preview :json-value="jsonVal" />
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="alarm-aggregation-detail-right">
            <el-button type="primary" @click="fullFlowSurveyHandle">全流量调查</el-button>
            <el-button type="primary" @click="addWhiteList">加入白名单</el-button>
            <el-button type="primary" @click="download">PCAP保存</el-button>
            <el-button plain type="primary" @click="handleDownloadExcel">下载告警详情</el-button>

            <el-button :disabled="!showPrevBtn" plain style="margin-top: 86px" @click="handlePrev">上一条</el-button>
            <el-button :disabled="!showNextBtn" plain @click="handleNext">下一条</el-button>
          </div>
        </div>
        <Whitelist v-if="show" :cur-data="selectAlert" :is-show="show" @on-closeEvent="show = false" />
      </template>
    </el-drawer>
    <template v-if="detailVisible">
      <Download class="download" :data="props.selectAlert" />
      <ToolBox class="toolBox" />
    </template>
  </div>
</template>

<style lang="scss">
  .selectAlert-attackDesc {
    transform: translate(0, 0);
    // transition: all 0s ease-in 0.5s;
    // transition: opacity 0.5s linear 0.5s;
    width: 230px;
  }
</style>

<style scoped lang="scss">
  .alert-detail {
    position: relative;
    .alarm-aggregation-detail-header {
      overflow: hidden;
      display: flex;
      flex-direction: column;
      margin-top: 28px;
      .word-content {
        // line-height: 24px;
        overflow: hidden;
        // max-height: 48px;
        text-overflow: ellipsis;
        word-break: break-all;
        word-wrap: break-word;
        display: -webkit-box;
        // background-clip: text;
        appearance: none;
        line-clamp: 2;
        -webkit-line-clamp: 2;
        /* autoprefixer: ignore next */
        -webkit-box-orient: vertical;
      }
      dt {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin-bottom: 30px;
        .alarm-btn {
          width: 56px;
          height: 24px;
          line-height: 20px;
          overflow: hidden;
          text-align: center;
          font-weight: 500;
          font-size: 14px;
          color: #ff4340;
          margin-top: 2px;
          border-radius: 12px;
          border: 2px solid transparent;
          background-image: linear-gradient(#fff0f0, #fff0f0),
            linear-gradient(166deg, rgba(255, 182, 182, 1), rgba(255, 184, 169, 1), rgba(255, 98, 98, 1));
          background-origin: border-box;
          background-clip: content-box, border-box;
        }
        .alarm-aggregation-detail-header-copy {
          font-size: 12px;
          color: #ffffff;
          width: 40px;
          height: 22px;
          border-radius: 4px;
          border: 1px solid #f8f7ff;
          display: inline-flex;
          position: absolute;
          top: 7px;
          right: 5px;
          text-indent: 0;
          text-indent: 8px;
          line-height: 20px;
        }
      }
      dd {
        display: flex;
        margin-left: 24px;
        text-align: left;
        label {
          font-weight: 500;
          font-size: 15px;
          color: #342e58;
          display: block;
          width: 80px;
          line-height: 26px;
        }
        div {
          flex: 1;
          overflow: hidden;
          display: -webkit-box;
          -webkit-box-orient: vertical;
          line-clamp: 8;
          -webkit-line-clamp: 8;
          word-break: break-all;
          text-overflow: ellipsis;
          font-weight: 400;
          font-size: 15px;
          color: #4a4759;
          line-height: 26px;
          margin-right: 20px;
          cursor: pointer;
          .attackResult {
            font-size: 14px;
            width: 58px;
            height: 24px;
            border-radius: 4px;
            display: block;
            line-height: 22px;
            text-align: center;
            .el-icon {
              margin-right: -4px;
              vertical-align: -3px;
              font-size: 17px;
            }
            &-0 {
              background: #fff5f5;
              color: #ff4340;
              border: 1px solid #ff4340;
            }
            &-1 {
              background: #f6f5ff;
              color: #9e9abd;
              border: 1px solid #9e9abd;
            }
            &-2 {
              background: #fff7ea;
              color: #ffa515;
              border: 1px solid #ffa515;
            }
            &-3 {
              background: #f5f4ff;
              border: 1px solid #6954f0;
              color: #6954f0;
            }
          }
        }
      }
    }
    .alarm-aggregation-detail-content {
      display: flex;
      height: 100%;
      .alarm-aggregation-detail-left {
        width: 316px;
      }
      .alarm-aggregation-detail-center {
        padding-top: 24px;
        background-color: #fff;
        width: calc(85vw - 516px);
        .alert-aggregation-detail-info {
          padding: 0 30px 25px;
          height: 100%;
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
              padding: 20px 0 15px;
              background: #fff;
              position: sticky;
              top: 0;
              z-index: 999;
            }
          }
        }
      }
      .alarm-aggregation-detail-right {
        width: 200px;
        background: #f8f8fd;
        padding: 20px 30px;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        :deep() {
          .el-button {
            width: 130px;
            margin-bottom: 14px;
            margin-left: 0;
          }
        }
      }
    }
    :deep() {
      .el-drawer {
        border-radius: 30px 0px 0px 30px;
        background: #534b89;
        .el-drawer__header {
          padding: 12px 30px 12px;
          margin-bottom: 0;
          .el-drawer__title {
            font-size: 20px;
            font-weight: 500;
            font-size: 20px;
            line-height: 34px;
            color: #fff;
          }
          .el-drawer__close-btn {
            width: 26px;
            height: 26px;
            border-radius: 50%;
            background: #eeedf9;
            display: flex;
            justify-content: center;
            align-items: center;
          }
        }
        .el-drawer__body {
          border-radius: 30px 0px 0px 30px;
          background: #f6f6fa;
          padding: 0;
        }
        .alert-aggregation-detail-info-content.active {
          .mask {
            height: calc(100vh - 362px);
          }
        }
      }
      .alarm-aggregation-detail-form {
        .el-form-item {
          width: 24%;
          margin-right: 1%;
          &:nth-of-type(4) {
            width: 25%;
            margin-right: 0;
          }
          &.search-btn {
            width: 50% !important;
            margin-right: 0;
            .el-form-item__content {
              justify-content: end;
            }
          }
        }
      }
    }
  }
  :deep() {
    .mask {
      width: 100%;
      height: 100%;
      overflow-y: auto;
      border: 1px solid rgb(230 231 240);
      background-color: #f8f7ff;
      border-radius: 8px;
    }
    .el-dialog {
      margin: var(--el-dialog-margin-top, 15vh) auto 0px;
    }
  }

  .download {
    position: fixed;
    left: 14px;
    bottom: 107px;
    z-index: 9999;
  }
  .toolBox {
    position: fixed;
    left: 14px;
    bottom: 66px;
    z-index: 9999;
  }
</style>
