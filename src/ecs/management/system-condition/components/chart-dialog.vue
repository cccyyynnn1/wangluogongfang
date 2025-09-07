<script lang="ts">
  export default {
    name: 'ChartDialog',
  }
</script>
<script setup lang="ts">
  // import TimeConfigBar from './time-config-bar.vue'

  import ChartsTotalY from './charts-total-y.vue'

  import ChartsTotalL from './charts-total-l.vue'

  import ChartsTotalC from './charts-total-c.vue'

  import ChartsTotalBL from './charts-total-b-l.vue'

  import CpuUsageRatesL from './cpu-usage-rates-l.vue'

  import CpuUsageRatesBL from './cpu-usage-rates-b-l.vue'
  import ChartDiskList from './chart-disk-list.vue'
  const props = defineProps<{
    dialogVisible: boolean
    mode: string
    resData: any
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent'): void
  }>()

  let resData = ref({
    cpuRate: {},
    diskTotal: {},
    disk: {},
    diskIo: {},
    total: {
      ldle: 0,
      uptime: 0,
    },
    memory: {},
    netRate: {},
  })

  type typeType = 'cpuRate' | 'diskTotal' | 'disk' | 'diskIo' | 'memory' | 'netRate' | 'total'

  let type: typeType = 'total'

  const visible = ref(false)

  let current = ChartsTotalC
  const title = ref('')
  onMounted(() => {
    visible.value = props.dialogVisible
    resData.value = props.resData
    initComponent()
  })

  const handleClose = () => {
    emit('on-closeEvent')
  }

  const initComponent = () => {
    switch (props.mode) {
      case 'ChartsTotalY':
        // @ts-ignore
        current = ChartsTotalY
        type = 'total'
        title.value = '系统总体状态'
        break
      case 'ChartsTotalL':
        // @ts-ignore
        current = ChartsTotalL
        type = 'memory'
        title.value = '内存使用率'
        break
      case 'ChartsTotalC':
        current = ChartDiskList as any
        type = 'disk'
        title.value = '磁盘使用详情'
        break
      case 'ChartsTotalBL':
        current = ChartsTotalBL
        type = 'netRate'
        title.value = '网卡流量'
        break
      case 'CpuUsageRatesL':
        // @ts-ignore
        current = CpuUsageRatesL
        type = 'cpuRate'
        title.value = 'CPU使用率'
        break
      case 'CpuUsageRatesBL':
        // @ts-ignore
        current = CpuUsageRatesBL
        type = 'diskIo'
        title.value = 'IO使用率'
        break
    }
  }
</script>

<template>
  <div v-if="visible">
    <el-dialog v-model="visible" :before-close="handleClose" title="Tips" width="60%">
      <!-- <template #header="{ titleId, titleClass }">
        <div class="my-header">
          <div :id="titleId" :class="titleClass">
            <time-config-bar :is-true="false" />
          </div>
        </div>
      </template> -->
      <template #header="{ titleId, titleClass }">
        <div class="my-header">
          <div :id="titleId" :class="titleClass">{{ title }}</div>
        </div>
      </template>
      <div class="content">
        <component :is="current" :total="resData[type]" />
      </div>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .my-header {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }

  .content {
    height: 600px;
  }

  :deep(.top-bar) {
    display: none;
  }
</style>
