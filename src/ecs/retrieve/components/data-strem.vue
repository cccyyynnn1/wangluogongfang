<script lang="ts">
  export default {
    name: 'DataStrem', //
  }
</script>

<script setup lang="ts">
  import { PacketsMagic } from '@/utils/magic'
  import { hexStringToArrayBuffer } from '@/utils/text'
  import { getPacketDecodeFlowDecodeApi } from '~/src/api-ecs/alert'
  import { formatNstime, formatTime } from '@/utils/time'
  import dayjs from 'dayjs'
  const decodingTime = ref('')

  onMounted(() => {
    const { startTimeNs } = props.infoVal
    const timeDate = dayjs(formatNstime(startTimeNs))
    const startDate = timeDate.subtract(1, 'minute')
    const endDate = timeDate.add(1, 'minute')
    decodingTime.value = `${formatTime(startDate)} - ${formatTime(endDate)}`
    getData()
  })

  const props = defineProps<{
    infoVal: any
  }>()
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
    getData()
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
  const getData = async () => {
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
    getData,
  })
</script>

<template>
  <div
    class="hostTable_title"
    style="display: flex; justify-content: space-between; align-items: center; display: none"
  >
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
  <div v-loading="loading" class="requests-and-responses" style="display: none">
    <div v-for="(stream, index) in showMagic ? magicStreamList : streamList" :key="index" class="stream">
      <h4>{{ stream.info }}</h4>
      <ul :class="{ isDst: stream.srcToDst }">
        <li v-for="(item, $index) in stream.hexs" :key="$index">{{ item }}</li>
      </ul>
    </div>
    <el-empty v-if="streamList.length === 0" description="暂无数据" />
  </div>
</template>

<style scoped lang="scss">
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
  }
</style>
