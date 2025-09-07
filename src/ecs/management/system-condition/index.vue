<script lang="ts">
  export default {
    name: 'SystemCondition',
  }
</script>

<script setup lang="ts">
  // import type { InjectionKey } from 'vue'

  // import { Ref } from 'vue'

  import TimeConfigBar from './components/time-config-bar.vue'

  import ChartsTotalY from './components/charts-total-y.vue'

  import ChartsTotalL from './components/charts-total-l.vue'

  import ChartsTotalC from './components/charts-total-c.vue'

  import ChartsTotalBL from './components/charts-total-b-l.vue'

  import CpuUsageRatesL from './components/cpu-usage-rates-l.vue'

  import CpuUsageRatesBL from './components/cpu-usage-rates-b-l.vue'

  import ChartDialog from './components/chart-dialog.vue'

  import { proxyNet } from '@/config/index'

  // import { useWebSocket } from '@vueuse/core'

  import { getToken } from '@/utils/token'

  const activeName = ref('management') // tabs选中项

  const showDialog = ref(false)

  const mode = ref('')

  let resData = ref({
    cpuRate: {},
    disk: {},
    diskTotal: {},
    diskIo: {},
    total: {
      ldle: 0,
      uptime: 0,
    },
    memory: {},
    netRate: {},
  })

  provide('activekey', showDialog)

  provide('modekey', mode)

  const closeEvent = () => {
    showDialog.value = false
  }

  const isShow = ref(true) // 是否加载

  const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname

  const cookie = getToken() as string

  const protocol = window.location.protocol == 'http:' ? 'wss' : 'wss'

  const { status, data, send, open, close } = useWebSocket(
    `${protocol}://${url}/v3/ecsPlatform/websocket/sysStatus?EcsSessionId=${cookie}`,
    // `${protocol}://${url}/v3/ecsPlatform/websocket/sysStatus`,
    {
      autoReconnect: {
        retries: 1,
        delay: 1000,
        onFailed() {
          ElMessage({ message: 'WebSocket链接失败', type: 'error' })
        },
      },
    }
  )

  const timer = ref()

  watchEffect(() => {
    if (data.value && data.value !== 'ws服务连接成功!') {
      const res = JSON.parse(data.value)
      if (res.node?.value) {
        const data = JSON.parse(res.node?.value)
        for (const key in resData.value) {
          if (key == 'total') {
            resData.value[key].ldle = data.ldle
            resData.value[key].uptime = data.uptime
          } else {
            // @ts-ignore
            resData.value[key] = data[key]
          }
        }
      }
    }
  })
  onBeforeRouteLeave((to, from, next) => {
    close()
    clearInterval(timer.value)
    next()
  })
</script>

<template>
  <div class="system-condition-container">
    <div v-if="isShow">
      <div class="header">
        <h3>本机系统状态</h3>
        <time-config-bar :is-true="false" />
      </div>

      <!-- 图表 -->
      <div class="chart-content">
        <div class="top">
          <div class="my-chart">
            <charts-total-y :total="resData.total" />
          </div>
          <div class="my-chart">
            <cpu-usage-rates-l :total="resData.cpuRate" />
          </div>
          <div class="my-chart">
            <charts-total-l :total="resData.memory" />
          </div>
          <div class="my-chart">
            <charts-total-c :total="resData.diskTotal" />
          </div>
        </div>
        <div class="bottom">
          <div class="my-chart">
            <charts-total-b-l :total="resData.netRate" />
          </div>
          <div class="my-chart">
            <cpu-usage-rates-b-l :total="resData.diskIo" />
          </div>
        </div>
      </div>
    </div>
    <chart-dialog
      v-if="showDialog"
      :dialog-visible="showDialog"
      :mode="mode"
      :res-data="resData"
      @on-close-event="closeEvent"
    />
  </div>
</template>

<style lang="scss" scoped>
  .system-condition-container {
    .header {
      display: flex;
      align-items: center;
      h3 {
        margin: 0 10px 0 0;
      }
    }
    .chart-content {
      height: 610px;
      width: 100%;

      .top,
      .bottom {
        display: flex;
        width: 100%;
        margin-top: 20px;
      }
    }

    .my-chart {
      width: 23.5%;
      margin-right: 2%;
      height: 295px;
      border: 1px solid #dcdfe6;

      &:nth-child(4) {
        margin-right: 0;
      }
    }
  }
</style>
