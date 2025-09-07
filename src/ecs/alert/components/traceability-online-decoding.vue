<script lang="ts">
  export default {
    name: 'OnlineDecoding',
  }
</script>

<script setup lang="ts">
  import { getPacketDecodeApi } from '@/api-ecs/alert'
  import { defaultTime } from '@/data/constant'
  import { downloadLogPacket } from '@/utils/download'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'
  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'
  import { PacketDecodeItem, PacketDecodeQuery } from '@/types'
  import dayjs from 'dayjs'
  const props = defineProps<{
    alertData?: any
    nodeData: any
    timeRange?: [string, string]
  }>()
  const option = [
    { label: 'TOP100', value: 100 },
    { label: 'TOP500', value: 500 },
    { label: 'TOP1000', value: 1000 },
    { label: 'TOP2000', value: 2000 },
    { label: 'TOP5000', value: 5000 },
  ]

  // 表单数据
  const queryForm = reactive({
    page: 1,
    limit: 10,
    total: 0,
    topCount: 100,
    lastTime: '',
    list: [] as PacketDecodeItem[],
    loading: true,
  })
  const currentData = ref<PacketDecodeQuery>()
  // 自定义时间
  const timeDate = ref()
  // 获取表格序号
  const curIndex = computed(() => (queryForm.page - 1) * queryForm.limit + 1)
  const decodingVisible = ref(false)
  const getData = async () => {
    queryForm.loading = true
    const { topCount, lastTime } = queryForm
    const { type, serverIp, snat, ip } = props.nodeData
    const query = type
      ? { serverIp: serverIp, clientIp: snat, timeRange: lastTime, top: topCount }
      : { timeRange: lastTime, ipAddr: ip, top: topCount }
    try {
      const { data } = await getPacketDecodeApi(query)
      queryForm.list = data
    } catch (error) {
      console.error(error)
    }
    queryForm.loading = false
  }

  async function download(row: any) {
    const startTimeNs = props.alertData ? props.alertData.startTimeNs : new Date(timeDate.value[0]).getTime() * 1000000
    try {
      const { clientIp, clientPort, serverIp, serverPort = '' } = row
      downloadLogPacket({ clientIp, clientPort, serverIp, serverPort, requestTimeNs: startTimeNs }, 2)
    } catch (error) {
      console.error(error)
    }
  }

  const decodingHandle = (row: any) => {
    const { clientIp, clientPort, probeId, serverIp, serverPort } = row
    currentData.value = {
      clientIp,
      clientPort,
      probeId,
      serverIp,
      serverPort,
    }
    currentData.value = row
    decodingVisible.value = true
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
    if (props.timeRange) {
      timeDate.value = props.timeRange
    }
  })
  watchEffect(() => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      queryForm.lastTime = `${formatTime(new Date(startDate).getTime())} - ${formatTime(new Date(endDate).getTime())}`
    }
  })
  onMounted(() => {
    getData()
  })
</script>

<template>
  <el-row :gutter="20">
    <el-col :span="8">
      <vab-date-time-picker v-model="timeDate" />
    </el-col>
    <el-col :span="12">
      <el-select v-model="queryForm.topCount">
        <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </el-col>
    <el-col :span="4">
      <el-button :loading="queryForm.loading" style="float: right; margin-left: 20px" type="primary" @click="getData">
        检索
      </el-button>
    </el-col>
  </el-row>
  <el-table v-loading="queryForm.loading" :data="queryForm.list" style="margin-top: 20px; width: 100%; height: 440px">
    <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="70" />
    <el-table-column align="center" label="源IP" prop="clientIp" />
    <el-table-column align="center" label="源端口" prop="clientPort" />
    <el-table-column align="center" label="目的IP" prop="serverIp" />
    <el-table-column align="center" label="目的端口" prop="serverPort" />
    <el-table-column align="center" fixed="right" label="操作" width="250">
      <template #default="{ row }">
        <el-button class="row_action" size="small" @click="decodingHandle(row)">PCAP解析</el-button>
        <el-button class="row_action" size="small" @click="download(row)">PCAP保存</el-button>
      </template>
    </el-table-column>
    <template #empty>
      <el-empty class="vab-data-empty" description="暂无数据" />
    </template>
  </el-table>

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
        <decoding-stream :current-data="currentData" :time-date="queryForm.lastTime" />
      </el-tab-pane>
      <el-tab-pane label="数据包" lazy name="decodingPacket">
        <decoding-packet :current-data="currentData" :time-date="queryForm.lastTime" />
      </el-tab-pane>
    </el-tabs>
  </vab-dialog>
</template>

<style scoped lang="scss">
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
