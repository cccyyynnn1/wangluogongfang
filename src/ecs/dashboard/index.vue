<script lang="ts">
  export default {
    name: 'DashboardScreen',
  }
  type CountType = {
    count?: number //总数
    tradeNum?: number // 上升趋势
  }
  type LineChartType = {
    count?: number[]
    date?: string[]
  }

  type scrrenDataType = {
    // 地图
    mapData: {
      domestic: {
        nodes?: any[]
        links?: any[]
      }
      foreign: {
        nodes?: any[]
        links?: any[]
      }
    }
    // 异常访问占比
    warnRequest: {
      normal?: number
      abnormal?: number
    }
    // 攻击类型
    threatType: {
      name: string
      value: number
    }[]
    // 异常访问
    abnormalWarn: LineChartType
    // 应用数据趋势
    index: LineChartType
    // 攻击趋势
    attackLine: LineChartType
    // 攻击源
    attackIp: {
      name: string
      value: number
    }[]
    // 资产总数
    topCount: {
      apiCount?: CountType
      assetCount?: CountType
      assetNotKnowCount?: CountType
      siteCount?: CountType
      warn_abonormal_count?: CountType
      warn_attack_count?: CountType
    }
    mapName: string
  }
</script>

<script setup lang="ts">
  import { ScreenView, ItemWrap } from '@/components/dashboard-screen'
  import DashboardLineChart from './components/line-chart.vue'
  import VisitProportion from './components/visit-proportion.vue'
  import Earth from './components/earth.vue'
  import AttackType from './components/attack-type.vue'
  import AttackSource from './components/attack-source.vue'
  import { useFavicon } from '@vueuse/core'
  import { useSettingsStore } from '@/store/modules/settings'
  import { useUserStore } from '@/store/modules/user'
  import { proxyNet } from '@/config/index'
  const $baseMessage: any = inject('$baseMessage')
  const settingsStore = useSettingsStore()
  const userStore = useUserStore()
  const { logo } = storeToRefs(settingsStore)
  const icon = useFavicon()
  const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname
  const dashboard_screen_url = `wss://${url}/v3/ecsPlatform/websocket/bigscreen?EcsSessionId=${userStore.token}`
  const { data, close } = useWebSocket(dashboard_screen_url, {
    autoReconnect: {
      retries: 1,
      delay: 1000,
      onFailed() {
        $baseMessage('WebSocket链接失败', 'error', 'vab-hey-message-error')
      },
    },
  })
  const loading = ref(true)
  const showFullScreen = ref(false)
  const scrrenData = reactive<scrrenDataType>({
    mapData: {
      domestic: {},
      foreign: {},
    },
    attackIp: [],
    warnRequest: {},
    topCount: {},
    threatType: [],
    abnormalWarn: {},
    index: {},
    attackLine: {},
    mapName: '态势感知大屏',
  })
  const header = ref()
  const scaleNum = ref()
  watchEffect(() => {
    if (data.value) {
      try {
        const {
          mapData,
          attackIp,
          warnRequest,
          topCount,
          threatType,
          abnormalWarn,
          index,
          attackLine,
          mapName,
        }: scrrenDataType = JSON.parse(data.value)
        scrrenData.mapData = mapData || []
        scrrenData.attackIp = attackIp || []
        scrrenData.threatType = threatType || []
        scrrenData.warnRequest = warnRequest || {}
        scrrenData.topCount = topCount || {}
        scrrenData.abnormalWarn = abnormalWarn || {}
        scrrenData.index = index || {}
        scrrenData.attackLine = attackLine || {}
        scrrenData.mapName = mapName || '态势感知大屏'

        const rootW = header.value.getBoundingClientRect().width
        const divElement = document.createElement('h4')
        divElement.textContent = scrrenData.mapName
        divElement.style.cssText = 'opacity: 0;width: max-content;'
        header.value.appendChild(divElement)
        const scale = divElement.offsetWidth > rootW ? divElement.offsetWidth / rootW : 1
        divElement.remove()
        scaleNum.value = scale === 1 ? 84 : 84 / scale

        loading.value = false
      } catch (error) {
        console.error(error)
      }
    }
  })
  const moveHandler = () => {
    if (showFullScreen.value) return
    showFullScreen.value = true
    setTimeout(() => {
      showFullScreen.value = false
    }, 2000)
  }
  onMounted(() => {
    icon.value = logo.value
  })
  onUnmounted(() => {
    close()
  })
