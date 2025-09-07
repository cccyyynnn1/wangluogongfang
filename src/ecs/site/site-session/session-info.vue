<script lang="ts">
  export default {
    name: 'SessionInfo',
  }
</script>
<script setup lang="ts">
  import { favoritesLogPageApi, appendCollectTraceApi, getIpLabelApi } from '~/src/api-ecs/retrieve'
  import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
  import { PacketDecodeQuery } from '~/src/types'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'
  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'
  import dayjs from 'dayjs'
  import { formatNstime, formatTime } from '@/utils/time'
  import { downloadLogPacket } from '@/utils/download'
  import ToolBox from '~/library/components/VabColumnBar/toolbox.vue'
  import Download from '~/library/components/VabColumnBar/download.vue'
  import { formatStrToJson } from '~/src/utils/transition'
  import { getAssetsPreviewListApi } from '~/src/api-ecs/assets-preview'
  import AttackHighlight from '@/components/attack-highlight.vue'
  import { PayloadDivider } from '~/src/hooks/useHttpTelegramMatch'
  import { removeLinerBeaksAndReturns } from '@/utils/text'
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    showSessionInfo: boolean
    infoData: any
    workspaceId?: number
    prev?: boolean
    next?: boolean
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-skipEvent', val: boolean): void
  }>()
  const magicIcon = require('@/assets/mofabang.svg')
  const showRequest = ref(false)
  const showResponse = ref(false)

  const decodingData = ref<PacketDecodeQuery>()
  const decodingVisible = ref(false)
  const decodingTime = ref('')
  const visible = ref(false)

  const allData = ref()

  const selects = ref<string[]>([])

  const requestList = ref([
    { Host: 'requestHost', label: '' },
    { Accept: 'requestAccept', label: '' },
    { Cookie: 'requestCookie', label: '' },
    { 'Proxy-Authorization': 'requestProxyAuthorization', label: '' },
    { 'User-Agent': 'requestUserAgent', label: '' },
    { ContentType: 'requestContentType', label: '' },
    { Referer: 'requestReferer', label: '' },
    { 'Content-Length': 'requestContentLength', label: '' },
    { Authorization: 'requestAuth', label: '' },
    { 'Accept-Language': 'requestAcceptLanguage', label: '' },
    { Connection: 'requestConnection', label: '' },
    { 'X-Forwarded-For': 'requestXForwardedFor', label: '' },
    { Via: 'requestVia', label: '' },
    { 'Accept-Charset': 'requestAcceptCharset', label: '' },
    { Rang: 'requestRang', label: '' },
    { 'If-Rang': 'requestIfRang', label: '' },
    { 'Accept-Encoding': 'requestAcceptEncoding', label: '' },
  ])

  const responseList = ref([
    { Server: 'responseServer', label: '' },
    { 'Content-Type': 'responseContentType', label: '' },
    { 'Content-Length': 'responseContentLength', label: '' },
    // { 'responseAuth': 'responseAuth', label: '' },
    { Connection: 'responseConnection', label: '' },
    { 'Set-Cookie': 'responseSetCookie', label: '' },
    { Authorization: 'responseProxyAuthorization', label: '' },
    { 'Content-Encoding': 'responseContentEncoding', label: '' },
    { Authorization: 'responseWwwAuthorization', label: '' },
    { 'Content-Disposition': 'responseContentDisposition', label: '' },
    { Via: 'responseVia', label: '' },
    { Location: 'responseLocation', label: '' },
  ])
  const attack_req_ref = ref<InstanceType<typeof AttackHighlight> | null>(null)
  const attack_res_ref = ref<InstanceType<typeof AttackHighlight> | null>(null)
  const codeMode = ref()
  const showPrevBtn = ref(false)

  const showNextBtn = ref(false)

  const getIpLabelHandle = async () => {
    if (props.infoData.requestXForwardedFor) {
      const { data } = await getIpLabelApi(props.infoData.requestXForwardedFor.split(','))
      const _ipLabel = Array.isArray(data) ? data : []
      setTimeout(() => {
        selects.value = [...selects.value, ..._ipLabel]
      })
    }
  }

  const isLoading = ref(false)

  const copyData = reactive({
    req: '',
    res: '',
  })

  const copyDataCopy = ref()

  const getCopyData = () => {
    let str = ''
    requestList.value.forEach((item: any) => {
      if (item.label) {
        str += `${item.label}\n`
      }
    })
    copyData.req = `${removeLinerBeaksAndReturns(
      `${allData.value.requestMethod} ${allData.value.requestUrl}`
    )}\n${removeLinerBeaksAndReturns(str)}\n${removeLinerBeaksAndReturns(
      allData.value.requestOtherData
    )}${PayloadDivider}${allData.value.requestPayload}`
    let str1 = ''
    responseList.value.forEach((item: any) => {
      if (item.label) {
        str1 += `${item.label}\n`
      }
    })
    copyData.res = `HTTP/${allData.value.responseVersion} ${
      allData.value.responseStatusCode
    }\n${removeLinerBeaksAndReturns(str1)}\n${removeLinerBeaksAndReturns(
      allData.value.responseOtherData
    )}${PayloadDivider}${allData.value.responsePayload}`
    copyDataCopy.value = copyData
  }

  onMounted(() => {
    visible.value = props.showSessionInfo
    allData.value = props.infoData
    showPrevBtn.value = props.prev
    showNextBtn.value = props.next
    window.addEventListener('keyup', handleKeyDown)
    getSelect()
    // getMetadata()
    changeData(requestList.value)
    changeData(responseList.value)
    getCopyData()
    getAssetName()
    // setTimeout(() => {
    //   initHeight()
    // }, 0)
  })

  const clientIPAssetName = ref()
  const serverIPAssetName = ref()
  const getAssetName = async () => {
    const { clientIp, clientPort, serverIp, serverPort, requestTimeNs: startTimeNs } = props.infoData
    const arr: string[] = []
    if (clientIp) {
      arr.push(clientIp)
    }
    if (serverIp) {
      arr.push(clientIp)
    }
    const {
      data: { records },
    } = await getAssetsPreviewListApi({
      pageNum: 1,
      pageSize: 10,
      searchIps: arr,
      dataSourceId: [] as number[],
      otherDataSource: false,
      searchStr: '',
    })
    if (clientIp) {
      clientIPAssetName.value = records[0]?.businessName || records[0]?.appName
    }
    if (serverIp) {
      serverIPAssetName.value = records[records.length - 1]?.businessName || records[records.length - 1]?.appName
    }
  }

  watch(
    () => props.infoData,
    () => {
      visible.value = props.showSessionInfo
      allData.value = props.infoData
      showRequest.value = false
      showResponse.value = false
      showPrevBtn.value = props.prev
      showNextBtn.value = props.next

      getSelect()
      changeData(requestList.value)
      changeData(responseList.value)
      getCopyData()
      getAssetName()
    },
    { deep: true }
  )
  const toMagic = (type: 'req' | 'res') => {
    if (type === 'req') {
      attack_req_ref.value?.setConvertsType('magic')
    } else {
      attack_res_ref.value?.setConvertsType('magic')
    }
    codeMode.value = ''
  }
  // 将数据转化成相应格式
  const changeData = (list: any) => {
    list.forEach((item: any) => {
      Object.keys(item).forEach((td) => {
        if (td !== 'label') {
          if (allData.value[item[td]]) {
            item['label'] = `${td}: ${allData.value[item[td]]}`
          }
        }
      })
    })
  }

  const getSelect = async () => {
    const { clientAssets, serverAssets, clientArea, serverArea, position } = props.infoData
    selects.value = []
    clientAssets && selects.value.push(clientAssets)
    serverAssets && selects.value.push(serverAssets)
    clientArea && selects.value.push(clientArea)
    serverArea && selects.value.push(serverArea)
    position && selects.value.push(position)
    getIpLabelHandle()
  }

  function formatDate(val: number) {
    const time = val / 1000000
    return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
  }

  const handleClose = () => {
    window.removeEventListener('keyup', handleKeyDown, false)
    emit('on-closeEvent', false)
  }

  const handleQueryByKey = async () => {
    downloadLogPacket(props.infoData, 1)
  }

  // 字段提取
  const fieldEextraction = () => {
    const origin = window.location.origin
    const url = `${origin}/#/config/data-extraction/extraction-rule?allData=allData`
    localStorage.setItem('fieldEextraction', JSON.stringify(allData.value))
    window.open(url)
  }
  const favoritesHandle = async () => {
    const { indexType, ...obj } = props.infoData
    const { data } = await favoritesLogPageApi({
      dataLog: JSON.stringify(obj),
      indexType,
      workspaceId: props.workspaceId!,
    })
    $baseMessage('添加收藏夹成功', 'success', 'vab-hey-message-success')
  }

  const traceabilityHandle = async () => {
    const { indexType, ...obj } = props.infoData
    const { data } = await appendCollectTraceApi({
      dataLog: JSON.stringify(obj),
      indexType,
      workspaceId: props.workspaceId!,
    })
    $baseMessage('添加溯源图成功', 'success', 'vab-hey-message-success')
  }
  // 下载数据
  const download = () => {
    handleQueryByKey()
  }
  const handleChange = (value: any) => {
    attack_req_ref.value?.setConvertsType(value)
    attack_res_ref.value?.setConvertsType(value)
  }
  const decodingHandle = () => {
    const { clientIp, clientPort, serverIp, serverPort, requestTimeNs: startTimeNs } = props.infoData
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
    decodingVisible.value = true
  }

  const handlePrev = () => {
    handleChange('default')
    codeMode.value = ''
    emit('on-skipEvent', false)
    isLoading.value = true
    setTimeout(() => {
      isLoading.value = false
    }, 400)
  }
  const handleNext = () => {
    handleChange('default')
    codeMode.value = ''
    emit('on-skipEvent', true)
    isLoading.value = true
    setTimeout(() => {
      isLoading.value = false
    }, 400)
  }

  const handleKeyDown = (e: any) => {
    if (e.keyCode == 37 && showPrevBtn.value == true) {
      handlePrev()
    } else if (e.keyCode == 39 && showNextBtn.value == true) {
      handleNext()
    }
  }
