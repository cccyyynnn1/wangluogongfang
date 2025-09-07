<script lang="ts">
  export default {
    name: 'Earth',
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import numberFormatter from '@/utils/number'
  type CountType = {
    count?: number //总数
    tradeNum?: number // 上升趋势
  }
  const router = useRouter()
  const props = defineProps<{
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
    countData: {
      apiCount?: CountType
      assetCount?: CountType
      assetNotKnowCount?: CountType
      siteCount?: CountType
      warn_abonormal_count?: CountType
      warn_attack_count?: CountType
    }
  }>()
  const chartOption = reactive({
    geo: {
      map: 'world',
      zoom: 1.2,
      scaleLimit: {
        min: 1.2,
      },
      emphasis: {
        label: {
          show: false,
        },
        itemStyle: {
          color: 'rgba(37, 43, 61, .5)', //悬浮背景
        },
      },
      roam: true, //是否允许缩放
      itemStyle: {
        areaColor: '#122445',
        borderWidth: 1,
        borderColor: '#1FC6FF',
      },
    },
    series: [] as any[],
  })
  const goToUrl = (url: string, query?: any) => {
    const resolveRouter = router.resolve({ path: url, query })
    window.open(resolveRouter.href, '_blank')
  }
  watch(
    () => props.mapData,
    () => {
      const { domestic, foreign } = props.mapData
      const domesticLinks = domestic.links || []
      const foreignLinks = foreign.links || []
      const lines = [...domesticLinks, ...foreignLinks].map((i: any) => ({
        coords: [
          [i.sLongitude, i.sLatitude],
          [i.tLongitude, i.tLatitude],
        ],
      }))
      const scatters_domestic = (domestic.nodes || [])?.map((i: any) => ({
        value: [i.longitude, i.latitude],
        name: i.localtionName,
        itemStyle: {
          show: false,
          color: '#f42525',
        },
      }))

      const scatters_foreign = (foreign.nodes || [])?.map((i: any) => ({
        value: [i.longitude, i.latitude],
        name: i.localtionName,
        itemStyle: {
          show: false,
          color: '#f1af20',
        },
      }))
      const series = [
        {
          type: 'lines',
          zlevel: 2,
          effect: {
            show: true,
            period: 4,
            trailLength: 0.02,
            symbol: 'arrow',
            symbolSize: 20,
          },
          lineStyle: {
            width: 1,
            opacity: 1,
            curveness: 0.3,
          },
          data: [] as any[],
        },
        {
          type: 'effectScatter',
          coordinateSystem: 'geo',
          zlevel: 2,
          rippleEffect: {
            period: 4,
            brushType: 'stroke',
            scale: 4,
          },
          label: {
            position: 'right',
            offset: [5, 0],
            fontSize: 33,
            show: true,
            formatter: '{b}',
            borderWidth: 0,
            color: '#fff',
          },
          symbolSize: 20,
          symbol: 'circle',
          data: [] as any[],
        },
      ]
      series[0].data = lines
      series[1].data = [...scatters_domestic, ...scatters_foreign]
      chartOption.series = series
    },
    {
      deep: true,
    }
  )
</script>

<template>
  <div class="earth">
    <vab-chart
      :init-options="{
        renderer: 'canvas',
        useDirtyRect: false,
      }"
      :option="chartOption"
      theme="vab-echarts-theme"
    />
  </div>
  <div class="navigations">
    <div class="navigation">
      <img class="total-assets" :src="require('@/assets/dashboard_images/total-assets.png')" />
      <div @click="() => goToUrl('/assets/index', { params: 'known' })">
        <h5>资产总数</h5>
        <h4>{{ countData?.assetCount?.count ? numberFormatter.format(+countData.assetCount.count) : '0' }}</h4>
      </div>
    </div>
    <div class="navigation">
      <img class="unknown-assets" :src="require('@/assets/dashboard_images/unknown-assets.png')" />
      <div @click="() => goToUrl('/assets/index', { params: 'unKnown' })">
        <h5>未知资产</h5>
        <h4>
          {{ countData?.assetNotKnowCount?.count ? numberFormatter.format(+countData.assetNotKnowCount.count) : '0' }}
        </h4>
      </div>
    </div>
    <div class="navigation">
      <img class="total-site" :src="require('@/assets/dashboard_images/total-site.png')" />
      <div @click="() => goToUrl('/site_index', { timeDuration: '6hours' })">
        <h5>站点总数</h5>
        <h4>{{ countData?.siteCount?.count ? numberFormatter.format(+countData.siteCount.count) : '0' }}</h4>
      </div>
    </div>
    <div class="navigation">
      <img class="api-site" :src="require('@/assets/dashboard_images/api-site.png')" />
      <div @click="() => goToUrl('/site_index', { timeDuration: '6hours', activeName: 'apiList' })">
        <h5>站点API</h5>
        <h4>{{ countData?.apiCount?.count ? numberFormatter.format(+countData.apiCount.count) : '0' }}</h4>
      </div>
    </div>
    <div class="navigation">
      <img class="total-attack" :src="require('@/assets/dashboard_images/total-attack.png')" />
      <div @click="() => goToUrl('/alerts/alarm-attack', { timeDuration: '6hours' })">
        <h5>攻击总量</h5>
        <h4>
          {{ countData?.warn_attack_count?.count ? numberFormatter.format(+countData.warn_attack_count.count) : '0' }}
        </h4>
      </div>
    </div>
    <div class="navigation">
      <img class="abnormal-request" :src="require('@/assets/dashboard_images/abnormal-request.png')" />
      <div>
        <h5>异常请求</h5>
        <h4>
          {{
            countData?.warn_abonormal_count?.count ? numberFormatter.format(+countData.warn_abonormal_count.count) : '0'
          }}
        </h4>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .earth {
    background: url('@/assets/dashboard_images/earth.png') no-repeat center;
    height: 1575px;
    border: 2px solid #0666be;

    :deep() {
      .echarts {
        width: 100%;
        height: 100%;
      }
    }
  }
  .navigations {
    height: 180px;
    margin-top: 46px;
    display: flex;
    justify-content: space-between;
    .navigation {
      width: 312px;
      height: 100%;
      display: flex;
      background: url('@/assets/dashboard_images/navigation.png') no-repeat center;
      justify-content: center;
      align-items: center;
      img {
        width: 105px;
        height: 113px;
        margin-left: 30px;
      }
      & > div {
        display: flex;
        flex-direction: column;
        text-align: center;
        justify-content: center;
        flex: 1;
        height: 110px;
        cursor: pointer;
        h4,
        h5 {
          margin: 0;
        }
        h5 {
          font-size: 26px;
          font-weight: 400;
          color: #9cbfd6;
        }
        h4 {
          font-size: 44px;
          font-weight: bold;
          color: #c8e2ff;
          text-shadow: 0px 3px 10px #144da1;
        }
      }
    }
  }
</style>