</script>

<template>
  <screen-view
    v-loading="loading"
    :box-style="{
      background: '#10101A',
      overflow: 'hidden',
    }"
    :delay="500"
    element-loading-background="rgba(16, 16, 26,.8)"
    element-loading-text="连接中..."
    :full-screen="false"
    height="2160"
    width="3840"
  >
    <div class="dashboard-screen" @mousemove="moveHandler">
      <header ref="header" class="screen-header">
        <h4 :style="{ fontSize: `${scaleNum}px` }">{{ scrrenData.mapName }}</h4>
      </header>
      <article class="screen-container">
        <div class="contetn_left">
          <ItemWrap class="contetn_left-top contetn-item" title="应用数据量趋势">
            <dashboard-line-chart color="blue" :data="scrrenData.index" />
          </ItemWrap>
          <ItemWrap class="contetn_left-center contetn-item" title="异常访问趋势">
            <dashboard-line-chart color="yellow" :data="scrrenData.abnormalWarn" />
          </ItemWrap>
          <ItemWrap class="contetn-item" title="异常访问占比">
            <visit-proportion :data="scrrenData.warnRequest" />
          </ItemWrap>
        </div>
        <div class="contetn_center">
          <earth :count-data="scrrenData.topCount" :map-data="scrrenData.mapData" />
        </div>
        <div class="contetn_right">
          <ItemWrap class="contetn-item" title="攻击趋势">
            <dashboard-line-chart color="red" :data="scrrenData.attackLine" />
          </ItemWrap>
          <ItemWrap class="contetn-item" title="攻击源">
            <attack-source :data="scrrenData.attackIp" />
          </ItemWrap>
          <ItemWrap class="contetn-item" title="攻击类型">
            <attack-type :data="scrrenData.threatType" />
          </ItemWrap>
        </div>
      </article>
      <div class="full-screen-btn" :class="{ show: showFullScreen }">
        <el-button>
          <vab-full-screen />
        </el-button>
      </div>
    </div>
  </screen-view>
</template>

<style scoped lang="scss">
  .dashboard-screen {
    box-sizing: border-box;
    width: 100%;
    height: 100%;
    background: url('@/assets/dashboard_images/bodyBox.svg') no-repeat center;
    .full-screen-btn {
      position: fixed;
      top: 50%;
      left: 20px;
      display: flex;
      align-items: center;
      width: 80px;
      height: 2000px;
      cursor: pointer;
      opacity: 0;
      transition: all 0.3s ease-in-out;
      transform: translateY(-50%);
      &.show {
        opacity: 1;
      }
      & > .el-button {
        width: 80px;
        height: 80px;
        background-color: rgb(16, 16, 26);
        border-radius: 100%;
        i {
          font-size: 50px;
        }
      }

      &::after {
        position: relative;
        z-index: -10;
        display: block;
        width: 100px;
        height: 100vh;
        padding: 0 200px;
        content: ' ';
      }
    }
    .screen-header {
      width: 1668px;
      height: 153px;
      padding: 0 450px;
      margin: 30px auto 0;
      font-size: 84px;
      font-weight: bold;

      line-height: 96px;
      text-align: center;
      text-align-last: justify;
      letter-spacing: normal;
      background: url('@/assets/dashboard_images/header.svg') no-repeat center;
      h4 {
        margin: 0 auto;
        color: #c2d6f3;
        text-align: center;
        white-space: nowrap;
      }
    }
    .screen-container {
      display: flex;
      justify-content: space-between;
      width: 100%;
      min-height: calc(100% - 155px);
      padding: 40px 110px;
      .contetn_left,
      .contetn_right {
        position: relative;
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        flex-shrink: 0;
        justify-content: start;
        width: 776px;
      }
      .contetn_center {
        display: flex;
        flex: 1;
        flex-direction: column;
        justify-content: flex-start;
        margin: 0 54px;
      }

      .contetn-item {
        height: 580px;
        margin-bottom: 35px;
      }
    }
  }
</style>
