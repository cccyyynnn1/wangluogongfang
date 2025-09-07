<script setup lang="ts">
  import { MonitoringItem, MonitoringType } from '~/src/types'

  import { VIEWOBJ } from '@/data/constant'

  import dayjs from 'dayjs'

  const viewRef = ref()

  const props = defineProps<{
    currentItem: MonitoringType
  }>()

  const formSave = reactive({
    id: 0,
    uuid: '',
    name: '',
    timeQuantum: '1hours',
    switchValue: false,
    timeDuration: '1.5min',
    timeDate: undefined,
  })

  const currentItem = ref() // 当前页面数据

  const containerData = ref<MonitoringItem[]>([]) // 要渲染的数组

  const toolData = reactive({
    timeDuration: '30s',
    switchValue: false,
    isFullscreen: false,
    showEditPage: false,
  })

  const timeDuratioOptions = [
    // {
    //   value: '30s',
    //   label: '30秒',
    // },
    {
      value: '1.5min',
      label: '90秒',
    },
    {
      value: '5min',
      label: '5分钟',
    },
    {
      value: '10min',
      label: '10分钟',
    },
  ]

  const timeSetting = ref({
    endTime: '',
    startTime: '',
  })

  const formatTimeDate = () => {
    if (formSave.timeDate) {
      // @ts-ignore
      const [startDate, endDate] = formSave.timeDate
      timeSetting.value.endTime = dayjs(endDate).format('YYYY-MM-DD HH:mm:ss')
      timeSetting.value.startTime = dayjs(startDate).format('YYYY-MM-DD HH:mm:ss')
    }
  }

  // 得到自定义时间
  watchEffect(() => {
    formatTimeDate()
  })

  const formatDayDate = () => {
    if (formSave.timeQuantum === 'user-defined') {
      formatTimeDate()
      return
    }
    const timeDate = dayjs()
    const endDate = timeDate.format('YYYY-MM-DD HH:mm:ss')
    const cur = timeDate.startOf('day')
    let startDate = null
    switch (formSave.timeQuantum) {
      case '1hours':
        startDate = timeDate.subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '3hours':
        startDate = timeDate.subtract(3, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '6hours':
        startDate = timeDate.subtract(6, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case '24hours':
        startDate = timeDate.subtract(24, 'hour').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'today':
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'yesterday':
        startDate = cur.subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-three-days':
        startDate = timeDate.subtract(3, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-week':
        startDate = timeDate.subtract(1, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      case 'last-two-week':
        startDate = timeDate.subtract(2, 'week').format('YYYY-MM-DD HH:mm:ss')
        break
      default:
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
    }
    timeSetting.value.endTime =
      formSave.timeQuantum === 'yesterday'
        ? timeDate.endOf('day').subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss')
        : endDate
    timeSetting.value.startTime = startDate
  }

  // 得到时间
  watchEffect(() => {
    formatDayDate()
  })

  // 自动刷新
  let timer: any
  const autoReflash = (time: string) => {
    if (formSave.switchValue) {
      const res = changeTime(time)
      timer = setInterval(() => {
        formatTimeDate()
        formatDayDate()
        viewRef.value.forEach((item: any) => {
          item.initData()
        })
      }, res)
    } else {
      clearInterval(timer)
    }
  }

  watch(
    () => formSave.switchValue,
    () => {
      autoReflash(formSave.timeDuration)
    },
    {
      immediate: true,
    }
  )

  onUnmounted(() => {
    clearInterval(timer)
  })

  const changeTime = (time: string) => {
    let timer = 0
    switch (time) {
      case '30s':
        timer = 30 * 1000
        break
      case '1min':
        timer = 60 * 1000
        break
      case '1.5min':
        timer = 90 * 1000
        break
      case '5min':
        timer = 5 * 60 * 1000
        break
      case '10min':
        timer = 10 * 60 * 1000
        break
    }
    return timer
  }

  const emit = defineEmits<{
    (e: 'on-back'): void
    (e: 'on-todesign'): void
    (e: 'on-update-event', val: MonitoringType): void
  }>()

  // 返回
  const goBack = () => {
    emit('on-back')
  }

  onMounted(() => {
    initData()
  })

  const initData = () => {
    currentItem.value = JSON.parse(JSON.stringify(props.currentItem))
    Object.keys(formSave).forEach((item) => {
      if (item !== 'timeDate') {
        // @ts-ignore
        formSave[item] = currentItem.value[item]
      } else {
        formSave[item] = currentItem.value[item] ? currentItem.value[item] : undefined
      }
    })
    containerData.value = JSON.parse(currentItem.value.list)
  }

  // 全屏
  const clickFullScreen = () => {
    toolData.isFullscreen = !toolData.isFullscreen
  }

  //  去保存编辑页
  const handleEdit = () => {
    emit('on-todesign')
  }
</script>

<script lang="ts">
  export default {
    name: 'DetailMonitoring',
  }
</script>
<template>
  <div class="detail-monitoring" :class="{ 'vab-fullscreen': toolData.isFullscreen }">
    <el-row class="toolbar" :gutter="10">
      <el-col :span="16">
        <!-- <el-button type="primary" @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button> -->
        <div class="el-toolbar-title" style="font-weight: 700; font-size: 16px">
          {{ currentItem && currentItem.name }}
        </div>
      </el-col>
      <el-col :span="8">
        <div class="toolbar-right">
          <el-select v-model="formSave.timeDuration" style="width: 100px">
            <el-option v-for="item in timeDuratioOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-switch v-model="formSave.switchValue" active-text="自动刷新" />
          <vab-icon
            class="icon"
            :icon="toolData.isFullscreen ? 'fullscreen-exit-fill' : 'fullscreen-fill'"
            @click="clickFullScreen"
          />
          <!-- </el-tooltip> -->
          <el-icon class="icon" @click="handleEdit"><EditPen /></el-icon>
        </div>
      </el-col>
    </el-row>
    <div class="content" :style="{ minHeight: '400px', height: containerData[0]?.style?.long + 'px' }">
      <template v-if="containerData.length > 0">
        <div v-for="item in containerData" :key="item.uuid" class="view" :style="{ ...item.style }">
          <div class="top-bar">
            <div class="top-bar-left">{{ item.data.name }}</div>
          </div>
          <div :class="{ viewContent: true, activeContent: item.hasError }">
            <component :is="VIEWOBJ[item.key]" ref="viewRef" :current-item="item" :time-setting="timeSetting" />
          </div>
          <div v-if="item.hasError" class="mark"></div>
          <div v-if="item.hasError" class="mark-boder">
            <div class="err-icon"><span>!</span></div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .detail-monitoring {
    background-color: #fff;
  }
  .toolbar {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 15px;

    .el-toolbar-title {
      max-width: 400px;
      min-width: 100px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      word-break: break-all;
      word-wrap: break-word;
    }
    .toolbar-right {
      width: 300px;
      display: flex;
      align-items: center;
      float: right;
      justify-content: space-around;
      .icon {
        &:hover {
          cursor: pointer;
        }
      }
    }
  }
  .content {
    position: relative;
    z-index: 0;
    // border: 1px solid #eee;
    // height: min;
    min-height: 400px;
    max-height: calc(100vh - 180px);
    overflow-y: auto;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    .view {
      position: absolute;
      z-index: 0;
      height: 200px;
      width: 50%;

      overflow: hidden;
      box-shadow: var(--el-box-shadow-light);
      border-radius: var(--el-card-border-radius);
      // border: 1px solid var(--el-card-border-color);
      background-color: #fff;
      border: 1px solid var(--el-border-color);
      color: var(--el-text-color-primary);
      padding: 5px 15px 15px 15px;

      .top-bar {
        height: 40px;
        padding: 10px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        // border-bottom: 1px solid var(--el-border-color);

        .icon {
          &:hover {
            cursor: pointer;
          }
        }
      }
      .viewContent {
        height: calc(100% - 40px);
        :deep() {
          .el-table {
            height: 100%;
            border: 0px solid #fff;
            border-bottom: 1px solid var(--el-border-color);
          }
        }
      }
      .mark,
      .mark-boder {
        opacity: 0.01;
        height: 100%;
        width: 100%;
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
        background-color: #ff3636;
      }
      .mark-boder {
        border: 1px solid #ff3636;
        background-color: inherit;
        opacity: 1;
        .err-icon {
          height: 24px;
          width: 24px;
          border-radius: 4px 4px 0 0;
          background-color: #ff3636;
          position: absolute;
          bottom: 0;
          right: 0;
          text-align: center;
          color: #fff;
        }
      }
      .activeContent {
        :deep() {
          .el-table {
            border: 1px solid #ff3636;
            border-top: 0px solid #fff;
          }
        }
      }
    }
  }
</style>
