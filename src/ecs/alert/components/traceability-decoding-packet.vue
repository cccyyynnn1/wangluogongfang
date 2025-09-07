<script lang="ts">
  export default {
    name: 'DecodingPacket',
  }
</script>

<script setup lang="ts">
  import { PacketDecodeQuery, PacketItem } from '@/types'
  import { getPacketDecodeListApi, getPacketDecodeDetailApi } from '@/api-ecs/alert'
  import { hexStringToArrayBuffer } from '@/utils/text'
  import { ElTable } from 'element-plus'
  const singleTableRef = ref<InstanceType<typeof ElTable>>()
  const props = defineProps<{
    currentData?: PacketDecodeQuery
    timeDate: string
  }>()
  const splitSize = 22
  const listLoading = ref(false)
  const option = [
    { label: 'TOP100', value: 100 },
    { label: 'TOP500', value: 500 },
    { label: 'TOP1000', value: 1000 },
    { label: 'TOP2000', value: 2000 },
    { label: 'TOP5000', value: 5000 },
  ]
  const queryForm = reactive({
    page: 1,
    limit: 10,
    topCount: 100,
    packeList: [] as PacketItem[],
    bytes: [] as string[][],
  })

  const currentRow = ref()
  const tableRowClassName = ({ row }: { row: any }) => {
    let val = ''
    if (row.protocol == 'HTTP') {
      val = 'http_row'
    }
    console.log(val)
    return val
  }
  // 获取表格序号
  const curIndex = computed(() => (queryForm.page - 1) * queryForm.limit + 1)
  const topCountChange = (topCount: number) => {
    getData()
  }
  const handleCurrentChange = (val: PacketItem | undefined) => {
    currentRow.value = val
    getByte(val)
  }
  const formatText = (txt: string[]) => {
    const textDecoder = new TextDecoder('utf-8')
    const str = txt.join('')
    return textDecoder.decode(hexStringToArrayBuffer(str))
  }
  const formatHex = (hexString: string) => {
    const stringArr = hexString.split(' ')
    const result = []
    while (stringArr.length > splitSize) {
      result.push(stringArr.splice(0, splitSize))
    }
    result.push(stringArr)
    return result
  }
  const getByte = async (data: any) => {
    const { data: byte } = await getPacketDecodeDetailApi(data.id)
    if (byte.data && byte.data?.hex) {
      queryForm.bytes = formatHex(byte.data?.hex[0] || '')
    }
    listLoading.value = false
  }
  const getData = async () => {
    if (!props.currentData) return
    try {
      listLoading.value = true
      const { clientPort, serverIp, serverPort, clientIp, probeId = '' } = props.currentData
      const { data } = await getPacketDecodeListApi({
        top: queryForm.topCount,
        query: {
          objectList: [
            {
              serverIp,
              probeId,
              serverPort: serverPort?.toString(),
              clientIp,
              clientPort: clientPort?.toString(),
            },
          ],
          timeStep: 'minuteStep',
          timeRange: props.timeDate,
        },
      })
      queryForm.packeList = data || []
      if (queryForm.packeList?.length > 0) singleTableRef.value?.setCurrentRow(data[0])
    } finally {
      listLoading.value = false
    }
  }
  onMounted(() => {
    getData()
  })
</script>

<template>
  <div class="DecodingPacketBox">
    <el-row :gutter="20">
      <el-col :offset="8" :span="16" style="text-align: end">
        <el-select v-model="queryForm.topCount" @change="topCountChange">
          <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-button :auto-insert-space="false" style="margin-left: 15px" type="danger" @click="getData">
          重新解码
        </el-button>
      </el-col>
    </el-row>
    <el-table
      ref="singleTableRef"
      v-loading="listLoading"
      class="my-table"
      :data="queryForm.packeList || [Array.of(10)]"
      element-loading-text="Loading..."
      highlight-current-row
      :row-class-name="tableRowClassName"
      style="width: 100%; margin-top: 20px"
      @current-change="handleCurrentChange"
    >
      <el-table-column :index="curIndex" label="序号" type="index" width="55px" />
      <el-table-column
        :formatter="({ dateTimeStr }) => dateTimeStr.split('.')[0]"
        label="日期"
        prop="dateTimeStr"
        show-overflow-tooltip
        width="180px"
      />
      <el-table-column
        :formatter="
          ({ diffTime }) => {
            return diffTime == '-' ? diffTime : (+diffTime).toFixed(2)
          }
        "
        label="时间差"
        prop="diffTime"
        show-overflow-tooltip
        width="100px"
      />
      <el-table-column
        :formatter="
          ({ relativeTime }) => {
            return relativeTime == '-' ? relativeTime : (+relativeTime).toFixed(2)
          }
        "
        label="相对时间"
        prop="relativeTime"
        show-overflow-tooltip
        width="100px"
      />
      <el-table-column label="源地址" prop="source" show-overflow-tooltip width="120px" />
      <el-table-column label="源端口" prop="sourcePort" show-overflow-tooltip width="90px" />
      <el-table-column label="目的IP" prop="target" show-overflow-tooltip width="120px" />
      <el-table-column label="目的端口" prop="targetPort" show-overflow-tooltip width="90px" />
      <el-table-column label="协议" prop="protocol" show-overflow-tooltip width="80px" />
      <el-table-column label="大小" prop="lengthStr" show-overflow-tooltip />
      <el-table-column label="信息" prop="info" show-overflow-tooltip />
      <template #empty>
        <el-empty class="vab-data-empty" description="暂无数据" />
      </template>
    </el-table>
    <ul class="hex2string">
      <li v-for="(byte, index) in queryForm.bytes" :key="index">
        <div class="index">{{ (index * splitSize).toString(16).padStart(8, '0') }}</div>
        <div class="hextxt">
          <span v-for="(item, index) in byte" :key="index" :span="1">{{ item }}</span>
        </div>
        <div class="strtxt">{{ formatText(byte) }}</div>
      </li>
    </ul>
  </div>
</template>
<style scoped lang="scss">
  .my-table {
    :deep(.http_row) {
      // tr {
      //   background-color: initial !important;
      // }
      background-color: #fff8ce !important;
    }
  }
  .hex2string {
    padding: 0;
    margin-top: 15px;
    position: sticky;
    bottom: 0;
    background: #ffffff;
    z-index: 9999;
    li {
      display: flex;
      padding: 5px 0;
      .index {
        width: 120px;
        text-align: center;
        margin-right: 20px;
      }
      .hextxt {
        flex: 1;
        display: flex;
        span {
          display: inline-block;
          flex-basis: 4.5%;
          text-align: center;
        }
      }
      .strtxt {
        width: 350px;
        text-align: center;
      }
    }
  }
  .DecodingPacketBox {
    height: calc(100vh - 150px);
    overflow: hidden auto;
    :deep() {
      .el-table__body-wrapper {
        height: max-content;
        min-height: auto;
        max-height: max-content;
      }
    }
  }
</style>
