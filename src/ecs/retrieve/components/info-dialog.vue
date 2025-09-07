<script setup lang="ts">
  import { getPacketDecodeApi } from '@/api-ecs/alert'
  import { downloadLogPacket } from '@/utils/download'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'
  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'
  import { PacketDecodeItem, PacketDecodeQuery } from '@/types'
  import dayjs from 'dayjs'
  import { useTableCopy } from '~/src/utils'
  const props = defineProps<{
    showInfoDialog: boolean
    currentRow: any
  }>()

  const visible = ref(false)

  const option = [
    { label: 'TOP100', value: 100 },
    { label: 'TOP500', value: 500 },
    { label: 'TOP1000', value: 1000 },
    { label: 'TOP2000', value: 2000 },
    { label: 'TOP5000', value: 5000 },
  ]

  const searchValue = ref('')

  // 表单数据
  const queryForm = reactive({
    page: 1,
    limit: 10,
    total: 0,
    topCount: 1000,
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
    if (searchValue.value) {
      if (queryForm.list.length == 0) return
      queryForm.list = queryForm.list.filter((item: any) => {
        return Object.keys(item).some((td: any) => {
          return item[td] == searchValue.value
        })
      })
      setTimeout(() => {
        queryForm.loading = false
      }, 400)
    } else {
      const { topCount, lastTime } = queryForm
      const { clientIp, clientPort, serverIp, serverPort } = props.currentRow
      const query = { clientIp, clientPort, serverIp, serverPort, timeRange: lastTime, top: topCount }
      try {
        const { data } = await getPacketDecodeApi(query)
        queryForm.list = data
      } finally {
        queryForm.loading = false
      }
    }
  }

  async function download(row: any) {
    const { flowBeginTimeNs } = props.currentRow
    try {
      const { clientIp, clientPort, serverIp, serverPort = '' } = row
      downloadLogPacket({ clientIp, clientPort, serverIp, serverPort, requestTimeNs: flowBeginTimeNs }, 1)
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
    if (props.currentRow) {
      const { flowBeginTimeNs } = props.currentRow
      const newDate = dayjs(formatNstime(flowBeginTimeNs))
      const startDate = newDate.subtract(30, 'second')
      const endDate = newDate.add(30, 'second')
      timeDate.value = [formatTime(startDate), formatTime(endDate)]
    }
  })
  watchEffect(() => {
    if (timeDate.value) {
      const [startDate, endDate] = timeDate.value
      queryForm.lastTime = `${formatTime(new Date(startDate).getTime())} - ${formatTime(new Date(endDate).getTime())}`
    }
  })
  onMounted(() => {
    visible.value = props.showInfoDialog
    getData()
  })

  const emit = defineEmits<{
    (e: 'on-closeEvent'): void
  }>()

  // 关闭
  const handleClose = () => {
    emit('on-closeEvent')
  }
</script>

<script lang="ts">
  export default {
    name: 'InfoDialog',
  }
</script>
<template>
  <div class="info-dialog">
    <el-dialog v-model="visible" :before-close="handleClose" :fullscreen="true" title="详情">
      <el-row :gutter="20">
        <el-col :span="8">
          <vab-date-time-picker v-model="timeDate" />
        </el-col>
        <el-col :span="8">
          <el-select v-model="queryForm.topCount">
            <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-col>
        <el-col :span="6"><el-input v-model="searchValue" clearable @keyup.enter="getData" /></el-col>
        <el-col :span="2">
          <el-button
            :loading="queryForm.loading"
            style="float: right; margin-left: 20px"
            type="primary"
            @click="getData"
          >
            检索
          </el-button>
        </el-col>
      </el-row>
      <div style="height: 80vh">
        <el-table
          v-loading="queryForm.loading"
          :data="queryForm.list"
          style="margin-top: 20px; width: 100%; height: 100%"
          @cell-contextmenu="useTableCopy"
        >
          <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="70" />
          <el-table-column align="center" label="源IP" prop="clientIp" />
          <el-table-column align="center" label="源端口" prop="clientPort" />
          <el-table-column align="center" label="目的IP" prop="serverIp" />
          <el-table-column align="center" label="目的端口" prop="serverPort" />
          <el-table-column align="center" fixed="right" label="操作" width="240">
            <template #default="{ row }">
              <el-button class="row_action" size="small" @click="decodingHandle(row)">PCAP解析</el-button>
              <el-button class="row_action" size="small" @click="download(row)">PCAP保存</el-button>
            </template>
          </el-table-column>
        </el-table>
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
            <decoding-stream :current-data="currentData" :time-date="queryForm.lastTime" />
          </el-tab-pane>
          <el-tab-pane label="数据包" lazy name="decodingPacket">
            <decoding-packet :current-data="currentData" :time-date="queryForm.lastTime" />
          </el-tab-pane>
        </el-tabs>
      </vab-dialog>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss"></style>
