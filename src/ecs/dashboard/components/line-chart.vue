<script lang="ts">
  export default {
    name: 'DashboardLineChart',
  }
</script>

<script setup lang="ts">
  import VabChart from '@/plugins/VabChart/index.vue'
  import numberFormatter from '@/utils/number'
  const props = defineProps<{
    color: 'blue' | 'red' | 'yellow'
    data: {
      count?: number[]
      date?: string[]
    }
  }>()
  const chartColor: {
    [key in 'blue' | 'red' | 'yellow']: {
      line: string
      area: string
      end: string
    }
  } = {
    blue: {
      line: '#23aefd',
      area: '#0a3b74',
      end: '#0f172f',
    },
    red: {
      line: '#ff3367',
      area: '#62203f',
      end: '#10152c',
    },
    yellow: {
      line: '#f9ec66',
      area: '#444539',
      end: '#12172b',
    },
  }

  const chartOption = reactive({
    grid: {
      top: '4%',
      left: '7%',
      right: '5%',
      bottom: '4%',
      containLabel: true,
    },
    xAxis: [
      {
        type: 'category',

        data: [] as any[],
        show: false,
      },
      {
        type: 'category',
        show: true,
        position: 'bottom',
        boundaryGap: false,
        axisLine: {
          lineStyle: {
            opacity: 0.3,
            color: '#3271D9',
          },
        },
        axisLabel: {
          color: '#5979B2',
          align: 'left',
          fontSize: 20,
          padding: [0, 0, 0, -80],
        },
        splitLine: {
          show: false,
        },
        data: [] as any[],
      },
    ],
    yAxis: [
      {
        axisLabel: {
          formatter: function (value: number) {
            return numberFormatter.format(+value)
          },
          color: '#5979B2',
          fontSize: 18,
        },
        axisTick: {
          show: false,
        },
        type: 'value',
        axisLine: {
          show: false,
        },
        splitLine: {
          lineStyle: {
            type: 'solid',
            width: 1,
            opacity: 0.3,
            color: '#3271D9',
          },
        },
      },
    ],
    series: {
      type: 'line',
      symbol: 'circle',
      smooth: true,
      yAxisIndex: 0,
      stack: 'Total',
      lineStyle: {
        width: 4,
        color: chartColor[props.color].line,
      },
      showSymbol: false,
      areaStyle: {
        opacity: 0.72,
        color: {
          colorStops: [
            {
              offset: 0,
              color: chartColor[props.color].area,
            },
            {
              offset: 1,
              color: chartColor[props.color].end,
            },
          ],
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          type: 'linear',
          global: false,
        },
      },
      emphasis: {
        focus: 'series',
      },
      data: [] as number[],
    },
  })

  watchEffect(() => {
    const { count, date } = props.data
    chartOption.series.data = count || []
    const _data = date || []
    chartOption.xAxis[1].data = _data.length > 3 ? [_data[0], _data[Math.floor(_data.length / 2)], _data.pop()] : _data

    chartOption.xAxis[0].data = date || []
  })
</script>

<template>
  <vab-chart
    class="dashboard-container"
    :init-options="{
      renderer: 'svg',
    }"
    :option="chartOption"
    theme="vab-echarts-theme"
  />
</template>

<style scoped lang="scss"></style>
