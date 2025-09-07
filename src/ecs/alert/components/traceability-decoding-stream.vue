<script lang="ts">
  export default {
    name: 'DecodingStream',
  }
</script>

<script setup lang="ts">
  import { hexStringToArrayBuffer } from '@/utils/text'
  import { PacketDecodeQuery } from '@/types'
  import { getPacketDecodeFlowDecodeApi } from '@/api-ecs/alert'
  import { PacketsMagic } from '@/utils/magic'
  const magicIcon = require('@/assets/mofabang.svg')
  const props = defineProps<{
    currentData?: PacketDecodeQuery
    timeDate: string
    protocol?: string
  }>()
  const loading = ref(false)
  const formatType = ref('utf-8')
  const topCount = ref(50)
  const protocol = ref(props.protocol || 'TCP')
  const showMagic = ref(false)
  const option = [
    { label: 'TOP50', value: 50 },
    { label: 'TOP100', value: 100 },
    { label: 'TOP200', value: 200 },
    { label: 'TOP500', value: 500 },
    { label: 'TOP1000', value: 1000 },
  ]
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

  const protocols = ['TCP', 'UDP']

  let STREAM_DATA = [] as { hexs: string[]; info: string; srcToDst: boolean }[]
  const streamList = ref<{ hexs: string[]; info: string; srcToDst: boolean }[]>([])
  const magicStreamList = ref<{ hexs: string[]; info: string; srcToDst: boolean }[]>([])

  const getData = async () => {
    showMagic.value = false
    if (!props.currentData) return
    try {
      loading.value = true
      const { probeId = '', clientPort, serverIp, serverPort, clientIp } = props.currentData
      const { data } = await getPacketDecodeFlowDecodeApi({
        top: topCount.value,
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
          timeRange: props.timeDate,
        },
        flowType: protocol.value,
      })
      STREAM_DATA = data || []
      toFormat(STREAM_DATA)
    } finally {
      loading.value = false
    }
  }

  const formatText = (txt: string) => {
    const textDecoder = new TextDecoder(formatType.value)
    const str = txt.replace(/\s*/g, '')
    return textDecoder.decode(hexStringToArrayBuffer(str))
  }
  const topCountChange = () => {
    getData()
  }
  const protocolChange = (protocolVal: string) => {
    protocol.value = protocolVal
    getData()
  }
  const toFormat = (streamData: any[]) => {
    streamList.value = streamData.map((stream) => ({
      ...stream,
      hexs: stream.hexs.map((hex: string) => formatText(hex)),
    }))
  }
  const toMagic = () => {
    if (showMagic.value) return (showMagic.value = false)
    magicStreamList.value = streamList.value.map((stream) => ({
      ...stream,
      hexs: stream.hexs.map((str: string) => PacketsMagic(str)),
    }))
    showMagic.value = true
  }
  watch(
    () => formatType.value,
    () => {
      showMagic.value = false
      toFormat(STREAM_DATA)
    }
  )
  watch(
    () => props.currentData,
    () => {
      getData()
    },
    {
      immediate: true,
      deep: true,
    }
  )
</script>

<template>
  <el-row :gutter="20">
    <el-col :span="14">
      <el-button-group>
        <el-button
          v-for="my_type of typeList"
          :key="my_type.value"
          :type="formatType === my_type.value ? 'primary' : 'default'"
          @click="formatType = my_type.value"
        >
          {{ my_type.lable }}
        </el-button>
      </el-button-group>
      <el-button-group style="margin-left: 20px">
        <el-button
          v-for="my_protocol of protocols"
          :key="my_protocol"
          :type="protocol === my_protocol ? 'primary' : 'default'"
          @click="protocolChange(my_protocol)"
        >
          {{ my_protocol }}
        </el-button>
      </el-button-group>
    </el-col>
    <el-col :span="10" style="text-align: end">
      <el-button
        :auto-insert-space="false"
        style="margin-left: 15px; border: 0; background-color: transparent"
        @click="toMagic()"
      >
        <el-image :src="magicIcon" style="width: 16px; height: 16px; cursor: pointer" />
      </el-button>
      <el-select v-model="topCount" style="width: 120px" @change="topCountChange">
        <el-option v-for="item in option" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-button :auto-insert-space="false" style="margin-left: 15px" type="danger" @click="getData">
        重新加载
      </el-button>
    </el-col>
  </el-row>
  <div
    v-loading="loading"
    style="max-height: calc(100vh - 202px); min-height: 300px; overflow-y: auto; margin-top: 20px"
  >
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
  .stream {
    padding: 0 20px;
    ul {
      list-style: none;
      padding: 0;
      &.isDst li {
        color: #38a046;
      }
      li {
        word-wrap: break-word;
        white-space: break-spaces;
        color: var(--el-color-primary);
        line-height: 30px;
      }
    }
  }
</style>
