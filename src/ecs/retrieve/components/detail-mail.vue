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
  async function getSelect() {
    const { clientAssets, serverAssets, clientArea, serverArea, position } = props.infoVal
    selects.value = []
    clientAssets && selects.value.push(clientAssets)
    serverAssets && selects.value.push(serverAssets)
    clientArea && selects.value.push(clientArea)
    serverArea && selects.value.push(serverArea)
    position && selects.value.push(position)
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
  const res = ref(`  <div class="datail-mail">
    <el-space :size="10">
      <el-button v-if="props.workspaceId" type="primary" @click="traceabilityHandle">添加到溯源图</el-button>
      <el-button v-if="props.workspaceId" type="primary" @click="favoritesHandle">添加到收藏夹</el-button>
      <el-button type="primary" @click="() => settingsStore.changeToolboxVisible(true)">工具箱</el-button>
      <el-button type="primary" @click="decodingHandle">PCAP解析</el-button>
      <el-button type="primary" @click="download">PCAP保存</el-button>
    </el-space>
    <el-descriptions border :column="2">
      <el-descriptions-item label="客户端IP">
        <div v-copy="infoVal.clientIp" style="padding-right: 30px">{{ infoVal.clientIp }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="客户端口">
        <div v-copy="infoVal.clientPort" style="padding-right: 30px">{{ infoVal.clientPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="服务端IP">
        <div v-copy="infoVal.serverIp" style="padding-right: 30px">{{ infoVal.serverIp }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="服务端口">
        <div v-copy="infoVal.serverPort" style="padding-right: 30px">{{ infoVal.serverPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="请求时间">
        <div v-copy="formatNstime(infoVal.startTimeNs)" style="padding-right: 30px">
          {{ formatNstime(infoVal.startTimeNs) }}
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="响应时长 ">
        <div v-copy="infoVal.responseTime" style="padding-right: 30px">{{ infoVal.responseTime }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="服务器">
        <div v-copy="infoVal.serverName" style="padding-right: 30px">{{ infoVal.serverName }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="版本">
        <div v-copy="infoVal.version" style="padding-right: 30px">{{ infoVal.version }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="类型">
        <div v-copy="infoVal.model" style="padding-right: 30px">{{ infoVal.model }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="操作类型">  <div class="datail-mail">
    <el-space :size="10">
      <el-button v-if="props.workspaceId" type="primary" @click="traceabilityHandle">添加到溯源图</el-button>
      <el-button v-if="props.workspaceId" type="primary" @click="favoritesHandle">添加到收藏夹</el-button>
      <el-button type="primary" @click="() => settingsStore.changeToolboxVisible(true)">工具箱</el-button>
      <el-button type="primary" @click="decodingHandle">PCAP解析</el-button>
      <el-button type="primary" @click="download">PCAP保存</el-button>
    </el-space>
    <el-descriptions border :column="2">
      <el-descriptions-item label="客户端IP">
        <div v-copy="infoVal.clientIp" style="padding-right: 30px">{{ infoVal.clientIp }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="客户端口">
        <div v-copy="infoVal.clientPort" style="padding-right: 30px">{{ infoVal.clientPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="服务端IP">
        <div v-copy="infoVal.serverIp" style="padding-right: 30px">{{ infoVal.serverIp }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="服务端口">
        <div v-copy="infoVal.serverPort" style="padding-right: 30px">{{ infoVal.serverPort }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="请求时间">
        <div v-copy="formatNstime(infoVal.startTimeNs)" style="padding-right: 30px">
          {{ formatNstime(infoVal.startTimeNs) }}
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="响应时长 ">
        <div v-copy="infoVal.responseTime" style="padding-right: 30px">{{ infoVal.responseTime }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="服务器">
        <div v-copy="infoVal.serverName" style="padding-right: 30px">{{ infoVal.serverName }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="版本">
        <div v-copy="infoVal.version" style="padding-right: 30px">{{ infoVal.version }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="类型">
        <div v-copy="infoVal.model" style="padding-right: 30px">{{ infoVal.model }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="操作类型">`)
  const ress = ref()
  const emailContent = ref()
  const loading = ref(false)
  onMounted(() => {
    emailContent.value = props.infoVal.content
    getSelect()
    setTimeout(() => {
      ress.value = res
    }, 1000)
    getAssetName()
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
  defineExpose({
    getAssetName,
  })
</script>

<script lang="ts">
  export default {
    name: 'MailDetail',
  }
</script>

<template>
  <div class="datail-mail">
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
      <el-descriptions-item label="服务器">
        <div v-copy="infoVal.serverName" style="padding-right: 30px">{{ infoVal.serverName }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="版本">
        <div v-copy="infoVal.version" style="padding-right: 30px">{{ infoVal.version }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="类型">
        <div v-copy="infoVal.model" style="padding-right: 30px">{{ infoVal.model }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="操作类型">
        <div v-copy="infoVal.operateType" style="padding-right: 30px">{{ infoVal.operateType }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="账号">
        <div v-copy="infoVal.account" style="padding-right: 30px">{{ infoVal.account }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="密码">
        <div v-copy="infoVal.passWord" style="padding-right: 30px">{{ infoVal.passWord }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="发件人" :span="2">
        <div v-copy="infoVal.sender" style="padding-right: 30px">{{ infoVal.sender }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="收件人" :span="2">
        <div v-copy="infoVal.receiver" style="padding-right: 30px">{{ infoVal.receiver }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="邮件标题" :span="2">
        <div v-copy="infoVal.title" style="padding-right: 30px">{{ infoVal.title }}</div>
      </el-descriptions-item>
      <!-- <el-descriptions-item label="邮件内容" :span="2">
        <div v-fold2>
          <div v-copy="infoVal.content" style="padding-right: 30px">{{ infoVal.content }}</div>
        </div>
      </el-descriptions-item> -->
      <el-descriptions-item label="附件" :span="2">
        <div v-copy="infoVal.attachment" v-fold style="padding-right: 30px">{{ infoVal.attachment }}</div>
      </el-descriptions-item>
      <el-descriptions-item label="执行结果" :span="2">
        <div v-copy="infoVal.result" style="padding-right: 30px">{{ infoVal.result }}</div>
      </el-descriptions-item>
    </el-descriptions>
    <div class="hostTable_title">标签信息</div>
    <el-card shadow="never">
      <div v-for="(item, $index) in selects" :key="$index" class="myElTag">
        {{ item }}
      </div>
      <div v-if="selects.length === 0" style="text-align: center">暂无数据</div>
    </el-card>
    <div class="hostTable_title" style="display: flex; justify-content: space-between; align-items: center">
      <span>邮件内容</span>
      <!-- <div class="my-bts">
        <el-button-group>
          <el-button
            v-for="my_type of typeList"
            :key="my_type.value"
            :type="formatType === my_type.value ? 'primary' : 'default'"
            @click="formatType = my_type.value"
            size="small"
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
        </el-button> -->
      <!-- </div> -->
    </div>
    <div v-loading="loading" class="requests-and-responses">
      {{ emailContent }}
      <el-empty v-if="!emailContent" description="暂无数据" />
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
    height: calc(100vh - 60px);
    width: 100%;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    overflow-y: auto;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    padding: 20px;
    .data-warp {
      overflow-y: auto;
      text-overflow: ellipsis;
      word-break: break-all;
      word-wrap: break-word;
    }
  }
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
