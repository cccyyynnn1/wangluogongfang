<script lang="ts">
  export default {
    name: 'AttackType',
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  const router = useRouter()
  const props = defineProps<{
    data: {
      name: string
      value: number
    }[]
  }>()
  const chartData = ref()
  const chartOption = reactive({
    legend: {
      show: false,
    },
    grid: {
      left: '2%',
      right: '2%',
      top: '5%',
      bottom: '3%',
    },
    xAxis: [
      {
        splitLine: 'none',
        type: 'value',
        show: false,
        axisLine: 'none',
      },
    ],
    yAxis: [
      {
        type: 'category',
        show: true,
        inverse: true,
        splitLine: 'none',
        axisTick: 'none',
        axisLine: 'none',
        axisLabel: 'none',
      },
      {
        type: 'category',
        axisTick: 'none',
        axisLine: 'none',
        splitLine: 'none',
        inverse: true,
        show: true,
        axisLabel: {
          inside: true,
          verticalAlign: 'bottom',
          lineHeight: 54,
          margin: 6,
          show: true,
          color: '#A2B0B8',
          fontWeight: 400,
          fontSize: 22,
        },
        data: [] as number[],
      },
    ],
    series: [
      {
        show: true,
        name: '',
        type: 'pictorialBar',
        data: [] as {
          name: string
          value: number
        }[],
        symbolSize: [8, 18],
        symbolRepeat: 'fixed',
        symbolMargin: 2,
        symbol: 'rect',
        symbolClip: true,
        symbolPosition: 'start',
        symbolOffset: [3, 0],
        label: {
          show: true,
          offset: [5, -30],
          color: '#A2B0B8',
          fontSize: 22,
          fontWeight: 400,
          position: 'left',
          align: 'left',
          formatter: function (params: any) {
            return params.data.name
          },
        },
        itemStyle: {
          color: '#1FAEFF',
        },
        backgroundStyle: {},
        z: 2,
        animationEasing: 'elasticOut',
      },
      {
        type: 'bar',
        data: [] as number[],
        sampling: 'lttb',
        barWidth: 25,
        barGap: '-100%',
        itemStyle: {
          color: '#282e44',
          borderColor: '#1FC6FF',
          borderWidth: 2,
        },
        z: -2,
        animationEasing: 'elasticOut',
      },
    ],
    dataZoom: [
      {
        yAxisIndex: [0, 1],
        show: false,
        type: 'slider',
        startValue: 0,
        endValue: 5,
      },
    ],
  })
  const goToAlert = (type: string) => {
    const resolveRouter = router.resolve({ path: '/alerts/alarm-attack', query: { threatType: type } })
    window.open(resolveRouter.href, '_blank')
  }
  watchEffect(() => {
    const { data } = props
    chartData.value = data
    const value = data.map((i) => i.value)
    chartOption.series[0].data = data
    chartOption.series[1].data = value.map(() => Math.max(...value))
    chartOption.yAxis[1].data = value
  })
  const handle = (params: any) => {
    const dataIndex = params.dataIndex
    const type = chartData.value[dataIndex].name
    goToAlert(type)
  }
</script>

<template>
  <vab-chart
    class="attack-type-container"
    :click="handle"
    :init-options="{
      renderer: 'svg',
    }"
    :option="chartOption"
    theme="vab-echarts-theme"
  />
</template>

<style scoped lang="scss"></style>
