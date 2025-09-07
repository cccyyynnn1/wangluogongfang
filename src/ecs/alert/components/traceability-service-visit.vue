<script setup lang="ts">
  import { getServiceAccessApi } from '@/api-ecs/assets'
  import { formatNstime, formatTime } from '~/src/utils/time'
  import { downloadLogPacket } from '@/utils/download'
  import dayjs from 'dayjs'
  import { defaultTime } from '@/data/constant'

  /* alertData 和 timeRange 必须有一个，alertData：{startTimeNs} */
  const props = defineProps<{
    alertData?: any
    timeRange?: [string, string]
    nodeData: any
  }>()
  // 自定义时间
  const timeDate = ref()

  // 表单数据
  const queryForm = reactive({
    page: 1,
    limit: 10,
    total: 0,
    lastTime: '',
    startTime: '',
    endTime: '',
    list: [] as any,
    loading: true,
  })

  const getData = async () => {
    queryForm.loading = true
    const { page, limit, lastTime, startTime, endTime } = queryForm
    const { type, serverIp, snat, ip } = props.nodeData
    const query = type ? { serverIp: serverIp, clientIp: snat, lastTime } : { lastTime, ipAddr: ip }
    const str = snat ? `源IP = "${snat}"` : ''
    const str1 = snat && serverIp ? ' and ' : ''
    const str2 = serverIp ? `目的IP = "${serverIp}"` : ''
    const searchSql = type ? str + str1 + str2 : `源IP = "${ip}" or 目的IP = "${ip}"`
    try {
      const {
        data: { resList, total },
      } = await getServiceAccessApi({
        indexType: 11,
        // tid: id,
        startTime: startTime,
        searchSql: searchSql,
        endTime: endTime,
        pageNum: queryForm.page,
        pageSize: queryForm.limit,
      })
      queryForm.list = resList
      queryForm.total = total
    } catch (error) {
      queryForm.list = []
      queryForm.total = 0
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
  // 获取表格序号
  const curIndex = computed(() => (queryForm.page - 1) * queryForm.limit + 1)
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
      queryForm.startTime = `${formatTime(new Date(startDate).getTime())}`
      queryForm.endTime = `${formatTime(new Date(endDate).getTime())}`
    }
  })
  onMounted(() => {
    getData()
  })
</script>

<script lang="ts">
  export default {
    name: 'TraceabilityServiceVisit',
  }
</script>

<template>
  <el-row :gutter="20">
    <el-col :span="20">
      <vab-date-time-picker v-model="timeDate" />
    </el-col>
    <el-col :span="4">
      <el-button :loading="queryForm.loading" style="float: right; margin-left: 20px" type="primary" @click="getData">
        检索
      </el-button>
    </el-col>
  </el-row>
  <el-table v-loading="queryForm.loading" :data="queryForm.list" style="margin-top: 20px; width: 100%">
    <el-table-column align="center" :index="(index) => curIndex + index" label="序号" type="index" width="70" />
    <el-table-column align="center" label="源IP" prop="clientIp" width="180" />
    <el-table-column align="center" label="目的IP" prop="serverIp" width="180" />
    <el-table-column align="center" label="目的端口" prop="serverPort" />
    <el-table-column align="center" label="协议" prop="appProtocolStr" />
    <el-table-column align="center" fixed="right" label="操作" width="120">
      <template #default="{ row }">
        <el-button class="row_action" size="small" @click="download(row)">PCAP保存</el-button>
      </template>
    </el-table-column>
    <template #empty>
      <el-empty class="vab-data-empty" description="暂无数据" />
    </template>
  </el-table>
  <el-pagination
    v-model:current-page="queryForm.page"
    v-model:page-size="queryForm.limit"
    background
    class="site_pagination"
    hide-on-single-page
    layout=" sizes, prev, pager, next, jumper"
    :page-sizes="[10, 20, 30]"
    :total="queryForm.total"
    @current-change="getData"
    @size-change="getData"
  />
</template>

<style scoped lang="scss">
  .timer {
    display: flex;
    align-items: center;
  }

  .m-2 {
    margin-right: 20px;
  }
  // :deep() {
  //   .el-table__body {
  //     height: 445px;
  //     overflow-y: auto;
  //   }
  // }
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
