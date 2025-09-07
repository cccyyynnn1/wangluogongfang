<script setup lang="ts">
  import { downloadLogPacket } from '@/utils/download'
  import { favoritesLogPageApi, appendCollectTraceApi } from '@/api-ecs/retrieve'
  import { formatNstime, formatTime } from '@/utils/time'
  import { useSettingsStore } from '@/store/modules/settings'
  import { PacketDecodeQuery } from '~/src/types'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'
  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'
  import dayjs from 'dayjs'
  import { getAssetsPreviewListApi } from '~/src/api-ecs/assets-preview'
  import { PacketsMagic } from '@/utils/magic'
  import { hexStringToArrayBuffer } from '@/utils/text'
  import { getPacketDecodeFlowDecodeApi } from '~/src/api-ecs/alert'
  const settingsStore = useSettingsStore()
  const props = defineProps<{
    infoVal: any
    workspaceId?: number
  }>()
  const $baseMessage: any = inject('$baseMessage')
  const selects = ref<string[]>([])
  const decodingData = ref<PacketDecodeQuery>()
  const decodingVisible = ref(false)
  const decodingTime = ref('')
  function download() {
    downloadLogPacket(props.infoVal, 1)
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
  // 数据流功能
  const formatType = ref('utf-8')
  const typeList = [
    {
      lable: 'UTF-8',
      value: 'utf-8',
    },
    {
      lable: 'UTF-16BE',
      value: 'utf-16be',
    },
    {
      lable: 'UTF-16LE',
      value: 'utf-16le',
    },
  ]
  const magicIcon = require('@/assets/mofabang.svg')
  const protocols = ['TCP', 'UDP']
  const protocol = ref('TCP')
  let STREAM_DATA = [] as { hexs: string[]; info: string; srcToDst: boolean }[]
  const protocolChange = (protocolVal: string) => {
    protocol.value = protocolVal
    getCode()
  }
  const showMagic = ref(false)
  const streamList = ref<{ hexs: string[]; info: string; srcToDst: boolean }[]>([])
  const magicStreamList = ref<{ hexs: string[]; info: string; srcToDst: boolean }[]>([])
  const toMagic = () => {
    if (showMagic.value) return (showMagic.value = false)
    magicStreamList.value = streamList.value.map((stream) => ({
      ...stream,
      hexs: stream.hexs.map((str: string) => PacketsMagic(str)),
    }))
    showMagic.value = true
  }
  const loading = ref(false)
  const getCode = async () => {
    showMagic.value = false
    if (!props.infoVal) return
    try {
      loading.value = true
      const { probeId = '', clientPort, serverIp, serverPort, clientIp, startTimeNs } = props.infoVal
      const timeDate = dayjs(formatNstime(startTimeNs))
      const startDate = timeDate.subtract(1, 'minute')
      const endDate = timeDate.add(1, 'minute')
      decodingTime.value = `${formatTime(startDate)} - ${formatTime(endDate)}`
      const { data } = await getPacketDecodeFlowDecodeApi({
        top: 50,
        query: {
          objectList: [
            {
              serverIp,
              serverPort: serverPort?.toString(),
              clientIp,
              probeId,
              clientPort: clientPort?.toString(),
            },
          ],
          timeStep: 'minuteStep',
          timeRange: decodingTime.value,
        },
        flowType: protocol.value,
      })
      STREAM_DATA = data || []
      toFormat(STREAM_DATA)
    } finally {
      loading.value = false
    }
  }
  const toFormat = (streamData: any[]) => {
    streamList.value = streamData.map((stream) => ({
      ...stream,
      hexs: stream.hexs.map((hex: string) => formatText(hex)),
    }))
  }
  const formatText = (txt: string) => {
    const textDecoder = new TextDecoder(formatType.value)
    const str = txt.replace(/\s*/g, '')
    return textDecoder.decode(hexStringToArrayBuffer(str))
  }
  defineExpose({
    getAssetName,
    getCode,
  })
</script>

<script lang="ts">
  export default {
    name: 'SmbDetail',
  }
</script>

<template>
  <div class="datail-smb">
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
      <el-descriptions-item label="客户端IP">
        <div v-copy="infoVal.clientIp" style="padding-right: 30px">
          {{ infoVal.clientIp }}
          <span style="margin-left: 6px; color: #4637a4">{{ clientIPAssetName }}</span>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="客户端口">
        <div v-copy="infoVal.clientPort" style="padding-right: 30px">{{ infoVal.clientPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="服务端IP">
        <div v-copy="infoVal.serverIp" style="padding-right: 30px">
          {{ infoVal.serverIp }}
          <span style="margin-left: 6px; color: #4637a4">{{ serverIPAssetName }}</span>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="服务端口">
        <div v-copy="infoVal.serverPort" style="padding-right: 30px">{{ infoVal.serverPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="开始时间">
        <div v-copy="infoVal?.startTimeNs" style="padding-right: 30px">
          {{ formatNstime(infoVal.startTimeNs) }}
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="响应时长">
        <div v-copy="infoVal?.responseTime" style="padding-right: 30px">
          {{ infoVal?.responseTime }}
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="服务器名称 ">
        <div v-copy="infoVal.serverName" style="padding-right: 30px">{{ infoVal.serverName }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="版本">
        <div v-copy="infoVal.version" style="padding-right: 30px">{{ infoVal.version }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="操作" :span="2">
        <div v-copy="infoVal.cmd" style="padding-right: 30px">{{ infoVal.cmd }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="账号" :span="2">
        <div v-copy="infoVal.account" style="padding-right: 30px">{{ infoVal.account }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="文件夹名称" :span="2">
        <div v-copy="infoVal.dirName" style="padding-right: 30px">{{ infoVal.dirName }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="文件名称" :span="2">
        <div v-copy="infoVal.fileName" style="padding-right: 30px">{{ infoVal.fileName }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="读写大小" :span="2">
        <div v-copy="infoVal.readWriteSize" style="padding-right: 30px">{{ infoVal.readWriteSize }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="执行结果" :span="2">
        <div v-copy="infoVal.result" style="padding-right: 30px">{{ infoVal.result }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="响应码" :span="2">
        <div v-copy="infoVal.responseCode" style="padding-right: 30px">{{ infoVal.responseCode }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="错误信息" :span="2">
        <div v-copy="infoVal.errorInfo" style="padding-right: 30px">{{ infoVal.errorInfo }}</div>
      </el-descriptions-item>
    </el-descriptions>
    <div class="hostTable_title">标签信息</div>
    <el-card shadow="never">
      <div v-for="(item, $index) in selects" :key="$index" class="myElTag">
        {{ item }}
      </div>
      <div v-if="selects.length === 0" style="text-align: center">暂无数据</div>
    </el-card>
    <!-- <div class="hostTable_title" style="display: flex; justify-content: space-between; align-items: center">
      <span>数据流</span>
      <div class="my-bts">
        <el-button-group>
          <el-button
            v-for="my_type of typeList"
            :key="my_type.value"
            size="small"
            :type="formatType === my_type.value ? 'primary' : 'default'"
            @click="formatType = my_type.value"
          >
            {{ my_type.lable }}
          </el-button>
        </el-button-group>
        <el-button-group style="margin-left: 15px; border: 0; background-color: transparent">
          <el-button
            v-for="my_protocol of protocols"
            :key="my_protocol"
            size="small"
            :type="protocol === my_protocol ? 'primary' : 'default'"
            @click="protocolChange(my_protocol)"
          >
            {{ my_protocol }}
          </el-button>
        </el-button-group>
        <el-button :auto-insert-space="false" style="border: 0; background-color: transparent" @click="toMagic()">
          <el-image :src="magicIcon" style="width: 16px; height: 16px; cursor: pointer" />
        </el-button>
      </div>
    </div>
    <div v-loading="loading" class="requests-and-responses">
      <div v-for="(stream, index) in showMagic ? magicStreamList : streamList" :key="index" class="stream">
        <h4>{{ stream.info }}</h4>
        <ul :class="{ isDst: stream.srcToDst }">
          <li v-for="(item, $index) in stream.hexs" :key="$index">{{ item }}</li>
        </ul>
      </div>
      <el-empty v-if="streamList.length === 0" description="暂无数据" />
    </div> -->
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
  .myElTag {
    background-color: #d9ecff;
    border-color: #d9ecff;
    color: var(--el-color-primary);
    background-color: rgba(102, 85, 231, 0.08);
    word-break: break-all;
    word-wrap: break-word;
    white-space: normal;
    display: inline-block;
    vertical-align: text-top;
    margin: 0 4px;
    padding: 3px 11px;
    border-radius: 4px;
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
