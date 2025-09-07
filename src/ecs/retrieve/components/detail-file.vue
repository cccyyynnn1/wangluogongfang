<script setup lang="ts">
  import { downloadFile, downloadLogPacket } from '@/utils/download'
  import { favoritesLogPageApi, appendCollectTraceApi, getFlowSearchShowFileAPI } from '@/api-ecs/retrieve'
  import { formatNstime, formatTime } from '@/utils/time'
  import { useSettingsStore } from '@/store/modules/settings'
  import { PacketDecodeQuery } from '~/src/types'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'
  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'
  import dayjs from 'dayjs'
  import { getAssetsPreviewListApi } from '~/src/api-ecs/assets-preview'
  import { proxyNet } from '@/config/index'
  import { FileRestoreDownloadFileAPI } from '~/src/api-ecs/unusual'
  const $baseMessage: any = inject('$baseMessage')
  const settingsStore = useSettingsStore()
  const props = defineProps<{
    infoVal: any
    workspaceId?: number
  }>()
  const selects = ref<string[]>([])
  const decodingData = ref<PacketDecodeQuery>()
  const decodingVisible = ref(false)
  const decodingTime = ref('')

  function download() {
    downloadLogPacket(props.infoVal, 1)
  }
  async function getSelect() {
    const { clientAssets, serverAssets, clientArea, serverArea, position } = props.infoVal
    selects.value = []
    clientAssets && selects.value.push(clientAssets)
    serverAssets && selects.value.push(serverAssets)
    clientArea && selects.value.push(clientArea)
    serverArea && selects.value.push(serverArea)
    position && selects.value.push(position)
  }
  const traceabilityHandle = async () => {
    const { indexType, ...obj } = props.infoVal
    const { data } = await appendCollectTraceApi({
      dataLog: JSON.stringify(obj),
      indexType,
      workspaceId: props.workspaceId!,
    })
    $baseMessage('添加溯源图成功', 'success', 'vab-hey-message-success')
  }
  const favoritesHandle = async () => {
    const { indexType, ...obj } = props.infoVal
    const { data } = await favoritesLogPageApi({
      dataLog: JSON.stringify(obj),
      indexType,
      workspaceId: props.workspaceId!,
    })
    $baseMessage('添加收藏夹成功', 'success', 'vab-hey-message-success')
  }
  const decodingHandle = () => {
    const { clientIp, clientPort, serverIp, serverPort, startTimeNs } = props.infoVal
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

  const mode = ref('Hex')

  const codeData = ref()
  const getCode = async () => {
    setTimeout(async () => {
      const query = props.infoVal
      let flag = true
      const res = await getFlowSearchShowFileAPI({ ...query })
      codeData.value = res?.data?.data
    }, 100)
  }

  onMounted(() => {
    getSelect()
    getAssetName()
    getCode()
  })

  const clientIPAssetName = ref()
  const serverIPAssetName = ref()
  const getAssetName = async () => {
    const { clientIp, serverIp } = props.infoVal
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

  function utf8ToUtf16be(str: any) {
    const utf16Buffer = new ArrayBuffer(str.length * 2)
    const view = new Uint16Array(utf16Buffer)
    for (let i = 0; i < str.length; i++) {
      view[i] = str.charCodeAt(i)
    }
    return new Uint8Array(utf16Buffer).buffer
  }

  function utf16beToUtf8(utf16beString: any) {
    // 将UTF-16BE字符串转换为ArrayBuffer
    const buffer = new ArrayBuffer(utf16beString.length * 2)
    const view = new Uint16Array(buffer)
    for (let i = 0; i < utf16beString.length; i++) {
      view[i] = utf16beString.charCodeAt(i)
    }

    // 使用TextDecoder转换为UTF-8
    const decoder = new TextDecoder('utf-8')
    return decoder.decode(buffer.slice(2)) // 假设我们要去掉字节序标记(BOM)
  }

  function utf8ToUtf16le(utf8String: string) {
    const utf16Buffer = new ArrayBuffer(utf8String.length * 2)
    const view = new Uint16Array(utf16Buffer)
    for (let i = 0; i < utf8String.length; i++) {
      view[i] = utf8String.charCodeAt(i)
    }
    return new Uint8Array(utf16Buffer).buffer
  }

  function utf16leToUtf8(str: any) {
    let utf16 = decodeURIComponent(escape(str))
    let utf8 = unescape(encodeURIComponent(utf16))
    return utf8
  }

  function utf8ToAscii(str: any) {
    return str.replace('/[^\x00-\x7F]/g', '')
  }

  function asciiToUtf8(asciiString: any) {
    // 将ASCII字符串转换为数组
    const asciiArray = asciiString.split('')

    // 转换ASCII数组到UTF-8数组
    const utf8Array = asciiArray.map((char: any) => {
      const code = char.charCodeAt(0)
      if (code <= 0x7f) {
        // 0xxxxxxx
        return char
      } else if (code <= 0x7ff) {
        // 110xxxxx 10xxxxxx
        const byte1 = 0xc0 | (code >> 6)
        const byte2 = 0x80 | (code & 0x3f)
        return String.fromCharCode(byte1, byte2)
      } else {
        // 1110xxxx 10xxxxxx 10xxxxxx
        const byte1 = 0xe0 | (code >> 12)
        const byte2 = 0x80 | ((code >> 6) & 0x3f)
        const byte3 = 0x80 | (code & 0x3f)
        return String.fromCharCode(byte1, byte2, byte3)
      }
    })

    // 将UTF-8数组合并为字符串
    return utf8Array.join('')
  }

  function utf8ToHex(str: any) {
    const encoder = new TextEncoder()
    const utf8Array = encoder.encode(str)
    let hexStr = ''
    for (const byte of utf8Array) {
      hexStr += byte.toString(16).padStart(2, '0')
    }
    return hexStr.replace(/../g, '$& ')
  }

  function hexToUtf8(hexString: any) {
    // 将十六进制字符串转换为ArrayBuffer
    const noSpacesStr = hexString.replace(/\s+/g, '')
    const buffer = new Uint8Array(noSpacesStr.match(/.{2}/g).map((byte: any) => parseInt(byte, 16))).buffer

    // 使用TextDecoder进行解码
    const decoder = new TextDecoder('utf-8')
    return decoder.decode(buffer)
  }

  const handleChangeCode = (val: string) => {
    if (mode.value == val) return
    const meta = JSON.parse(JSON.stringify(mode.value))
    mode.value = val
    let temp = ''
    let res = undefined
    // 'UTF-8','UTF-16BE','UTF-16LE','Hex','ASCII'
    switch (meta) {
      case 'UTF-16BE':
        temp = utf16beToUtf8(codeData.value)
        break
      case 'UTF-16LE':
        temp = utf16leToUtf8(codeData.value)
        break
      case 'Hex':
        temp = hexToUtf8(codeData.value)
        break
      case 'ASCII':
        temp = asciiToUtf8(codeData.value)
        break

      default:
        temp = codeData.value
        break
    }
    switch (val) {
      case 'UTF-16BE':
        res = utf8ToUtf16be(temp)
        break
      case 'UTF-16LE':
        res = utf8ToUtf16le(temp)
        break
      case 'Hex':
        res = utf8ToHex(temp)
        break
      case 'ASCII':
        res = utf8ToAscii(temp)
        break

      default:
        res = temp
        break
    }
    codeData.value = res
  }

  // 下载
  const disabled = ref(false)
  const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname
  const URL = `https://${url}/v3/admin/ci_proxy/rest/v1/flowSearch/downloadFile`
  const handleDownload = async () => {
    const { dataLens, fileName, dataOffsets } = props.infoVal
    const { data, headers } = await FileRestoreDownloadFileAPI({
      query: JSON.stringify({
        dataLens,
        fileName,
        dataOffsets,
      }),
    })
    try {
      disabled.value = true
      const downloadElement = document.createElement('a')
      const blob = new Blob([data], { type: headers['content-type'] })
      const href = window.URL.createObjectURL(blob)
      downloadElement.href = href
      let fileNames = ''
      const contentDisposition = headers['content-disposition']
      if (contentDisposition) {
        fileNames = window.decodeURI(headers['content-disposition'].split('=')[1])
        downloadElement.download = fileNames
      }
      document.body.appendChild(downloadElement)
      downloadElement.click()
      document.body.removeChild(downloadElement)
      window.URL.revokeObjectURL(href)
      disabled.value = false
    } catch (error) {
      disabled.value = false
    }

    // downloadFile()
  }
  defineExpose({
    getCode,
    getAssetName,
  })
</script>

<script lang="ts">
  export default {
    name: 'DetailFile',
  }
</script>
<template>
  <div class="datail-file">
    <div class="marginB15" style="display: flex; justify-content: space-between; align-items: center">
      <div class="hostTable_title">基本信息</div>
      <div>
        <el-space :size="10">
          <el-button v-if="props.workspaceId" type="primary" @click="traceabilityHandle">添加到溯源图</el-button>
          <el-button v-if="props.workspaceId" type="primary" @click="favoritesHandle">添加到收藏夹</el-button>
          <!-- <el-button type="primary" @click="() => settingsStore.changeToolboxVisible(true)">工具箱</el-button> -->
          <el-button type="primary" @click="decodingHandle">PCAP解析</el-button>
          <el-button type="primary" @click="download">PCAP保存</el-button>
        </el-space>
      </div>
    </div>
    <el-descriptions border :column="2">
      <el-descriptions-item label="日期">
        <div v-copy="infoVal?.statTimeSec" style="padding-right: 30px">
          {{ infoVal?.date }}
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="文件名">
        <div v-copy="infoVal?.fileName" style="padding-right: 30px">{{ infoVal?.fileName }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="文件大小">
        <div v-copy="infoVal?.fileBytes" style="padding-right: 30px">{{ infoVal?.fileBytes }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="源IP">
        <div v-copy="infoVal?.clientIp" style="padding-right: 30px">
          {{ infoVal?.clientIp }}
          <span style="margin-left: 6px; color: #4637a4">{{ clientIPAssetName }}</span>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="源端口">
        <div v-copy="infoVal?.clientPort" style="padding-right: 30px">{{ infoVal?.clientPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="目的IP">
        <div v-copy="infoVal?.serverIp" style="padding-right: 30px">
          {{ infoVal?.serverIp }}
          <span style="margin-left: 6px; color: #4637a4">{{ serverIPAssetName }}</span>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="目的端口">
        <div v-copy="infoVal?.serverPort" style="padding-right: 30px">{{ infoVal?.serverPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="文件信息">
        <div v-copy="infoVal?.fileInfo" style="padding-right: 30px">{{ infoVal?.fileInfo }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="文件类型">
        <div v-copy="infoVal?.fileType" style="padding-right: 30px">{{ infoVal?.fileType }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="链路">
        <div v-copy="infoVal?.probeIds[0]" style="padding-right: 30px">{{ infoVal?.probeIds[0] }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="SHAI值">
        <div v-copy="infoVal?.sha1" style="padding-right: 30px">{{ infoVal?.sha1 }}</div>
      </el-descriptions-item>
    </el-descriptions>
    <div class="hostTable_title">标签信息</div>
    <el-card shadow="never">
      <div v-for="(item, $index) in selects" :key="$index" class="myElTag">
        {{ item }}
      </div>
      <div v-if="selects.length === 0" style="text-align: center">暂无数据</div>
    </el-card>
    <div class="hostTable_title">
      <span>文件传输</span>
      <div class="my-bts">
        <el-button-group>
          <!-- TODO -->
          <!-- <el-button size="small" type="primary" @click="handleChangeCode('UTF-8')">UTF-8</el-button> -->
          <!-- <el-button size="small" @click="handleChangeCode('UTF-16BE')">UTF-16BE</el-button>
          <el-button size="small" @click="handleChangeCode('UTF-16LE')">UTF-16LE</el-button> -->
          <!-- <el-button size="small" style="margin-left: 10px" type="primary" @click="handleChangeCode('Hex')">
            十六进制
          </el-button> -->
          <!-- <el-button size="small" @click="handleChangeCode('ASCII')">ASCII</el-button> -->
          <el-button
            :disabled="infoVal?.fileBytes == 0 || disabled"
            size="small"
            style="margin-left: 20px"
            type="primary"
            @click="handleDownload"
          >
            下载
          </el-button>
        </el-button-group>
      </div>
    </div>
    <div class="requests-and-responses">
      <div class="requests line">
        <div class="top">
          <span style="color: #4b4764; font-weight: 500">{{ infoVal?.fileName }}</span>
        </div>
        <div class="data-warp" style="padding: 20px">
          {{ codeData }}
        </div>
      </div>
      <!-- {{ codeData }} -->
    </div>
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
</template>

<style scoped lang="scss">
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
      width: 100%;
      height: 100%;
      overflow-y: auto;
      &::-webkit-scrollbar {
        width: 0;
        height: 0;
      }

      .data-warp {
        overflow-y: auto;
        text-overflow: ellipsis;
        word-break: break-all;
        word-wrap: break-word;
      }

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
  .el-descriptions {
    margin-bottom: 15px;
    clear: both;

    :deep() {
      .el-descriptions__label {
        width: 172px;
      }

      .el-descriptions__content {
        width: 305px;
      }
    }
  }

  .el-space--horizontal {
    float: right;
    margin: 5px -10px 15px 0;
  }

  .hostTable_title {
    height: 14px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    line-height: 14px;
    position: relative;
    text-indent: 0.8em;
    font-weight: 700;
    color: #303133;
    margin-bottom: 15px;

    &::before {
      content: ' ';
      display: inline-block;
      position: absolute;
      width: 3px;
      height: 100%;
      left: 0;
      background: var(--el-color-primary);
      margin-right: 5px;
    }
  }
</style>