</script>

<template>
  <div class="session-info">
    <div v-loading="isLoading">
      <el-drawer v-model="visible" :before-close="handleClose" size="85%" title="会话详情">
        <div v-loading="isLoading">
          <div
            v-dialogBackTop
            style="
              display: flex;
              justify-content: space-between;
              align-items: center;
              position: sticky;
              top: 0;
              background-color: #fff;
              z-index: 99;
              padding-block: 20px;
            "
          >
            <div class="hostTable_title">连接信息</div>
            <div>
              <el-button :disabled="!showPrevBtn" :icon="ArrowLeft" plain type="primary" @click="handlePrev" />
              <el-button :disabled="!showNextBtn" :icon="ArrowRight" plain type="primary" @click="handleNext" />
              <el-button v-if="props.workspaceId" type="primary" @click="traceabilityHandle">添加到溯源图</el-button>
              <el-button v-if="props.workspaceId" type="primary" @click="favoritesHandle">添加到收藏夹</el-button>
              <!-- <el-button type="primary" @click="() => settingsStore.changeToolboxVisible(true)">工具箱</el-button> -->
              <el-button type="primary" @click="decodingHandle">PCAP解析</el-button>
              <ElButton type="primary" @click="download">PCAP保存</ElButton>
              <ElButton v-permissions="['Admin']" type="primary" @click="fieldEextraction">字段提取</ElButton>
            </div>
          </div>
          <ElDescriptions border :column="2">
            <ElDescriptionsItem label="源IP">
              <div v-copy="allData?.clientIp" style="padding-right: 30px">
                {{ allData?.clientIp || '-' }}
                <span style="margin-left: 6px; color: #4637a4">{{ clientIPAssetName }}</span>
              </div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="源端口">
              <div v-copy="allData?.clientPort" style="padding-right: 30px">{{ allData?.clientPort || '-' }}</div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="目的IP">
              <div v-copy="allData?.serverIp" style="padding-right: 30px">
                {{ allData?.serverIp || '-' }}
                <span style="margin-left: 6px; color: #4637a4">{{ serverIPAssetName }}</span>
              </div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="目的端口">
              <div v-copy="allData?.serverPort" style="padding-right: 30px">{{ allData?.serverPort || '-' }}</div>
            </ElDescriptionsItem>
            <!-- <ElDescriptionsItem label="请求时间">
          <div v-copy="allData?.requestTimeNs" style="padding-right: 30px">
            {{ formatDate(allData.requestTimeNs) || '-' }}
          </div>
        </ElDescriptionsItem> -->
            <!-- <ElDescriptionsItem label="响应时间">
          <div v-copy="allData?.responseTimeNs" style="padding-right: 30px">
            {{ formatDate(allData.responseTimeNs) || '-' }}
          </div>
        </ElDescriptionsItem> -->
            <ElDescriptionsItem label="响应标题">
              <div v-copy="allData?.title" style="padding-right: 30px">{{ allData?.title || '-' }}</div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="状态码">
              <div v-copy="allData?.responseStatusCode" style="padding-right: 30px">
                {{ allData?.responseStatusCode || '-' }}
              </div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="源IP资产">
              <div v-copy="allData?.clientAssets" style="padding-right: 30px">{{ allData?.clientAssets || '-' }}</div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="目的IP资产">
              <div v-copy="allData?.serverAssets" style="padding-right: 30px">{{ allData?.serverAssets || '-' }}</div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="源区域">
              <div v-copy="allData?.clientArea" style="padding-right: 30px">{{ allData?.clientArea || '-' }}</div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="目的区域">
              <div v-copy="allData?.serverArea" style="padding-right: 30px">{{ allData?.serverArea || '-' }}</div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="方向">
              <div v-copy="allData?.position" style="padding-right: 30px">{{ allData?.position || '-' }}</div>
            </ElDescriptionsItem>
            <ElDescriptionsItem label="XFF">
              <div v-copy="allData?.requestXForwardedFor" style="padding-right: 30px">
                {{ allData?.requestXForwardedFor || '-' }}
              </div>
            </ElDescriptionsItem>
          </ElDescriptions>
          <ElDescriptions border class="marginB15" :column="1" style="margin-top: -1px">
            <ElDescriptionsItem label="完整URL">
              <div v-copy="allData?.requestFullUrl" style="padding-right: 30px">
                {{ allData?.requestFullUrl || '-' }}
              </div>
            </ElDescriptionsItem>
          </ElDescriptions>
          <div class="hostTable_title marginB15">标签信息</div>
          <ElCard class="marginB15" shadow="never">
            <div v-for="(item, $index) in selects" :key="$index" class="myElTag">
              {{ item }}
            </div>
            <div v-if="selects.length === 0" style="text-align: center">暂无数据</div>
          </ElCard>
          <div class="hostTable_title marginB15">
            请求和响应
            <div class="fomart-code">
              <el-select v-model="codeMode" placeholder="请选择编码模式" style="width: 140px" @change="handleChange">
                <el-option key="default" label="恢复原始编码" value="default" />
                <el-option key="GB2312" label="GB2312" value="gb2312" />
                <el-option key="GBK" label="GBK" value="gbk" />
              </el-select>
            </div>
          </div>
          <div class="requests-and-responses">
            <div class="requests line">
              <div class="top">
                <span style="color: #4b4764; font-weight: 500">请求</span>
                <span style="color: #888b91">{{ formatDate(allData.requestTimeNs) || '-' }}</span>
                <el-tooltip content="魔法棒" effect="light" placement="top" :show-arrow="false">
                  <el-image
                    :src="magicIcon"
                    style="width: 16px; height: 16px; cursor: pointer"
                    @click="toMagic('req')"
                  />
                </el-tooltip>
              </div>
              <div class="content">
                <attack-highlight ref="attack_req_ref" :attack-http-msg="copyData.req" attack-type="request" />
              </div>
            </div>
            <div class="requests">
              <div class="top">
                <span style="color: #4b4764; font-weight: 500">响应</span>
                <span style="color: #888b91">{{ formatDate(allData.responseTimeNs) || '-' }}</span>
                <el-tooltip content="魔法棒" effect="light" placement="top" :show-arrow="false">
                  <el-image
                    :src="magicIcon"
                    style="width: 16px; height: 16px; cursor: pointer"
                    @click="toMagic('res')"
                  />
                </el-tooltip>
              </div>
              <div class="content">
                <attack-highlight ref="attack_res_ref" :attack-http-msg="copyData.res" attack-type="response" />
              </div>
            </div>
          </div>
        </div>
        <Download class="download" :data="props.infoData" />
        <ToolBox class="toolBox" />
      </el-drawer>
      <vab-dialog
        v-model="decodingVisible"
        :close-on-click-modal="false"
        destroy-on-close
        show-fullscreen
        title="详情"
        width="1175px"
      >
        <el-tabs :model-value="'decodingStream'">
          <el-tab-pane label="数据流" name="decodingStream">
            <decoding-stream :current-data="decodingData" :time-date="decodingTime" />
          </el-tab-pane>
          <el-tab-pane label="数据包" lazy name="decodingPacket">
            <decoding-packet :current-data="decodingData" :time-date="decodingTime" />
          </el-tab-pane>
        </el-tabs>
      </vab-dialog>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .hostTable_title {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .session-info {
    position: relative;
    :deep() {
      .el-dialog {
        margin: var(--el-dialog-margin-top, 15vh) auto 0px;
      }
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
          background: #fff;
          padding: 0 20px 20px;
        }
      }
    }
  }
  .toolBox {
    position: fixed;
    left: 15px;
    bottom: 65px;
  }
  .download {
    position: fixed;
    left: 15px;
    bottom: 105px;
  }
  .requests-and-responses {
    height: calc(100vh - 20px);
    width: 100%;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));

    display: flex;

    .line {
      border-right: 1px solid var(--el-border-color-lighter);
    }
    .requests {
      width: 50%;

      .top {
        padding: 0 20px;
        height: 40px;
        background: #f5f7fa;
        border-bottom: 1px solid var(--el-border-color-lighter);
        display: flex;
        align-items: center;
        justify-content: space-between;
      }
      .content {
        height: calc(100vh - 65px);
        overflow-y: auto;
        :deep() {
          .attack-highlight {
            background-color: transparent;
            border-radius: 0;
            border: 0;
            height: 100%;
            // background-color: red;
            // overflow-y: auto;
            .cm-scroller {
              min-height: 300px;
              .cm-content {
                width: 100%;
                white-space: break-spaces;
                word-wrap: break-word;
              }
            }
          }
        }
        &::-webkit-scrollbar {
          width: 0;
          height: 0;
        }
        .content-item {
          padding: 20px 20px 0;
          line-height: 20px;
          overflow: hidden;
          word-break: break-all;
          word-wrap: break-word;
          // box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color-lighter)) inset;
          // border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
          transition: var(--el-transition-box-shadow);
          white-space: pre-wrap;
          margin: 0;
          &:hover {
            cursor: text;
          }
        }
        :deep() {
          .el-textarea {
            margin: 0;
            border: none !important;
            background-color: #fff;
            .el-textarea__inner {
              padding: 15px 20px;
              background-color: #fff;
              border: none !important;
              box-shadow: none;
            }
          }
        }
      }
    }
  }
  .marginB15 {
    margin-bottom: 15px !important;
  }
  .my-item {
    background: #f8fbff;
    padding: 20px;
    line-height: 20px;
    overflow: hidden;
    word-break: break-all;
    word-wrap: break-word;
    box-shadow: 0 0 0 1px var(--el-input-border-color, var(--el-border-color-lighter)) inset;
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    transition: var(--el-transition-box-shadow);
    white-space: pre-wrap;
    margin: 0;
    &:hover {
      cursor: text;
    }
  }

  // .el-textarea {
  //   color: #303133;

  // :deep() {
  // .el-textarea__inner {
  // background: #f8fbff;
  // padding: 20px;

  // &:hover {
  //   box-shadow: 0 0 0 1px var(--el-border-color-lighter) inset !important;
  // }

  // &:focus {
  //   box-shadow: 0 0 0 1px var(--el-border-color-lighter) inset !important;
  // }
  // }
  // }
  // }

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
      top: 0px;
      background: var(--el-color-primary);
      margin-right: 5px;
    }
  }

  .el-descriptions {
    clear: both;
    :deep() {
      .el-descriptions__label {
        width: 172px;
      }

      .el-descriptions__content {
        width: 305px;
      }
    }

    &:nth-child(3) {
      background-color: red;

      :deep() {
        .el-descriptions__label {
          width: 172px;
        }

        .el-descriptions__content {
          width: 782px;
          overflow: hidden;
          word-break: break-all;
          // white-space: nowrap;
          word-wrap: break-word;
        }
      }
    }
  }

  .el-space--horizontal {
    float: right;
    margin: 15px -10px 20px 0;
  }

  .myElTag {
    background-color: #d9ecff;
    border-color: #d9ecff;
    color: var(--el-color-primary);
    background-color: rgba(102, 85, 231, 0.08);
    white-space: normal;
    display: inline-block;
    vertical-align: text-top;
    margin: 0 4px;
    padding: 3px 11px;
    border-radius: 4px;
    word-break: break-all;
    word-wrap: break-word;
  }
</style>
